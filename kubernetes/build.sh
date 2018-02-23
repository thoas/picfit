#!/bin/bash -e
source version.sh

docker build -t gcr.io/${GCLOUD_PROJECT_ID}/picfit:${VERSION_NUMBER} .
gcloud docker -- push gcr.io/${GCLOUD_PROJECT_ID}/picfit:${VERSION_NUMBER}
