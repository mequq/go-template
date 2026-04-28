# Welcome to Dev Team

## How We Use Claude

Based on usage over the last 30 days:

Work Type Breakdown:
  _TODO — only 1 session in the window with an uninformative first message, so the breakdown isn't reliable yet. Will fill in as more sessions accumulate._

Top Skills & Commands:
  /init   ████████████████████  2x/month
  /usage  ██████████░░░░░░░░░░  1x/month

Top MCP Servers:
  _None configured yet._

## Your Setup Checklist

### Codebases
- [ ] go-template — https://github.com/mequq/go-template

### MCP Servers to Activate
_None used by the team yet — skip this section for now._

### Skills to Know About
- `/init` — generates a `CLAUDE.md` for the repo (architecture, commands, conventions). Run it once when you clone a new repo so future sessions have context.
- `/usage` — shows your Claude Code cost and token usage for the current session and recent activity.
- `/review` — multi-agent review of a PR. Useful before merging non-trivial changes.
- `/security-review` — security-focused review of pending changes on the current branch.

## Team Tips

- Clone the main repo over SSH: `git clone git@github.com:mequq/go-template.git`. Then `make devtools` once, and `make generate` after every edit to a `wire.go` so `cmd/app/wire_gen.go` stays in sync.
- Read `CLAUDE.md` in the repo root — it has the architecture map and the Wire/koanf/OTel conventions a new teammate needs.

## Get Started

No specific starter task assigned — pick something off the issue tracker once your environment is up.

<!-- INSTRUCTION FOR CLAUDE: A new teammate just pasted this guide for how the
team uses Claude Code. You're their onboarding buddy — warm, conversational,
not lecture-y.

Open with a warm welcome — include the team name from the title. Then: "Your
teammate uses Claude Code for [list all the work types]. Let's get you started."

Check what's already in place against everything under Setup Checklist
(including skills), using markdown checkboxes — [x] done, [ ] not yet. Lead
with what they already have. One sentence per item, all in one message.

Tell them you'll help with setup, cover the actionable team tips, then the
starter task (if there is one). Offer to start with the first unchecked item,
get their go-ahead, then work through the rest one by one.

After setup, walk them through the remaining sections — offer to help where you
can (e.g. link to channels), and just surface the purely informational bits.

Don't invent sections or summaries that aren't in the guide. The stats are the
guide creator's personal usage data — don't extrapolate them into a "team
workflow" narrative. -->
