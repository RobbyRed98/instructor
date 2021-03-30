# instructor
A small tool to create scope based shortcuts for shell commands.

## Description

A cli-tool to create shortcuts for specific shell comands. The commands are scope based. A scope is represented by a directory. Currently the usage of the parent directory scope is not supported in a sub-directory. The tool allows to create, use, list, remove, rename and reorganize the shortcuts.

## Build 

### On Linux

#### Linux
```bash
$ go mod download
$ pandoc doc/instructor.1.md -s -t man | gzip | tee doc/instructor.1.gz > doc/ins.1.gz
$ goreleaser release --skip-publish --rm-dist --snapshot
$ go build -o ins
```

#### Windows
```bash
$ go mod download
$ pandoc doc/instructor.1.md -s -t man | gzip | tee doc/instructor.1.gz > doc/ins.1.gz
$ goreleaser release --skip-publish --rm-dist --snapshot
$ GOOS=windows GOARCH=amd64 go build -o ins.exe
```

## Usage
Executes a shortcut command.
```bash
$ ins <shortcut>
```

Creates a shortcut command which runs a specific instruction.
```bash
ins add <shortcut> <instruction>
```

Replaces the name of the old shortcut by the a new shortcut name.
```bash
ins mv <shortcut-old> <shortcut-new>
```

Removes the shortcut with the passed name in the current scope.
```bash
ins rm <shortcut>
```

Reorganizes the file in which the shortcuts and instructions are stored.
```bash
ins reorganize
```

Lists all existing shortcuts. 
```bash
**ins** list
```
