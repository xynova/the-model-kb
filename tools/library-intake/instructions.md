# Library intake CLI

## Project Overview

`library-intake` prepares YouTube AI roundups for filing into the models repo
`library/` catalog. It orchestrates installed tools (`fabric`, `yt-dlp`,
`ffmpeg`) and Polypus OpenAI-compatible chat for candidate proposal (GLM Flash)
and frame curation (Gemma vision). It does not write library entries or invent
project links; WebSearch enrichment stays in the `ai-library-intake` skill after
each keep decision.

## Core Functionalities

1. `transcript` — grab YouTube transcript via fabric (no chat pattern).
2. `candidates` — propose ordered keep candidates from a transcript via Polypus.
3. `frames` — download video and extract sparse JPEG frames via yt-dlp + ffmpeg.
4. `curate` — map frames to candidates with Gemma vision via Polypus.
5. `prepare` — run transcript + candidates, then soft-fail frames + curation.
6. `index refresh` — regenerate `library/INDEX.md`, each `library/<category>/INDEX.md`, and `library/additions/` week files (section tables per category; open/closed from `added` dates).

## Docs and Libraries

- Polypus gateway: `http://127.0.0.1:1320` (`POLYPUS_BASE_URL`)
- Default chat model: `cf_local/@cf/zai-org/glm-4.7-flash` (`LIBRARY_INTAKE_POLYPUS_MODEL`)
- Default vision model: `cf_local/@cf/google/gemma-4-26b-a4b-it` (`LIBRARY_INTAKE_VISION_MODEL`)
- Skill overlay: `../../.cursor/skills/ai-library-intake/SKILL.md`

## Current File Structure

```text
tools/library-intake/
  cmd/library-intake/main.go
  internal/cli/
  internal/config/
  internal/errors/
  internal/clients/polypus/
  internal/execx/
  internal/transcript/
  internal/frames/
  internal/candidates/
  internal/workflow/
  internal/youtubeid/
  Makefile
  instructions.md
  go.mod
```

## Build and run

```bash
cd tools/library-intake
make build
./bin/library-intake prepare --url 'https://www.youtube.com/watch?v=VIDEO_ID'
```

## Frame download notes

YouTube often returns HTTP 403 for default `android_vr` / SABR-only formats.
`frames` / `prepare` pass `--extractor-args youtube:player_client=android,web`
by default and prefer progressive `mp4` under 720p. Override with
`LIBRARY_INTAKE_YTDLP_EXTRACTOR_ARGS` when needed. Soft-fail messages include
yt-dlp/ffmpeg stderr fields from the typed error.
