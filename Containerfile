FROM golang:1.22.3-bullseye AS build
WORKDIR /usr/app
ENV CGO_ENABLED=0

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go fmt ./... && \
    go mod tidy -v

ENV GOCACHE=/usr/app/.cache/go-build
RUN --mount=type=cache,target="/usr/app/.cache/go-build" go build -ldflags='-s -w' -o=./build/api ./cmd/api

FROM gcr.io/distroless/static AS final
ENV APP_HOME=/home/app
WORKDIR $APP_HOME
COPY --from=build /usr/app/build/api ./api

EXPOSE 8080

CMD ["./api","-port","8080"]