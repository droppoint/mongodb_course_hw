rs.initiate({
  _id: "shard2",
  members: [
    { _id: 0, host: "mongodb-shard2-replica0:27017" },
    { _id: 1, host: "mongodb-shard2-replica1:27017" },
    { _id: 2, host: "mongodb-shard2-replica2:27017" }
  ]
});