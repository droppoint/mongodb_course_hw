#!/bin/bash
set -e

echo "Initializing config servers..."
docker compose exec mongodb-configdb-replica0 mongosh /scripts/init-configserver.js

echo "Initializing shard0..."
docker compose exec mongodb-shard0-replica0 mongosh /scripts/init-shard0.js
sleep 10
docker compose exec mongodb-shard0-replica0 mongosh /scripts/create-users.js

echo "Initializing shard1..."
docker compose exec mongodb-shard1-replica0 mongosh /scripts/init-shard1.js
sleep 10
docker compose exec mongodb-shard1-replica0 mongosh /scripts/create-users.js

sleep 20

echo "Adding shards to cluster..."
docker compose exec mongos mongosh --port 27017 /scripts/init-router.js

echo "Creating root user via mongos..."
docker compose exec mongos mongosh --port 27017 /scripts/create-users.js

echo "Applying config to pbm-agents"
docker compose exec pbm-agent-configdb-replica0 bash -c "pbm config --file /etc/pbm-config.yaml"
docker compose exec pbm-agent-configdb-replica1 bash -c "pbm config --file /etc/pbm-config.yaml"
docker compose exec pbm-agent-configdb-replica2 bash -c "pbm config --file /etc/pbm-config.yaml"
docker compose exec pbm-agent-shard0-replica0 bash -c "pbm config --file /etc/pbm-config.yaml"
docker compose exec pbm-agent-shard0-replica1 bash -c "pbm config --file /etc/pbm-config.yaml"
docker compose exec pbm-agent-shard0-replica2 bash -c "pbm config --file /etc/pbm-config.yaml"
docker compose exec pbm-agent-shard1-replica0 bash -c "pbm config --file /etc/pbm-config.yaml"
docker compose exec pbm-agent-shard1-replica1 bash -c "pbm config --file /etc/pbm-config.yaml"
docker compose exec pbm-agent-shard1-replica2 bash -c "pbm config --file /etc/pbm-config.yaml"


echo "Cluster initialization complete!"