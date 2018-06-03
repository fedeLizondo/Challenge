FROM golang:latest

RUN go get github.com/gin-gonic/gin
RUN go get github.com/mattn/go-sqlite3
RUN go get github.com/go-sql-driver/mysql
RUN go get -u google.golang.org/api/drive/v3
RUN go get -u google.golang.org/api/googleapi
RUN go get -u golang.org/x/oauth2/...

ADD ./src/api /go/src/api
COPY ./src/api/client_secret.json /go/
RUN go install api
CMD ["api"]
