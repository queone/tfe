# tfe Architecture

## Purpose

tfe answers everyday Terraform Cloud and Terraform Enterprise questions from the shell, and clones a workspace, which the web interface cannot do.

## System Summary

tfe reads `TF_ORG`, `TF_DOMAIN`, and `TF_TOKEN`, calls the Terraform Cloud API through the official `go-tfe` client, and prints lists or details. `clone` creates a workspace and its variables.

## Current Platform

- Go

## Major Components

- `cmd/tfe/main.go`: help and command dispatch
- `cmd/tfe/config.go`: credentials from the environment or the config file
- `cmd/tfe/api.go`: the `go-tfe` client and list paging
- `cmd/tfe/orgs.go`, `modules.go`, `workspaces.go`: the commands

## Core Files

- `AGENTS.md`: base governance contract
- `plan.md`: prioritized roadmap and approved direction
- `build.sh`: self-contained build / release-prep / release script (Bash 3.2+, no external tools)
- `govna/development-cycle.md`: workflow from roadmap through release
- `govna/ac-template.md`: acceptance-criteria template for new work
- `govna/build-release.md`: build, test, and release rules

## Data And Control Flow

The command line selects a command. tfe then resolves credentials, builds the client, reads the API page by page, filters by the optional substring, and prints the result. `clone` reads the source workspace and its variables, then creates the destination workspace and copies each variable.

## AC Lifecycle Control Flow

The governed change path is `Draft → Audit → Refine → Implement → Ratify → Package`. Draft creates the AC; Audit, Refine, Implement, and Ratify are the four AC phases; Package is post-Ratify release preparation and is not a fifth phase.

Integrated audit adoption is the only command-mediated phase exception. It can advance one emitted adoption AC through immediate Audit and no-edit Refine, but it cannot enter Implement. Every unpackaged AC with implementation in the unreleased state enters the pending release batch, including work awaiting Ratify. A private pre-Implement calculation prevents that complete batch from growing beyond one 80-byte prefix-plus-summary message. Package requires every member to be Ratified, rejects excluded implemented work, and rechecks the complete batch before prep. A named request such as `Package AC70+AC71` establishes a fitting multi-AC batch; a standalone Package alias reuses the complete batch already established in the active session.

## Architecture Notes

- record stable system decisions here
- prefer durable structure and interfaces over transient implementation detail
- `clone` is the only call that changes live state.
- Credentials come from the environment when all three are set, otherwise from `$XDG_CONFIG_HOME/tfe/config.yaml` or `~/.config/tfe/config.yaml`, which is created as a 0600 skeleton when missing. The token is never printed.
- The API never returns sensitive values, so `clone` creates sensitive variables with empty values and prints their names.
- Help renders through `github.com/queone/gkit/help`, and colors come from `github.com/queone/gkit/color`.

## Conventions

- update this document when architecture or major workflow changes materially
- keep implementation detail in code and stable architecture here
