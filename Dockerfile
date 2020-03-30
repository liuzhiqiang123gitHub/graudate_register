FROM golang

ARG app_env
ENV APP_ENV $app_env

COPY . /go/src/backend/graduate_registrator
WORKDIR /go/src/backend/graduate_registrator
RUN export GOPROXY="http://goproxy.cn"
RUN go build

CMD if [ ${APP_ENV} = production ]; \
        then \
        app; \
        else \
        go get github.com/pilu/fresh && \
        fresh; \
        fi

EXPOSE 9370