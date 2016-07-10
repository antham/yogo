Yogo [![Build Status](https://travis-ci.org/antham/yogo.svg?branch=master)](https://travis-ci.org/antham/yogo) [![codecov](https://codecov.io/gh/antham/yogo/branch/master/graph/badge.svg)](https://codecov.io/gh/antham/yogo) [![codebeat badge](https://codebeat.co/badges/c561682e-6834-4325-8725-6167c16214b1)](https://codebeat.co/projects/github-com-antham-yogo)
====

Check yopmail mails and inboxes from command line.

## Install

Download binaries here : [yogo](https://github.com/antham/yogo/releases/).

## Usage ##

```
Check yopmail mails and inboxes from command line.

Usage:
  yogo [command]

Available Commands:
  inbox       Handle inbox messages
  version     App version

Flags:
  -h, --help   help for yogo

Use "yogo [command] --help" for more information about a command.
```

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
