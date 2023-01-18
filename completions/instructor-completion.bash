__ins_completions() {
  declare -a INSTRUCTIONS
  if test -f .instructions; then
    INSTRUCTIONS=$(grep "${PWD}|" ~/.instructions | cut -d '|' -f2- | cut -d '-' -f1)
  fi
 
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

complete -F __ins_completions ins
