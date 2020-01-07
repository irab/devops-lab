# Simple Go API endpoint using HTTProuter to provide API routing

## Build

```
CGO_ENABLED=0 go build -ldflags="-s -w"
```

Syntax/linting with VSCode plugins

[HTTProuter](https://github.com/julienschmidt/httprouter) used as it provides [lower latency and better performance](https://github.com/julienschmidt/go-http-routing-benchmark) than the default http.ServeMux and all other HTTP/API router frameworks.

Added dependency tracking with Go Modules (go.mod)

Added Cloudbuild.yaml to build the binary and store it in a Google Cloud Storage bucket called *anz-test*.

Added go build compiler flags -s and -w to reduce binary size by remove symbol table and debugging information. This reduces the ability to debug the application but it does reduce the binary size from 7.4MB to 5.8MB, a 22% reduction. Also added 

Updated build process to remove dynamic links

## Troubleshooting

You can access the docker binary container by attaching it to the a debugging container, as shown below

Run the container:
```bash
docker run -e VERSION=1.2.3 -e LASTCOMMITSHA=asdas12312 gcr.io/ira-nz/anz-test:latest
```

Attach an alpine container:
```bash
docker run --rm -it --pid container:$(docker container ls -l -q) alpine
```

check files:
```bash
ls /proc/1/root/
anz-test  dev       etc       proc      sys
```

References:
 - https://www.ianlewis.org/en/building-go-applications-google-container-builder
 - https://hub.docker.com/_/scratch
 - https://stackoverflow.com/questions/55106186/no-such-file-or-directory-with-docker-scratch-image
