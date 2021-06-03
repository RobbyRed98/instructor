# instructor
Cli-tool to create scope based shortcuts for shell commands.

## Description

Cli-tool to create shortcuts for specific shell commands. The commands are scope based. A scope is represented by a directory. Currently, the usage of the parent directory scope is not supported in a sub-directory. The tool allows to create, use, list, remove, rename, edit and reorganize the shortcuts.

## Build

The instructions only refer to builds on linux systems. 
```bash
$ go mod download
$ pandoc doc/instructor.1.md -s -t man | gzip | tee doc/instructor.1.gz > doc/ins.1.gz
$ goreleaser release --skip-publish --rm-dist --snapshot
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

Replaces the name of the old shortcut by the new shortcut name.
```bash
ins mv <shortcut-old> <shortcut-new>
```

Edits the instruction of the shortcut by a replacing it with a new one.
```bash
ins edit <shortcut> <instruction-new>
```

Removes the shortcut with the passed name in the current scope.
```bash
ins rm <shortcut>
```

Reorganizes the file in which the shortcuts and instructions are stored.
```bash
ins reorganize
```

Lists existing shortcuts. 
```bash
ins list
```

### Options
The options `-b` and `--bash ` allows instructions to be executed in bash mode. 
This happens by running an instruction like `bash -c [instruction]`.
Warning this mode is **not stable** it is just experimental. 
The option can only be used in combination with a shortcut.
```bash
ins <shortcut> --bash
```
