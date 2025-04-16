#!/bin/bash

# Configuration
GO_SERVICE_URL="http://localhost:8080"  # Go service URL
PY_SERVICE_URL="http://localhost:5000"   # Python service URL
TEST_OWNER="testuser_$(date +%s)"       # Unique owner name
MAX_WAIT=30                             # Max wait time for services (seconds)

echo "=== Banking System Integration Test ==="

# Helper function with error handling
api_request() {
    local method=$1
    local endpoint=$2
    local data=$3
    local service_url=${4:-$GO_SERVICE_URL}
    
    response=$(curl -s -X "$method" "$service_url$endpoint" \
        -H "Content-Type: application/json" \
        -w "\nHTTP_STATUS:%{http_code}" \
        -d "$data")
    
    http_status=$(echo "$response" | grep "HTTP_STATUS:" | cut -d':' -f2)
    response_body=$(echo "$response" | grep -v "HTTP_STATUS:")
    
    if [ "$http_status" -ge 400 ]; then
        echo "ERROR: API request failed" >&2
        echo "Endpoint: $method $endpoint" >&2
        echo "Status: $http_status" >&2
        echo "Response:" >&2
        echo "$response_body" | jq . >&2
        return 1
    fi
    
    echo "$response_body"
}

# Wait for services to be ready
echo -e "\n[1/12] Waiting for services to be ready..."
for i in $(seq 1 $MAX_WAIT); do
    if curl -s -f "$PY_SERVICE_URL/health" >/dev/null && 
       curl -s -f "$GO_SERVICE_URL/health" >/dev/null; then
        break
    fi
    sleep 1
    if [ "$i" -eq "$MAX_WAIT" ]; then
        echo "ERROR: Services not ready after $MAX_WAIT seconds" >&2
        exit 1
    fi
done
echo "Services ready!"

# ====== Account Tests ======
echo -e "\n[2/12] Creating test account:"
CREATE_RESPONSE=$(api_request "POST" "/accounts" "{\"owner\":\"$TEST_OWNER\"}")
echo "$CREATE_RESPONSE" | jq .
ACCOUNT_ID=$(echo "$CREATE_RESPONSE" | jq -r '.id')

echo -e "\n[3/12] Verifying account exists:"
api_request "GET" "/accounts/$ACCOUNT_ID" "" | jq .

# ====== Transaction Tests ======
echo -e "\n[4/12] Testing invalid account rejection:"
api_request "POST" "/transactions" \
    '{"account_id":"invalid-account-id","amount":100.0,"type":"deposit"}'

echo -e "\n[5/12] Testing deposit transaction:"
DEPOSIT_RESPONSE=$(api_request "POST" "/transactions" \
    '{"account_id":"'"$ACCOUNT_ID"'","amount":200.0,"type":"deposit"}')
echo "$DEPOSIT_RESPONSE" | jq .
DEPOSIT_ID=$(echo "$DEPOSIT_RESPONSE" | jq -r '.transaction_id')

echo -e "\n[6/12] Testing withdrawal transaction:"
WITHDRAWAL_RESPONSE=$(api_request "POST" "/transactions" \
    '{"account_id":"'"$ACCOUNT_ID"'","amount":50.0,"type":"withdrawal"}')
echo "$WITHDRAWAL_RESPONSE" | jq .
WITHDRAWAL_ID=$(echo "$WITHDRAWAL_RESPONSE" | jq -r '.transaction_id')

echo -e "\n[7/12] Checking transactions for account:"
api_request "GET" "/transactions?account_id=$ACCOUNT_ID" "" | jq .

echo -e "\n[8/12] Listing all transactions:"
api_request "GET" "/transactions" "" | jq .

echo -e "\n[9/12] Getting single transaction:"
api_request "GET" "/transactions/$DEPOSIT_ID" "" | jq .

echo -e "\n[10/12] Testing analytics:"
api_request "GET" "/analytics?account_id=$ACCOUNT_ID" "" | jq .

# ====== Cleanup ======
echo -e "\n[11/12] Deleting test account:"
api_request "DELETE" "/accounts/$ACCOUNT_ID" ""

echo -e "\n[12/12] Verifying account deletion:"
api_request "GET" "/accounts/$ACCOUNT_ID" "" | jq .

echo -e "\n=== Test Completed Successfully ==="