# syntax=docker/dockerfile:1.7
FROM golang:1.25-bookworm AS build
WORKDIR /src
COPY . .
ARG TARGETOS TARGETARCH
RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -o /out/server ./cmd/server

FROM gcr.io/distroless/static:nonroot
COPY --from=build /out/server /server
USER 65532:65532
EXPOSE 8080
ENTRYPOINT ["/server"]
