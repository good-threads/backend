db = db.getSiblingDB('goodthreads')
db['users'].createIndex({ ['name']: 1 }, { unique: true });

print('mongo-init-script.js successfully ran.');