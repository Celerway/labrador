#!/usr/bin/env bash



# Philips Hue Bridge IP address
BRIDGE_IP="ecb5faac4a4c.local"

# The data payload
DATA='{"devicetype":"labrador#lab"}'

# Make a POST request to create a new user
RESPONSE=$(curl -s -X POST -d "$DATA" http://$BRIDGE_IP/api)

# Check if the response contains an error
if echo "$RESPONSE" | grep -q "error"; then
  echo "Failed to create user. Make sure you pressed the bridge button."
  echo "Response from Hue Bridge: $RESPONSE"
else
  # If success, dump the relevant information
  echo "User created successfully."
  echo "Response from Hue Bridge: $RESPONSE"
fi