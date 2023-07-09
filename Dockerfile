# https://klotzandrew.com/blog/smallest-golang-docker-image/
# Dockerfile.debian
FROM golang:1.20-bullseye

RUN adduser \
  --disabled-password \
  --gecos "" \
  --home "/nonexistent" \
  --shell "/sbin/nologin" \
  --no-create-home \
  --uid 65532 \
  small-user

WORKDIR $GOPATH/src/smallest-golang/app/

COPY . .

RUN go mod download
RUN go mod verify

RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o /chibi .

USER small-user:small-user

CMD ["/chibi"]

