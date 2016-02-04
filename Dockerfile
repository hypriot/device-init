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

RUN gem update --system && \
    gem install --no-document serverspec \
    pry-byebug \
    bundler

RUN sudo wget --quiet --output-document - ${TARBALL} | sudo tar xfz - -C ${UNTARPATH}

RUN go get github.com/tools/godep
RUN go get github.com/mitchellh/gox

CMD bash
