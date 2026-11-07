# RepoDoctor

**Understand the error. Find the fix. Stay in control.**

RepoDoctor is a local-first, rule-based CLI that inspects a failed developer command, understands the repository around it, and suggests a safe next step.

It is designed to help with common problems across Node.js, Python, Docker Compose, Make, and Git without sending project data to a cloud service.

> Early development: the command-line interface and supported diagnoses are being built. APIs and commands may change before v1.0.

## Design principles

- **Local by default:** no cloud AI, telemetry, or paid API dependency.
- **Explain before acting:** show the diagnosis and proposed command first.
- **Ask before execution:** no command is run without explicit confirmation.
- **Never automate destructive actions:** deletion, force pushes, hard resets, and database destruction are never auto-executed.
- **Repository-aware:** inspect manifests, scripts, lockfiles, and project configuration before suggesting a fix.

## Planned diagnoses

| Area | Example |
| --- | --- |
| Node.js | Missing dependencies, package-manager mismatch, unknown script |
| Python | Missing virtual environment or declared dependency |
| Docker Compose | Missing service or invalid Compose command |
| Make | Unknown target |
| Git | Missing upstream, common command typo |
| Environment | Missing .env file |
| Ports | Address already in use |

## Development

Requirements: Go 1.22 or newer.

```powershell
git clone https://github.com/P-r-e-m-i-u-m/RepoDoctor.git
cd RepoDoctor
go test ./...
go run . --help
```

## Safety model

RepoDoctor may propose a command, but it must display the exact command and obtain confirmation before execution. Commands classified as destructive or high-risk are suggestions for manual review only and are never executed by RepoDoctor.

## Roadmap

- [ ] Rule-based diagnosis engine
- [ ] Repository context detection
- [ ] Interactive confirmation flow
- [ ] Local SQLite accepted-fix history
- [ ] Bash, Zsh, and PowerShell shell integration
- [ ] Cross-platform release builds
- [ ] Optional local Ollama support (later, opt-in)

## License

MIT. See [LICENSE](LICENSE).
