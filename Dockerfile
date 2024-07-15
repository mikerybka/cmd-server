FROM alpine:latest

RUN apk update
RUN apk add go
COPY . .

CMD go run github.com/mikerybka/cmd-server/cmd/server
