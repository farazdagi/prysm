#!/bin/sh -l
set -e

cd $GITHUB_WORKSPACE

# Check if any files are not formatted.
set +e
test -z "$(gofmt -l -d -e $1)"
SUCCESS=$?
set -e

# Exit if `go fmt` passes.
if [ $SUCCESS -eq 0 ]; then
    echo "::set-output name=result::gofmt OK"
    exit 0
fi

echo "::set-output name=result::gofmt FAIL"
exit 1
