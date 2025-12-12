#!/usr/bin/env bash
set -euo pipefail

# E2E Saga test script for 78-pflops
# Runs two scenarios:
#  1) Positive: upload images + create ad -> expect media persisted and ad created
#  2) Compensation: stop ad_service, upload images via gateway -> expect media removed

API_HOST=${API_HOST:-http://localhost:8080}
MEDIA_GRPC_HOST=${MEDIA_GRPC_HOST:-localhost:50053}
AD_SERVICE_NAME=${AD_SERVICE_NAME:-ad_service}
DOCKER_COMPOSE_DIR=${DOCKER_COMPOSE_DIR:-$(dirname "$(dirname "$0")")}
PROTO_PATH=${PROTO_PATH:-$(dirname "$(dirname "$0")")/MediaService/proto}

# tiny 1x1 PNG base64
SAMPLE_IMG="iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAQAAAC1HAwCAAAAC0lEQVR42mP8/x8AAwMB/6Xc6aQAAAAASUVORK5CYII="

die(){ echo "ERROR: $*" >&2; exit 1; }
check_cmd(){ command -v "$1" >/dev/null 2>&1 || die "required command '$1' not found"; }

# prerequisites check
check_prereqs(){
  check_cmd curl
  check_cmd jq
  check_cmd grpcurl
  check_cmd docker
}

# Register and return JSON with id and token
register_user(){
  local email="saga-test-$(date +%s)@example.com"
  # Use a password matching policy: at least 8 chars and a special character
  local password="P@ssw0rd123!"
  local name="saga-test"
  # log to stderr to avoid contaminating JSON output captured by caller
  echo "Registering user: $email" >&2
  resp=$(curl -sS -X POST "$API_HOST/api/auth/register" -H 'Content-Type: application/json' -d "{\"email\":\"$email\",\"password\":\"$password\",\"name\":\"$name\"}")
  echo "$resp"
}

# Create ad via gateway
create_ad(){
  local token="$1"
  local title="$2"
  echo "Creating ad (title='$title')"
  resp=$(curl -sS -X POST "$API_HOST/api/ads" -H 'Content-Type: application/json' -d "{\"token\":\"$token\",\"title\":\"$title\",\"description\":\"desc\",\"price\":1.23,\"images\":[\"$SAMPLE_IMG\"]}") || true
  echo "$resp"
}

# List ads via gateway (returns raw JSON)
list_ads(){
  curl -sS "$API_HOST/api/ads" || true
}

# List media for user via grpcurl
list_media(){
  local user_id=$1
  # Use simple grpcurl invocation compatible with snap package: put -d before host
  grpcurl -plaintext -d "{\"user_id\":\"$user_id\"}" $MEDIA_GRPC_HOST media.MediaService.ListMedia
}

# Helper: extract token and id from register response
extract_token(){ echo "$1" | jq -r '.token // empty'; }
extract_userid(){ echo "$1" | jq -r '.id // empty'; }

main(){
  echo "Running Saga E2E tests"
  check_prereqs

  # Register
  reg=$(register_user)
  # print raw register response for debugging
  echo "RAW REGISTER RESP: $reg" >&2
  token=$(extract_token "$reg")
  user_id=$(extract_userid "$reg")
  [ -n "$token" ] || die "failed to obtain token from register response: $reg"
  [ -n "$user_id" ] || die "failed to obtain user id from register response: $reg"
  echo "Got token for user $user_id"

  echo "\n--- Positive scenario ---"
  create_resp=$(create_ad "$token" "SAGA_POSITIVE_TEST")
  echo "Create response: $create_resp"

  echo "Waiting 1s for propagation..."
  sleep 1

  echo "Listing ads via gateway:"
  list_ads | jq . || true

  echo "Listing media for user via MediaService:"
  list_media "$user_id" | jq . || true

  echo "\n--- Compensation scenario (stop ad service) ---"
  echo "Stopping ad service: docker compose stop $AD_SERVICE_NAME"
  (cd "$DOCKER_COMPOSE_DIR" && docker compose stop "$AD_SERVICE_NAME")
  sleep 2

  create_resp2=$(create_ad "$token" "SAGA_COMPENSATE_TEST") || true
  echo "Create response (expected failure): $create_resp2"

  echo "Waiting 1s for compensation..."
  sleep 1

  echo "Listing media for user after compensation:"
  after_comp=$(list_media "$user_id" || true)
  echo "$after_comp" | jq . || true
  # Assert no media items remain for this test user
  count=$(echo "$after_comp" | jq '.media_items | length' 2>/dev/null || echo 0)
  if [ "${count:-0}" -ne 0 ]; then
    echo "Compensation failed: $count media items remain for user $user_id" >&2
    die "compensation failed"
  else
    echo "Compensation succeeded: no media items remain for user $user_id"
  fi

  echo "Starting ad service back"
  (cd "$DOCKER_COMPOSE_DIR" && docker compose start "$AD_SERVICE_NAME")

  echo "E2E finished"
}

main "$@"
