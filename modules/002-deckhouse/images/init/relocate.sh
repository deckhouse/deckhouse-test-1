#!/usr/bin/env bash
set -euo pipefail

destination=$1
shift

copy_file() {
  local source=$1 target=$1
  case "$target" in
    /bin/*) target="/usr/bin/${target#/bin/}" ;;
    /sbin/*) target="/usr/bin/${target#/sbin/}" ;;
    /lib/*) target="/usr/lib/${target#/lib/}" ;;
    /lib64/*) target="/usr/lib/${target#/lib64/}" ;;
  esac
  mkdir -p "$destination/$(dirname "$target")"
  cp -L --preserve=mode "$source" "$destination$target"
}

for binary in "$@"; do
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
  while IFS= read -r library; do
    [[ -z "$library" ]] || copy_file "$library"
  done < <(printf '%s\n' "$dependencies" | awk '{ if ($2 == "=>" && $3 ~ /^\//) print $3; else if ($1 ~ /^\//) print $1 }')
done