# Сборка статического бинарника (SQLite через modernc — CGO не нужен).
FROM golang:1.25-alpine AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /out/server .

FROM ubuntu:24.04
RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates \
	&& rm -rf /var/lib/apt/lists/*
WORKDIR /app
COPY --from=build /out/server /app/server
COPY web /app/web
EXPOSE 7540
ENV TODO_PORT=7540
ENV TODO_DBFILE=/app/data/scheduler.db
RUN mkdir -p /app/data
CMD ["/app/server"]
