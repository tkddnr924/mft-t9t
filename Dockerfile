FROM golang:1.19.6

ENV GOOS=windows \
    GOARCH=amd64 \
    CGO_ENABLED=0 \
    GO111MODULE=auto

WORKDIR /go/src/app

RUN apt-get update
RUN apt-get upgrade -y

# github clone
RUN git clone https://github.com/tkddnr924/mft-t9t.git
RUN cd mft-t9t

# download module
RUN go get -u github.com/t9t/gomft/mft
RUN go get -u github.com/t9t/gomft/bootsect

# build
RUN go build /go/src/app/mft-t9t/main.go