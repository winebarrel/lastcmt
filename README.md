# lastcmt

lastcmt is a CLI that comments on the PR and minimizes previous comments.

![](https://github.com/user-attachments/assets/88dfcef4-e52d-4ac7-9619-174dcbbf32c6)

## Usage

```
Usage: lastcmt --owner=STRING --repo=STRING --number=INT --token=STRING <body-file> [flags]

Arguments:
  <body-file>    Comment body file. '-' is accepted for stdin.

Flags:
  -h, --help             Show help.
  -o, --owner=STRING     Owner name.
  -r, --repo=STRING      Repo name.
  -n, --number=INT       Issues/PR number.
  -k, --key="lastcmt"    Identification key.
      --token=STRING     Auth token ($GITHUB_TOKEN).
      --left=0           Number of comments not minimized.
      --version
```

## Installation

```sh
go install github.com/winebarrel/lastcmt/cmd/lastcmt@latest
```
