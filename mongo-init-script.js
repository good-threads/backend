goodthreads = db.getSiblingDB('goodthreads');

const users = goodthreads['users']
const sessions = goodthreads['sessions']

users.createIndex({ name: 1 }, { unique: true });
sessions.createIndex({ last_update_date: 1 }, { expireAfterSeconds: 8*3600 });

print('mongo-init-script.js successfully ran.');