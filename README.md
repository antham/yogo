# Yogo [![CircleCI](https://circleci.com/gh/antham/yogo.svg?style=svg)](https://circleci.com/gh/antham/yogo) [![codecov](https://codecov.io/gh/antham/yogo/branch/master/graph/badge.svg)](https://codecov.io/gh/antham/yogo)

Check yopmail mails and inboxes from command line.

## Install

Download binaries here : [yogo](https://github.com/antham/yogo/releases/).

A package is available in aur for archlinux.

## Usage

```
Check yopmail mails and inboxes from command line.

Usage:
  yogo [command]

Available Commands:
  help        Help about any command
  inbox       Handle inbox messages
  version     App version

Flags:
  -h, --help   help for yogo
      --json   Dump the output as json

Use "yogo [command] --help" for more information about a command.
```

## Flag

Use the `--json` output flag to get the output as JSON.

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

### Delete a mail

Delete first message from inbox helloworld@yopmail.com

```bash
yogo inbox delete helloworld 1
```
