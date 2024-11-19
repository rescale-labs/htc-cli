# htc v2 CLI

Goal: a single, static binary for doing the most useful things in HTC,
including with the API and inside the job runtime.

## Build and test

```
make build
make test
```

## Building for distribation

```
make dist
```

and poke around in `build/dist`. I guess we should zip these things at
some point.

## Developing

Docs are useful, especially API docs, which you can view under
[localhost:6060/pkg/github.com/rescale/htc-storage-cli/v2/api/\_oas/](http://localhost:6060/pkg/github.com/rescale/htc-storage-cli/v2/api/_oas/)
after running:

```
make godoc
```

### Updating dependencies

Dependabot will occasionally want updates, or you'll want to do them
yourself. To do that, from the top directory of the repo:

```
go get -u
go get -u ./...
go mod tidy
```

## Releasing

Steps are pretty simple:

1. In `CHANGES.md`, rename `## Pending next release` -> `## {VERSION}`
1. Commit and push the branch for review.
1. After merge, pull main and tag the version as an annoted tag:
   ```
   git pull origin main
   git tag -a v<LATEST_VERSION_HERE>
   ```
1. Build the dist archives:
   ```
   make dist
   ```
1. Push the tag:
   ```
   git push --tags
   ```
1. Create a new release at
   https://github.com/rescale/htc-storage-cli/releases/new
   1. Copy in the text from `CHANGES.md`.
   1. Upload the dist archives you built in `build/dist`.
