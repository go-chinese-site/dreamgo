FROM golang

RUN mkdir /dreamgo
RUN mkdir /dreamgo/log
ADD ./bin /dreamgo/bin
ADD ./config /dreamgo/config
ADD ./static /dreamgo/static
ADD ./template /dreamgo/template

EXPOSE 2017

WORKDIR /dreamgo

# Define default command.
CMD [ "/dreamgo/bin/dreamgo" ]
