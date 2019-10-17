############################
# STEP 1 test the code
############################

FROM golang:alpine as tester

ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GO111MODULE=on

# Install git + SSL ca certificates.
# Git is required for fetching the dependencies.
# Ca-certificates is required to call HTTPS endpoints.
RUN apk update && apk add --no-cache git ca-certificates tzdata && update-ca-certificates

WORKDIR $GOPATH/src/github.com/cthulhu/tw-trend/

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

CMD [ "go", "test", "./...", "-v" ]

############################
# STEP 2 build the binaries
############################
FROM tester AS builder

RUN mkdir /opt/bin
RUN go build -ldflags="-w -s" -a -installsuffix cgo -o /opt/bin/server cmd/server/main.go

############################
# STEP 3 build a small image
############################
FROM scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /opt/bin/server /bin/

ENTRYPOINT [ "/bin/server" ]