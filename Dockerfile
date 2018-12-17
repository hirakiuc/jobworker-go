# build stage
FROM golang:1.11 AS build-stage

RUN mkdir -p /go/src/github.com/hirakiuc/jobworker-go
WORKDIR /go/src/github.com/hirakiuc/jobworker-go

COPY . /go/src/github.com/hirakiuc/jobworker-go

# install dep
RUN curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh \
  && dep ensure \
  && make build

# build binary => /go/src/github.com/hirakiuc/jobworker-go/jobdaemon

# runtime stage
FROM alpine:3.8

COPY --from=build-stage /go/src/github.com/hirakiuc/jobworker-go/jobdaemon /app/jobdaemon

ENTRYPOINT ["/app/jobdaemon"]

CMD ["start"]
