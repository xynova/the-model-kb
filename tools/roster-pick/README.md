# roster-pick

Rank Model Picks host offers for a natural-language task using committed `roster/scores.json`.

## Agents

- Read [AGENTS.md](AGENTS.md) and [.cursor/skills/roster-pick/SKILL.md](../../.cursor/skills/roster-pick/SKILL.md).
- Prefer MCP tool `roster_pick` on stdio (`roster-pick mcp`).

## Scores refresh

After roster markdown changes:

```bash
make -C site generate-roster
make -C dbt build
```

## Build and test

From this directory:

```bash
go tool task build
go tool task test
```

Flag details: `bin/roster-pick help` or `go run ./cmd/roster-pick help`.

## MCP launcher (local)

```bash
go tool task build
/path/to/the-model-kb/tools/roster-pick/bin/roster-pick mcp
```

Wire that command in your Cursor MCP settings; do not commit secrets.
