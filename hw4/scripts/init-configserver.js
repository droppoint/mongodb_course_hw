rs.initiate({
  _id: "cfg",
  configsvr: true,
  members: [
    { _id: 0, host: "mongodb-configdb-replica0:27017" },
    { _id: 1, host: "mongodb-configdb-replica1:27017" },
    { _id: 2, host: "mongodb-configdb-replica2:27017" }
  ]
});
  