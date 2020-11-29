#!/usr/bin/env bash
set -ex

#FOLDER_ID=${FOLDER_ID:}
#YC_PROFILE=private
: ${YC_FOLDER_ID?Need to define FOLDER_ID}
: ${YC_PROFILE?Need to define YC_PROFILE}


# check build status before upload
function build {
  go build -buildmode=plugin .
}

# create new version of function upload
function deploy {
  yc --profile ${YC_PROFILE} --folder-id ${YC_FOLDER_ID} serverless function create --name=$1 || true

  [[ -s build.zip ]] && rm -r

  zip -r build.zip main.go go.mod go.sum audio.go util.go

  yc --profile ${YC_PROFILE} \
      --folder-id ${YC_FOLDER_ID} serverless function version create \
      --function-name=$1 \
      --runtime golang114 \
      --entrypoint main.Handler \
      --memory 256m \
      --execution-timeout 60s \
      --source-path ./build.zip
}

build
deploy "alice-english"
