FROM golang:1.22-alpine as build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY . .

FROM build as migration
RUN go build -v -o migrator-tool ./cmd/migrator
ENV CONFIG_PATH=/app/configs/local.docker.json
RUN ./migrator-tool -migrations-path=/app/migrations -direction=up as migration

FROM migration as run
RUN go build -v -o app ./cmd/app

EXPOSE 22313

CMD ["/app/app"]