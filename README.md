# lastcmt

lastcmt is a CLI that comments on a pull request and minimizes previous comments.

![](https://github.com/user-attachments/assets/2a1b8bb1-dd0d-4dc1-a1d3-3923925dc0c5)

## Usage

```
Usage: lastcmt --owner=STRING --repo=STRING --number=INT --token=STRING <body-file> [flags]

Arguments:
  <body-file>    Comment body file. '-' is accepted for stdin.

Flags:
  -h, --help             Show help.
  -o, --owner=STRING     Owner name.
  -r, --repo=STRING      Repo name.
  -n, --number=INT       Pull Request number.
  -k, --key="lastcmt"    Identification key.
      --token=STRING     Auth token ($GITHUB_TOKEN).
      --version
```

## Installation

```sh
go install github.com/winebarrel/lastcmt/cmd/lastcmt@latest
```
