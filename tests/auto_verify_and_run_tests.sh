#!/bin/bash

# Auto-verify email tokens found in tmp/mock_emails.json while tests run
BASE_DIR="$(cd "$(dirname "$0")/.." && pwd)"
MOCK_FILE="$BASE_DIR/tmp/mock_emails.json"
PROCESSED_TOKENS_FILE="$BASE_DIR/tmp/processed_tokens.txt"
HOST="http://localhost:8080"

mkdir -p "$BASE_DIR/tmp"
touch "$PROCESSED_TOKENS_FILE"

process_new_tokens() {
  if [ ! -f "$MOCK_FILE" ]; then
    return
  fi

  # Extract any verification tokens from HTML content in mock emails
  jq -r '.[] | .HTML' "$MOCK_FILE" 2>/dev/null | while IFS= read -r html; do
    # look for verify-email?token=TOKEN using awk to avoid nested-quote issues
    tokens=$(echo "$html" | awk '{
      while(match($0,"verify-email\\?token=[^\"&]+")) {
        s = substr($0, RSTART, RLENGTH);
        sub("verify-email\\?token=","",s);
        print s;
        $0 = substr($0, RSTART+RLENGTH);
      }
    }') || true
    for t in $tokens; do
      if [ -z "$t" ]; then
        continue
      fi
      if ! grep -qx "$t" "$PROCESSED_TOKENS_FILE" 2>/dev/null; then
        echo "[auto-verify] Verifying token: $t"
        # Call verify endpoint (ignore result)
        curl -s -o /dev/null -w "%{http_code}" "$HOST/api/v1/auth/verify-email/$t" || true
        echo "$t" >> "$PROCESSED_TOKENS_FILE"
      fi
    done
  done
}

# Start watcher in background
( while true; do process_new_tokens; sleep 0.5; done ) &
WATCHER_PID=$!

echo "Started auto-verify watcher (PID=$WATCHER_PID), running tests..."

cd "$BASE_DIR" || exit 1
./run-hurl-tests.sh || true

# Stop watcher
kill "$WATCHER_PID" 2>/dev/null || true
wait "$WATCHER_PID" 2>/dev/null || true

echo "Auto-verify watcher stopped. Tests finished."
