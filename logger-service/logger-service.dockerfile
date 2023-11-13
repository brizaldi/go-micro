# build a tiny docker image
FROM alpine:latest

RUN mkdir /app

COPY /build/loggerServiceApp /app

CMD [ "/app/loggerServiceApp" ]