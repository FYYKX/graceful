FROM scratch
COPY graceful /graceful
ENTRYPOINT [ "/graceful" ]