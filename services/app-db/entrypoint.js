let db = connect("mongodb://admin:pass@localhost:27017/admin");

const USER = process.env.MONGO_USER
const PASSWORD = process.env.MONGO_PASSWORD
const DATABASE = process.env.MONGO_DATABASE

db = db.getSiblingDB(DATABASE); // we can not use "use" statement here to switch db

db.createUser(
    {
        user: USER,
        pwd: PASSWORD,
        roles: [ { role: "readWrite", db: "admin"} ],
        passwordDigestor: "server",
    }
)