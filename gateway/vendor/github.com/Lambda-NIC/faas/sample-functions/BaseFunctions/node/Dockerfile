FROM node:6.9.1-alpine

RUN apk --no-cache add curl \
    && curl -sL https://github.com/Lambda-NIC/faas/releases/download/0.9.6/fwatchdog > /usr/bin/fwatchdog \
    && chmod +x /usr/bin/fwatchdog
WORKDIR /root/

COPY package.json .

RUN npm i
COPY handler.js .

ENV fprocess="node handler.js"

HEALTHCHECK --interval=1s CMD [ -e /tmp/.lock ] || exit 1

CMD ["fwatchdog"]

