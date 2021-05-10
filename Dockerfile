FROM alpine:3
WORKDIR /clazz
COPY bin/app ./app
CMD [ "/clazz/app" ]