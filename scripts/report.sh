#!/bin/bash
ROOT="$(git rev-parse --show-toplevel)"
cdiff='diff'
if [[ ! -z "$(command -v colordiff)" ]]; then
  cdiff='colordiff'
fi

GO_REPORT=$(go test $ROOT | grep FAIL: | sort | grep -o -e 'test_[^ ]*')
RB_REPORT=$(ruby $ROOT/fnmatch_test.rb | grep Failure | sort | grep -o -e 'test_[^(]*')

$cdiff -U 8 <(echo "$RB_REPORT") <(echo "$GO_REPORT") | sed 's|--- /dev/[^ ]*|--- ruby  |' | sed 's|+++ /dev/[^ ]*|+++ golang|'
