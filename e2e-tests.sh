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
        },
        {
            "id": "c3",
            "datetime": "'"$(date --rfc-3339=seconds | sed -e 's/ /T/')"'",
            "type": "createThread",
            "payload": {
                "id": "t2",
                "name": "the pepest thread"
            }
        },
        {
            "id": "c4",
            "datetime": "'"$(date --rfc-3339=seconds | sed -e 's/ /T/')"'",
            "type": "createThread",
            "payload": {
                "id": "t3",
                "name": "the amazing thread"
            }
        },
        {
            "id": "c5",
            "datetime": "'"$(date --rfc-3339=seconds | sed -e 's/ /T/')"'",
            "type": "hideThread",
            "payload": {
                "id": "t2"
            }
        },
        {
            "id": "c6",
            "datetime": "'"$(date --rfc-3339=seconds | sed -e 's/ /T/')"'",
            "type": "relocateThread",
            "payload": {
                "id": "t3",
                "newIndex": 0
            }
        },
        {
            "id": "c7",
            "datetime": "'"$(date --rfc-3339=seconds | sed -e 's/ /T/')"'",
            "type": "createKnot",
            "payload": {
                "threadID": "t1",
                "knotID": "k2",
                "knotBody": "no edit here"
            }
        },
        {
            "id": "c8",
            "datetime": "'"$(date --rfc-3339=seconds | sed -e 's/ /T/')"'",
            "type": "deleteKnot",
            "payload": {
                "threadID": "t1",
                "knotID": "k1"
            }
        },
        {
            "id": "c9",
            "datetime": "'"$(date --rfc-3339=seconds | sed -e 's/ /T/')"'",
            "type": "editKnot",
            "payload": {
                "threadID": "t1",
                "knotID": "k2",
                "knotBody": "the edit is real"
            }
        },
        {
            "id": "c10",
            "datetime": "'"$(date --rfc-3339=seconds | sed -e 's/ /T/')"'",
            "type": "createKnot",
            "payload": {
                "threadID": "t1",
                "knotID": "k3",
                "knotBody": "the last knot"
            }
        }
    ]
}'
curl http://localhost:8080/ -b cookies