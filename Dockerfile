FROM golang:1.10.8

ENV GOOS=windows \
    GOARCH=386 \
    CGO_ENABLED=1 \
    GO111MODULE=auto \
    CC=i686-w64-mingw32-gcc

# golang 1.11 이상
# ENV GO111MODULE=auto
# ENV GOOS=windows
# ENV GOARCH=amd64
# ENV CGO_ENABLED=0

WORKDIR /go/src/app

RUN apt-get update
RUN apt-get upgrade -y

# install mingw-w64
RUN apt-get install -y mingw-w64

# github clone
RUN git clone https://github.com/tkddnr924/mft-t9t.git
RUN cd mft-t9t

# download module
RUN go get -u github.com/tkddnr924/gomft/mft
RUN go get -u github.com/tkddnr924/gomft/bootsect

# build
RUN go build /go/src/app/mft-t9t/main.go
