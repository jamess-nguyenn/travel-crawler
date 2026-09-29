FROM golang:1.27

RUN apt-get update && apt-get install -y \
    git \
    build-essential \
    curl \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*

# Install Delve debugger
RUN go install github.com/go-delve/delve/cmd/dlv@latest

# Install Goose migration tool
RUN go install github.com/pressly/goose/v3/cmd/goose@v3.26.0

WORKDIR /app

EXPOSE 8000
EXPOSE 2345
