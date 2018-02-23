ROOT_DIR=$(dirname "${BASH_SOURCE[0]}")
REPOSITORY_DIR="${ROOT_DIR}/.."

VERSION=$(awk '/Version/ { gsub("\"", ""); print $NF }' ${REPOSITORY_DIR}/constants/constants.go)
VERSION_NUMBER=$(echo ${VERSION} | cut -d' ' -f2)
