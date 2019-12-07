FROM golang:1.13.1 as build

WORKDIR /build

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY . .

RUN go build -o kafql ./cmd/kafql

FROM build as test

CMD go test -race -coverprofile=/artifacts/coverage.txt -covermode=atomic ./...

FROM ubuntu:18.04 as final

WORKDIR /app

COPY --from=0 /build/kafql /app/

CMD ["/app/kafql"]
