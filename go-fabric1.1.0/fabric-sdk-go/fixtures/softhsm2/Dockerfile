FROM hyperledger/fabric-baseimage

ENV GOPATH=/opt/gopath \
    GOROOT=/opt/go \
    GO_VERSION=1.9.2 \
    PATH=$PATH:/opt/go/bin:/opt/gopath/bin

COPY test/fixtures/softhsm2/install-softhsm2.sh /tmp
COPY scripts/_go/src/pkcs11helper /opt/gopath/src/pkcs11helper
RUN bash /tmp/install-softhsm2.sh
CMD ["/bin/bash"]
