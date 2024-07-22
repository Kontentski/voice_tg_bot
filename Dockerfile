FROM golang:latest AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -tags netgo -ldflags "-extldflags '-static'" -o /bot .

FROM alpine:latest

# Ffmpeg for converting .mp4 to .oga
RUN apk add --no-cache ffmpeg

COPY --from=builder /bot /bot
EXPOSE 8080
CMD ["/bot"]
