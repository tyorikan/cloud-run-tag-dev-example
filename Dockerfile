FROM golang:1.21 as builder
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -v -o api ./cmd/api

FROM gcr.io/distroless/static-debian12
COPY --from=builder /app/api /api
CMD ["/api"]