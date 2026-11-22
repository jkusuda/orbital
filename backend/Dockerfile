# Build stage: compile the API and goose (for migrations) as static binaries.
FROM golang:1.27-alpine AS build
WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/api ./cmd/api
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/goose github.com/pressly/goose/v3/cmd/goose

# Run stage: no shell or package manager, just the binaries and CA certs (for NASA HTTPS).
FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=build /out/api /out/goose /usr/local/bin/
COPY migrations /migrations

EXPOSE 8080
ENTRYPOINT ["/usr/local/bin/api"]
