#!/bin/bash
set -e

echo "Initializing config servers..."
docker compose exec mongodb-configdb-replica0 mongosh /scripts/init-configserver.js

echo "Initializing shard0..."
docker compose exec mongodb-shard0-replica0 mongosh /scripts/init-shard0.js

echo "Initializing shard1..."
docker compose exec mongodb-shard1-replica0 mongosh /scripts/init-shard1.js

echo "Initializing shard2..."
docker compose exec mongodb-shard2-replica0 mongosh /scripts/init-shard2.js

echo "Adding shards to cluster..."
docker compose exec mongos mongosh --port 27017 /scripts/init-router.js

echo "Creating root user via mongos..."
docker compose exec mongos mongosh --port 27017 /scripts/create-admin.js

echo "Cluster initialization complete!"