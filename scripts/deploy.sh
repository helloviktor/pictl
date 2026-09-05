#!/usr/bin/env bash
# Builds the arm64 binaries and rsyncs them to the remote host defined in .env.
set -euo pipefail

cd "$(dirname "${BASH_SOURCE[0]}")/.."

ENV_FILE=".env"
if [[ ! -f "$ENV_FILE" ]]; then
	echo "error: $ENV_FILE not found. Copy .env.example to $ENV_FILE and fill in your connection details." >&2
	exit 1
fi

# shellcheck disable=SC1090
source "$ENV_FILE"

: "${DEPLOY_HOST:?DEPLOY_HOST is not set in $ENV_FILE}"
: "${DEPLOY_USER:?DEPLOY_USER is not set in $ENV_FILE}"
DEPLOY_PATH="${DEPLOY_PATH:-/usr/local/bin}"
DEPLOY_SSH_PORT="${DEPLOY_SSH_PORT:-22}"

OUTPUT_DIR="out"

echo "Building arm64 binaries..."
GOARCH=arm64 TAGS="${TAGS:-rpi}" make build

echo "Deploying to ${DEPLOY_USER}@${DEPLOY_HOST}:${DEPLOY_PATH}..."
rsync -avz --progress \
	-e "ssh -p ${DEPLOY_SSH_PORT}" \
	"${OUTPUT_DIR}/dashboard" \
	"${OUTPUT_DIR}/pictl" \
	"${DEPLOY_USER}@${DEPLOY_HOST}:${DEPLOY_PATH}/"

echo "Deploy complete."
