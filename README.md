# tfe

Terraform cloud utility

## Usage

```bash
tfe v1.1.0
Terraform cloud utility - https://github.com/queone/tfe
Usage: tfe [options]
  -o [filter]              List orgs, filter option
  -m[j] [filter]           List only latest version of modules, filter option; JSON option
  -ma [filter]             List all version of modules, filter option
  -w [filter]              List workspaces (100 limit), filter option
  -ws NAME                 Show workspace details
  -wc SRC DES              Clone workspace named SRC as DES
  -?, -h, --help           Print this usage page

  Note: This utility relies on below 3 critical environment variables:
    TF_ORG       TFE Organization name (MYORG, etc)
    TF_DOMAIN    TFE domain name (https://app.terraform.io, etc)
    TF_TOKEN     Security token to access the respective TFE instance

  Current values:
    TF_ORG="QUE1"
    TF_DOMAIN="https://app.terraform.io"
    TF_TOKEN="__DELIBERATELY_REDACTED__"
```
