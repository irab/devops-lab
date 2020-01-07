# Simple Go API endpoint using HTTProuter to provide API routing

Syntax/linting with VSCode plugins

[HTTProuter](https://github.com/julienschmidt/httprouter) used as it provides [lower latency and better performance](https://github.com/julienschmidt/go-http-routing-benchmark) than the default http.ServeMux and all other HTTP/API router frameworks.

Added dependency tracking with Go Modules (go.mod)

Added Cloudbuild.yaml to build the binary and store it in a Google Cloud Storage bucket called *anz-test*.

Added go build compiler flags -s and -w to reduce binary size by remove symbol table and debugging information. This reduces the ability to debug the application but it does reduce the binary size from 7.4MB to 5.8MB, a 22% reduction. Also added 

Updated build process to remove dynamic links

References:
 - https://www.ianlewis.org/en/building-go-applications-google-container-builder
 - https://hub.docker.com/_/scratch
 - https://stackoverflow.com/questions/55106186/no-such-file-or-directory-with-docker-scratch-image
