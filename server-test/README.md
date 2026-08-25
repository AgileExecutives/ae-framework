# server-test

This folder contains integration test harness and Hurl tests for the server.

Development notes
- CI sets `SKIP_REDUNDANT_HURL=true` so redundant Hurl tests are skipped when Go integration tests cover the same scenarios.
- Archived processed Hurl files: `tests/hurl/processed/archived/` — restore any archived `.hurl` file into `tests/hurl/processed/` to re-enable it.

Running tests locally
- Run Go integration tests:

```bash
cd server-test
go test ./integration -v
```

- Run Hurl tests (if not skipped):

```bash
cd server-test
./run-hurl-tests.sh
```

If you need CI to run Hurl tests again, remove or unset `SKIP_REDUNDANT_HURL` in the workflow.
