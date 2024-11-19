# Changelog

## v.0.03

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
