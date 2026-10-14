// Copyright 2026 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package main

import (
	"archive/zip"
	"bytes"
	"crypto/sha256"
	"database/sql"
	_ "embed"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"net/url"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"regexp"
	"strings"
	"syscall"
	"time"
	"unsafe"

	"golang.org/x/sys/windows"
	_ "modernc.org/sqlite"
)

//go:embed payload.zip
var payload []byte

func main() {
	if e := run(); e != nil {
		fmt.Fprintln(os.Stderr, e)
		os.Exit(1)
	}
}

func run() error {
	data := flag.String("data", "", "Dedicated protected Gitea data directory")
	port := flag.Int("port", 8994, "Loopback HTTP port")
	flag.Parse()
	if !filepath.IsAbs(*data) || *port < 1024 || *port > 65535 {
		return errors.New("run: absolute data directory and valid port required")
	}
	exe, e := os.Executable()
	if e != nil {
		return e
	}
	runtime := filepath.Join(filepath.Dir(exe), "gitea-runtime-28.1.0")
	if e = extract(runtime); e != nil {
		return e
	}
	if e = os.MkdirAll(*data, 0700); e != nil {
		return e
	}
	if e = protectData(*data); e != nil {
		return e
	}
	git := filepath.Join(runtime, "git", "cmd")
	env := append(os.Environ(), "PATH="+git+string(os.PathListSeparator)+os.Getenv("PATH"))
	binary := filepath.Join(runtime, "gitea.exe")
	config := filepath.Join(*data, "app.ini")
	if _, e = os.Stat(config); os.IsNotExist(e) {
		if e = writeConfig(config, *data, *port, git); e != nil {
			return e
		}
	}
	command := func(args ...string) *exec.Cmd {
		args = append([]string{"--work-path", *data, "--config", config}, args...)
		c := exec.Command(binary, args...)
		c.Dir = *data
		c.Env = env
		c.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
		return c
	}
	marker := filepath.Join(*data, "workspace-bootstrap.json")
	if _, e = os.Stat(marker); os.IsNotExist(e) {
		dbPath := filepath.Join(*data, "gitea.db")
		if _, e = os.Stat(dbPath); e == nil {
			return errors.New("run: refusing to bootstrap an existing Gitea database")
		}
		if out, err := command("migrate").CombinedOutput(); err != nil {
			_ = out
			return errors.New("run: Gitea migration failed")
		}
		db, err := sql.Open("sqlite", dbPath)
		if err != nil {
			return err
		}
		defer db.Close()
		cfg := map[string]any{"Provider": "openidConnect", "ClientID": os.Getenv("WORKSPACE_OIDC_CLIENT_ID"), "ClientSecret": os.Getenv("WORKSPACE_OIDC_CLIENT_SECRET"), "OpenIDConnectAutoDiscoveryURL": strings.TrimRight(os.Getenv("WORKSPACE_OIDC_ISSUER"), "/") + "/.well-known/openid-configuration", "Scopes": []string{"openid", "profile", "email"}, "FullNameClaimName": "name", "ExternalIDClaim": "sub"}
		blob, _ := json.Marshal(cfg)
		now := time.Now().Unix()
		if _, err = db.Exec("INSERT INTO login_source (type,name,is_active,is_sync_enabled,two_factor_policy,cfg,created_unix,updated_unix) VALUES (6,'myaccounts',1,0,'',?,?,?)", string(blob), now, now); err != nil {
			return fmt.Errorf("run: provision authentication source: %w", err)
		}
		// Capture the generated recovery credential in memory, then protect it with Windows DPAPI.
		out, err := command("admin", "user", "create", "--username", "workspace-recovery", "--email", "workspace-recovery@localhost.invalid", "--admin", "--random-password", "--random-password-length", "48", "--must-change-password=false").CombinedOutput()
		if err != nil {
			return errors.New("run: Gitea recovery user creation failed")
		}
		match := regexp.MustCompile(`generated random password is '([^']+)'`).FindSubmatch(out)
		if len(match) != 2 {
			return errors.New("run: Gitea recovery handoff unavailable")
		}
		if err = protect(filepath.Join(*data, "recovery.dpapi"), match[1]); err != nil {
			return err
		}
		if e = os.WriteFile(marker, []byte("{\"schema\":1,\"gitea\":\"28.1.0\"}\n"), 0600); e != nil {
			return e
		}
	}
	child := command("web")
	child.Stdout = os.Stdout
	child.Stderr = os.Stderr
	return child.Run()
}

