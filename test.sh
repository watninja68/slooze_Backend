#!/bin/bash

# --- Configuration ---
BASE_URL="http://localhost:8080"
# In a real scenario, you'd get this token from your login endpoint
AUTH_TOKEN="YOUR_PLACEHOLDER_JWT_TOKEN"
HEADER_AUTH="Authorization: Bearer $AUTH_TOKEN"
HEADER_CONTENT="Content-Type: application/json"

# --- Helper for Pretty Printing ---
function print_json {
  if command -v jq &> /dev/null; then
    jq .
  else
    cat
  fi
}

echo "======================================="
echo " Starting API Endpoint Tests"
echo " Target: $BASE_URL"
echo "======================================="

# --- Public Routes ---

echo "Testing GET / (HelloWorld)"
curl -s -X GET "$BASE_URL/" | print_json
echo -e "\n---------------------------------------\n"

echo "Testing GET /health"
curl -s -X GET "$BASE_URL/health" | print_json
echo -e "\n---------------------------------------\n"

# --- Authenticated Routes (/api/v1) ---
# Note: These requests will use the placeholder auth middleware.
# You might need to adjust based on the hardcoded user in your middleware.
# If your placeholder middleware returns 403/401, that means it's working as expected.

echo "Testing GET /api/v1/restaurants (All Roles)"
curl -s -X GET -H "$HEADER_AUTH" "$BASE_URL/api/v1/restaurants" | print_json
echo -e "\n---------------------------------------\n"

echo "Testing GET /api/v1/restaurants/1/menu (All Roles)"
curl -s -X GET -H "$HEADER_AUTH" "$BASE_URL/api/v1/restaurants/1/menu" | print_json
echo -e "\n---------------------------------------\n"

echo "Testing POST /api/v1/orders (Create Order - All Roles)"
curl -s -X POST -H "$HEADER_AUTH" -H "$HEADER_CONTENT" -d '{"restaurant_id": 1, "items": [{"menu_item_id": 101, "quantity": 2}]}' "$BASE_URL/api/v1/orders" | print_json
echo -e "\n---------------------------------------\n"

echo "Testing GET /api/v1/orders (List Orders)"
curl -s -X GET -H "$HEADER_AUTH" "$BASE_URL/api/v1/orders" | print_json
echo -e "\n---------------------------------------\n"

echo "Testing GET /api/v1/orders/1 (Get Order Details)"
curl -s -X GET -H "$HEADER_AUTH" "$BASE_URL/api/v1/orders/1" | print_json
echo -e "\n---------------------------------------\n"

echo "Testing POST /api/v1/orders/1/items (Add Item - All Roles)"
curl -s -X POST -H "$HEADER_AUTH" -H "$HEADER_CONTENT" -d '{"menu_item_id": 102, "quantity": 1}' "$BASE_URL/api/v1/orders/1/items" | print_json
echo -e "\n---------------------------------------\n"

echo "Testing POST /api/v1/orders/1/checkout (Manager/Admin Only)"
curl -s -X POST -H "$HEADER_AUTH" -H "$HEADER_CONTENT" -d '{}' "$BASE_URL/api/v1/orders/1/checkout" | print_json
echo -e "\n---------------------------------------\n"

echo "Testing POST /api/v1/orders/1/cancel (Manager/Admin Only)"
curl -s -X POST -H "$HEADER_AUTH" -H "$HEADER_CONTENT" -d '{}' "$BASE_URL/api/v1/orders/1/cancel" | print_json
echo -e "\n---------------------------------------\n"

echo "Testing PUT /api/v1/payment-methods/1 (Admin Only - EXPECT 403 if placeholder user isn't Admin)"
curl -s -X PUT -H "$HEADER_AUTH" -H "$HEADER_CONTENT" -d '{"details": "new_details"}' "$BASE_URL/api/v1/payment-methods/1" | print_json
echo -e "\n---------------------------------------\n"

echo "Testing GET /api/v1/users/5/payment-methods (Admin Only - EXPECT 403 if placeholder user isn't Admin)"
curl -s -X GET -H "$HEADER_AUTH" "$BASE_URL/api/v1/users/5/payment-methods" | print_json
echo -e "\n---------------------------------------\n"

echo "Testing POST /api/v1/users/5/payment-methods (Admin Only - EXPECT 403 if placeholder user isn't Admin)"
curl -s -X POST -H "$HEADER_AUTH" -H "$HEADER_CONTENT" -d '{"method_type": "upi", "details": "user5@upi"}' "$BASE_URL/api/v1/users/5/payment-methods" | print_json
echo -e "\n---------------------------------------\n"

echo "======================================="
echo " API Endpoint Tests Finished"
echo "======================================="
