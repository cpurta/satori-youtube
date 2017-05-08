FROM ubuntu:latest

MAINTAINER Chris Purta cpurta@gmail.com

RUN apt-get update && \
    apt-get install -y ca-certificates && \
    mkdir -p /opt

ADD ./bin/crawler /opt

RUN chmod +x /opt/crawler

WORKDIR /opt

ENTRYPOINT ["/opt/crawler"]
