#!/bin/bash
# Script to add test domains to /etc/hosts
# Run with: sudo ./setup-hosts.sh

HOSTS_FILE="/etc/hosts"
DOMAINS=(
    "127.0.0.1 fider.local"
    "127.0.0.1 app.local"
    "127.0.0.1 multi.local"
    "127.0.0.1 tenant1.multi.local"
    "127.0.0.1 tenant2.multi.local"
)

echo "Adding test domains to hosts file..."

# Check if running as root
if [ "$EUID" -ne 0 ]; then
    echo "ERROR: This script must be run as root"
    echo "Run with: sudo $0"
    exit 1
fi

# Add domains if not already present
modified=false
for domain in "${DOMAINS[@]}"; do
    if ! grep -q "^$domain$" "$HOSTS_FILE"; then
        echo "$domain" >> "$HOSTS_FILE"
        echo "Added: $domain"
        modified=true
    else
        echo "Already exists: $domain"
    fi
done

if [ "$modified" = true ]; then
    echo ""
    echo "Hosts file updated successfully!"
    echo "Location: $HOSTS_FILE"
else
    echo ""
    echo "All domains already configured."
fi

echo ""
echo "You can now run: docker-compose -f docker-compose-test.yml up -d"
