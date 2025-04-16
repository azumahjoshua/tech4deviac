#!/bin/bash

# Configuration
GO_SERVICE_URL="http://localhost:8080"  # Your Go service URL
TEST_OWNER="testuser_$(date +%s)"       # Unique owner name for testing

echo "=== Comprehensive API Test ==="

# Helper function for API calls
api_request() {
    local method=$1
    local endpoint=$2
    local data=$3
    
    curl -s -X "$method" "$GO_SERVICE_URL$endpoint" \
        -H "Content-Type: application/json" \
        -d "$data"
}

# ====== Account Tests ======
echo -e "\n=== Account Service Tests ==="

# 1. Create test account
echo -e "\n1. Creating test account:"
CREATE_RESPONSE=$(api_request "POST" "/accounts" "{\"owner\":\"$TEST_OWNER\"}")
echo "$CREATE_RESPONSE" | jq .
ACCOUNT_ID=$(echo "$CREATE_RESPONSE" | jq -r '.id')

# 2. Verify account exists
echo -e "\n2. Getting created account:"
api_request "GET" "/accounts/$ACCOUNT_ID" "" | jq .

# ====== Transaction Tests ======
echo -e "\n=== Transaction Service Tests ==="

# 3. Test invalid account rejection
echo -e "\n3. Testing invalid account rejection:"
api_request "POST" "/transactions" \
    '{"account_id":"invalid-account-id","amount":100.0,"type":"deposit"}' | jq .

# 4. Test successful deposit
echo -e "\n4. Testing valid deposit transaction:"
DEPOSIT_RESPONSE=$(api_request "POST" "/transactions" \
    '{"account_id":"'"$ACCOUNT_ID"'","amount":200.0,"type":"deposit"}')
echo "$DEPOSIT_RESPONSE" | jq .
DEPOSIT_ID=$(echo "$DEPOSIT_RESPONSE" | jq -r '.transaction_id')

# 5. Test successful withdrawal
echo -e "\n5. Testing valid withdrawal transaction:"
WITHDRAWAL_RESPONSE=$(api_request "POST" "/transactions" \
    '{"account_id":"'"$ACCOUNT_ID"'","amount":50.0,"type":"withdrawal"}')
echo "$WITHDRAWAL_RESPONSE" | jq .
WITHDRAWAL_ID=$(echo "$WITHDRAWAL_RESPONSE" | jq -r '.transaction_id')

# 6. Verify initial deposit exists (from account creation)
echo -e "\n6. Checking initial deposit transaction:"
api_request "GET" "/transactions?account_id=$ACCOUNT_ID" "" | jq .

# 7. Test transaction listing
echo -e "\n7. Listing all transactions:"
api_request "GET" "/transactions" | jq .

# 8. Test getting single transaction
echo -e "\n8. Getting single transaction:"
api_request "GET" "/transactions/$DEPOSIT_ID" "" | jq .

# 9. Test analytics endpoint
echo -e "\n9. Testing analytics endpoint:"
api_request "GET" "/analytics?account_id=$ACCOUNT_ID" "" | jq .

# ====== Cleanup ======
echo -e "\n=== Cleanup ==="

# 10. Delete test account
echo -e "\n10. Deleting test account:"
api_request "DELETE" "/accounts/$ACCOUNT_ID" ""
echo "Account deleted (status 204 expected)"

# 11. Verify account deletion
echo -e "\n11. Verifying account deletion:"
api_request "GET" "/accounts/$ACCOUNT_ID" "" | jq .

echo -e "\n=== All Tests Completed ==="