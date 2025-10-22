const username = process.env.MONGO_INITDB_ROOT_USERNAME || "root";
const password = process.env.MONGO_INITDB_ROOT_PASSWORD || "password";

db = db.getSiblingDB("admin");

db.auth(username, password);

rs.initiate({
  _id: "rs",
  members: [
    { _id: 0, host: "mongodb-rs0:27017" },
    { _id: 1, host: "mongodb-rs1:27017" },
    { _id: 2, host: "mongodb-rs2:27017" }
  ],
  settings:{"keyFile": "/etc/secrets/keyfile"}
});
