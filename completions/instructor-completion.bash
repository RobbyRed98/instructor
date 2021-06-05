#!/usr/bin/env bash
_ins_completions() {
  INSTRUCTIONS=$(grep "${PWD}|" ~/.instructions | cut -d '|' -f2- | cut -d '-' -f1)
  if [ "${#COMP_WORDS[@]}" == "2" ]; then
    COMMANDS="add edit list mv rm rename reorganize help"
    WORDS="$COMMANDS $INSTRUCTIONS"
    COMPREPLY=($(compgen -W "$WORDS" "${COMP_WORDS[1]}"))
  elif [ "${#COMP_WORDS[@]}" == "3" ]; then
    case "${COMP_WORDS[1]}" in
      list)
        COMPREPLY=($(compgen -W "all" "${COMP_WORDS[2]}"))
        ;;
      edit|mv|rename|rm)
        COMPREPLY=($(compgen -W "$INSTRUCTIONS" "${COMP_WORDS[2]}"))
        ;;
      *)
        return
        ;;
    esac
  fi
}

complete -F _ins_completions ins