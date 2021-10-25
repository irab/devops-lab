# Simple Go API endpoint using HTTProuter
[![Go Report Card](https://goreportcard.com/badge/github.com/irab/devops-lab?v=lastCommitID)](https://goreportcard.com/report/github.com/irab/devops-lab)

This repository contains a Go app that dynamically returns the version of the application and SHA commit hash along with a description.

The app uses [HTTProuter](https://github.com/julienschmidt/httprouter) for API routing, which has [lower latency and better performance](https://github.com/julienschmidt/go-http-routing-benchmark) than supplied by the default **http.ServeMux** router and other HTTP/API router frameworks. This provides an RESTful API that's cleaner and more easily extensible than custom code.

The final Docker image is running on Distroless, as building 'FROM Scratch' had some issues on Cloud Run. Distroless is Google's [suggested approach](https://github.com/GoogleContainerTools/distroless/blob/master/base/README.md) to running static binaries on Google Cloud.


## Pre-requisites

1. A GCP Project Setup
1. Installing the [Cloud Build App to Github](https://cloud.google.com/build/docs/automating-builds/build-repos-from-github)
1. Cloud Build Triggers setup

A CI/CD pipeline has been created with a Cloudbuild Github Action integration that triggers for each push to the master branch and also for every tag.

The dynamic variables (**LASTCOMMITSHA**, **VERSION**) in the app are passed to the Go app via environment variables. Theses are derived automatically from [Cloudbuild builtin substitutions](https://cloud.google.com/cloud-build/docs/configuring-builds/substitute-variable-values) - **VERSION** from **build.Source.RepoSource.Revision.TagName** and **LASTCOMMITSHA** from **build.SourceProvenance.ResolvedRepoSource.Revision.CommitSha**. This is to avoid hardcoding and ensure changes to these values come via changes to the git repository.

The cloudbuild.yaml file has the following steps:

- **go_build**: Builds a statically linked Go binary with the flags -s and -w to reduce binary size by removing the symbol table and debugging information.
- **go_test**: Run 'go test', which currently covers testing that the function returns an expected JSON output.
- **go_security**: Runs [gosec](https://github.com/securego/gosec), a code scanner that checks the AST (Abstract Syntax Tree) for security issues
- **docker_image_build**: Builds a minimal Docker image using the Dockerfile in the root of the repository (FROM gcr.io/distroless/static)
- **docker_image_push**: pushes the newly built image to the container registry
- **cloud_run_deploy**: deploys the image to Cloud Run
- **gke_deploy**: deploys the image to a GKE Cluster

## Risks and Benefits

This pipeline has no different environments for managing changes safely. But this could easily be added by creating Cloud Run services and K8s deployments based on branch names or tag regex.

No real failure logic with Cloudbuild, it either fails completely or all steps are successfull. The current couldbuild.yaml can get into a state where the Cloud Run deployment is successful but the GKE fails. There's little logic that can be added to resolve this in a single yaml file. Could be better resolved by splitting out into seperate yaml files.

The application can also be deployed from a commit that hasn't got a tag, so the version field is missing.

External third party dependencies expose the pipeline and application to supply chain attacks. Mitigation of this would be to review and build these components from source and deploy them only from a known and audited repository.

## Local development

### Building

The go build command has been modified to build a binary with a minimal footprint as below:

```bash
CGO_ENABLED=0 go build -ldflags="-s -w"
```

This reduces the binary size from 7.4MB to 5.8MB, a 22% reduction, but the trade off is that the capability to debug/trace the application is degraded somewhat.

### Releasing a new version

Currently the build runs every time a push is made to master and also for each time a new tag is detected. Ideally there would be different environments with a single release branch or there would be a branch created for each release version (depends on the teams gitflow approach). To create a new release that correctly assigns the version dynamically, create a new tag and push to the repo as below:

```bash
git tag 0.0.2 && git push --tags
```

### Updating Dependencies with Go Modules

The go.mod will automatically updated when a new version of external dependencies are updated after a new build. Make sure it these changes are committed.

### Load testing

I added some [loadtests with k6](https://blog.loadimpact.com/load-testing-with-postman-collections) for fun under /loadtests. K6 can be installed on OSX with **brew install k6**.

Tests with 1000 virtual users can be run as follows:

```bash
k6 run -u 1000 loadtests/k6-kubernetes.js
k6 run -u 1000 loadtests/k6-cloudrun.js
```

Performance seems to be much slower on Cloud Run than on Kubernetes:

```bash
k6 run -u 1000 loadtests/k6-cloudrun.js
http_reqs..................: 10000  764.054335/s

k6 run -u 1000 loadtests/k6-kubernetes.js
http_reqs..................: 10000  2154.421445/s
```

### Troubleshooting

You can access the docker binary container by attaching it to the a debugging container, as shown below.

Run the container:

```bash
docker run -e VERSION=1.2.3 -e LASTCOMMITSHA=asdas12312 gcr.io/ira-nz/devops-lab:latest
```

Attach an alpine container (ideally from another terminal session):

```bash
docker run --rm -it --pid container:$(docker container ls -l -q) alpine
```

You should now be able to check the app container filesystem under **/proc/1/root**:

```bash
ls /proc/1/root/
devops-lab  dev       etc       proc      sys
```

## Handy References

- https://www.ianlewis.org/en/building-go-applications-google-container-builder
- https://stackoverflow.com/questions/55106186/no-such-file-or-directory-with-docker-scratch-image
- https://github.com/GoogleCloudPlatform/cloud-builders/tree/master/go
- https://cloud.google.com/solutions/best-practices-for-building-containers
- https://medium.com/digio-australia/building-a-robust-ci-pipeline-for-golang-with-google-cloud-build-4b5029617bc9
- https://blog.loadimpact.com/load-testing-with-postman-collections
- https://ordina-jworks.github.io/cloud/2019/03/28/Building-With-Google-Cloud-Build.html
