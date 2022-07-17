# Digit

## Concepts

Here is our report of Digit: https://docs.google.com/document/d/1g7YeATxRtSnpZQw8roqUBm9e1ZnmsCmPhwIqURbqOFA/edit?usp=sharing

## Installation

This is an example of how you may give instructions on setting up your project locally. To get a local copy up and running follow these simple example steps.

Before started: Make sure you have Go installed, and that go is in your path.

### Clone this repository

```
git clone https://github.com/SinHanYang/Digit.git
cd Digit
```

### Run installation

```
# for Digit CLI
go install

# for Digit Core
go install ./digit-server
```

## Usage

### Start Digit Server

```
digit-server

# Or to start in background:
digit-server &
```

### Start playing with Digit!

```
$ digit -h
Digit is a super fancy version control system for your database.

Usage:
  digit [flags]
  digit [command]

Available Commands:
  add         Stage table
  commit      commit
  diff        Show the changes between two commits
  init        Initialize a empty Digit data repository
  log         show commit log
  reset       Reset the database to specified commit
  sql         Runs a SQL query

Flags:
  -h, --help   help for digit

Use "digit [command] --help" for more information about a command.
```
