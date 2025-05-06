# Changelog

## v0.0.10
* Add `job logs -f [JOB_UUID]` flag to live tail job logs.
* Add `job events [JOB_UUID]` to view job lifecycle events.
* Add `task stats` to get task summary statistics.
* Add `workspace clusters` to view details about all GCP clusters that can run jobs for the specified HTC workspace.
* Add `workspace retention-policy get` to view current task retention policy of a specific Workspace
* Add `workspace retention-policy apply` (admin) to define or update task retention policy for the HTC workspace
* Add `workspace dimensions get` to view the various hardware configurations and environments available within a specific workspace.
* Add `project retention-policy get` to view the current task retention policy of a specific project.
* Add `project retention-policy apply` (admin) to define or update the task retention policy for a specific project.
* Add `project limits get` to view all resource limitations associated with this project.
* Add `project limits apply` to add a new limit to this project or overwrite an existing limit if one already exists with the provided `modifierRole`
* Add `project dimensions get` to view the current set of dimension combinations configured for a specific project
* Add `project dimensions apply` to create/update/delete the dimension combinations for a project


## v0.0.9

* Add `job logs [JOB_UUID]` to support viewing job logs
* Add `job cancel` to support attempting cancellation of all jobs in a task
* Add option `job submit --env var1=val1` to support adding environment variables to a job
* Add option `job submit -w $(pwd)` to support passing the current working directory to a job

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
