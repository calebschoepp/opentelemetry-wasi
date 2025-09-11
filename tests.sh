#!/bin/bash

set -e

build_and_run() {
    # This will let us capture STDERR and STDOUT from the background process
    OUTPUT=$(mktemp)
    spin build -f $1
    # This is run in the background so we can `curl` it later on
    spin otel up -- -f $1 &> $OUTPUT &
    SPIN_PID=$!
    sleep 2

    # If the process doesn't exist, `spin otel up` failed
    if ! kill -0 $SPIN_PID 2>/dev/null; then
        # Allow final output to be written
        sleep 0.5
        cat $OUTPUT
        rm $OUTPUT
        exit 1
    fi

    # TODO: Is there a way to catch if the process panics but doesn't exit?
    curl localhost:3000 && echo
    cat $OUTPUT

    # There are multiple processes to kill
    pkill -9 spin
    rm $OUTPUT
}

build_and_run "rust/examples/spin-basic/spin.toml"
build_and_run "rust/examples/spin-tracing/spin.toml"
