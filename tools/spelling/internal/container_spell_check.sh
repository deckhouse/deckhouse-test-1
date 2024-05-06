#!/bin/sh

set -e

ex_result=0

DICTIONARIES="/app/dictionaries/ru_RU,/app/dictionaries/en_US,/app/dictionaries/dev_OPS"

if [ -n "$1" ]; then
  FILENAME=$1
fi

echo "Checking docs..."

if [ -n "$1" ]; then
    check=1
    if test -f "/app/internal/filesignore"; then
      while read file_name_to_ignore; do
        if [[ "${FILENAME}" =~ ${file_name_to_ignore} ]]; then
          unset check
          check=0
        fi
      done <<-__EOF__
  $(cat /app/internal/filesignore | grep -vE '^#\s*|^\s*$')
__EOF__
      if [ "$check" -eq 1 ]; then
        echo "Possible typos in $(echo ${FILENAME} | sed '#^\./#d'):"
        result=$(python3 /app/internal/clean-files.py ${FILENAME} | sed '/^\s*$/d' | hunspell -d ${DICTIONARIES} -l)
        if [ -n "$result" ]; then
          echo $result | sed 's/\s\+/\n/g'
        fi
      else
        echo "Ignoring ${FILENAME}..."
      fi
    fi
else
  for file in `find ./ -type f`
  do
    check=1
    if test -f "/app/internal/filesignore"; then
      while read file_name_to_ignore; do
        if [[ "$file" =~ ${file_name_to_ignore} ]]; then
          unset check
          check=0
        fi
      done <<-__EOF__
  $(cat /app/internal/filesignore | grep -vE '^#\s*|^\s*$')
__EOF__
      if [ "$check" -eq 1 ]; then
        result=$(python3 /app/internal/clean-files.py $file | sed '/^\s*$/d' | hunspell -d ${DICTIONARIES} -l)
        if [ -n "$result" ]; then
          unset ex_result
          ex_result=1
          echo "Possible typos in $(echo ${file} | sed '#^\./#d'):"
          echo $result | sed 's/\s\+/\n/g'
          echo
        fi
      else
        echo "Ignoring $file..."
      fi
    fi
  done
  if [ "$ex_result" -eq 1 ]; then
    exit 1
  fi
fi
