# Workspace/MyAccounts package

Requested new Gitea fork as part of Workspace installation/hosting and MyAccounts
login. Fork fahimc/gitea, I:/Projects/gitea, upstream development HEAD preserved.
Independent Windows launcher module integrations/workspace embeds unmodified
official Gitea28.1.0 and MinGit2.56.0.2 with pinned checksums, rather than falsely
claiming a fork core rebuild. build.ps1 and README document exact provenance.

Fresh NTFS data only; protected ACL owner/SYSTEM/admins and machine DPAPI recovery.
Runs native migrations, inserts OIDC source in pinned fresh SQLite schema and
uses native CLI random administrator captured in memory, never credentials in
arguments. Existing/interrupted DB with no completion marker refused. All app
state/repositories outside FileBrowser files. Runtime extracted beside launcher,
each file checked against embedded payload on start. Native OIDC at
ROOT_URL/user/oauth2/myaccounts/callback, no email auto-link/admin mapping.
HTTP Git tokens supported; SSH and mail disabled. Explicitly disable authorized
keys-file creation and set SSH root into private data: initial test exposed Gitea
writing an empty global ~/.ssh/authorized_keys even with SSH disabled. Only that
new empty test-created file was removed; existing SSH keys/config unchanged.

Build and native fresh bootstrap/bundled Git succeeded. Real isolated Workspace
MyAccounts browser login/callback and ordinary Gitea user verified; DB confirmed
external account is_admin=false. Anonymous API denied. Gitea28 API intentionally
rejects web-session cookie auth (403); tokens required, not a login failure.
Wrapper gofmt/build passed; make help unavailable (make not installed). No Gitea
core source changed. Full final installer/checksum/restart verification still
being finalized; see one-workspace agents/history/2026-10-14-git-music-compacted.md.
No existing live MyAccounts/FileBrowser installation changed, no production DNS.

Final package with SSH isolation passed full suite restart and real native callback.
No global authorized_keys file created after the corrected start. Music/Drive/portal
end-to-end checks also passed. Final installer catalog holds exact launcher SHA256.

