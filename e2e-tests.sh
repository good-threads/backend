# this basic test script / request list can be used as a means to understand
# the way the API works. it's meant to be read from top to botseba,
# as the requests are ordered in a way that shows how the state of the
# application evolves.

# get board
#
# as no session cookie is set,
# a 401 should be returned
#
# expected response body:
# {"message":"Invalid session"}
curl http://localhost:8000/ -b cookies

# create session
#
# as the credentials don't match any existing user,
# a 401 should be returned
#
# expected response body:
# {"message":"Wrong credentials"}
curl http://localhost:8000/session -d '{"username":"seba","password":"pepeasdasd"}' -c cookies

# create user
#
# as the request was made just fine,
# a 201 should be returned
#
# expected response body:
# {"message":"User created"}
curl http://localhost:8000/user -d '{"username":"seba","password":"pepe123"}'

# create session
#
# as the request was made just fine,
# a 201 should be returned
#
# expected response body:
# {"message":"Session created"}
curl http://localhost:8000/session -d '{"username":"seba","password":"pepe123"}' -c cookies

# get board
#
# as the request was made just fine (it now includes a valid session cookie),
# a 200 should be returned, with an empty board
#
# expected response body:
# {"activeThreads":[],"hiddenThreads":[],"lastProcessedCommandID":null}
curl http://localhost:8000/ -b cookies

# patch board
#
# as the request was made just fine,
# a 200 should be returned, indicating the last processed command id (should be c10)
#
# expected response body:
# {"lastProcessedCommandID":"c10"}
curl http://localhost:8000/ -b cookies -X PATCH -d '{
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
            "type": "editKnotBody",
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

# get board
#
# as the request was made just fine (and changes were made to the board),
# a 200 should be returned, showing the new state of the board
#
# expected response body:
# {"activeThreads":[{"id":"t3","name":"the amazing thread","knots":[]},{"id":"t1","name":"the pepe thread","knots":[{"id":"k2","body":"the edit is real"},{"id":"k3","body":"the last knot"}]}],"hiddenThreads":["t2"],"lastProcessedCommandID":"c10"}
curl http://localhost:8000/ -b cookies

# get thread
#
# as the request was made just fine,
# a 200 should be returned, showing the hidden thread
#
# expected response body:
# {"id":"t2","name":"the pepest thread","knots":[]}
curl http://localhost:8000/thread/t2 -b cookies
