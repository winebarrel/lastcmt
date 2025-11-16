# lastcmt

lastcmt is a CLI that comments on a issue/pull request and minimizes previous comments.

![demo](https://github.com/user-attachments/assets/0da3da2f-f649-4563-84bf-b7c7cb7d9f82)

## Usage

```
Usage: lastcmt --owner=STRING --repo=STRING --number=INT --token=STRING [<body-file>] [flags]

Arguments:
  [<body-file>]    Comment body file. '-' is accepted for stdin.

Flags:
  -h, --help                  Show help.
  -o, --owner=STRING          Owner name ($GITHUB_OWNER).
  -r, --repo=STRING           Repo name ($GITHUB_REPO).
  -n, --number=INT            Issue/Pull Request number.
  -k, --key="lastcmt"         Identification key.
  -m, --[no-]minimize-only    Minimize only.
      --token=STRING          Auth token ($GITHUB_TOKEN).
      --version
```

## Download

https://github.com/winebarrel/lastcmt/releases/latest

## Installation

```sh
go install github.com/winebarrel/lastcmt/cmd/lastcmt@latest
```
