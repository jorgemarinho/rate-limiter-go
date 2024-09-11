FROM golang:1.22.4 as build
WORKDIR /app
COPY . .
RUN GOOS=linux CGO_ENABLED=0 go build -ldflags="-w -s" -o rate-limiter-go cmd/server/main.go

FROM scratch
WORKDIR /app
COPY --from=build /app/rate-limiter-go .
ENTRYPOINT ["./rate-limiter-go"]