#!/bin/bash

set -eu

err() {
  echo "[ERR] $*" >&2
}

info() {
  echo "[INF] $*"
}

script_err() {
  err "failed the cyclomatic complexity test"
}

if ! which gocyclo > /dev/null; then
  err "please install gocyclo (https://github.com/fzipp/gocyclo)"
  script_err
  exit 1
fi

gocyclo . | awk '21<=$1' > result.txt
err_count="$(wc -l < result.txt)"
if [[ 0 -lt "$err_count" ]]; then
  err "the cyclomatic complexity is 21 or higher. please fix those funcitons"
  err "---------------------------------------------------------------------"
  cat result.txt >&2
  err "---------------------------------------------------------------------"
  script_err
  rm result.txt
  exit 1
fi

rm result.txt
info "passed the cyclomatic complexity test"
