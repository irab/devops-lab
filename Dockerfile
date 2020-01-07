FROM gcr.io/distroless/static
COPY anz-test /
EXPOSE 8080
ENTRYPOINT [ "/anz-test" ]