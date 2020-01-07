# Simple Go API endpoint using HTTProuter to provide API routing

## Intro

This repository contains a Go app that dynamically returns the version of the application and SHA commit hash along with a description.

The app uses [HTTProuter](https://github.com/julienschmidt/httprouter) for API routing, which has [lower latency and better performance](https://github.com/julienschmidt/go-http-routing-benchmark) than supplied by the default http.ServeMux and other HTTP/API router frameworks. This provides an RESTful API that's clearer and more easily extensible than custom code.

The final Docker image is running on Distroless, as building 'FROM Scratch' did not work on Cloud Run. Distroless is Google's [suggested approach](https://github.com/GoogleContainerTools/distroless/blob/master/base/README.md) to running static binaries on Google Cloud and switching to this works.

## CI/CD Pipeline

A CI/CD pipeline has been created with a Cloudbuild Github Action integration that triggers for each push to the master branch and also for every tag.

The dynamic variables (**LASTCOMMITSHA**, **VERSION**) in the app are passed to the Go app via environment variables. Theses are derived automatically from [Cloudbuild builtin substitutions](https://cloud.google.com/cloud-build/docs/configuring-builds/substitute-variable-values) - **VERSION** from **build.Source.RepoSource.Revision.TagName** and **LASTCOMMITSHA** from **build.SourceProvenance.ResolvedRepoSource.Revision.CommitSha**. This is to avoid hardcoding and ensure changes to these values come via changes to the git repository.

The cloudbuild.yaml file has the following steps:

- go_build: Builds a statically linked Go binary with the flags -s and -w to reduce binary size by removing the symbol table and debugging information.
- go_test: Run go tests
- go_security: Runs [gosec](https://github.com/securego/gosec), a code scanner that checks the AST (Abstract Syntax Tree) for security issues
- docker_image_build: Builds a minimal Docker image using the Dockerfile in the root of the repository (FROM gcr.io/distroless/static)
- cloud_run_deploy: deploys the image above to Cloud Run in us-central1

## Local development

### Building

The build command can be modified to build a binary with a minimal footprint as below:

```bash
CGO_ENABLED=0 go build -ldflags="-s -w"
```

This reduces the ability to debug/trace the application but it does reduce the binary size from 7.4MB to 5.8MB, a 22% reduction.

### Updating Dependencies with Go Modules

The go.mod will automatically updated when a new version of external dependencies are updated after a new build. Make sure it these changes are committed.

### Troubleshooting

You can access the docker binary container by attaching it to the a debugging container, as shown below.

Run the container:

```bash
docker run -e VERSION=1.2.3 -e LASTCOMMITSHA=asdas12312 gcr.io/ira-nz/anz-test:latest
```

Attach an alpine container (ideally from another terminal session):

```bash
docker run --rm -it --pid container:$(docker container ls -l -q) alpine
```

You should now be able to check the app container filesystem under **/proc/1/root**:

```bash
ls /proc/1/root/
anz-test  dev       etc       proc      sys
```

## References

- https://www.ianlewis.org/en/building-go-applications-google-container-builder
- https://stackoverflow.com/questions/55106186/no-such-file-or-directory-with-docker-scratch-image
