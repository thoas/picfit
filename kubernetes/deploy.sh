#!/bin/bash -e
source version.sh

echo ${VERSION_NUMBER}

# CONFIG=$(envsubst < picfit-config.yml)
# cat <<EOF | kubectl replace --force -f -
# ${CONFIG}
# EOF
#
# DEPLOYMENT=$(envsubst < picfit-deployment.yml)
# cat <<EOF | kubectl replace --force -f -
# ${DEPLOYMENT}
# EOF
#
# SERVICE=$(envsubst < picfit-service.yml)
# cat <<EOF | kubectl replace --force -f -
# ${SERVICE}
# EOF
