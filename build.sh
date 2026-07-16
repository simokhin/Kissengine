#!/bin/bash
set -e

label="$1"
if [ -z "$label" ]; then
	echo "Usage: ./build.sh <label>"
	exit 1
fi

hash=$(git rev-parse --short HEAD)
if ! git diff --quiet; then
	hash="${hash}-dirty"
fi

output_dir="$HOME/Kissengine"
mkdir -p "$output_dir"

output="$output_dir/kissengine-${label}-${hash}"
go build -o "$output" ./uci/

echo "Built: $output"
