FROM golang:1.24-bookworm AS builder

RUN apt-get -y update && apt-get -y install \
    libgl1-mesa-dev \
    libxcursor-dev \
    libxi-dev \
    libxinerama-dev \
    libxrandr-dev \
    libxxf86vm-dev \
    libasound2-dev \
    pkg-config

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=1 go build -o /app/snake-game -ldflags="-s -w" ./cmd/snake-game

FROM debian:bookworm-slim

RUN apt-get -y update && apt-get -y install \
    libgl1-mesa-glx \
    libxcursor1 \
    libxi6 \
    libxinerama1 \
    libxrandr2 \
    libxxf86vm1 \
    libasound2 \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY --from=builder /app/snake-game .
COPY --from=builder /app/levels ./levels/

CMD ["./snake-game"]