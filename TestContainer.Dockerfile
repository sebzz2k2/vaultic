FROM golang:1.23

RUN curl -sSL https://taskfile.dev/install.sh | sh -s -- -d && \
    mv ./bin/task /usr/local/bin/task && \
    rm -rf ./bin
# RUN -c 'sh -c "$(curl --location https://taskfile.dev/install.sh)" -- -d'
# RUN 
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY . .
RUN task build
EXPOSE 5381

CMD ["task", "run"]
# CMD ["ls"]
