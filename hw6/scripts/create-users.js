const username = process.env.MONGO_ROOT_USERNAME || "root";
const password = process.env.MONGO_ROOT_PASSWORD || "password";
const pbmUsername = process.env.PBM_USERNAME || "pbm";
const pbmPassword = process.env.PBM_PASSWORD || "password";

db = db.getSiblingDB("admin");

db.createUser({
  user: username,
  pwd: password,
  roles: [
    { role: "root", db: "admin" }
  ]
});

db.auth(username, password);

// Роль и user для Percona Backup MongoDB
db.createRole({
  role: "pbmAnyAction",
  privileges: [
    { resource: { anyResource: true }, actions: [ "anyAction" ] }
  ],
  roles: []
});

db.createUser({
  user: pbmUsername,
  pwd: pbmPassword,
  roles: [
    { role: "backup", db: "admin" },
    { role: "restore", db: "admin" },
    { role: "clusterMonitor", db: "admin" },
    { role: "pbmAnyAction", db: "admin" }
  ]
});