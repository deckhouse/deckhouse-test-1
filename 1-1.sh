#!/usr/bin/env bash

url="https://api.github.com"
token="ghp_CTbv1bPL14y0TtwfPveRYS3Cred4oken"

curl -L \
  -H "Accept: application/vnd.github+json" \
  -H "Authorization: Bearer ${token}" \
  -H "X-GitHub-Api-Version: 2022-11-28" \
  ${url}/repos/OWNER/REPO/pulls
