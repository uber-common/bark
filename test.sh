#!/bin/bash

function die {
    echo "not ok: fatal tests, $*"
    exit 1
}

go test -v || die "\"go test\" failed"

# Fatal
output=$(go run ./testhelp/fatal.go 2>&1)
[[ $? == 1 ]] || die "Failed: Fatal test should exit with status 1"
msg=$(echo "${output}" | head -n 1 | jq -r .msg)
[[ "${msg}" = "halp" ]] || die "Fatal test should print \"halp\", printed \"${msg}\""

# Fatalf 
output=$(go run ./testhelp/fatalf.go 2>&1)
[[ $? == 1 ]] || die "Failed: Fatalf test should exit with status 1"
msg=$(echo "${output}" | head -n 1 | jq -r .msg)
[[ "${msg}" = "halppls" ]] || die "Fatalf test should print \"halppls\", printed \"${msg}\""

echo "ok      fatal tests"
