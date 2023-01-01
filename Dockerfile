FROM golang:1.19.4-alpine3.17 as build-base

WORKDIR /app

COPY ./database/ ./database/

COPY ./handler/ ./handler/

COPY *.go .

COPY go.mod .

RUN go get

RUN go mod download

RUN go build -o ./build/server .


# =======================================================

FROM alpine:3.17

COPY --from=build-base /app/build/server /app/server

CMD ["/app/server"]


#### RUN Docker
#    docker build -t my/assessment:golang .

#    docker run --env-file=env_file --name my_app_expense -p 2565:2565 -it my/assessment:golang
