FROM golang:1.23.9 as build
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o cryple-api ./cmd

FROM scratch
WORKDIR /app
COPY --from=build /app/cryple-api .
COPY --from=build /app/migrations ./migrations
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
ENTRYPOINT ["./cryple-api"]