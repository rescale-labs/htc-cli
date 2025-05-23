# TODO

Roughly sorted task list while building out core functionality of `htc`.

As things check off, please remember to copy to
[CHANGES.md # Pending next release](./CHANGES.md#pending-next-release).

## Soon

Essential API operations:
  * Project limits, dimensions, and task retention policy (ENK-2469):
    * `project dimensions apply`
    * `project limits apply`
    * `project retention-policy apply`
  * Workspace limits, dimensions, and task retention policy (ENK-2450):
    * `workspace dimensions get`
    * `workspace limits get`
    * `workspace retention-policy get`
    * `workspace dimensions apply`
    * `workspace limits apply`
    * `workspace retention-policy apply`
  * `clusters get` (ENK-2651)
  * Jobs:
    * GET /jobs/{jobId}/events  (ENK-2652)
    * POST /jobs/cancel (ENK-2653)
  * Tasks:
    * GET /tasks/summary-statistics (ENK-2654)
	* Fix JSON output:
		* Make runner.PrintResult should require `tabler.Tabler` interface
		* Make tabler.Tabler interface include `MarshalJSON([]byte) error`
    * Document in the interface why this is important (because for non-slice
      types, we need to use oapi's `MarshalJSON()` implementation, otherwise
      JSON encoding can fail, as it did quite unpleasantly with tabler.HTCJob.

Warts:
  * We now have `--project-id` on CLI, `project_id` in TOML, and `projectId` in some (but not all!) JSON. Will need to ask for Enkis' thoughts on this.
  * Some commands display global flags in help even though they don't use them. And, calling cobra.Command.ResetFlags() on them doesn't seem to help. So. Yuck.
  * Deleting context does not delete credentials

Testing:
  * Identify the most important things to test. Test them. Somehow.

Usability:
  * Print the env vars expected, if any, in usage.
  * Are we setting User-Agent?
  * For `htc job submit --group`, it might be nice if it automatically ran `htc config set group GROUP`
  * Decide on whether we need the YAML encoder. It doesn't handle OptString, etc properly, e.g.:
    ```
    $ go run cmd/htc/main.go  task get --project-id e8dee146-5606-407e-8e51-3f28070ece6e -l 1  -o yaml
    2024/09/25 13:03:40 HtcProjectsProjectIdTasksGet: pageIndex= pageSize=500
    -   archivedat:
            value: 0001-01-01T00:00:00Z
            set: true
            "null": true
    ```
		* Brian would like per-record (transposed table output) - this may / may not be the same as YAML output.

Automation:
    * Run `make dist` and upload them to Github releases. Maybe S3, too.
        * https://docs.github.com/en/actions/use-cases-and-examples/building-and-testing/building-and-testing-go
        * https://docs.github.com/en/actions/writing-workflows/choosing-what-your-workflow-does/storing-and-sharing-data-from-a-workflow
        * https://github.com/actions/upload-artifact?tab=readme-ov-file#where-does-the-upload-go
    * prepare releases when tagged with https://github.com/softprops/action-gh-release
    * Brian would like to be able to add extra columns to tabular output.
    * Brian would like `htc config edit` that drops into an editor.

## stuff from Alistair 2024-10-22

things we run a lot in my team are:
  * list out all the workspaces - with the name / ID and (possibly) API key
  * [done] list out all the projects that are configured as part of a given WS
  * describe the vcpu-limits of a project
    * and set the vcpu-limits of a project
  * describe the dimensions of a project
    * and change the dimensions of a project
  * describe the task-lifecycle of a project
    * and set this
  * show the details of 50 jobs from a Task

separate to these simpler scripts is:
  * run a 'task-report', which gives a summary of all Tasks within any given Project
    * task ID
    * task owner
    * task created / last modified time
    * job-status for all jobs (ie. how many in 'submitted-to-rescale', how many in 'submitted-to-provider', 'running', 'completed', and so on)
  * this script allows us to get a feel for what the users are up to, as a complement to graphical metrics views like our dashboards
    * the dashboards give an amalgamated view of many users/tasks/jobs, and the corresponding infrastructure
    * but sometimes we want more specifics about what's actually driving those graphed metrics, so a task-by-task breakdown can be useful

I can give example output from any of the above / existing scripting, if that helps ?


## Later

* merge in upload/download from htccli
* scope out htcctl functionality to implement here
* task deletion

* All the methods:
    ```
    enkictl --help
    Please enable `enable-htc-api` for your user to fetch bearer
    Unknown option: '--help'
    Usage: enkictl [-o=<outputFormat>] [COMMAND]
      -o, --output=<outputFormat>
             json, yaml
    Commands:
      help                    Displays help information about the specified command
      get
        health
          ready               Get readiness of HTC
        auth
          jwt                 Get Roles, Claim, Subject for current bearer token
          me                  Make a who am I call
          token               Get Bearer token for further authentication
        regions               Region commands
          health              Region health command
        projects              Project commands
        registry
          images              Get project images
          token               Get token for authentication to repo
        tasks                 Task commands
        jobs                  Jobs commands
        job-details
          events              Get events of a job
          logs                Get logs of a job
        summary               Summary statistics for task
        storage               Storage commands
          presigned-url
          token
          tokens
        infra                 Only Admin can access those to check BYO workspaces
                                infra
          ecr-replication     Get rule to replicate across ECR regions
          event-rules         Get rules in EventBridge
          job-queues          Get job queues
          compute-env         Get compute environment
          update-compute-env  Get launchTemplates
        metrics               Metrics commands
      create
        project               Create project
        task                  Create task
        batch                 Submit batch of jobs
        registry              Create registry repo
          repo                Create registry repo
        infra                 Only Superuser can access those to create
                                infrastructure
          compute-env         Create compute environment in BYO region
          ecr-replication     Create rule to replicate across ECR regions
          event-rule          Create rule in EventBridge
          job-queues          Create job queues
          update-compute-env  Create compute environment in BYO region
      delete
        registry              Delete image from project registry
        task                  Delete task from project
      cancel
        jobs                  Cancel submitted jobs
      completion              bash/zsh completion:  source <(enkictl completion)
      config
        set                   Change configuration in config file
          new                 Create new config entry
          baseuri             Set baseuri
          current-config      Set current configuration
        get                   Get current active configuration
        delete                Delete config entry from config file
    ```
* Eventually support project deletion along the lines of
  https://rescale.atlassian.net/browse/ENK-2095 ?
* Show the last job run in a task/project/workspace.

Security:
  * Before we can deprecate htcctl, we'll need to do security certification for `htc`. (See Rescale internal sync notes from 2024-11-26.)


<!-- vim: set tw=999999 sts=0 ts=2 sw=2: -->
