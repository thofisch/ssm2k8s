#!/bin/bash

owner=${1}
repo=${2}
version=${3}
name=${4}
file=${5}
token=$(git config --global github.token)

echo "Uploading asset ${file} to ${version}"

upload_url=$(curl -s "https://api.github.com/repos/${owner}/${repo}/releases/tags/${version}" | \
    grep "upload_url"  | \
    cut -d '"' -f 4 | \
    sed -e "s/{?name,label}//")

curl --netrc \
    --header "Content-Type:application/gzip" \
    --header "Authorization: token ${token}" \
    --data-binary "@${file}" \
    "${upload_url}?name=${name}"

