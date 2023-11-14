# build a tiny docker image
FROM alpine:latest

RUN mkdir /app

COPY /build/mailerApp /app

CMD [ "/app/mailerApp" ]