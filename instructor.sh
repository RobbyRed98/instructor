#!/bin/bash

SCOPE="$PWD"
INSTRUCTIONS_FILE=~/.instructions
PARAMETER=$1
VERBOSE=0

function error() {
    echo "Error: $1"
}

function verbose() {
    if [[ $VERBOSE -eq 1 ]]; then
        echo "$1"
    fi
}

function label-scope-exists() {
    label=$1
    scope=$2
    if [[ ! -f "$INSTRUCTIONS_FILE" ]]; then 
        return 0
    fi
    
    if grep "^$scope|$label->" $INSTRUCTIONS_FILE > /dev/null 2>&1; then 
        return 1
    fi

    return 0
}

function sed-escape() {
    text=$1
    echo $(echo "$text" | sed "s|/|\\\/|g")
}

function help() {
    echo "Usage: 
instructor <command> <args>
ins <command> <args>

Allows the creation and usage of scope-bound shell shortcuts.

<shortcut>      Executes a created shortcut.
add             Creates a scope-bound shortcut for a shell command.
mv              Renames a shortcut.
rename          Also renames a shortcut.
rm              Removes a shortcut.
list            Lists all existing shortcuts.
reorganize      Reorganizes the file in which the shortcuts and commands are stored.

help            Prints this help text."
}

function list() {
    echo "Scope | Label -> Instruction"
    cat $INSTRUCTIONS_FILE | sort | sed "s/|/ | /g" | sed "s/->/ -> /g"
}

function add() {
    label=$1
    instruction=$2
    verbose "Adding instruction for \"$label\" -> \"$instruction\" in scope \"$SCOPE\"."

    if [[ -z "$label" ]]; then 
        error "An empty instruction label is not allowed!"
        exit 1
    fi

    if [[ -z "$instruction" ]]; then 
        error "An empty instruction is not allowed!"
        exit 1
    fi

    label-scope-exists "$label" "$SCOPE"
    if [[ $? -ne 0 ]]; then 
        error "A instruction for \"$label\" already exists in scope \"$SCOPE\"!"
        exit 1
    fi

    echo "$SCOPE|$label->$instruction" >> $INSTRUCTIONS_FILE
    exit 0
}

function remove() {
    label=$1
    scope=$2

    if [[ -z "$instruction" ]]; then
        verbose "No scope has been passed. Thus using local scope: \"$SCOPE\""
        scope="$SCOPE"
    fi

    verbose "Removing instruction for \"$label\" in \"$scope\"."
    
    label-scope-exists "$label" "$scope"
    if [[ $? -eq 0 ]]; then
        error "Cannot remove label-scope combination. The combination does not exist."
    fi

    sed -i "/^$(sed-escape $scope)|$(sed-escape $label)->/d" $INSTRUCTIONS_FILE

    exit 0
}

function rename() {
    current_label=$1
    new_label=$2
    scope=$3

    if [[ -z $scope ]]; then
        scope=$SCOPE
    fi

    label-scope-exists "$current_label" "$scope"
    if [[ $? -eq 0 ]]; then
        error "Cannot rename \"$scope|$current_label\" label-scope combination."
        error "Because the combination does not exist!"
        exit 0
    fi 

    label-scope-exists "$new_label" "$scope"
    if [[ $? -ne 0 ]]; then
        error "Cannot rename \"$scope|$current_label\" to \"$scope|$new_label\" label-scope combination." 
        error "Because the combination \"$scope|$new_label\" already exists!"
        exit 0
    fi 

    verbose "Renaming \"$scope|$current_label\" -> \"$scope|$new_label\""
    sed -i "s/^$(sed-escape $scope)|$(sed-escape $current_label)->/$(sed-escape $scope)|$(sed-escape $new_label)->/g" $INSTRUCTIONS_FILE
}

function reorganize() {
    verbose "Reorganizing the instruction file: $INSTRUCTIONS_FILE"
    sort "$INSTRUCTIONS_FILE" -o "$INSTRUCTIONS_FILE.sorted"
    mv "$INSTRUCTIONS_FILE.sorted" "$INSTRUCTIONS_FILE"
}

function run() {
    label=$1
    verbose "Instructing: \"$label\" in scope \"$SCOPE\""
    instruction=$(grep -m1 "^$SCOPE|$label" $INSTRUCTIONS_FILE | sed 's/.*->//')
    if [[ -z "$instruction" ]]; then 
        error "No instruction found for scope-label combination."
        exit 1
    fi
    verbose "Running: $instruction"
    bash -c "$instruction"
}

case $1 in

  list)
    list
    ;;

  add)
    add "$2" "$3"
    ;;

  rm)
    remove "$2" "$3"
    ;;

  mv|rename)
    rename "$2" "$3" "$4"
    ;;

  reorganize)
    reorganize
    ;;

  help)
    help    
    ;;

  *)
    run "$1"
    ;;
esac