#!/usr/bin/env bash

TEST=./...;
FMT="*.go"

echo "Running tests...";
go test -v $TEST;
go test -v -race $TEST;

echo "Checking gofmt..."
fmtRes=$(gofmt -l -s $FMT)
if [ -n "${fmtRes}" ]; then
	echo -e "gofmt checking failed:\n${fmtRes}"
	exit 255
fi

echo "Success";