func writeConfig(path, data string, port int, git string) error {
	origin := strings.TrimRight(os.Getenv("WORKSPACE_GITEA_ROOT_URL"), "/") + "/"
	issuer := os.Getenv("WORKSPACE_OIDC_ISSUER")
	u, e := url.Parse(origin)
	if e != nil || u.Hostname() == "" || strings.ContainsAny(origin+data, "\r\n") {
		return errors.New("writeConfig: invalid origin or data path")
	}
	if u.Scheme != "https" && !(u.Scheme == "http" && (u.Hostname() == "localhost" || u.Hostname() == "127.0.0.1")) {
		return errors.New("writeConfig: HTTPS or loopback required")
	}
	if issuer == "" || os.Getenv("WORKSPACE_OIDC_CLIENT_ID") == "" || os.Getenv("WORKSPACE_OIDC_CLIENT_SECRET") == "" {
		return errors.New("writeConfig: MyAccounts client configuration required")
	}
	slash := func(s string) string { return filepath.ToSlash(s) }
	body := fmt.Sprintf(`APP_NAME = Workspace Git
RUN_MODE = prod
[server]
HTTP_ADDR = 127.0.0.1
HTTP_PORT = %d
ROOT_URL = %s
APP_DATA_PATH = %s/data
DISABLE_SSH = true
SSH_CREATE_AUTHORIZED_KEYS_FILE = false
SSH_ROOT_PATH = %s/ssh
OFFLINE_MODE = true
[database]
DB_TYPE = sqlite3
PATH = %s/gitea.db
[repository]
ROOT = %s/repositories
[security]
INSTALL_LOCK = true
REVERSE_PROXY_LIMIT = 1
REVERSE_PROXY_TRUSTED_PROXIES = 127.0.0.0/8,::1/128
[service]
DISABLE_REGISTRATION = false
ALLOW_ONLY_EXTERNAL_REGISTRATION = true
SHOW_REGISTRATION_BUTTON = false
REQUIRE_SIGNIN_VIEW = true
REGISTER_EMAIL_CONFIRM = false
ENABLE_NOTIFY_MAIL = false
DEFAULT_ALLOW_CREATE_ORGANIZATION = true
[oauth2_client]
ENABLE_AUTO_REGISTRATION = true
USERNAME = userid
ACCOUNT_LINKING = disabled
REGISTER_EMAIL_CONFIRM = false
[mailer]
ENABLED = false
[session]
COOKIE_NAME = workspace_gitea
COOKIE_SECURE = %t
[git]
PATH = %s/git.exe
[log]
MODE = console
LEVEL = Warn
`, port, origin, slash(data), slash(data), slash(data), slash(data), u.Scheme == "https", slash(git))
	return os.WriteFile(path, []byte(body), 0600)
}

func extract(dest string) error {
	z, e := zip.NewReader(bytes.NewReader(payload), int64(len(payload)))
	if e != nil {
		return e
	}
	for _, entry := range z.File {
		rel := filepath.FromSlash(entry.Name)
		if filepath.IsAbs(rel) || strings.HasPrefix(filepath.Clean(rel), "..") || entry.Mode()&os.ModeSymlink != 0 {
			return errors.New("extract: invalid embedded runtime path")
		}
		target := filepath.Join(dest, rel)
		if entry.FileInfo().IsDir() {
			continue
		}
		in, e := entry.Open()
		if e != nil {
			return e
		}
		content, e := io.ReadAll(in)
		in.Close()
		if e != nil {
			return e
		}
		if existing, err := os.ReadFile(target); err == nil {
			if sha256.Sum256(existing) != sha256.Sum256(content) {
				return errors.New("extract: runtime integrity mismatch; restore the application package")
			}
			continue
		}
		if e = os.MkdirAll(filepath.Dir(target), 0700); e != nil {
			return e
		}
		f, e := os.OpenFile(target, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0700)
		if e != nil {
			return e
		}
		_, e = f.Write(content)
		f.Close()
		if e != nil {
			return e
		}
	}
	return nil
}

func protect(path string, value []byte) error {
	in := windows.DataBlob{Size: uint32(len(value)), Data: &value[0]}
	var out windows.DataBlob
	if e := windows.CryptProtectData(&in, nil, nil, 0, nil, windows.CRYPTPROTECT_UI_FORBIDDEN|windows.CRYPTPROTECT_LOCAL_MACHINE, &out); e != nil {
		return e
	}
	defer windows.LocalFree(windows.Handle(unsafe.Pointer(out.Data)))
	return os.WriteFile(path, unsafe.Slice(out.Data, int(out.Size)), 0600)
}

func protectData(dir string) error {
	path, e := windows.UTF16PtrFromString(dir)
	if e != nil {
		return e
	}
	root := make([]uint16, 32768)
	if e = windows.GetVolumePathName(path, &root[0], uint32(len(root))); e != nil {
		return e
	}
	var flags uint32
	if e = windows.GetVolumeInformation(&root[0], nil, 0, nil, nil, &flags, nil, 0); e != nil {
		return e
	}
	if flags&windows.FILE_PERSISTENT_ACLS == 0 {
		return errors.New("protectData: choose an NTFS data folder")
	}
	current, e := user.Current()
	if e != nil {
		return e
	}
	sid := current.Uid
	marker := filepath.Join(dir, ".owner-sid")
	if b, err := os.ReadFile(marker); err == nil {
		sid = strings.TrimSpace(string(b))
	} else if b, err = os.ReadFile(filepath.Join(filepath.Dir(dir), ".owner-sid")); err == nil {
		sid = strings.TrimSpace(string(b))
	}
	if _, e = windows.StringToSid(sid); e != nil {
		return e
	}
	sd, e := windows.SecurityDescriptorFromString("D:P(A;OICI;FA;;;" + sid + ")(A;OICI;FA;;;SY)(A;OICI;FA;;;BA)")
	if e != nil {
		return e
	}
	acl, _, e := sd.DACL()
	if e != nil {
		return e
	}
	if e = windows.SetNamedSecurityInfo(dir, windows.SE_FILE_OBJECT, windows.DACL_SECURITY_INFORMATION|windows.PROTECTED_DACL_SECURITY_INFORMATION, nil, nil, acl, nil); e != nil {
		return e
	}
	return os.WriteFile(marker, []byte(sid), 0600)
}
