FROM golang:1.22.2-bullseye AS build
WORKDIR /usr/app
ENV CGO_ENABLED=0

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go fmt ./... && \
    go mod tidy -v

ENV GOCACHE=/root/.cache/go-build
RUN --mount=type=cache,target="/root/.cache/go-build" go build -ldflags='-s -w' -o=./build/benktech ./cmd/web

FROM gcr.io/distroless/static AS final
ENV APP_HOME=/home/app
WORKDIR $APP_HOME
COPY --from=build /usr/app/build/benktech ./benktech
COPY --from=build /usr/app/assets ./assets

EXPOSE 8080

CMD ["./benktech","-port","8080"]