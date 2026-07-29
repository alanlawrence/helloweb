# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## What this project is

A Go web server that serves an interactive maths calculator. The browser sends AJAX requests to the server, which computes results and returns HTML fragments that are injected into the page without a reload. Previously deployed to GKE (Google Kubernetes Engine) via Docker. More recently deployed to Google Serverless Cloud Run via a container image registered to Google Artefact Registry from location europe-docker.pkg.dev/alans-gcp-project/helloweb-repo.

## Commands

All commands run from `helloweb/` (the directory containing `go.mod` and `Dockerfile`).

**Run locally:**
```bash
go run helloweb/webserver
```
The server starts on port 8080 (or `$PORT`). It runs self-tests for `IsPrime` at startup and exits if they fail. View in a browser by visiting localhost:8080 (or localhost:$PORT). To run the webserver in the background, add a trailing &.

**Lint:**
```bash
go vet ./...
```

**Formatting:**
Use 4 spaces for indentions, not tabs.


**Test Driven Development:**
Write enough function code for tests to compile.
Write the tests.
Write the remaining code so that the tests pass.


**Design Philosophy:**
Simple to maintain UI, fast to load and process.
Using HTML without styling and minimal third party javascript packages allows fast loading and simple dependency management.
Injecting HTML response fragments into the DOM which are formatted as plain text ensures super fast update responses.
Using the minimal Alpine container image allows fast cold startups.


**Run all tests in verbose mode:**
```bash
go test -v ./...
```

**Run all tests for a package:**
```bash
go test helloweb/series
go test helloweb/longDiv
go test helloweb/quadratic
go test helloweb/digits
```

**Run a single test:**
```bash
go test helloweb/series -run TestARCalc
```

**Build and push Docker image** (run from `helloweb/`):
```bash
./build.sh hello-app <version>         # build + optional push to Google Artifact Registry
```

**Deploy to Cloud Run**
```bash
./cloud-run-deploy.sh <image path:version> <region>  # deploy image from Google Artifact Registry
```

**Build, push Docker Image, and deploy to GKE** (after first-time setup):
```bash
export PROJECT_ID=alans-gcp-project
./build-deploy.sh <version>            # build + push + kubectl apply
```

**Deploy to GKE** (first-time setup):
```bash
cd helloweb/webserver/manifests
./launch-helloweb-app.sh mutate        # omit "mutate" for a dry run
```

## Architecture

```
helloweb/
├── go.mod                      # module: helloweb, go 1.26
├── Dockerfile                  # multi-stage: golang:alpine build → alpine run
├── webserver/
│   ├── webserver.go            # main: HTTP router + inline prime and GCD logic; longmult and division logic included via main package
│   └── index.html              # single-page UI; AJAX calls hit the API endpoints
├── series/series.go            # arithmetic series sum: S = n/2 * (2a + (n-1)d)
├── quadratic/quadratic.go      # quadratic formula, real and complex roots
├── longDiv/longDiv.go          # long division with step-by-step working
├── digits/digits.go            # Digits struct: digit-level manipulation of integers
├── longmult/longmult.go        # long multiplication package (included in main package)
└── division/division.go        # division utilities package (included in main package)
```

**Request flow:** `index.html` → AJAX `XMLHttpRequest` → Go handler in `webserver.go` → package function → HTML fragment string → response body → injected into `<span>` in page.

Each endpoint reads URL query params, calls a computation function, and returns an HTML fragment. Every package that generates output has a paired `GenerateHtml()` function.

**Two testing approaches coexist:**
- `IsPrime` has an inline `TestIsPrime()` called from `main()` at startup (process exits on failure).
- All packages use standard `go test` with `_test.go` files.

**`digits.Digits`** is a shared struct used by `longDiv` to represent numbers as digit slices, enabling the step-by-step working display (each intermediate row in the long division layout is a `Digits` value).

**Deployment to Cloud Run:** `build.sh` builds the image. If the push to Artifact Registry is accepted then the image is tagged with the supplied version tag and pushed to `europe-docker.pkg.dev/alans-gpc-project/helloweb-repo/`. The app name is `hello-app` and is typically deployed to `europe-west2` using `cloud-run-deploy.sh`.

**Deployment to GKE:** `build-deploy.sh` builds the image, updates the version tag in `helloweb-deployment.yaml` via `sed`, and runs `kubectl apply`. The GKE cluster is `hello-cluster` in `europe-west2-a` under GCP project `alans-gcp-project`, with the image stored in `europe-docker.pkg.dev/alans-gcp-project/helloweb-repo/`.
