#!/usr/bin/env bash

# Run this script with the Makefile
# make create-test-resources
#
# This tool generates json files in the inputs/ for `cfn test`.
#
set -o errexit
set -o nounset
set -o pipefail
set -x

# setting projectName
projectName="ct-cloud-backup-restore-jobs-$((1 + RANDOM % 10000))"

./test/cfn-test-create-inputs.sh $projectName