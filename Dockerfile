FROM scratch
COPY anz-test /anz-test
EXPOSE 8080
ENTRYPOINT [ "/anz-test" ]