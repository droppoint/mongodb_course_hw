rs.initiate({
  _id: "shard0",
  members: [
    { _id: 0, host: "mongodb-shard0-replica0:27017" },
    { _id: 1, host: "mongodb-shard0-replica1:27017" },
    { _id: 2, host: "mongodb-shard0-replica2:27017" }
  ]
});