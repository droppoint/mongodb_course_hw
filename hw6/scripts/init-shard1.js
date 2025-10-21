rs.initiate({
  _id: "shard1",
  members: [
    { _id: 0, host: "mongodb-shard1-replica0:27017" },
    { _id: 1, host: "mongodb-shard1-replica1:27017" },
    { _id: 2, host: "mongodb-shard1-replica2:27017" }
  ]
});