<div align="center">

<!-- logo -->
<img src=".github/logo.png" alt="ChainRun" width="120" />

# ChainRun

**A lightweight, YAML-driven workflow automation engine for the command line.**

<!-- placeholder -->
![Build](https://img.shields.io/github/actions/workflow/status/vikingpingvin/chainrun/ci.yml?branch=main&label=build)
![Go Version](https://img.shields.io/badge/go-1.24-blue)
![License](https://img.shields.io/badge/license-MIT-green)
![Release](https://img.shields.io/github/v/release/vikingpingvin/chainrun)

</div>

---
## About   

ChainRun lets you define multi-step automation workflows in a single YAML file and run them from any terminal — on demand, on a schedule, or in response to external events. It ships as a single static binary with no runtime dependencies, making it trivial to drop into CI pipelines, developer machines, or server cron jobs.

The engine evaluates each step in sequence, passes outputs forward via Go `text/template` expressions, and enforces per-step timeouts and retry policies so your workflows fail fast and recover cleanly.

### Core Use Case

ChainRun is a local, single-binary workflow automation engine.   
The primary use case is: "I want n8n/Zapier-style step   chaining, but as a CLI tool I can drop anywhere without cloud accounts, Docker, or a web UI."  

**Strongest concrete targets:**
  - Personal developer automation (daily summaries, repo digests, build reports)
  - CI side-effects on a server/VPS without a full CI system
  - LLM-augmented scripting without writing custom Go/Python glue code
  - Small team shared workflows via a config file in a repo
---

## Features

- YAML-defined workflows — one file per project or one file for everything
- Multiple trigger types: cron, manual, webhook, watch, git
- Multiple step types: shell, HTTP, LLM, file, notify
- Go `text/template` support in all string fields (env vars, prior step output, trigger metadata)
- Per-step timeout, retry with configurable backoff, and `continue_on_error`
- Workflow-level error policies: stop, continue, or notify
- Single static binary, no runtime dependencies

---

## Quick Start

```bash
# Install (placeholder — module path not yet published)
go install github.com/vikingpingvin/chainrun@latest
```

Create a minimal workflow file:

```yaml
# chainrun.yaml
workflows:
  - name: hello
    trigger:
      type: manual
    steps:
      - id: greet
        type: shell
        command: echo "Hello from ChainRun"
        timeout: 10s
```

Run it:

```bash
chainrun run hello
```

---

## CLI Reference

| Command | Description |
|---|---|
| `chainrun run [name]` | Execute a named workflow once |
| `chainrun daemon` | Start all trigger sources and keep running |
| `chainrun validate [file]` | Parse and validate a config file |
| `chainrun list` | List all loaded workflows |

**Persistent flag:** `--config <path>` — path to the YAML config file (default: `chainrun.yaml`)

---

## Configuration Reference

TODO

---

### Common step fields

| Field | Type | Default | Description |
|---|---|---|---|
| `id` | string | — | Unique step identifier (required) |
| `type` | string | — | Step executor type (required) |
| `timeout` | duration | `30s` | Maximum run time for this step |
| `retry` | object | — | Retry policy: `attempts`, `delay`, `backoff` |
| `continue_on_error` | bool | `false` | Proceed to the next step even if this one fails |
| `env` | map | — | Step-scoped environment variables |

---

### Templating

All string fields in step definitions are rendered as Go `text/template` expressions before the step runs. The following data is available:

| Expression | Description |
|---|---|
| `{{ .Env.VAR_NAME }}` | Environment variable from workflow, step, or OS env |
| `{{ .Steps.<id>.Stdout }}` | Captured stdout of a previously completed step |
| `{{ .Steps.<id>.Stderr }}` | Captured stderr of a previously completed step |
| `{{ .Steps.<id>.ExitCode }}` | Exit code of a previously completed step |
| `{{ .Trigger.Type }}` | Trigger type that fired this run |
| `{{ .Trigger.Metadata }}` | Trigger-specific metadata (webhook payload, etc.) |

---

## Full Example Workflow

A daily stand-up digest: pull a git log, summarise it with an LLM, and post to Slack.

```yaml
workflows:
  - name: daily-standup
    description: Summarise yesterday's commits and post to Slack
    trigger:
      type: cron
      schedule: "0 9 * * 1-5"
    env:
      REPO_PATH: /home/user/projects/myapp
    on_error: notify
    steps:
      - id: git_log
        type: shell
        command: git -C {{ .Env.REPO_PATH }} log --oneline --since="24 hours ago"
        timeout: 10s

      - id: summarise
        type: llm
        provider: anthropic
        model: claude-opus-4-6
        system: You are a concise engineering lead. Summarise git commits into a short stand-up update.
        prompt: |
          Commits from the last 24 hours:
          {{ .Steps.git_log.Stdout }}
        timeout: 30s

      - id: post_slack
        type: notify
        channel: slack
        to: "#standup"
        message: |
          *Daily Stand-up Digest*
          {{ .Steps.summarise.Stdout }}
        timeout: 10s
```

---

## Contributing

Contributions are welcome. Fork the repository, create a feature branch, and open a pull request against `main`. Please read [CONTRIBUTING.md](CONTRIBUTING.md) before submitting.

---

## License

<!-- LICENSE placeholder -->

MIT License. See [LICENSE](LICENSE) for details.
