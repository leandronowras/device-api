#!/bin/bash

set -e

API_URL="${API_URL:-http://localhost:8080/v1/devices}"

# Check if API is running
if ! curl -s -f -o /dev/null --connect-timeout 3 "$API_URL" 2>/dev/null; then
  echo "❌ Error: API is not running!"
  echo "Please start the API first with: go run cmd/api/main.go"
  exit 1
fi

echo "📖 Testing read-only operations..."
echo ""

# List all devices
echo "1️⃣  List all devices:"
curl -s "$API_URL" | jq '.'
echo ""

# Get first device ID
DEVICE_ID=$(curl -s "$API_URL" | jq -r '.[0].id // empty')

if [ -n "$DEVICE_ID" ]; then
  echo "2️⃣  Get device by ID ($DEVICE_ID):"
  curl -s "$API_URL/$DEVICE_ID" | jq '.'
  echo ""
fi

# Filter by brand
echo "3️⃣  Filter by brand (Apple):"
curl -s "$API_URL?brand=Apple" | jq '.'
echo ""

# Filter by state
echo "4️⃣  Filter by state (available):"
curl -s "$API_URL?state=available" | jq '.'
echo ""

# Pagination
echo "5️⃣  Pagination (page=1, limit=3):"
curl -s "$API_URL?page=1&limit=3" | jq '.'
echo ""

echo "6️⃣  Pagination (page=2, limit=3):"
curl -s "$API_URL?page=2&limit=3" | jq '.'
echo ""

# Combined filters
echo "7️⃣  Combined filters (brand=Samsung, state=available):"
curl -s "$API_URL?brand=Samsung&state=available" | jq '.'
echo ""

# Statistics
echo "📊 Statistics:"
echo ""
echo "Total devices:"
curl -s "$API_URL" | jq 'length'
echo ""

echo "Devices by brand:"
curl -s "$API_URL" | jq -r 'group_by(.brand) | map({brand: .[0].brand, count: length}) | sort_by(-.count) | .[] | "  \(.brand): \(.count)"'
echo ""

echo "Devices by state:"
curl -s "$API_URL" | jq -r 'group_by(.state) | map({state: .[0].state, count: length}) | .[] | "  \(.state): \(.count)"'
echo ""

echo "✅ Read-only operations complete!"
