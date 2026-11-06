---
name: providers-mint
description: >-
  Add a new inference host to the roster setup: taxonomy host id, roster-mill
  allowlists, and retrieve-limits LinkKeys. Use when the user says add provider,
  providers mint, mint host, add this host, or names a serving company (for
  example Xinference) to wire into nightly mine. MUST NOT file Perplexity output
  as prices, URLs, or capacity.
---

# Providers mint

Wire a **new host** into the roster compiler and nightly mine. This skill is a **consumer overlay** (not cursor-packs). It does not pick models (use `roster-pick`) and does not add model cards (use `model-roster`).

## When to load

Load when the user asks to add a provider or host to the setup, or pastes a public pricing URL for a serving company that is not yet a `hosts[].id`.

## Core constraints

**CONSTRAINT:** One candidate keep-gate before any allowlist write

- MUST: Ask about **exactly one** host per turn.
- MUST: WebFetch a **public** pricing or model page first (no login dashboard).
- MUST: Show proposed `id`, `runtime`, `link_key`, hostname, page URL, short cold read, then `keep? y/n`.
- MUST NOT: Batch-keep or write taxonomy / Python allowlists without `y`.
- MUST: On `n`, skip with no file change; continue with the next named host if any.
- MUST: On `y`, apply the allowlist edits in this skill, then stop unless the user also asked to attach a `hosts` row on a named card.

Enforcement: One ask → one reply → at most one new `host_id`.
Violation: STOP, revert unconfirmed writes.

CORRECT:
```text
[1/1] xinference | runtime: xinference | link_key: xinference
hostname: xinference.co
url: https://xinference.co/pricing.html
AU Model API; published $ for some models; Gemma listed on request.

keep? y/n
```

PROHIBITED:
```text
Adding Fireworks, Together, and Xinference in one write.
```

**CONSTRAINT:** Perplexity and aggregators are discovery only

- MUST: Treat Perplexity, chat recap, or a models.dev/OpenRouter provider list as a **name to fetch**, not as filed facts.
- MUST: File `links` URLs, `price_in` / `price_out`, and `limits_note` only from WebSearch / WebFetch of the host's own public page (or omit / `unknown`).
- MUST NOT: Call Perplexity MCP to fill mint fields.
- MUST NOT: Copy Perplexity prices, RPM, or URLs onto taxonomy, cards, or `LINK_KEYS`.

Enforcement: Every filed URL hostname matches the fetched page; prices quote that page or stay unset.
Violation: STOP, drop Perplexity-derived fields, re-fetch the public page.

CORRECT:
```text
User: Perplexity says Xinference hosts Gemma
→ WebFetch https://xinference.co/pricing.html
→ propose mint from that page
```

PROHIBITED:
```text
price_in: 0.30 from Perplexity summary with no host URL
```

**CONSTRAINT:** Public page gate before keep? y/n

Before proposing keep, complete this gate in working memory:

```text
Candidate: [proposed host id]
Public HTTP URL fetched this turn: YES: [exact URL] / NO
Hostname is the host's own site (not openrouter.ai, models.dev, artificialanalysis.ai, lmarena.ai, aider.chat): YES: [hostname] / NO
Login or dashboard-only: YES: [quote] / NO
Decision: PROPOSE | SKIP
If PROPOSE: id: [kebab]
            runtime: [same kebab unless user said otherwise]
            link_key: [kebab, not provider_pricing unless the page is cloud.google.com]
            hostname: [no www.]
If SKIP: reason: [one line]
```

- MUST NOT: Write the keep prompt until Decision is PROPOSE.
- MUST NOT: Reuse `provider_pricing` for a non-Google hostname (`provider_pricing` is bound to `cloud.google.com` in retrieve-limits).
- MUST NOT: Mint `cursor` or `local` (those ids already exist and are skipped by retrieve).

Enforcement: Gate slots filled; SKIP never reaches taxonomy.
Violation: STOP, complete the gate, skip or re-fetch.

CORRECT:
```text
Candidate: xinference
Public HTTP URL fetched this turn: YES: https://xinference.co/pricing.html
Hostname is the host's own site: YES: xinference.co
Login or dashboard-only: NO
Decision: PROPOSE
```

PROHIBITED:
```text
Mint host_id fireworks from an OpenRouter provider table with no fireworks.ai fetch.
```

**CONSTRAINT:** Dedup against the live allowlists

- MUST: Read `roster/taxonomy.md` Host row ids, `AllowedHostIDs` / `AllowedRuntimes` in `tools/roster-mill/internal/card`, and `LinkKeys` in `tools/roster-mill/internal/limits` before proposing.
- MUST: If the `id` already exists, STOP mint. Tell the user to use `model-roster` to append a `hosts` row on a card.
- MUST: Propose a new kebab `id` that is not in those lists.

Enforcement: Grep the three allowlists for the proposed `id` before `keep? y/n`.
Violation: STOP, do not duplicate; point at model-roster.

CORRECT:
```text
id deepinfra already in ALLOWED_HOST_IDS
→ do not mint; offer model-roster host row
```

PROHIBITED:
```text
Second runtime google_vertex beside google for the same Vertex pricing page.
```

