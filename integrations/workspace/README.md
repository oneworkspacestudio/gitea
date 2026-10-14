<img src="branding/logo.svg" alt="My Git" width="80" height="80">

# Gitea for Workspace

This fork adds a self-contained Windows launcher for M8E Workspace. Gitea's
native OpenID Connect handles MyAccounts login, account creation and stable
external subject links. HTTP Git uses Gitea's own access tokens; Workspace does
not put a portal-cookie gate in front of Git clients.

The launcher embeds the **unmodified, checksum-pinned Gitea 28.1.0 release**
and MinGit 2.56.0.2. The root fork retains upstream source; this package does
not claim to be a rebuild of the fork's development HEAD. MinGit includes its
licenses. Gitea is MIT licensed. The independent launcher module is here.

`build.ps1 -Go PATH_TO_GO` downloads and verifies both pinned payloads and builds
`dist/Gitea.exe`; Go 1.27 and Windows x64 are required for building. Workspace
embeds that EXE, so installing needs no Go, Docker, Python, Node or existing Git.

Workspace sets these environment variables automatically:

- `WORKSPACE_GITEA_ROOT_URL`: the full public URL, e.g. `https://workspace.example/git/`.
- `WORKSPACE_OIDC_ISSUER`, `WORKSPACE_OIDC_CLIENT_ID`, `WORKSPACE_OIDC_CLIENT_SECRET`.

Start with `Gitea.exe --data C:\Private\Gitea --port 8994`. Data requires NTFS;
the launcher restricts its ACL to the installing owner, SYSTEM and administrators.
For standalone use register the exact callback
`ROOT_URL/user/oauth2/myaccounts/callback` (one slash between components), set the
same environment variables, and strip the public URL prefix at your reverse
proxy. ROOT_URL can also be a loopback HTTP origin for local use.

On the **first start in a new data directory**, the launcher runs Gitea's database
migrations, provisions the `myaccounts` authentication source in the fresh SQLite
database using the pinned release's schema, and creates a random local recovery
administrator. This account is stored only in machine DPAPI at `recovery.dpapi`.
An existing database without our completion marker is refused, including after
an interrupted bootstrap. Preserve it and investigate; use fresh directories for
another test. Existing installations are never silently adopted or overwritten.

Normal MyAccounts users are normal Gitea users, with no automatic administrator
group mapping or email-based account linking. Registration requires the external
provider. Email sending and SSH are off by default; clone/push use HTTPS and a
Gitea access token. Config, repositories, database and recovery file live under
the protected data folder, outside FileBrowser's shared files. Back up the whole
data folder. Client-secret rotation requires Gitea authentication-source upkeep.

Runtime payload files are extracted beside the launcher and checked against the
embedded archive each start; modified files cause a startup failure. First start
needs writable program storage. This is an unsigned evaluation integration.
