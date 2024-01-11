curl http://localhost:8080/ -b cookies
curl http://localhost:8080/session -d '{"username":"tom","password":"pepe123"}' -c cookies
curl http://localhost:8080/user -d '{"username":"tom","password":"pepe123"}'
curl http://localhost:8080/session -d '{"username":"tom","password":"pepe123"}' -c cookies
curl http://localhost:8080/ -b cookies
curl http://localhost:8080/ -b cookies -X PATCH -d '{
    "lastProcessedCommandID": null,
    "commands": [
        {
            "id": "c1",
            "datetime": "'"$(date --rfc-3339=seconds | sed -e 's/ /T/')"'",
            "type": "createThread",
            "payload": {
                "id": "t1",
                "name": "the pepe thread"
            }
        },
        {
            "id": "c2",
            "datetime": "'"$(date --rfc-3339=seconds | sed -e 's/ /T/')"'",
            "type": "createKnot",
            "payload": {
                "threadID": "t1",
                "knotID": "k1",
                "knotBody": "the pepe knot"
            }
        }
    ]
}'
curl http://localhost:8080/ -b cookies