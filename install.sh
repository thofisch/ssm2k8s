#!/bin/bash

set -o errexit
set -o pipefail

if [[ "${OS}" = "Windows_NT" ]]
then
    os="windows"
    if [[ "${PROCESSOR_ARCHITECTURE}" = "AMD64" ]]
    then
        arch="amd64"
    elif [[ "${PROCESSOR_ARCHITECTURE}" = "x86" ]]
    then
        arch="386"
    fi
else
    os=$(uname -s)
    if [[ "${os}" = "Linux" ]]
    then
        os="linux"
    elif [[ "${os}" = "Darwin" ]]
    then
        os="darwin"
    fi
    arch=$(uname -m)
    if [[ "${arch}" = "x86_64" ]]
    then
        arch="amd64"
    elif [[ "${arch}" = *"86"* ]]
    then
        arch="386"
    fi
fi

echo "Downloading latest version of secrets (OS=${os},ARCH=${arch})"

[[ ${os} = "windows" ]] && ext=".exe" || ext=""

target="secrets${ext}"

curl -s https://api.github.com/repos/thofisch/ssm2k8s/releases/latest | \
    grep "browser_download_url" | \
    grep "${os}-${arch}" | \
    cut -d '"' -f 4 | \
    xargs curl -o ${target} -sSL

chmod u+x ${target}
