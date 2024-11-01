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