**CONSTRAINT:** Allowlist edits on `y` (setup only)

- MUST: Add the `id` to taxonomy **Runtime values** (if sovereign/provider runtime) and **Host row ids**.
- MUST: Add the link key under **Allowed `links` keys** and a row in **Allowed enrichment sources**.
- MUST: Add the `id` to `AllowedHostIDs` and, when it is a runtime, `AllowedRuntimes` in `tools/roster-mill/internal/card`.
- MUST: Add `link_key -> hostname` to `LinkKeys` in `tools/roster-mill/internal/limits`.
- MUST: Run `make -C site test-roster` after those edits.
- MUST NOT: Set `capacity: industrial` during mint.
- MUST NOT: Append a card `hosts` row unless the user named a model slug and asked to attach this host. Then load `model-roster` for that row (`capacity: unknown` unless they filed it).
- MUST NOT: Edit `.github/workflows/roster-mine.yml` for a new host (mine already reads `LINK_KEYS`).

Enforcement: After `y`, Grep taxonomy, `AllowedHostIDs`, and `LinkKeys` for the new id; `test-roster` exits 0.
Violation: STOP, finish the three allowlists, re-run tests.

CORRECT:
```text
LinkKeys["xinference"] = "xinference.co"
AllowedHostIDs includes xinference
taxonomy Host row ids includes xinference
```

PROHIBITED:
```text
Only a README mention; mill LinkKeys still has four entries.
```

**CONSTRAINT:** Prices and capacity stay conservative

- MUST: Leave card prices unset during mint unless attaching a row and the fetched page has an explicit USD-per-million pair for that model.
- MUST: File `capacity: unknown` on any new host row. Operator later uses `roster/retrieval/capacity-decisions.yaml`.
- MUST NOT: Invent RPM, neurons, or AUD-only figures as `price_in` / `price_out`.
- MUST NOT: Treat "on request" catalog lines as a numeric price.

Enforcement: New host rows have `capacity: unknown`; numeric prices quote the fetched page.
Violation: STOP, clear invented numbers.

CORRECT:
```text
Xinference Gemma on request → no price_in on a Gemma host row
```

PROHIBITED:
```text
capacity: industrial because the marketing page said production
```

## Flow

```
- [ ] Read taxonomy Host row ids, generate-roster allowlists, retrieve-limits LINK_KEYS
- [ ] Resolve one candidate name or URL (Perplexity names are pointers only)
- [ ] WebFetch the host's public page; complete the public-page gate
- [ ] Dedup; if id exists, stop and point at model-roster
- [ ] keep? y/n; wait
- [ ] On y: taxonomy + ALLOWED_* + LINK_KEYS; make -C site test-roster
- [ ] Attach a card hosts row only if the user named a slug
```

1. **Load allowlists:** Fail if those files are missing.
2. **One candidate:** Name, URL, proposed `id` / `runtime` / `link_key` / hostname, cold read (what they serve; residency if stated; published vs on-request prices). End with `keep? y/n`.
3. **On `n`:** Skip; next named host if any.
4. **On `y`:** Edit taxonomy, mill card allowlists, and mill `LinkKeys`. Run `make -C site test-roster`.
5. **Done:** Summarize minted id vs skipped. Nightly mine will fetch that hostname on the next run.

## Pre-completion checklist

- [ ] **Allowlists read:** taxonomy, `ALLOWED_HOST_IDS`, `LINK_KEYS` were Read this run
      Method: this turn includes those Reads
      Pass: files Read before keep prompt
      Fail: STOP, Read, then propose
- [ ] **Public page fetched:** keep prompt cites a host-owned URL fetched this turn
      Method: URL hostname is not an aggregator
      Pass: gate Decision is PROPOSE
      Fail: STOP, WebFetch or skip
- [ ] **Keep gate held:** one ask → one reply → at most one new `host_id`
      Method: count writes of new ids this turn
      Pass: zero writes before `y`; at most one id after `y`
      Fail: STOP, revert extras
- [ ] **No Perplexity metadata:** no price or link taken from Perplexity
      Method: each filed URL/price traces to WebFetch
      Pass: all traces are host pages
      Fail: STOP, drop those fields
- [ ] **Three-file wire:** taxonomy + generate-roster + retrieve-limits
      Method: Grep the new id in all three
      Pass: all match
      Fail: STOP, finish the missing file
- [ ] **Tests:** `make -C site test-roster` exits 0
      Method: run after `y`
      Pass: exit 0
      Fail: STOP, fix tests
- [ ] **No silent industrial:** no `capacity: industrial` from mint
      Method: Grep the diff
      Pass: no industrial added
      Fail: STOP, set `unknown`

## Additional resources

- Taxonomy: [roster/taxonomy.md](../../roster/taxonomy.md)
- Mine spider: [tools/roster-mill/internal/limits/limits.go](../../tools/roster-mill/internal/limits/limits.go)
- Compiler allowlists: [tools/roster-mill/internal/card/card.go](../../tools/roster-mill/internal/card/card.go)
- Card host rows: skill `model-roster`
- Capacity confirm: [roster/retrieval/capacity-decisions.yaml](../../roster/retrieval/capacity-decisions.yaml)
