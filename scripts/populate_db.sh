#!/bin/bash

set -e

API_URL="${API_URL:-http://localhost:8080/v1/devices}"

# Check if API is running
if ! curl -s -f -o /dev/null --connect-timeout 3 "$API_URL" 2>/dev/null; then
  echo "❌ Error: API is not running!"
  echo "Please start the API first with: go run cmd/api/main.go"
  exit 1
fi

echo "🔧 Populating device database..."

devices=(
  '{"name":"iPhone 15 Pro","brand":"Apple"}'
  '{"name":"Galaxy S24","brand":"Samsung"}'
  '{"name":"Pixel 8","brand":"Google"}'
  '{"name":"MacBook Pro","brand":"Apple"}'
  '{"name":"Surface Laptop","brand":"Microsoft"}'
  '{"name":"ThinkPad X1","brand":"Lenovo"}'
  '{"name":"XPS 15","brand":"Dell"}'
  '{"name":"iPad Pro","brand":"Apple"}'
  '{"name":"Galaxy Tab","brand":"Samsung"}'
  '{"name":"Pixel Tablet","brand":"Google"}'
)

count=0
for device in "${devices[@]}"; do
  if curl -sf -X POST "$API_URL" -H "Content-Type: application/json" -d "$device" > /dev/null; then
    ((count++))
    echo "✓ Created device $count"
  else
    echo "✗ Failed to create device"
  fi
done

echo "✅ Done! Created $count devices"
