#!/usr/bin/env bash

set -euo pipefail

binary="${1}"

archives=$(curl -s -L \
  -H "Accept: application/vnd.github+json" \
  -H "Authorization: Bearer ${GITHUB_TOKEN}" \
  -H "X-GitHub-Api-Version: 2022-11-28" \
  "https://api.github.com/repos/cattlecloud/${binary}/actions/artifacts")

best=$(echo "${archives}" | jq -r '.artifacts | sort_by(.created_at) | last | .archive_download_url')
dest="/tmp/${binary}.zip"
file="/tmp/${binary}"

echo "downloading ${best} into ${dest}"
curl -L -H "Authorization: Bearer ${GITHUB_TOKEN}" -H "X-GitHub-Api-Version: 2022-11-28" -o "${dest}" -s "${best}"
pushd /tmp && unzip "${dest}" && popd
echo "download complete."

scp "${file}" "${HOST}:~/${binary}"
ssh -t "${HOST}" "sudo mv ~/${binary} /opt/bin/${binary}"
ssh -t "${HOST}" "sudo systemctl restart ${binary}"
ssh -t "${HOST}" "sudo systemctl status  ${binary}"

rm "${dest}" "${file}"
echo "done."

