FROM golang:1.24-alpine AS BUILDER

ENV GOCACHE=/app/.cache/go-build
ENV GOMODCACHE=/app/.cache/go-mod

WORKDIR /usr/src/app

COPY go.mod go.sum ./

RUN --mount=type=cache,target=/app/.cache/go-build \
    --mount=type=cache,target=/app/.cache/go-mod \
        go mod download && go mod verify

COPY . .

RUN --mount=type=cache,target=/app/.cache/go-build \
    --mount=type=cache,target=/app/.cache/go-mod \
    GOOS=linux go build \
    -ldflags='-w -s -extldflags "-static"' \
    -o wordlee-app-backend .

FROM scratch

ENV SERVICE_PORT=8080

COPY --from=BUILDER /usr/src/app/wordlee-app-backend /wordlee-app-backend

EXPOSE $SERVICE_PORT

CMD [ "./wordlee-app-backend" ]
