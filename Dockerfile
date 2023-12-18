# syntax=docker/dockerfile:1

FROM golang:1.21  

WORKDIR /app 

COPY . .

RUN go mod download            

RUN CGO_ENABLED=0 GOOS=linux go build -o /good-threads-backend

EXPOSE 3000

CMD ["/good-threads-backend"]

