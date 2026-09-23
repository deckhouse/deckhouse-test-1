#!/usr/bin/env bash
set -euo pipefail

destination=$1
shift

copy_file() {
  local source=$1 target
  case "$source" in
    /bin/*|/sbin/*|/usr/sbin/*) target=/usr/bin/${source##*/} ;;
    /lib/*) target=/usr/lib/${source#/lib/} ;;
    /lib64/*) target=/usr/lib/${source#/lib64/} ;;
    /usr/lib64/*) target=/usr/lib/${source#/usr/lib64/} ;;
    *) target=$source ;;
  esac
  mkdir -p "$destination$(dirname "$target")"
  cp -L --preserve=mode,timestamps "$source" "$destination$target"
}

for pattern in "$@"; do
  matched=false
  for binary in $pattern; do
    [[ -f "$binary" ]] || continue
    matched=true
    copy_file "$binary"
    dependencies=$(ldd "$binary" 2>&1) || {
      case "$dependencies" in
        *"not a dynamic executable"*|*"statically linked"*) continue ;;
        *) printf '%s\n' "$dependencies" >&2; exit 1 ;;
      esac
    }
    if [[ "$dependencies" == *"not found"* ]]; then
      printf '%s\n' "$dependencies" >&2
      exit 1
    fi
    while read -r dependency; do
      copy_file "$dependency"
    done < <(printf '%s\n' "$dependencies" | awk '$2 == "=>" && $3 ~ /^\// {print $3} NF == 2 && $1 ~ /^\// && $2 ~ /^\(0x[[:xdigit:]]+\)$/ {print $1}')
  done
  if [[ "$matched" == false ]]; then
    printf 'Missing runtime artifact: %s\n' "$pattern" >&2
    exit 1
  fi
done