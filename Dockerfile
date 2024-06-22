FROM golang:1.22.4 AS build-stage

WORKDIR /api

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /youtoob ./cmd/api/main.go

FROM gcr.io/distroless/base-debian11 AS release-stage

WORKDIR /

COPY --from=build-stage /youtoob /youtoob

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["/youtoob"]
