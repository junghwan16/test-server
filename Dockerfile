FROM golang:1.25.1 AS builder
WORKDIR /app
COPY . . 
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o server main.go

FROM gcr.io/distroless/static:nonroot
COPY --from=builder /app/server /server
USER 65532:65532
EXPOSE 8080
ENTRYPOINT [ "/server" ]