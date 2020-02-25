FROM gcr.io/distroless/static
COPY devops-lab /
EXPOSE 8080
ENTRYPOINT [ "/devops-lab" ]