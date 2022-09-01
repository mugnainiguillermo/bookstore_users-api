# TODO Make it work

# Start from base image 1.19:
FROM golang:1.19

ENV ELASTIC_HOSTS=localhost:9200
ENV LOG_LEVEL=info

# Configure the repo url so we can configure our work directory:
ENV REPO_URL=github.com/mugnainiguillermo/bookstore_users-api

# Setup out $GOPATH
ENV GOPATH=/app

ENV APP_PATH=$GOPATH/src/$REPO_URL

# /app/src/github.com/mugnainiguillermo/bookstore_users-api/src

# Copy the entire source code from the current directory to $WORKPATH
ENV WORKPATH=$APP_PATH/src
COPY src $WORKPATH
WORKDIR $WORKPATH

RUN go build -o items-api .

EXPOSE 8080

CMD ["./users-api"]