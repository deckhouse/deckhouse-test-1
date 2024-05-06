#!/bin/bash

set -e

HELP_STRING=$(cat <<EOF
Usage: spell_check.sh [OPTION]
Run CMD on all kubernetes clusters

Optional arguments:
  FILENAME           the name of the file with a path (relative from the Deckhouse repo)
  -h, --help         output this message
EOF
)

if [ -n "$1" ]; then
  if [ "$1" = "-h" ] || [ "$1" = "--help" ]; then
    echo "$HELP_STRING"; exit 1;
  else
    FILENAME=$1;
  fi
fi

cd ../../docs/site/

werf run docs-spell-checker --config=werf-spell-check.yaml --env development --docker-options="--entrypoint=sh" -- /app/internal/container_spell_check.sh $FILENAME
