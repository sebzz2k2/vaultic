FROM golang:1.23

RUN apt-get update && apt-get install -y netcat-openbsd && rm -rf /var/lib/apt/lists/*

RUN curl -sSL https://taskfile.dev/install.sh | sh -s -- -d && \
    mv ./bin/task /usr/local/bin/task && \
    rm -rf ./bin

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY . .
RUN task build
EXPOSE 5381

CMD ["task", "run"]