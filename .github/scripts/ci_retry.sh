retry() {
  local max_attempts="${RETRY_MAX_ATTEMPTS:-3}"
  local base_sleep_seconds="${RETRY_BASE_SLEEP_SECONDS:-5}"
  local attempt=1
  local exit_code=0

  while true; do
    "$@"
    exit_code=$?

    if [ "${exit_code}" -eq 0 ]; then
      return 0
    fi

    if [ "${attempt}" -ge "${max_attempts}" ]; then
      echo "[retry] failed after ${attempt}/${max_attempts}: $* (exit=${exit_code})"
      return "${exit_code}"
    fi

    local wait_seconds=$((attempt * base_sleep_seconds))
    echo "[retry] attempt ${attempt}/${max_attempts} failed (exit=${exit_code}), retry in ${wait_seconds}s: $*"
    sleep "${wait_seconds}"
    attempt=$((attempt + 1))
  done
}

print_retry_policy() {
  local max_attempts="${RETRY_MAX_ATTEMPTS:-3}"
  local base_sleep_seconds="${RETRY_BASE_SLEEP_SECONDS:-5}"
  echo "[retry] policy: max_attempts=${max_attempts}, base_sleep_seconds=${base_sleep_seconds}"
}

run_with_retry() {
  local step_name="$1"
  shift
  echo "${step_name}"
  retry "$@"
}
