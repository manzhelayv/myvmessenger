FROM golang:1.22 as builder

ENV CGO_ENABLED=0

COPY . /app
WORKDIR /app

#RUN  go build -v -a -tags server -installsuffix server -ldflags "-X main.version=${version} -s -w" -o /go/bin/server -mod vendor
RUN go build -a -installsuffix cgo -o main /app/*.go
#    && apt install curl && apt install bash

# контейнер окружения
#FROM scratch
#COPY --from=builder /go/bin/server /bin/server