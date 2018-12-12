FROM openjdk:8u121-jdk-alpine

RUN apk --no-cache add curl \
    && curl -sL https://github.com/Lambda-NIC/faas/releases/download/0.9.6/fwatchdog > /usr/bin/fwatchdog \
    && chmod +x /usr/bin/fwatchdog

WORKDIR /root/

COPY Handler.java .
RUN javac Handler.java

ENV fprocess="java Handler"

HEALTHCHECK --interval=1s CMD [ -e /tmp/.lock ] || exit 1

CMD ["fwatchdog"]

