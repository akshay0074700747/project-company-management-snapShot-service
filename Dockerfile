FROM golang:1.21.5-bullseye AS build

RUN apt-get update && apt-get install -y curl libcurl-dev

WORKDIR /app

COPY . .

RUN go mod download

WORKDIR /app/cmd

RUN go build -o snapshot-service

FROM busybox:latest

WORKDIR /snapshot-service

COPY --from=build /app/cmd/snapshot-service .

COPY --from=build /app/cmd/.env .

EXPOSE 50005

CMD [ "./snapshot-service" ]