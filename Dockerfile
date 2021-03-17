FROM golang:1.16 as build

COPY ./cmd /usr/src/app/cmd
COPY go.* /usr/src/app/
COPY .git /usr/src/app/

ENV GOOS=linux
ENV GOARCH=amd64
ENV CGO_ENABLED=0
ENV GOFLAGS="-trimpath"

RUN cd /usr/src/app \
  && go mod download \
  && go mod verify \
  && go build -v -o node-label-to-pod -ldflags \
  "-X main.gitVersion=$(git describe --tags `git rev-list --tags --max-count=1`)-$(date +%Y%m%d%H%M%S)-$(git log -n1 --pretty='%h')" \
  ./cmd \
  && /usr/src/app/node-label-to-pod -version

FROM alpine:3

COPY --from=build /usr/src/app/node-label-to-pod /usr/local/bin/node-label-to-pod

RUN addgroup -g 101 -S app \
&& adduser -u 101 -D -S -G app app

USER 101

CMD /usr/local/bin/node-label-to-pod