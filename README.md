# RepoDoctor

**Your terminal says something broke. RepoDoctor helps you figure out what.**

[![CI](https://github.com/P-r-e-m-i-u-m/RepoDoctor/actions/workflows/ci.yml/badge.svg)](https://github.com/P-r-e-m-i-u-m/RepoDoctor/actions/workflows/ci.yml)
[![Website](https://img.shields.io/badge/website-live-f04b3a)](https://p-r-e-m-i-u-m.github.io/RepoDoctor/)
[![License: MIT](https://img.shields.io/badge/license-MIT-3159d9.svg)](LICENSE)

RepoDoctor is an open-source Go CLI for the moments when a command fails and the error message leaves you guessing. Give it the command and its output. It checks the project around that failure, explains what it found, and can offer a next step.

No mystery chatbot. No uploading your repository. Just practical, local-first debugging help.

**[Visit the website](https://p-r-e-m-i-u-m.github.io/RepoDoctor/) · [Explore the code](https://github.com/P-r-e-m-i-u-m/RepoDoctor/issues) · [Report a bug](https://github.com/P-r-e-m-i-u-m/RepoDoctor/issues/new)**

> **Project status:** early development. The diagnosis engine is taking shape, but the CLI, supported cases, and command interface are still evolving. Expect rough edges. Issues and thoughtful contributions are welcome.

## The idea

A terminal error rarely tells the whole story.

`npm` says a script is missing. Is it a typo, or does this project use a different script? Python cannot import a module. Is it missing from the environment, or was it never added to the requirements? Git says there is no upstream. Which branch and remote are actually configured?

RepoDoctor connects the error to the files and settings that can explain it.

1. **Show it the failure.** Pass the command and the error output.
2. **Let it inspect the project.** It checks relevant manifests, scripts, lockfiles, and configuration.
3. **Review the diagnosis.** Get an explanation and, where appropriate, a suggested command.
4. **You decide.** RepoDoctor asks before running a suggested command. High-risk actions stay manual.

## Try it locally

You'll need **Go 1.22 or newer**.

```powershell
git clone https://github.com/P-r-e-m-i-u-m/RepoDoctor.git
cd RepoDoctor

go test ./...
go run . --help
```

### Give it a failed command

For example, if `npm run dev` reports that the `dev` script is missing:

```powershell
go run . --command "npm run dev" --error 'npm ERR! Missing script: "dev"' --dir "C:\path\to\your\project"
```

Or check a Python import error:

```powershell
go run . --command "python app.py" --error "ModuleNotFoundError: No module named 'requests'" --dir "C:\path\to\your\project"
```

Use `--dir` to point RepoDoctor at the project you want it to inspect. The command and error are required.

## What it can look for

The current rule-based engine includes checks for these situations:

| Project area | Examples |
| --- | --- |
| Node.js | Missing dependencies, unknown scripts, package-manager clues |
| Python | Missing imports, declared requirements, missing virtual environments |
| Docker Compose | Unknown services and Compose project detection |
| Git | Missing upstream configuration, possible command typos |
| Make | Unknown targets |
| Environment | Missing `.env` when an `.env.example` is available |
| Ports | Address already in use |

Coverage is still being tested and expanded. A listed check does not mean every variation of that error is supported.

## Safety is part of the feature

A debugging tool should not turn one problem into three more.

- **No surprise execution.** Suggested commands are shown before RepoDoctor asks for confirmation.
- **Risky commands stay manual.** Destructive actions are not automatically executed.
- **Local-first.** Diagnosis happens on your machine. No cloud AI or paid API is required.
- **Project-aware.** It checks relevant files instead of guessing from the error text alone.

There is also a `--yes` flag for skipping confirmation on commands RepoDoctor considers safe. Use it only when you understand what the suggested commands will do.

## Built with

- Go and [Cobra](https://github.com/spf13/cobra) for the CLI
- SQLite for local accepted-fix history
- Rule-based diagnosis, with no external AI service required

## Where the project is going

The goal is simple: make the next step after a failed command less frustrating.

- [x] Initial rule-based diagnosis engine
- [x] Repository context checks for supported project types
- [x] Confirmation before suggested command execution
- [x] Local SQLite fix history
- [ ] More real-world error patterns and regression tests
- [ ] Bash, Zsh, and PowerShell shell integration
- [ ] Cross-platform release binaries and simpler installation
- [ ] Optional local Ollama support, only if it genuinely helps

## Want to help?

RepoDoctor is young, and there is plenty of room to shape it. You do not need to be a Go expert to make a useful contribution. A good bug report, a real-world error example, or a regression test can be just as valuable as a new feature.

**A few good places to jump in:**

- [Browse the open issues](https://github.com/P-r-e-m-i-u-m/RepoDoctor/issues).
- [Pick a feature or bug to work on](https://github.com/P-r-e-m-i-u-m/RepoDoctor/issues?q=is%3Aissue+is%3Aopen).
- [Read the contribution guide](CONTRIBUTING.md) before opening a PR.
- [Report a bug or suggest an idea](https://github.com/P-r-e-m-i-u-m/RepoDoctor/issues/new/choose).

If you want to take an issue, leave a comment saying which part you plan to handle. That gives everyone a chance to coordinate before work overlaps. If an issue feels too large, suggest a smaller first slice in the comments.

We aim to review focused PRs quickly. Please include tests and run the checks listed in the contribution guide. Green checks and a clear diff make it much easier to get your work reviewed and merged.

First contribution? That's welcome too. Ask questions in the issue instead of getting stuck guessing.

## License

MIT. See [LICENSE](LICENSE).
