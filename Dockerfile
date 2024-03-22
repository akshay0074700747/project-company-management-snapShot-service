FROM golang:1.21.5-alpine AS build

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