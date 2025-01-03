# Changelog

## v0.0.8

* Add `auth logout` for clearing existing stored API key and token.
* Add `project create` with sample JSON in usage string.
* Add `region get` to print the global list of HTC regions.
* Extend `dimensions get` to support tabular text output by default.
* Extend `job get` to support sorting by several fields.
* Document that `job submit` can be fed from stdin and include a
  sample JSON payload in usage string.
* Fix top level help text so that all subcommands have applicable
  headings.

## v0.0.7

* Add new experimental cloudfilesystems section from swagger to
  HTCJobSubmitRequest.
* Future proof authentication against method renames.

## v0.0.6

* Update `htc image login-repo` to feed token to docker/podman over
  stdin instead of as a command arg.
* Update `htc image push` to get ECR name using `GET
  /.../projects/{projectId}` instead of fetching all images in a
  project.

## v0.0.5

* Add `htc image push`

## v0.0.4

* Extend `htc job get` so it takes optional job UUID
* Fix `htc job get` so it lists jobs that have not yet started or
  completed. (ENK-2318)
* Fix JSON output issue in `htc job get`.

## v.0.0.3

* Extend `htc config context get` to report workspace name, workspace
  id, and the email address associated with this context's credential.
* Add `htc project retention-policy get`
* Add `htc image create-repo/login-repo/get [IMAGE_NAME:TAG]`.
  `login-repo` is particularly nice since it will log docker or podman
  into a given HTC project's private registry.

## v0.0.2

* Add `htc version`
* Add `htc config set/unset` and `htc config context get/use/delete`.
* Add README to release tarballs including docs on several basic
  commands.

## v0.0.1

* Initial functionality for the following commands:
  ```
  auth
    login
    whoami
  image
    get
  job
    get
    submit
  metrics
    get
  project
    get
  task
    create
    get
  ```
* Automate universal OS X binaries from Makefile: see [GitHub - randall77/makefat: A tool for making fat OSX binaries (a portable lipo)](https://github.com/randall77/makefat)
* Default to table output for things that are Tabler interface
