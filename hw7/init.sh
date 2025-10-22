#!/bin/bash
set -e

echo "Initializing replicaset..."
docker compose exec mongodb-rs0 mongosh /scripts/init-replicaset.js


echo "Initialization complete!"