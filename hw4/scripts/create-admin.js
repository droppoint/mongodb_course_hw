const username = process.env.MONGO_ROOT_USERNAME || "root";
const password = process.env.MONGO_ROOT_PASSWORD || "password";

db = db.getSiblingDB("admin");

db.createUser({
  user: username,
  pwd: password,
  roles: [
    { role: "root", db: "admin" }
  ]
});
