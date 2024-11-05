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

## Releasing

Steps are pretty simple:

1. Update `VERSION` in `Makefile`:
   ```
   VERSION := v0.0.1
   ```
1. In `CHANGES.md`, rename `## Pending next release` -> `## {VERSION}`
   and create a new pending next release section.
1. Commit.
1. Build the dist archives:
   ```
   make dist
   ```
1. Create the tag and push:
   ```
   git tag $(make echo-version)
   ```
1. Create a new release at
   https://github.com/rescale/htc-storage-cli/releases/new
   1. Copy in the text from `CHANGES.md`.
   1. Upload the dist archives.
