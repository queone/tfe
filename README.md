## tfe

Terraform Cloud CLI utility

### Why?

Just a simple utility that allows:

- **Useful Functionalities**: Quickly viewing workspaces and modules registered in a particular TFE or Terraform Cloud instance. Also allows **cloning** workspaces, which is usually not possible via the Web UI.
- **Cross-Platform**: Easily compile on Linux, Mac, or Windows.
- **Simplicity**: Covers 99% of typical use cases with a minimal feature set.
- **Learning Opportunity**: A great way to practice coding in Go.

### Getting Started

To compile, clone this repository and run:
```bash
./build_go
```

### Usage

```bash
tfe v1.1.4
Terraform Cloud CLI utility
=======================
Usage
  tfe [options] [arguments]

Overview
  This utility provides basic functionality for interacting with a Terraform Enterprise or TF Cloud account instance.

Authentication
  You can authenticate using one of two methods:
  1. Set the following environment variables:
     • TF_ORG=<value>     TFE Organization name (e.g. MYORG)
     • TF_DOMAIN=<value>  TFE domain name (e.g. https://app.terraform.io)
     • TF_TOKEN=<value>   Security token to access the TFE instance
  2. Create a YAML configuration file at ~/.tfe/config.yaml with those variables:
     • TF_ORG:    <value>
     • TF_DOMAIN: <value>
     • TF_TOKEN:  <value>

Options
  The following options are available:
  -o [filter]       List organizations; filter option
  -m[j] [filter]    List only latest version of modules; filter option; JSON option
  -ma [filter]      List all versions of modules; filter option
  -w [filter]       List workspaces (100 limit); filter option
  -ws NAME          Show workspace details
  -wc SRC DEST      Clone workspace named SRC as a new workspace named DEST
  -?, -h, --help    Show this help message and exit

Examples
  tfe -o myorg
  tfe -m mymodule
  tfe -ws myworkspace
  tfe -wc source_ws_name new_ws_name
  tfe -h
```
