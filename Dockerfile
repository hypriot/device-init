FROM debian:jessie

ENV TARBALL "https://storage.googleapis.com/golang/go1.5.3.linux-amd64.tar.gz"
ENV UNTARPATH "/opt"
ENV GOROOT "${UNTARPATH}/go"
ENV GOPATH "${UNTARPATH}/gopath"
ENV PATH "${GOROOT}/bin:${GOPATH}/bin:/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin"
RUN apt-get update && \
    apt-get -y install \
    build-essential \
    sudo \
    file \
    wget \
    git-core \
    ruby \
    ruby-dev \
    --no-install-recommends && \
    rm -rf /var/lib/apt/lists/*

# install Docker and dependencies
RUN apt-get update && \
    apt-get install -y apt-transport-https ca-certificates && \
    apt-key adv --keyserver hkp://p80.pool.sks-keyservers.net:80 --recv-keys 58118E89F3A912897C070ADBF76221572C52609D && \
    echo "deb https://apt.dockerproject.org/repo debian-jessie main" > /etc/apt/sources.list.d/docker.list && \
    apt-get update && \
    apt-get install -y docker-engine=1.9.1-0~jessie

RUN gem update --system && \
    gem install --no-document serverspec \
    pry-byebug \
    bundler

RUN sudo wget --quiet --output-document - ${TARBALL} | sudo tar xfz - -C ${UNTARPATH}

RUN go get github.com/tools/godep
RUN go get github.com/mitchellh/gox

CMD bash
