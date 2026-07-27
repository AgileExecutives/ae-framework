#!/usr/bin/env bash
set -euo pipefail

# Starts base-server with shared in-memory SQLite, waits for readiness, runs HURL tests, then stops the server.
# Usage: ./scripts/run-hurl-inmemory.sh [port]

PORT=${1:-:8080}
LOG=/tmp/base-server.log
export SERVER_PORT=${PORT}
export USE_IN_MEMORY_DB=true
# Optional override for DSN: export IN_MEMORY_DB_DSN="file:ae_saas?mode=memory&cache=shared"

echo "Starting base-server on ${SERVER_PORT} with in-memory DB (logs: ${LOG})"
# Kill any existing server on this port
p=$(lsof -tiTCP:${SERVER_PORT#*:} -sTCP:LISTEN || true)
if [ -n "$p" ]; then
  echo "Killing existing process(es): $p"
  kill -9 $p || true
fi

nohup go run . > "$LOG" 2>&1 &
PID=$!
sleep 1

echo "Waiting for server readiness..."
for i in {1..60}; do
  if curl -sSf "http://localhost:${SERVER_PORT#*:}/api/v1/health" >/dev/null 2>&1; then
    echo "Server ready"
    break
  fi
  echo "waiting ${i}..."
  sleep 1
done

if ! curl -sSf "http://localhost:${SERVER_PORT#*:}/api/v1/health" >/dev/null 2>&1; then
  echo "Server did not become ready; dumping last 200 lines of log:" >&2
  tail -n 200 "$LOG" >&2 || true
  kill -9 $PID || true
  exit 1
fi

# Run HURL tests
echo "Running HURL tests"
./run-hurl-tests.sh
RC=$?

echo "HURL tests finished with exit code $RC"

# Stop server
echo "Stopping server (pid $PID)"
kill -9 $PID || true
sleep 1

echo "Saved logs to $LOG"
exit $RC
