FROM golang

ADD . /dreamgo
ENV GO111MODULE=on
ENV GOPROXY=https://goproxy.cn
# install dreamgo
WORKDIR /dreamgo
RUN ./build.sh
EXPOSE 2017

ENTRYPOINT [ "./bin/dreamgo" ]
