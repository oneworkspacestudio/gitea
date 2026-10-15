# <img src="public/assets/img/logo.svg" alt="" width="40" height="40"> My Git

Gitea repository hosting with MyAccounts sign-in and a portable Windows launcher for One Workspace. Keep repositories, issues and collaboration in your own installation, with a dedicated Git app in the shared portal.

[Download Windows preview](https://github.com/oneworkspacestudio/gitea/releases) · [Launcher and standalone setup](integrations/workspace/README.md) · [One Workspace](https://github.com/oneworkspacestudio/one-workspace) · [Upstream Gitea](https://github.com/go-gitea/gitea)

## Repositories in your workspace

| Capability | How it works |
| --- | --- |
| Browser sign-in | Gitea's native OpenID Connect client uses MyAccounts |
| Workspace address | Open Git from the portal at `/git/` |
| Git clients | Clone and push over HTTP(S) with Gitea access tokens |
| Account linking | Stable provider subjects; ordinary users are not automatically administrators |
| Storage | Repositories, SQLite and configuration live in a separate private data folder |
| Branding | The My Git logo, favicon and Windows icon are bundled with the launcher |

Workspace shares identity with the other apps. Git retains its own app session and access tokens, so command-line Git traffic is handled by Gitea's normal authentication.

## Install with One Workspace

1. Visit [One Workspace](https://github.com/oneworkspacestudio/one-workspace) for the installer and setup instructions.
2. Complete the Workspace setup wizard and launch the suite.
3. Open **Git** from the portal and sign in with MyAccounts.
4. Create a Gitea access token when using an HTTP Git client.

The installer supplies the launcher, Git runtime and dedicated OIDC client configuration. No separate Docker or Git installation is required for that package.

## Run independently

Download `Gitea.exe` from this fork's [preview release](https://github.com/oneworkspacestudio/gitea/releases). Supply your MyAccounts issuer/client settings and the public Git URL in a protected launcher, then run:

```powershell
.\Gitea.exe --data 'D:\Private\MyGit' --port 8994
```

Use a fresh **NTFS** data folder and an exact registered callback. Public hosting uses an HTTPS proxy; the process itself remains bound to loopback. The [launcher guide](integrations/workspace/README.md) lists all variables, callback rules and first-start behaviour.

## What's in the Windows package

The independent launcher embeds the **unmodified, checksum-pinned Gitea 28.1.0** release and **MinGit 2.56.0.2**. It applies branding through Gitea's supported custom public-assets directory. This repository retains upstream source; the package is not presented as a rebuild of development HEAD.

On first start, the launcher verifies/extracts the runtime, runs migrations, registers MyAccounts and creates a separate recovery administrator protected by Windows DPAPI. It refuses an existing unmarked database. SSH and email sending are off by default. Back up the protected data folder and retain recovery access.

## Build locally

```powershell
cd integrations/workspace
./build.ps1 -Go 'go'
```

Requires Windows x64 and Go 1.27. The script downloads pinned upstream payloads, checks their hashes and produces `integrations/workspace/dist/Gitea.exe`. GitHub Actions is disabled for this fork.

For development of Gitea core itself, follow [development documentation](docs/development.md) and [contributing guidelines](CONTRIBUTING.md).

## Preview status and licenses

This is an unsigned evaluation integration. Current branding source can be newer than the published EXE; read release notes before deploying. Existing installations are not adopted or migrated automatically.

Gitea and the launcher retain their [MIT license](LICENSE). MinGit carries its component licenses in the runtime. Upstream documentation, authors, translations, supporters and UI references are preserved in the [archived Gitea README](docs/README-archive.md).
