# build a tiny docker image
FROM alpine:latest

RUN mkdir /app

COPY /build/brokerApp /app

CMD [ "/app/brokerApp" ]