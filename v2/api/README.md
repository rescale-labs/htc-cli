# htcapi (to eventually be part of htccli!)

This is very early work in progress to generate the go client code for
talking to the HTC API.

## Roadmap

Sami would like first to have a basic CLI only for reading things from the
API. Limits, tasks, etc.

So, steps here are:

* Create the initial main.go, import `_oas`, and implement auth and at
  least one read method.
* Get other read methods working here in  own directory (here).
* Sketch the long term CLI design.
* Integrate this down into the main htccli.
* Document and issue as the tool for talking to HTC API.
* Sort out remaining roadmap for htcctl feature completion and write
  methods for HTC API.

## Dependencies

To patch the swagger file, you will need `jsonnet`.

## Generating the client

From the parent directory, run:

```
make oapi
```

And commit the updated files into the repo. Note that the above does 3
things:

1. Downloads `swagger.json` from https://htc.rescale.com/
2. Patches the swagger into `swagger-patched.json` using jsonnet, so
   that it works for our needs.
3. Generates the code in `_oas` using [ogen](https://ogen.dev/docs/intro).

## Updating the original swagger

To download the latest swagger JSON from the prod HTC API and regenaret
the client, from the parent directory, run:

```
make refresh-swagger
make oapi
```

## Comparing original and patched swagger

In general, reading `swagger.jsonnet` is probably the best way to see
what we've changed from the prod API. But to view the real diff, for
instance when debugging your jsonnet, use this make target from the
parent directory:

```
make diff-swagger
```

This will compare a sorted version of the prod JSON to the
jsonnet-patched output, `swagger-patched.json`.

## Other options for OpenAPI codegen.

ogen-go/ogen has nice instrumentation for traces and timing, even if it
also creates a bunch of server code. So, we're going there first. But
other things looked at follow.

* go-swagger (OpenAPI 2.0 only)
* oapi-codegen (works!)

### go-swagger

go-swagger only supports OpenAPI 2.0

but if we could get a 2.0 response we'd run:

```
go run github.com/go-swagger/go-swagger/cmd/swagger@v0.30.5 generate cli -f swagger.json
```

## oapi-codegen

Works!

```
mkdir -p oapi
go run github.com/deepmap/oapi-codegen/v2/cmd/oapi-codegen@latest \
    --config oapi-config.yaml swagger-patched.yaml
```

## Readings

* https://medium.com/julotech/implementing-swagger-in-go-projects-8579a5fb955
* https://www.reddit.com/r/golang/comments/15h8q9l/openapiswagger_in_go_anyone_implementing_it/?rdt=33785
* https://www.reddit.com/r/golang/comments/10rlp31/toolsgo_pattern_still_valid_today_i_want_to/

