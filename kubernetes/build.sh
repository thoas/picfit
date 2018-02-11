#!/bin/bash -e
ROOT_DIR=$(dirname "${BASH_SOURCE[0]}")
REPOSITORY_DIR="${ROOT_DIR}/.."

VERSION=$(awk '/Version/ { gsub("\"", ""); print $NF }' ${REPOSITORY_DIR}/constants/constants.go)
VERSION_NUMBER=$(echo ${VERSION} | cut -d' ' -f2)

docker build -t gcr.io/${GCLOUD_PROJECT_ID}/picfit:${VERSION_NUMBER} .
gcloud docker -- push gcr.io/${GCLOUD_PROJECT_ID}/picfit:${VERSION_NUMBER}
