# roster-mill

Go CLI that compiles roster cards and retrieves host limits plus Artificial
Analysis scores. Replaces the former `site/scripts/*.py` mill.

```text
roster-mill generate
roster-mill retrieve-limits
roster-mill apply-mine
roster-mill retrieve-aa
roster-mill apply-aa
roster-mill polypus-diff
roster-mill init
roster-mill store-aa-key
roster-mill review-body
roster-mill version
```

From the repo root:

```text
go run ./tools/roster-mill/cmd/roster-mill generate
make -C site generate-roster
make -C site test-roster
```

Allowlists live in `internal/card` (`AllowedHostIDs`, `AllowedRuntimes`) and
`internal/limits.LinkKeys`. Keep them in sync with `roster/taxonomy.md`.
