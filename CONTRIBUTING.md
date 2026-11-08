# Contributing to RepoDoctor

Thanks for taking a look. RepoDoctor is early, so useful contributions can have a real impact.

You do not need to build a huge feature to help. A reproducible bug, a realistic error fixture, a focused test, or a clearer diagnosis is worth sending.

## Before you start

1. Browse the [open issues](https://github.com/P-r-e-m-i-u-m/RepoDoctor/issues).
2. If you want to work on an issue, leave a short comment saying what you plan to change.
3. For a larger change, wait for a maintainer reply before spending time on it. We want to avoid two people solving the same thing.
4. If the issue is broad, suggest a smaller first step. A focused PR is easier to review.

If you have found a bug that does not have an issue yet, open one with the command you ran, the relevant error output, what you expected, and what actually happened. Remove secrets, tokens, and private project data before posting logs.

## Local setup

You need Go 1.22 or newer.

```powershell
git clone https://github.com/P-r-e-m-i-u-m/RepoDoctor.git
cd RepoDoctor
go test ./...
go vet ./...
go build ./...
```

## Pull requests

- Keep each PR focused on one issue or closely related change.
- Explain the problem first, then describe the approach.
- Add or update tests for behavior changes. Prefer table-driven tests and temporary project fixtures where they make sense.
- Do not add network calls, cloud AI requirements, or paid services to the diagnosis path without discussing the design first.
- Do not make RepoDoctor execute commands during diagnosis. Suggested commands must remain subject to the existing safety and confirmation behavior.
- Update the README when user-facing behavior or supported cases change.
- Include screenshots only when a change affects the website or another visual interface.

Before opening the PR, run:

```powershell
go test ./...
go vet ./...
go build ./...
```

If one of these cannot pass, say why in the PR and include the output. Please do not hide failing checks.

## Review and communication

We want to keep the PR queue moving. The maintainer aims to review focused contributions quickly and merge them when the checks and implementation look good. That is a goal, not a promise of automatic approval: correctness, safety, and maintainability still come first.

Be kind, explain trade-offs, and ask when something is unclear. Thanks for helping make RepoDoctor more useful.
