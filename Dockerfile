FROM 092574054921.dkr.ecr.eu-west-1.amazonaws.com/remote/dockerhub/library/golang:1.24-bookworm

ENV GOPROXY http://goproxy.appsflyer.com,https://proxy.golang.org,direct
ENV GONOSUMDB gitlab.appsflyer.com,github.com/appsflyerrnd/*

RUN useradd --create-home --shell /bin/bash docker && \
    apt-get update && \
    apt-get -y install ca-certificates curl && \
    rm -rf /var/lib/apt/lists/*

USER docker
WORKDIR /home/docker

COPY go.mod ./

# by running the go mod tidy command before copying the other Go files, subsequent build will be faster when there are only changes in Code and not in Dependencies
RUN go mod download

# now we shuold copy the code files
COPY cmd ./cmd
COPY docs ./docs
COPY internal ./internal

# change the path here to point to your main.go file
RUN go build -o server cmd/server/main.go && chmod +x server

EXPOSE 11666 8090

CMD ["./server"]
