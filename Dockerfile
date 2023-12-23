# syntax=docker/dockerfile:1

FROM golang:1.21  

WORKDIR /app 

COPY go.mod go.sum main.go ./
COPY internal ./internal

RUN ls

RUN go mod download            

RUN CGO_ENABLED=0 GOOS=linux go build -o /good-threads-backend

EXPOSE 3000

CMD ["/good-threads-backend"]

