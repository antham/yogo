# Yogo [![Go Report Card](https://goreportcard.com/badge/github.com/antham/yogo)](https://goreportcard.com/report/github.com/antham/yogo) [![codecov](https://codecov.io/gh/antham/yogo/branch/master/graph/badge.svg)](https://codecov.io/gh/antham/yogo) [![GitHub tag](https://img.shields.io/github/tag/antham/yogo.svg)]() [![Go Reference](https://pkg.go.dev/badge/github.com/antham/yogo.svg)](https://pkg.go.dev/github.com/antham/yogo)

Check yopmail mails from command line.

## Install

Download binaries here : [yogo](https://github.com/antham/yogo/releases/).

Or run:
`go install github.com/antham/yogo/v4@latest`

A package is available in aur for archlinux.

## Usage

```
Check yopmail mails from command line.

Usage:
  yogo [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  inbox       Handle inbox messages
  version     App version

Flags:
      --debug   Log all requests/responses
  -h, --help    help for yogo
      --json    Dump the output as json

Use "yogo [command] --help" for more information about a command.

```

⚠️ Performing too much calls will trigger a CAPTCHA that you will need to solve through a browser. Add a delay to prevent this.

## Flag

Use the `--json` output flag to get the output as JSON.

In case of an issue with `yogo`, use the `--debug` flag to log the requests/responses.

## Inbox

### List

Retrieve 10 messages from mailbox test1@yopmail.com :

```bash
yogo inbox list test1 10
```

### Flush

Flush inbox test1@yopmail.com :

```bash
yogo inbox flush test1
```

### Read a mail

Retrieve first message from inbox helloworld@yopmail.com

```bash
yogo inbox show helloworld 1
```

Retrieve second message from inbox helloworld@yopmail.com

```bash
yogo inbox show helloworld 2
```

### Read the source of the mail with all headers

```bash
yogo inbox source helloworld 1
```

### Delete a mail

Delete first message from inbox helloworld@yopmail.com

```bash
yogo inbox delete helloworld 1
```
