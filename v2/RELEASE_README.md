# Rescale HTC cli

`htc` is a single, static binary for doing the most useful things in
HTC, including with the API and inside the job runtime.

**Please note**: this version of `htc` is pre-release, alpha software
with very limited functionality.

## Quickstart

Set `RESCALE_API_TOKEN` in your environment, taking care to disable
history if you don't want tokens in your shell history (thus the leading
spaces in `  export` below).

```
  export RESCALE_API_TOKEN=<YOUR_API_TOKEN>
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
