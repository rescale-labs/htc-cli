# Rescale HTC cli

`htc` is a single, static binary for doing the most useful things in
HTC, including with the API and inside the job runtime.

## DISCLAIMER

**Please note**: this version of `htc` is pre-release, alpha software
with very limited functionality. It is provided with NO WARRANTY. Bugs
are likely and features are mostly incomplete.

Bug reports and feature requests are welcome, though at this stage,
Rescale's HTC team is mostly just working to fill out basic
functionality to parity with the Rescale HTC API, and on a best effort
basis. For any problems or questions, please contact `htc`'s current
maintainer, Hunter Blanks, at hblanks@rescale.com or in any appropriate,
existing Slack channel.

## Quickstart

Set `RESCALE_API_KEY` in your environment, taking care to disable
history if you don't want tokens in your shell history (thus the leading
spaces in `  export` below).

```
  export RESCALE_API_KEY=<YOUR_API_KEY>
```

From then, run commands and use `htc help` to see what commands you can
run. A few examples:


Get all projects in your workspace:

```
htc projects get
```

Get all tasks in a project with json output:

```
htc task get --project-id=<PROJECT_ID> -o json
```

Submit a batch of jobs you've prepared in `batch.json`:

```
htc job submit --project-id=<PROJECT_ID> --task-id=<TASK_ID> batch.json
```
 
## Setting up multiple config contexts

`htc` also supports saving configuration and auth tokens for multiple
workspaces through "contexts."

Initial context is always `default`. Example with output:

```
$   export RESCALE_API_KEY=<YOUR_API_KEY>
$ htc auth login
2024/11/13 11:50:17 Bearer token: ExpiresIn=21600
$ htc config context get
    NAME                                             PROJECT ID              TASK ID
 *  default
```

Switch context with `htc config context use` and set variables with `htc
config set`. Example with output:

```
$ htc config context use kaiju
$ htc config set project_id 8f9db624-62de-44da-942f-edcc244f4fcb
$ htc config context get -o json
[
  {
    "name": "default",
    "selected": false
  },
  {
    "name": "kaiju",
    "selected": true,
    "project_id": "8f9db624-62de-44da-942f-edcc244f4fcb"
  }
]
$ htc task get  -l 1 -o json | jq .[].projectId
2024/11/13 12:26:27 HtcProjectsProjectIdTasksGet: projectId=8f9db624-62de-44da-942f-edcc244f4fcb pageIndex= pageSize=500
"8f9db624-62de-44da-942f-edcc244f4fcb"
```
