#!/usr/bin/env bash
set -euo pipefail

# CI-friendly HURL runner: builds server-test, starts it, waits for health, runs hurl tests, then cleans up
BIN=server-test-bin
PORT=8080
HOST="http://localhost:${PORT}"
LOGFILE="/tmp/server-test-ci.log"

# Resolve script directory robustly so the script works from any cwd
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" >/dev/null 2>&1 && pwd)"
# If the script is placed in server-test/ (contains go files), build there.
# If the script is in server-test/scripts/, build from parent.
if compgen -G "${SCRIPT_DIR}/*.go" > /dev/null; then
  SERVER_DIR="${SCRIPT_DIR}"
else
  SERVER_DIR="${SCRIPT_DIR}/.."
fi

echo "Building server-test..."
(
  cd "${SERVER_DIR}" && go build -o ${BIN} .
)

echo "Starting server-test binary..."
pushd "${SERVER_DIR}" > /dev/null
./${BIN} > ${LOGFILE} 2>&1 &
SERVER_PID=$!
popd > /dev/null
echo "Server PID: ${SERVER_PID}"

cleanup() {
  echo "Stopping server (pid ${SERVER_PID})" || true
  kill ${SERVER_PID} 2>/dev/null || true
  wait ${SERVER_PID} 2>/dev/null || true
}
trap cleanup EXIT

# Wait for health endpoint
echo "Waiting for health endpoint ${HOST}/api/v1/health..."
RETRIES=30
SLEEP=1
for i in $(seq 1 ${RETRIES}); do
  if curl -s --max-time 2 "${HOST}/api/v1/health" > /dev/null 2>&1; then
    echo "Server healthy"
    break
  fi
  if ! kill -0 ${SERVER_PID} 2>/dev/null; then
    echo "Server process exited prematurely. See log: ${LOGFILE}" >&2
    cat ${LOGFILE} >&2 || true
    exit 1
  fi
  sleep ${SLEEP}
done

# Run HURL tests using existing script
echo "Running HURL tests..."
(
  cd "${SERVER_DIR}" && ./run-hurl-tests.sh
)

echo "HURL run completed"

# cleanup will run via trap
