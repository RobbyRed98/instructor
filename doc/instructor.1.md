% instructor(1) 0.1.0
% RobbyRed98
% March 2021

# NAME
instructor - a small cli-tool to add scope bound shortcuts for regular tasks/commands. 

# SYNOPSIS

**instructor** <shortcut>

**instructor** <command> <args>

# DESCRIPTION

A cli-tool to create shortcuts for specific shell comands. The commands are scope based. A scope is represented by a directory. Currently the usage of the parent directory scope is not supported in a sub-directory. The tool allows to create, use, list, remove, rename and reorganize the shortcuts.

# Commands

**instructor** <shortcut>
: Executes the command of the shortcut.

**instructor** add <shortcut> <instruction>
: Adds a shortcut for the passed instruction in the current scope.

**instructor** mv <shortcut-old> <shortcut-new>
: Renames replaces the old shortcut name by the new shortcut name.

**instructor** rm <shortcut>
: Removes the passed shortcut name.

**instructor** reorganize
: Reorganizes the files in which the shortcuts and instructions are stored.

**instructor** list
: Lists all existing shortcuts. 

# EXAMPLE

**instructor** test
: Run the command for the test shortcut.

**instructor** add test "echo 'test'"
: Adds a shotcut named test with the instruction "echo 'test'".

**instructor** mv test new-test
: Renames the test shortcut to new-test.

**instructor** rm test
: Deletes the test shortcut.

# COPYRIGHT
Copyright (c) 2021 RobbyRed98 | MIT License