# Yogo #

Interact with yopmail from command-line

## Install

Download binary here : [yogo](https://github.com/antham/yogo/releases/tag/v1.0.3).

## Usage ##

```
usage: yogo [<flags>] <command> [<args> ...]

Interact with yopmail from command line

Flags:
--help  Show help (also see --help-long and --help-man).

Commands:
help [<command>...]
Show help.

mailbox [<flags>] <mail> [<action>]
Manage mailbox

mail <mail> [<position>] [<action>]
Manage mail

```

## Mailbox command

### Read

Display first sum up from mailbox test1@yopmail.com :

```bash
yogo mailbox test1
```

Retrieve 10 messages from mailbox test1@yopmail.com :

```bash
yogo mailbox test1 --limit 10
```

### Flush

Flush mailbox test1@yopmail.com :

```bash
yogo mailbox test1 flush
```

## Mail command

### Read

Retrieve first message from mailbox helloworld@yopmail.com

```bash
yogo mail helloworld
```

Retrieve second message from mailbox helloworld@yopmail.com

```bash
yogo mail helloworld 2
```

### Delete

Delete first message from mailbox helloworld@yopmail.com

```bash
yogo mail helloworld 1 delete
```
