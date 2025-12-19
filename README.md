# lastcmt

lastcmt is a CLI that comments on a issue/pull request and minimizes previous comments.

![demo](https://github.com/user-attachments/assets/0da3da2f-f649-4563-84bf-b7c7cb7d9f82)

## Usage

```
Usage: lastcmt --repo=REPO --token=STRING <number> [<body-file>] [flags]

Arguments:
  <number>         Issue/Pull Request number.
  [<body-file>]    Comment body file. If not specified, read from stdin.

Flags:
  -h, --help                  Show help.
  -R, --repo=REPO             OWNER/REPO ($GITHUB_REPOSITORY)
  -k, --key="lastcmt"         Identification key.
  -m, --[no-]minimize-only    Minimize only.
      --token=STRING          Auth token ($GH_TOKEN, $GITHUB_TOKEN).
      --version
```

## Download

https://github.com/winebarrel/lastcmt/releases/latest

## Installation

```sh
go install github.com/winebarrel/lastcmt/cmd/lastcmt@latest
```
