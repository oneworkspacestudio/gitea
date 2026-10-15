# Packaged Windows Git hooks and display fixes

Wrapper prepends packaged git/cmd and git/usr/bin; SCRIPT_TYPE=sh ensures new
hooks use the bundled shell. Trusted sh.exe is also supplied as bash.exe for old
hooks. This avoids falling through to Windows WSL bash when MinGit has no bash.
Scoped INI helper preserves unrelated settings; idempotence/preservation tests
pass. Official pinned Gitea/MinGit payload provenance unchanged.

UI DEFAULT_SHOW_FULL_NAME and pinned28.1 template overrides present friendly
profile/owner labels. Canonical usernames, repository URLs and ownership remain
unchanged. Browser initialized a private synthetic QA repository, then committed
README through the editor. Bare-repo history confirms both commits. Dashboard
friendly name verified after final update. Clean-host portability not tested.
Local wrapper rebuilt/applied; no publication.
