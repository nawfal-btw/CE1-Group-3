FROM golang:1.26.8-bookworm AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /api ./api

FROM debian:bookworm-slim
WORKDIR /app
COPY --from=build /api /app/api
COPY sql/migrations /app/sql/migrations
EXPOSE 8080
CMD ["./api"]
