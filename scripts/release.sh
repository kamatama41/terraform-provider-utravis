#!/usr/bin/env bash

set -eu

if [[ "${TRAVIS}" != "true" ]]; then
  echo "This script is allowed to run on Travis CI"
  exit 1
fi

git config --global user.email "shiketaudonko41@gmail.com"
git config --global user.name "kamatama41"
git checkout master
git reset --hard origin/master

PROJECT_ROOT=$(cd $(dirname $0)/..; pwd)
VERSION_FILE=${PROJECT_ROOT}/version
VERSION=$(cat ${VERSION_FILE})

echo "## Create the new release"
go get github.com/aktau/github-release
github-release release \
  --user kamatama41 \
  --repo terraform-provider-unofficial-travis \
  --tag ${VERSION}


echo "## Build and upload release binaries"
PLATFORMS="darwin/amd64"
PLATFORMS="${PLATFORMS} windows/amd64"
PLATFORMS="${PLATFORMS} linux/amd64"
for PLATFORM in ${PLATFORMS}; do
  GOOS=${PLATFORM%/*}
  GOARCH=${PLATFORM#*/}
  BIN_FILENAME="terraform-provider-utravis_${VERSION}_x4"
  CMD="GOOS=${GOOS} GOARCH=${GOARCH} go build -o ${BIN_FILENAME}"
  rm -f ${BIN_FILENAME}
  echo "${CMD}"
  eval ${CMD}

  ZIP_FILENAME="terraform-provider-utravis_${VERSION}_${GOOS}_${GOARCH}.zip"
  CMD="zip ${ZIP_FILENAME} ${BIN_FILENAME}"
  echo "${CMD}"
  eval ${CMD}

  github-release upload \
    --user kamatama41 \
    --repo terraform-provider-unofficial-travis \
    --tag ${VERSION} \
    --name ${ZIP_FILENAME} \
    --file ${ZIP_FILENAME}
done

gem install semantic
script=$(cat << EOS
require 'semantic'
puts "v#{Semantic::Version.new(gets[1..-1]).increment!(:patch)}"
EOS
)
NEXT_VERSION=$(cat version | ruby -e "${script}")
echo ${NEXT_VERSION} > ${VERSION_FILE}

echo "## Bump up the version to ${NEXT_VERSION}"
git add ${VERSION_FILE}
git commit -m "Bump up to the next version"
git push origin master
