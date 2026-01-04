#!/usr/bin/env bash
set -euo pipefail

BITCOIN_CLI="${BITCOIN_CLI:-bitcoin-cli}"
BITCOIN_CLI_ARGS="${BITCOIN_CLI_ARGS:-}"
DESC_FILE="${1:-${BTC_CORE_COMPAT_DESCRIPTOR_FILE:-}}"

if [[ -z "${DESC_FILE}" ]]; then
  echo "Usage: $0 <descriptor-file>"
  echo "Or set BTC_CORE_COMPAT_DESCRIPTOR_FILE to the descriptor JSON/text file exported from keygen or Bitcoin Core."
  exit 1
fi

if [[ ! -f "${DESC_FILE}" ]]; then
  echo "Descriptor file not found: ${DESC_FILE}"
  exit 1
fi

run_cli() {
  # shellcheck disable=SC2086
  "${BITCOIN_CLI}" ${BITCOIN_CLI_ARGS} "$@"
}

echo "=== Bitcoin Core descriptor compatibility smoke test ==="
echo "bitcoin-cli: ${BITCOIN_CLI} ${BITCOIN_CLI_ARGS}"
run_cli -version
echo "Using descriptor file: ${DESC_FILE}"

echo "--- Importing descriptors into Bitcoin Core ---"
run_cli importdescriptors "$(cat "${DESC_FILE}")"
echo "✓ importdescriptors succeeded"

echo "--- Deriving first two addresses from first descriptor ---"
first_desc=$(
  python - <<'PY'
import json,sys
from pathlib import Path
data=Path(sys.argv[1]).read_text()
try:
    recs=json.loads(data)
    if isinstance(recs,list) and recs and "desc" in recs[0]:
        print(recs[0]["desc"])
        sys.exit(0)
except Exception:
    pass
for line in data.splitlines():
    line=line.strip()
    if line and not line.startswith("#"):
        print(line)
        sys.exit(0)
sys.exit(1)
PY
  "${DESC_FILE}"
)

if [[ -z "${first_desc}" ]]; then
  echo "Could not extract a descriptor from ${DESC_FILE}"
  exit 1
fi

run_cli deriveaddresses "${first_desc}" "[0,1]"
echo "✓ deriveaddresses succeeded"

echo "Descriptor compatibility checks completed."
