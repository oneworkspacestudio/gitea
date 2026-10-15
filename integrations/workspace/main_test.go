// Copyright 2026 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT
package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestConfigureHookShellPreservesRepositoryAndOtherSettings(t *testing.T) {
	for _, input := range []string{"[repository]\nROOT = C:/private/repos\nSCRIPT_TYPE = bash\n[other]\nSCRIPT_TYPE = untouched\n", "[repository]\nROOT = C:/private/repos\n", "[server]\nHTTP_PORT = 8994\n"} {
		path := filepath.Join(t.TempDir(), "app.ini")
		if err := os.WriteFile(path, []byte(input), 0600); err != nil {
			t.Fatal(err)
		}
		if err := configureHookShell(path); err != nil {
			t.Fatal(err)
		}
		first, _ := os.ReadFile(path)
		if !strings.Contains(string(first), "SCRIPT_TYPE = sh") {
			t.Fatal("packaged interpreter was not selected")
		}
		if strings.Contains(input, "ROOT") && !strings.Contains(string(first), "ROOT = C:/private/repos") {
			t.Fatal("repository path changed")
		}
		if strings.Contains(input, "untouched") && !strings.Contains(string(first), "SCRIPT_TYPE = untouched") {
			t.Fatal("unrelated section changed")
		}
		if err := configureHookShell(path); err != nil {
			t.Fatal(err)
		}
		second, _ := os.ReadFile(path)
		if string(first) != string(second) {
			t.Fatal("configuration update is not idempotent")
		}
	}
}
