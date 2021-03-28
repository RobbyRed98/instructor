% instructor(1) 0.1.0
% RobbyRed98
% March 2021

# NAME
instructor - a small cli-tool to add scope bound shortcuts for regular tasks/commands. 

# SYNOPSIS

**ins** <shortcut>

**ins** <command> <args>

# DESCRIPTION

A cli-tool to create shortcuts for specific shell comands. The commands are scope based. A scope is represented by a directory. Currently the usage of the parent directory scope is not supported in a sub-directory. The tool allows to create, use, list, remove, rename and reorganize the shortcuts.

# COMMANDS

**ins** <shortcut>
: Executes the command of the shortcut.

**ins** add <shortcut> <instruction>
: Adds a shortcut for the passed instruction in the current scope.

**ins** mv <shortcut-old> <shortcut-new>
: Renames replaces the old shortcut name by the new shortcut name.

**ins** rm <shortcut>
: Removes the passed shortcut name.

**ins** reorganize
: Reorganizes the files in which the shortcuts and instructions are stored.

**ins** list
: Lists all existing shortcuts. 

# EXAMPLE

**ins** test
: Run the command for the test shortcut.

**ins** add test "echo 'test'"
: Adds a shotcut named test with the instruction "echo 'test'".

**ins** mv test new-test
: Renames the test shortcut to new-test.

**ins** rm test
: Deletes the test shortcut.

# COPYRIGHT
Copyright (c) 2021 RobbyRed98 | MIT License