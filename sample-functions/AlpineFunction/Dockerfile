FROM alpine:3.8

RUN apk --no-cache add curl \
    && curl -sL https://github.com/Lambda-NIC/faas/releases/download/0.9.6/fwatchdog > /usr/bin/fwatchdog \
    && chmod +x /usr/bin/fwatchdog

# Populate example here
# ENV fprocess="wc -l"

HEALTHCHECK --interval=5s CMD [ -e /tmp/.lock ] || exit 1
CMD ["fwatchdog"]
