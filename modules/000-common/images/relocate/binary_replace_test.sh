#!/bin/bash

set -euo pipefail

eval "$(sed -n '/^function relocate()/,/^function get_binary_path/ { /^function get_binary_path/!p; }' "$(dirname "$0")/binary_replace.sh")"
declare -F relocate relocate_item >/dev/null || exit 1

RDIR=/relocate

function mkdir() { :; }
function cp() {
  [[ "$1" == "-L" && "$2" == "--preserve=mode,timestamps" ]] || exit 1
  printf '%s\n' "$4"
}
function ldd() {
  if [[ "$1" == /usr/bin/static ]]; then
    return 1
  fi
  printf '%s\n' \
    'linux-vdso.so.1 (0x0000)' \
    'libc.so.6 => /lib64/libc.so.6 (0x0000)' \
    'libexample.so => /usr/lib64/libexample.so (0x0000)' \
    '/lib64/ld-linux-x86-64.so.2 (0x0000)'
}

for directory in /bin /sbin /usr/sbin /usr/bin; do
  [[ "$(relocate_item "$directory/tool")" == /relocate/usr/bin/tool ]] || exit 1
done
for directory in /lib /lib64 /usr/lib64 /usr/local/lib64 /usr/lib; do
  [[ "$(relocate_item "$directory/libexample.so")" == /relocate/usr/lib/libexample.so ]] || exit 1
done
[[ "$(relocate_item /usr/lib64/gconv/example.so)" == /relocate/usr/lib/gconv/example.so ]] || exit 1
[[ "$(relocate /usr/bin/static)" == /relocate/usr/bin/static ]] || exit 1
[[ "$(relocate /bin/tool)" == "$(printf '%s\n' /relocate/usr/bin/tool /relocate/usr/lib/libc.so.6 /relocate/usr/lib/libexample.so /relocate/usr/lib/ld-linux-x86-64.so.2)" ]] || exit 1

printf '%s\n' 'PASS: relocation paths, dereferencing, static binaries and dynamic dependencies'