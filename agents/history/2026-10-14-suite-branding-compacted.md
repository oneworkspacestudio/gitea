# Suite branding

New my_git artwork updates assets/logo.svg, assets/favicon.svg and public/assets/img
logo/favicon/PNG/touch fallback. Independent integrations/workspace launcher embeds
branding/* and installs to protected data/workspace-branding/public/assets/img, passing
explicit --custom-path. This is Gitea's supported asset overlay; bundled upstream core
remains unmodified. Existing identity/database/SSH isolation is unchanged.
Windows amd64 launcher resource embeds the same generated ICO. Go source gofmt applied;
native launcher vet/build passed. Output: sibling one-workspace/dist/branding/Gitea.exe. No Gitea core Go/API/locale changes,
so core make fmt/tidy/swagger steps do not apply; make unavailable on this host.
No existing Gitea data/runtime changed and no new source/release publication.
Canonical artwork and sync/resource scripts live in sibling one-workspace.

## Source publication follow-up

User requested pushing the affected repositories and adding logos to their main
READMEs. This supersedes the earlier local-only publication scope. All five main
READMEs now reference their local SVG assets. Drive publication includes the
previous Photos/Studio prerequisite source so the branded app imports build from
the pushed branch; unrelated preview files and concurrent Workspace Mail changes
remain excluded. Existing executable releases and live installations are unchanged.
GitHub Actions is disabled before publication; builds remain local.
