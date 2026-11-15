# Migrations

This directory contains immutable `golang-migrate` schema versions. Every shared schema change has paired `up` and `down` SQL files; application queries do not belong here.

## Rules
- File names: `NNNNNN_description.up.sql` and `NNNNNN_description.down.sql` (e.g. `000002_create_resources.up.sql`).
- Never edit a migration that has already been merged; add a new numbered one instead.
- `down` must exactly undo `up`.
- Tables must match the GORM row structs in `internal/repository/postgres/`.

## Current versions
| Version | What it does |
|---|---|
| 000001 | Creates the `incidents` table |
