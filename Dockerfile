FROM golang

RUN go get github.com/polaris1119/gvt
RUN ln -sf /go/bin/gvt /usr/local/bin/
ADD . /dreamgo

# install dreamgo
WORKDIR /dreamgo
RUN ./getpkg.sh
RUN ./install.sh
EXPOSE 2017

ENTRYPOINT [ "./bin/dreamgo" ]
