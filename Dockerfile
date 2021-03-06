FROM golang:latest as BUILDER

MAINTAINER TommyLike<tommylikehu@gmail.com>

# build binary
COPY . /go/src/github.com/opensourceways/app-cla-server
RUN cd /go/src/github.com/opensourceways/app-cla-server && CGO_ENABLED=1 go build -v -o ./cla-server main.go

# copy binary config and utils
FROM golang:latest
RUN apt-get update && apt-get install -y python-pip && mkdir -p /opt/app/
COPY ./conf /opt/app/conf
COPY ./util/merge-signature.py /opt/app/util
# overwrite config yaml
COPY ./deploy/app.conf /opt/app/conf
COPY  --from=BUILDER /go/src/github.com/opensourceways/app-cla-server/cla-server /opt/app

WORKDIR /opt/app/
ENTRYPOINT ["/opt/app/cla-server"]
