#!/bin/bash
set -e

# Create private key file from environment variable
if [ -n "$NODE_PRIVATE_KEY" ]; then
    echo "Creating private key file from environment variable..."
    echo "$NODE_PRIVATE_KEY" > /config/node-key
    chmod 600 /config/node-key
    echo "✅ Private key file created"
else
    echo "❌ ERROR: NODE_PRIVATE_KEY environment variable not set"
    exit 1
fi

# Execute Besu with all arguments
exec besu "$@"
