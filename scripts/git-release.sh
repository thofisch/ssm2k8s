#!/bin/bash

owner=${1}
repo=${2}
version=${3}
branch=$(git rev-parse --abbrev-ref HEAD)
token=$(git config --global github.token)

generate_post_data()
{
  cat <<EOF
{
  "tag_name": "${version}",
  "target_commitish": "${branch}",
  "name": "${version}",
  "draft": false,
  "prerelease": false
}
EOF
}

echo "Creating release for ${version} on ${branch}"

curl -s --data "$(generate_post_data)" "https://api.github.com/repos/${owner}/${repo}/releases?access_token=${token}"
