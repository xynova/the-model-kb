---
name: roster-pick
description: >-
  Rank Model Picks host offers for a user task via MCP roster_pick or the
  roster-pick CLI. Use when the user asks which model or host to use for a job,
  fit, cost, or industrial limits. MUST NOT grep roster markdown to invent a pick.
---

# Roster pick

Agents choose models from scored data, not from reading every card.

## When to load

Load when the user asks which model, host, or offer to use for a task, budget, or capacity (for example industrial batch vs hobby tier).

## Core constraints

**CONSTRAINT:** Call MCP `roster_pick` with the user's task string.

- MUST: Pass the full user task as the `task` argument.
- MUST: Surface `matched_jobs`, `matched_capacity`, and the top offer (`name`, `host_id`, `cost_per_m`).
- MAY: Mention that offers in `scores.json` also carry `quality` and `capability_cohort` for inspect; pick sort does not use them in v1.
- MUST NOT: Open `roster/*.md` or `catalog.json` to invent a pick.
- MUST NOT: Run `make -C dbt build` or edit roster facts during a pick request.

Enforcement: The pick path is MCP or CLI JSON only.
Violation: STOP, call `roster_pick`, show structured results.

**CONSTRAINT:** Fallback when MCP is unwired.

- MAY: Run `go run ./tools/roster-pick/cmd/roster-pick pick --task "..." --json` from the repo root when `roster-pick` MCP is missing.
- MUST: Tell the user MCP is unwired when using the fallback.

CORRECT:

```text
User: industrial cheap vision on DeepInfra
→ roster_pick task="..."
→ Top offer: gemma-4-26b-a4b @ deepinfra
```

PROHIBITED:

```text
User: pick a model
→ rg 'vision' roster/
→ guess from markdown
```

## Operator refresh

Operators rebuild scores after catalog changes:

```bash
make -C site generate-roster
make -C dbt build
```

See [roster/README.md](../../roster/README.md) and [tools/roster-pick/README.md](../../tools/roster-pick/README.md).
