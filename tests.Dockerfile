FROM bash:4
ENV BATS_VERSION=v1.7.0
RUN apk update && apk add git
RUN git clone --branch $BATS_VERSION https://github.com/bats-core/bats-core.git /tmp/bats-core
RUN cd /tmp/bats-core && ./install.sh /usr/local && bats --version
RUN mkdir /helpers
RUN git clone https:g/github.com/bats-core/bats-assert.git /helpers/bats-assert

ENTRYPOINT [ "bats" ]
