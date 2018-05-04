#!/bin/sh

script_dir=$(dirname $0)

log_dir=$script_dir/.testdata/run/$(date +"%s")

$script_dir/datamesh -logtostderr -v 8 -dir $log_dir server
