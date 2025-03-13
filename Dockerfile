FROM golang:1.24-alpine AS builder
WORKDIR /App
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o app
FROM scratch
WORKDIR /root/
COPY --from=builder /App/app .
ENV STORAGE_MODE=inmemory
ENV DATABASE_URL=""
EXPOSE 3000
CMD ["./app"]