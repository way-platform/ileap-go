#!/bin/bash

#
# What this script does:
# 1. find out the architecture of the machine
# 2. download the conformance binary for the target architecture
# 3. make the binary executable and executes is
#

set -xeuo pipefail

####
# Step 1: Find out the architecture of the machine or default to `arm64`
####

arch=$(uname -m)

case ${arch} in
arm64 | x86_64)
    echo "Downloading ${arch} binary."
    ;;
*)
    echo "The architecture specified is not supported, yet."
    exit 5    
    ;;
esac


####
# Step 2: Download the conformance binary for the target architecture
####
conformance_bin=$(mktemp)
trap "rm ${conformance_bin}" EXIT
curl -o ${conformance_bin} "https://actbin.blob.core.windows.net/act-bin/conformance_${arch}"


####
# Step 3: Make the binary executable and execute it
####
chmod +x ${conformance_bin}
exec ${conformance_bin} "$@"
