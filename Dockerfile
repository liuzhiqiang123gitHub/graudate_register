FROM golang

ARG app_env
ENV APP_ENV $app_env

COPY . /go/src/backend/graduate_registrator
WORKDIR /go/src/backend/graduate_registrator

RUN go build
RUN export GOPROXY="http://goproxy.cn"
CMD if [ ${APP_ENV} = production ]; \
        then \
        app; \
        else \
        go get github.com/pilu/fresh && \
        fresh; \
        fi

EXPOSE 9370