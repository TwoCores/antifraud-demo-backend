FROM golang:1.25-alpine AS backend

RUN apk add upx

WORKDIR /app

ADD . .

RUN go mod tidy

RUN CGO_ENABLED=0 go build -o server -ldflags="-s -w" -buildvcs=false cmd/server/main.go
RUN upx ./server

# Create a group and user
RUN addgroup -S user && adduser -S user -G user -u 1001

FROM scratch

COPY --from=backend /app/server /server

USER 1001

ENTRYPOINT ["/server"]