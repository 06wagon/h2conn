#! /usr/bin/env bash

set -e

if [ -v COVER ]
then
    echo "Running tests with coverage"
    FLAGS="${FLAGS} -coverprofile=/tmp/coverage.txt -covermode=atomic"
fi

if [ -v RACE ]
then
    echo "Running tests with race"
    FLAGS="${FLAGS} -race"
fi

function append-coverage {
    if [ -f /tmp/coverage.txt ]
    then
        cat /tmp/coverage.txt >> coverage.txt
    fi
}

rm -f coverage.txt /tmp/coverage.txt

# Run all tests in package
go test ${FLAGS} ./...
append-coverage

# Run only passing netconn.TestConn tests
TEST_CONN=1 go test -run "TestConn/(BasicIO|PingPong)" ${FLAGS}
append-coverage
