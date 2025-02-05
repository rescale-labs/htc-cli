// Code generated by ogen, DO NOT EDIT.

package _oas

import (
	"context"
)

// Handler handles operations described by OpenAPI v3 specification.
type Handler interface {
	AuthHandler
	ImageHandler
	JobHandler
	MetricsHandler
	ProjectHandler
	TaskHandler
	// AuthTokenWhoamiGet implements GET /auth/token/whoami operation.
	//
	// This endpoint will get a JWT token payload given a bearer token.
	//
	// GET /auth/token/whoami
	AuthTokenWhoamiGet(ctx context.Context) (AuthTokenWhoamiGetRes, error)
	// HtcGcpClustersWorkspaceIdGet implements GET /htc/gcp/clusters/{workspaceId} operation.
	//
	// This endpoint returns details about all GCP clusters that can run jobs for the specified HTC
	// workspace.
	//
	// GET /htc/gcp/clusters/{workspaceId}
	HtcGcpClustersWorkspaceIdGet(ctx context.Context, params HtcGcpClustersWorkspaceIdGetParams) (HtcGcpClustersWorkspaceIdGetRes, error)
	// HtcProjectsProjectIdDimensionsPut implements PUT /htc/projects/{projectId}/dimensions operation.
	//
	// This endpoint allows _workspace_, _organization_, and _Rescale administrators_ to _create_,
	// _update_, or _delete_ the dimension combinations for a project. It accepts a list of dimension
	// combinations, each specifying a unique set of computing environment attributes to tailor the
	// computing environment(s) of a project to match specific job requirements.
	// For example, a project’s dimensions can be configured to require jobs to run on a particular
	// type of processor architecture, within a certain region, and with or without hyperthreading.
	// It's important to note that the dimensions set through this endpoint must align with the available
	// dimensions at the workspace level.
	// **In the event that a project's dimension requirements need to be reset to allow for a broader
	// range of job types, sending an empty list to this endpoint will remove all existing dimension
	// restrictions, returning the project to a state where it can accommodate any dimension available in
	// the workspace.**.
	//
	// PUT /htc/projects/{projectId}/dimensions
	HtcProjectsProjectIdDimensionsPut(ctx context.Context, req []HTCComputeEnvironment, params HtcProjectsProjectIdDimensionsPutParams) (HtcProjectsProjectIdDimensionsPutRes, error)
	// HtcProjectsProjectIdLimitsDelete implements DELETE /htc/projects/{projectId}/limits operation.
	//
	// This endpoint will remove all resource limits associated with this project.
	// Any jobs `SUBMITTED_TO_RESCALE` will transition to `SUBMITTED_TO_PROVIDER` if no other limits apply.
	//
	// DELETE /htc/projects/{projectId}/limits
	HtcProjectsProjectIdLimitsDelete(ctx context.Context, params HtcProjectsProjectIdLimitsDeleteParams) (HtcProjectsProjectIdLimitsDeleteRes, error)
	// HtcProjectsProjectIdLimitsIDDelete implements DELETE /htc/projects/{projectId}/limits/{id} operation.
	//
	// This endpoint will remove a single resource limit associated with this project if it exists.
	//
	// DELETE /htc/projects/{projectId}/limits/{id}
	HtcProjectsProjectIdLimitsIDDelete(ctx context.Context, params HtcProjectsProjectIdLimitsIDDeleteParams) (HtcProjectsProjectIdLimitsIDDeleteRes, error)
	// HtcProjectsProjectIdLimitsIDGet implements GET /htc/projects/{projectId}/limits/{id} operation.
	//
	// This endpoint will get either the `PROJECT_ADMIN` or `WORKSPACE_ADMIN` limit for this project.
	//
	// GET /htc/projects/{projectId}/limits/{id}
	HtcProjectsProjectIdLimitsIDGet(ctx context.Context, params HtcProjectsProjectIdLimitsIDGetParams) (HtcProjectsProjectIdLimitsIDGetRes, error)
	// HtcProjectsProjectIdLimitsIDPatch implements PATCH /htc/projects/{projectId}/limits/{id} operation.
	//
	// This endpoint will update one of the existing resource limits associated with this project.
	// Any user who belongs the project's workspace can modify the `PROJECT_ADMIN` limit. Higher
	// permissions are required to modify the `WORKSPACE_ADMIN` limit.
	//
	// PATCH /htc/projects/{projectId}/limits/{id}
	HtcProjectsProjectIdLimitsIDPatch(ctx context.Context, req OptHTCLimitUpdate, params HtcProjectsProjectIdLimitsIDPatchParams) (HtcProjectsProjectIdLimitsIDPatchRes, error)
	// HtcProjectsProjectIdLimitsPost implements POST /htc/projects/{projectId}/limits operation.
	//
	// This endpoint will add a new limit to this project or overwrite an existing limit if one already
	// exists with the provided `modifierRole`.
	// Jobs submitted to this project will only run when the active resource count falls below the
	// minimum of all limits associated with this project.
	// Any user who belongs the project's workspace can modify the `PROJECT_ADMIN` limit. Higher
	// permissions are required to modify the `WORKSPACE_ADMIN` limit.
	//
	// POST /htc/projects/{projectId}/limits
	HtcProjectsProjectIdLimitsPost(ctx context.Context, req OptHTCLimitCreate, params HtcProjectsProjectIdLimitsPostParams) (HtcProjectsProjectIdLimitsPostRes, error)
	// HtcProjectsProjectIdPatch implements PATCH /htc/projects/{projectId} operation.
	//
	// This endpoint allows for updating a project's regions.
	//
	// PATCH /htc/projects/{projectId}
	HtcProjectsProjectIdPatch(ctx context.Context, req OptHTCProjectUpdate, params HtcProjectsProjectIdPatchParams) (HtcProjectsProjectIdPatchRes, error)
	// HtcProjectsProjectIdStoragePresignedURLGet implements GET /htc/projects/{projectId}/storage/presigned-url operation.
	//
	// This endpoint will get a presigned url for project storage.
	//
	// GET /htc/projects/{projectId}/storage/presigned-url
	HtcProjectsProjectIdStoragePresignedURLGet(ctx context.Context, params HtcProjectsProjectIdStoragePresignedURLGetParams) (HtcProjectsProjectIdStoragePresignedURLGetRes, error)
	// HtcProjectsProjectIdStorageTokenGet implements GET /htc/projects/{projectId}/storage/token operation.
	//
	// This endpoint will get temporary access information for a project storage.
	//
	// GET /htc/projects/{projectId}/storage/token
	HtcProjectsProjectIdStorageTokenGet(ctx context.Context, params HtcProjectsProjectIdStorageTokenGetParams) (HtcProjectsProjectIdStorageTokenGetRes, error)
	// HtcProjectsProjectIdStorageTokenRegionGet implements GET /htc/projects/{projectId}/storage/token/{region} operation.
	//
	// This endpoint will get temporary access information for a project storage given a region.
	//
	// GET /htc/projects/{projectId}/storage/token/{region}
	HtcProjectsProjectIdStorageTokenRegionGet(ctx context.Context, params HtcProjectsProjectIdStorageTokenRegionGetParams) (HtcProjectsProjectIdStorageTokenRegionGetRes, error)
	// HtcProjectsProjectIdStorageTokensGet implements GET /htc/projects/{projectId}/storage/tokens operation.
	//
	// This endpoint will get temporary access information for all project storages.
	//
	// GET /htc/projects/{projectId}/storage/tokens
	HtcProjectsProjectIdStorageTokensGet(ctx context.Context, params HtcProjectsProjectIdStorageTokensGetParams) (HtcProjectsProjectIdStorageTokensGetRes, error)
	// HtcProjectsProjectIdTaskRetentionPolicyDelete implements DELETE /htc/projects/{projectId}/task-retention-policy operation.
	//
	// This endpoint allows users to delete the task retention policy for the specified project. When a
	// project-level policy is deleted, the auto-archival and auto-deletion behavior for tasks within the
	// project will fall back to the workspace-level policy (if any). If no workspace-level policy is set,
	//  tasks within the project will not be subject to any auto-archival or auto-deletion.
	//
	// DELETE /htc/projects/{projectId}/task-retention-policy
	HtcProjectsProjectIdTaskRetentionPolicyDelete(ctx context.Context, params HtcProjectsProjectIdTaskRetentionPolicyDeleteParams) (HtcProjectsProjectIdTaskRetentionPolicyDeleteRes, error)
	// HtcProjectsProjectIdTaskRetentionPolicyGet implements GET /htc/projects/{projectId}/task-retention-policy operation.
	//
	// This endpoint is used to retrieve the current task retention policy of a specific project. The
	// task retention policy is necessary in managing the lifecycle of tasks within a project. The task
	// retention policy includes two key aspects:
	// * **Deletion Grace Period**: The `deleteAfter` field represents the duration (in hours) after
	// which an archived task is automatically deleted. Archived tasks can be unarchived during this
	// period, protecting users from prematurely deleting task resources.
	// * **Auto-Archive After Inactivity**: The `archiveAfter` field represents the duration (in hours)
	// of inactivity after which an active task is automatically archived. This feature helps in keeping
	// the project organized by archiving active tasks, ensuring that storage resources are freed
	// optimistically.
	// Setting either value to `0` will result in disabling of that feature. For example, a project's
	// task retention policy with `deleteAfter` set to `0` will result in tasks within that project never
	// auto-deleting.
	// If no policy is set at the project level (i.e., the response is a 404), the policy at the
	// workspace level will apply. If the policy has archiveAfter or deleteAfter set to 0, it means that
	// auto-archival or auto-deletion is disabled at the project level and any workspace level policy is
	// ignored.
	//
	// GET /htc/projects/{projectId}/task-retention-policy
	HtcProjectsProjectIdTaskRetentionPolicyGet(ctx context.Context, params HtcProjectsProjectIdTaskRetentionPolicyGetParams) (HtcProjectsProjectIdTaskRetentionPolicyGetRes, error)
	// HtcProjectsProjectIdTaskRetentionPolicyPut implements PUT /htc/projects/{projectId}/task-retention-policy operation.
	//
	// This endpoint enables project administrators to define or update the task retention policy for a
	// specific project. The task retention policy includes two key aspects:
	// * **Deletion Grace Period**: The `deleteAfter` field allows administrators to set the duration (in
	// hours) after which an archived task is automatically deleted. This control allows for flexibility
	// in managing the lifecycle of tasks, ensuring that data is retained for an adequate period before
	// being permanently deleted. Archived tasks can be unarchived during this period, protecting users
	// from prematurely deleting task resources
	// * **Auto-Archive After Inactivity**: The `archiveAfter` field allows administrators to specify the
	// duration (in hours) of inactivity after which an active task is automatically archived. This
	// feature helps in keeping the project organized by archiving active tasks, ensuring that storage
	// resources are freed optimistically.
	// Setting either value to `0` will result in disabling of that feature. For example, a project's
	// task retention policy with `deleteAfter` set to `0` will result in tasks within that project never
	// auto-deleting.If no policy is set at the project level, the workspace-level policy (if any) will
	// be applied to the project.
	//
	// PUT /htc/projects/{projectId}/task-retention-policy
	HtcProjectsProjectIdTaskRetentionPolicyPut(ctx context.Context, req OptTaskRetentionPolicy, params HtcProjectsProjectIdTaskRetentionPolicyPutParams) (HtcProjectsProjectIdTaskRetentionPolicyPutRes, error)
	// HtcProjectsProjectIdTasksTaskIdDelete implements DELETE /htc/projects/{projectId}/tasks/{taskId} operation.
	//
	// This endpoint will delete a task by ID.
	//
	// DELETE /htc/projects/{projectId}/tasks/{taskId}
	HtcProjectsProjectIdTasksTaskIdDelete(ctx context.Context, params HtcProjectsProjectIdTasksTaskIdDeleteParams) (HtcProjectsProjectIdTasksTaskIdDeleteRes, error)
	// HtcProjectsProjectIdTasksTaskIdGet implements GET /htc/projects/{projectId}/tasks/{taskId} operation.
	//
	// This endpoint will get a task by ID.
	//
	// GET /htc/projects/{projectId}/tasks/{taskId}
	HtcProjectsProjectIdTasksTaskIdGet(ctx context.Context, params HtcProjectsProjectIdTasksTaskIdGetParams) (HtcProjectsProjectIdTasksTaskIdGetRes, error)
	// HtcProjectsProjectIdTasksTaskIdGroupSummaryStatisticsGet implements GET /htc/projects/{projectId}/tasks/{taskId}/group-summary-statistics operation.
	//
	// This endpoint will get job status summary statistics for each group in a task.
	//
	// GET /htc/projects/{projectId}/tasks/{taskId}/group-summary-statistics
	HtcProjectsProjectIdTasksTaskIdGroupSummaryStatisticsGet(ctx context.Context, params HtcProjectsProjectIdTasksTaskIdGroupSummaryStatisticsGetParams) (HtcProjectsProjectIdTasksTaskIdGroupSummaryStatisticsGetRes, error)
	// HtcProjectsProjectIdTasksTaskIdGroupsGet implements GET /htc/projects/{projectId}/tasks/{taskId}/groups operation.
	//
	// This endpoint will get task groups.
	//
	// GET /htc/projects/{projectId}/tasks/{taskId}/groups
	HtcProjectsProjectIdTasksTaskIdGroupsGet(ctx context.Context, params HtcProjectsProjectIdTasksTaskIdGroupsGetParams) (HtcProjectsProjectIdTasksTaskIdGroupsGetRes, error)
	// HtcProjectsProjectIdTasksTaskIdJobsJobIdEventsGet implements GET /htc/projects/{projectId}/tasks/{taskId}/jobs/{jobId}/events operation.
	//
	// This endpoint will get events for a job.
	//
	// GET /htc/projects/{projectId}/tasks/{taskId}/jobs/{jobId}/events
	HtcProjectsProjectIdTasksTaskIdJobsJobIdEventsGet(ctx context.Context, params HtcProjectsProjectIdTasksTaskIdJobsJobIdEventsGetParams) (HtcProjectsProjectIdTasksTaskIdJobsJobIdEventsGetRes, error)
	// HtcProjectsProjectIdTasksTaskIdJobsJobIdLogsGet implements GET /htc/projects/{projectId}/tasks/{taskId}/jobs/{jobId}/logs operation.
	//
	// This endpoint will get job logs.
	//
	// GET /htc/projects/{projectId}/tasks/{taskId}/jobs/{jobId}/logs
	HtcProjectsProjectIdTasksTaskIdJobsJobIdLogsGet(ctx context.Context, params HtcProjectsProjectIdTasksTaskIdJobsJobIdLogsGetParams) (HtcProjectsProjectIdTasksTaskIdJobsJobIdLogsGetRes, error)
	// HtcProjectsProjectIdTasksTaskIdPatch implements PATCH /htc/projects/{projectId}/tasks/{taskId} operation.
	//
	// This endpoint allows for managing the lifecycle of tasks. Users may set the `LifecycleStatus` of
	// an HTCTask in accordance with their data retention requirements.
	// Archiving a Task: To archive an active task, submit a PATCH request with "lifecycleStatus":
	// "ARCHIVED". This action is permissible only if the task is currently active and has no running
	// jobs. Once archived, the task enters a state where it is no longer operational, but its data is
	// retained. An archived task will be automatically scheduled for deletion after a period defined in
	// the project's task retention policy.
	// Unarchiving a Task: If a task is in an archived state and you wish to defer its automatic deletion,
	//  you can restore it to an active state. To unarchive a task, PATCH it with "lifecycleStatus":
	// "ACTIVE". This action reactivates the task, making it modifiable and operational again. Note that
	// this action is only applicable to tasks in the ARCHIVED state.
	// Restrictions: Tasks in a DELETED state are immutable and cannot be transitioned to any other state
	// using this endpoint. Similarly, tasks can only be archived if they are in an ACTIVE state and do
	// not have any running jobs.
	//
	// PATCH /htc/projects/{projectId}/tasks/{taskId}
	HtcProjectsProjectIdTasksTaskIdPatch(ctx context.Context, req OptHTCTaskUpdate, params HtcProjectsProjectIdTasksTaskIdPatchParams) (HtcProjectsProjectIdTasksTaskIdPatchRes, error)
	// HtcProjectsProjectIdTasksTaskIdStoragePresignedURLGet implements GET /htc/projects/{projectId}/tasks/{taskId}/storage/presigned-url operation.
	//
	// This endpoint will get a presigned url for task storage.
	//
	// GET /htc/projects/{projectId}/tasks/{taskId}/storage/presigned-url
	HtcProjectsProjectIdTasksTaskIdStoragePresignedURLGet(ctx context.Context, params HtcProjectsProjectIdTasksTaskIdStoragePresignedURLGetParams) (HtcProjectsProjectIdTasksTaskIdStoragePresignedURLGetRes, error)
	// HtcProjectsProjectIdTasksTaskIdStorageRegionalStorageGet implements GET /htc/projects/{projectId}/tasks/{taskId}/storage/regional-storage operation.
	//
	// This endpoint will get temporary access information for all task storages.
	//
	// GET /htc/projects/{projectId}/tasks/{taskId}/storage/regional-storage
	HtcProjectsProjectIdTasksTaskIdStorageRegionalStorageGet(ctx context.Context, params HtcProjectsProjectIdTasksTaskIdStorageRegionalStorageGetParams) (HtcProjectsProjectIdTasksTaskIdStorageRegionalStorageGetRes, error)
	// HtcProjectsProjectIdTasksTaskIdStorageTokenGet implements GET /htc/projects/{projectId}/tasks/{taskId}/storage/token operation.
	//
	// This endpoint will get temporary access information for a task storage.
	//
	// GET /htc/projects/{projectId}/tasks/{taskId}/storage/token
	HtcProjectsProjectIdTasksTaskIdStorageTokenGet(ctx context.Context, params HtcProjectsProjectIdTasksTaskIdStorageTokenGetParams) (HtcProjectsProjectIdTasksTaskIdStorageTokenGetRes, error)
	// HtcProjectsProjectIdTasksTaskIdStorageTokenRegionGet implements GET /htc/projects/{projectId}/tasks/{taskId}/storage/token/{region} operation.
	//
	// This endpoint will get temporary access information for a task storage given a region.
	//
	// GET /htc/projects/{projectId}/tasks/{taskId}/storage/token/{region}
	HtcProjectsProjectIdTasksTaskIdStorageTokenRegionGet(ctx context.Context, params HtcProjectsProjectIdTasksTaskIdStorageTokenRegionGetParams) (HtcProjectsProjectIdTasksTaskIdStorageTokenRegionGetRes, error)
	// HtcProjectsProjectIdTasksTaskIdStorageTokensGet implements GET /htc/projects/{projectId}/tasks/{taskId}/storage/tokens operation.
	//
	// This endpoint will get temporary access information for all task storages.
	//
	// GET /htc/projects/{projectId}/tasks/{taskId}/storage/tokens
	HtcProjectsProjectIdTasksTaskIdStorageTokensGet(ctx context.Context, params HtcProjectsProjectIdTasksTaskIdStorageTokensGetParams) (HtcProjectsProjectIdTasksTaskIdStorageTokensGetRes, error)
	// HtcRegionsGet implements GET /htc/regions operation.
	//
	// This endpoint will get HTC region settings for all regions.
	//
	// GET /htc/regions
	HtcRegionsGet(ctx context.Context, params HtcRegionsGetParams) (HtcRegionsGetRes, error)
	// HtcRegionsRegionGet implements GET /htc/regions/{region} operation.
	//
	// This endpoint will get HTC region settings for a specific region.
	//
	// GET /htc/regions/{region}
	HtcRegionsRegionGet(ctx context.Context, params HtcRegionsRegionGetParams) (HtcRegionsRegionGetRes, error)
	// HtcStorageGet implements GET /htc/storage operation.
	//
	// This endpoint will get storages in all enabled regions.
	//
	// GET /htc/storage
	HtcStorageGet(ctx context.Context) (HtcStorageGetRes, error)
	// HtcStorageRegionRegionGet implements GET /htc/storage/region/{region} operation.
	//
	// This endpoint will get a storage for a region.
	//
	// GET /htc/storage/region/{region}
	HtcStorageRegionRegionGet(ctx context.Context, params HtcStorageRegionRegionGetParams) (HtcStorageRegionRegionGetRes, error)
	// HtcWorkspacesWorkspaceIdDimensionsGet implements GET /htc/workspaces/{workspaceId}/dimensions operation.
	//
	// This endpoint provides a comprehensive view of the various hardware configurations and
	// environments available within a specific workspace. This read-only API is primarily designed for
	// users who need to understand the different "dimensions" or attributes that describe the hardware
	// and other aspects of job runs within their workspace. By offering insights into available
	// environments, it aids users in selecting the most suitable configuration for their jobs,
	// especially when performance testing across different hardware setups.
	// Normal users can access this endpoint for the workspace they belong to
	// Rescale personnel are required in order to modify any of these dimensions.
	//
	// GET /htc/workspaces/{workspaceId}/dimensions
	HtcWorkspacesWorkspaceIdDimensionsGet(ctx context.Context, params HtcWorkspacesWorkspaceIdDimensionsGetParams) (HtcWorkspacesWorkspaceIdDimensionsGetRes, error)
	// HtcWorkspacesWorkspaceIdLimitsGet implements GET /htc/workspaces/{workspaceId}/limits operation.
	//
	// This endpoint will get the resource limit applied to this workspace.
	//
	// GET /htc/workspaces/{workspaceId}/limits
	HtcWorkspacesWorkspaceIdLimitsGet(ctx context.Context, params HtcWorkspacesWorkspaceIdLimitsGetParams) (HtcWorkspacesWorkspaceIdLimitsGetRes, error)
	// HtcWorkspacesWorkspaceIdTaskRetentionPolicyGet implements GET /htc/workspaces/{workspaceId}/task-retention-policy operation.
	//
	// This endpoint is used to retrieve the current task retention policy of a specific Workspace. The
	// task retention policy is necessary in managing the lifecycle of tasks within a Workspace. The task
	// retention policy includes two key aspects:
	// * **Deletion Grace Period**: The `deleteAfter` field represents the duration (in hours) after
	// which an archived task is automatically deleted. Archived tasks can be unarchived during this
	// period, protecting users from prematurely deleting task resources.
	// * **Auto-Archive After Inactivity**: The `archiveAfter` field represents the duration (in hours)
	// of inactivity after which an active task is automatically archived. This feature helps in keeping
	// the project organized by archiving active tasks, ensuring that storage resources are freed
	// optimistically.
	// Setting either value to `0` will result in disabling of that feature. For example, a project's
	// task retention policy with `deleteAfter` set to `0` will result in tasks within that project never
	// auto-deleting.
	//
	// GET /htc/workspaces/{workspaceId}/task-retention-policy
	HtcWorkspacesWorkspaceIdTaskRetentionPolicyGet(ctx context.Context, params HtcWorkspacesWorkspaceIdTaskRetentionPolicyGetParams) (HtcWorkspacesWorkspaceIdTaskRetentionPolicyGetRes, error)
	// HtcWorkspacesWorkspaceIdTaskRetentionPolicyPut implements PUT /htc/workspaces/{workspaceId}/task-retention-policy operation.
	//
	// This endpoint enables Workspace administrators to define or update the task retention policy for a
	// specific workspace. The task retention policy includes two key aspects:
	// * **Deletion Grace Period**: The `deleteAfter` field allows administrators to set the duration (in
	// hours) after which an archived task is automatically deleted. This control allows for flexibility
	// in managing the lifecycle of tasks, ensuring that data is retained for an adequate period before
	// being permanently deleted. Archived tasks can be unarchived during this period, protecting users
	// from prematurely deleting task resources
	// * **Auto-Archive After Inactivity**: The `archiveAfter` field allows administrators to specify the
	// duration (in hours) of inactivity after which an active task is automatically archived. This
	// feature helps in keeping the project organized by archiving active tasks, ensuring that storage
	// resources are freed optimistically.
	// Setting either value to `0` will result in disabling of that feature. For example, a workspace's
	// task retention policy with `deleteAfter` set to `0` will result in tasks within that project never
	// auto-deleting. The policy applies to all projects within the workspace that do not have their own
	// project-level policy defined. If a project within the workspace has its own retention policy
	// defined, the project-level policy takes precedence over the workspace-level policy.
	//
	// PUT /htc/workspaces/{workspaceId}/task-retention-policy
	HtcWorkspacesWorkspaceIdTaskRetentionPolicyPut(ctx context.Context, req OptWorkspaceTaskRetentionPolicy, params HtcWorkspacesWorkspaceIdTaskRetentionPolicyPutParams) (HtcWorkspacesWorkspaceIdTaskRetentionPolicyPutRes, error)
	// OAuth2TokenPost implements POST /oauth2/token operation.
	//
	// This endpoint will get an OAuth access token.
	//
	// POST /oauth2/token
	OAuth2TokenPost(ctx context.Context) (OAuth2TokenPostRes, error)
	// WellKnownJwksJSONGet implements GET /.well-known/jwks.json operation.
	//
	// This endpoint will get the public keys used to verify JWT.
	//
	// GET /.well-known/jwks.json
	WellKnownJwksJSONGet(ctx context.Context) (WellKnownJwksJSONGetRes, error)
}

// AuthHandler handles operations described by OpenAPI v3 specification.
//
// x-ogen-operation-group: Auth
type AuthHandler interface {
	// GetToken implements getToken operation.
	//
	// This endpoint will get a JWT token given an API key.
	//
	// GET /auth/token
	GetToken(ctx context.Context) (GetTokenRes, error)
	// WhoAmI implements whoAmI operation.
	//
	// This endpoint will get Rescale user information given a Rescale API key.
	//
	// GET /auth/whoami
	WhoAmI(ctx context.Context) (WhoAmIRes, error)
}

// ImageHandler handles operations described by OpenAPI v3 specification.
//
// x-ogen-operation-group: Image
type ImageHandler interface {
	// CreateRepo implements createRepo operation.
	//
	// This endpoint will create a private container repository belonging to this project
	// Private container registries are collections of repositories, and private repositories are
	// collections of container images. These images are referenced when running jobs within this project.
	//  In order to upload an image to a repository, you will need the `registryURI`, the
	// `repositoryName`, and the token (see `/htc/projects/:projectId/container-registry/token`).
	//
	// POST /htc/projects/{projectId}/container-registry/repo/{repoName}
	CreateRepo(ctx context.Context, params CreateRepoParams) (CreateRepoRes, error)
	// GetImage implements getImage operation.
	//
	// Retrieves the current status of an image across cloud providers. The status indicates whether the
	// image is ready for use or still being processed. Returns READY when the image is available in all
	// cloud providers, PENDING while the image is being replicated, and a 404 if the image does not
	// exist.
	//
	// GET /htc/projects/{projectId}/container-registry/images/{imageName}
	GetImage(ctx context.Context, params GetImageParams) (GetImageRes, error)
	// GetImages implements getImages operation.
	//
	// This endpoint will list all images for a project.
	//
	// GET /htc/projects/{projectId}/container-registry/images
	GetImages(ctx context.Context, params GetImagesParams) (GetImagesRes, error)
	// GetRegistryToken implements getRegistryToken operation.
	//
	// This endpoint will get a container registry authorization token.
	// To use this token run `docker login --username AWS --password {TOKEN} {CONTAINER_REGISTRY_DOMAIN}`.
	// e.g. `docker login --username AWS --password "eyJwYXlsb2FkIjoiZHhtSzJuQ0x..." 183929446192.dkr.ecr.
	// us-west-2.amazonaws.com`.
	//
	// GET /htc/projects/{projectId}/container-registry/token
	GetRegistryToken(ctx context.Context, params GetRegistryTokenParams) (GetRegistryTokenRes, error)
}

// JobHandler handles operations described by OpenAPI v3 specification.
//
// x-ogen-operation-group: Job
type JobHandler interface {
	// CancelJobs implements cancelJobs operation.
	//
	// This endpoint will attempt to cancel submitted jobs.
	// Note a 200 response status code does not mean all jobs were cancelled.
	//
	// POST /htc/projects/{projectId}/tasks/{taskId}/jobs/cancel
	CancelJobs(ctx context.Context, params CancelJobsParams) (CancelJobsRes, error)
	// GetJob implements getJob operation.
	//
	// This endpoint will get a job by id.
	//
	// GET /htc/projects/{projectId}/tasks/{taskId}/jobs/{jobId}
	GetJob(ctx context.Context, params GetJobParams) (GetJobRes, error)
	// GetJobs implements getJobs operation.
	//
	// This endpoint will get all jobs for a task.
	//
	// GET /htc/projects/{projectId}/tasks/{taskId}/jobs
	GetJobs(ctx context.Context, params GetJobsParams) (GetJobsRes, error)
	// SubmitJobs implements submitJobs operation.
	//
	// This endpoint will submit a batch of jobs for a task.
	//
	// POST /htc/projects/{projectId}/tasks/{taskId}/jobs/batch
	SubmitJobs(ctx context.Context, req []HTCJobSubmitRequest, params SubmitJobsParams) (SubmitJobsRes, error)
}

// MetricsHandler handles operations described by OpenAPI v3 specification.
//
// x-ogen-operation-group: Metrics
type MetricsHandler interface {
	// GetMetrics implements getMetrics operation.
	//
	// Get all HTC Metrics for a workspace.
	//
	// GET /htc/metrics
	GetMetrics(ctx context.Context, params GetMetricsParams) (GetMetricsRes, error)
}

// ProjectHandler handles operations described by OpenAPI v3 specification.
//
// x-ogen-operation-group: Project
type ProjectHandler interface {
	// CreateProject implements createProject operation.
	//
	// This endpoint will create a project. A project is a collection of tasks and container images used
	// to run jobs. Several projects can belong to a single workspace.
	//
	// POST /htc/projects
	CreateProject(ctx context.Context, req OptHTCProject) (CreateProjectRes, error)
	// GetDimensions implements getDimensions operation.
	//
	// This endpoint is designed to retrieve the current set of dimension combinations configured for a
	// specific project so that users can understand the existing computing environment constraints of a
	// project. It returns a list of dimension combinations such as pricing priority, geographical region,
	//  compute scaling policy, and hyperthreading options.
	// Any user who _belongs to the workspace this project belongs to_ can use this endpoint to verify or
	// audit the current configuration of a project. This can be helpful in ensuring that the project's
	// settings align with expectations.
	// The payload also includes a read-only set of `derived` dimensions which help describe the
	// currently configured `machineType`.
	//
	// GET /htc/projects/{projectId}/dimensions
	GetDimensions(ctx context.Context, params GetDimensionsParams) (GetDimensionsRes, error)
	// GetLimits implements getLimits operation.
	//
	// This endpoint will list all resource limitations associated with this project.
	// A job running in this project will be subject to all resulting limits as well as any associated
	// with the workspace (see `/htc/workspaces/{workspaceId}/limits`).
	//
	// GET /htc/projects/{projectId}/limits
	GetLimits(ctx context.Context, params GetLimitsParams) (GetLimitsRes, error)
	// GetProject implements getProject operation.
	//
	// This endpoint will get a project by id.
	//
	// GET /htc/projects/{projectId}
	GetProject(ctx context.Context, params GetProjectParams) (GetProjectRes, error)
	// GetProjects implements getProjects operation.
	//
	// This endpoint will get all projects.
	//
	// GET /htc/projects
	GetProjects(ctx context.Context, params GetProjectsParams) (GetProjectsRes, error)
}

// TaskHandler handles operations described by OpenAPI v3 specification.
//
// x-ogen-operation-group: Task
type TaskHandler interface {
	// CreateTask implements createTask operation.
	//
	// This endpoint will create a task for a project.
	//
	// POST /htc/projects/{projectId}/tasks
	CreateTask(ctx context.Context, req OptHTCTask, params CreateTaskParams) (CreateTaskRes, error)
	// GetTasks implements getTasks operation.
	//
	// This endpoint will get all tasks in a project.
	//
	// GET /htc/projects/{projectId}/tasks
	GetTasks(ctx context.Context, params GetTasksParams) (GetTasksRes, error)
}

// Server implements http server based on OpenAPI v3 specification and
// calls Handler to handle requests.
type Server struct {
	h   Handler
	sec SecurityHandler
	baseServer
}

// NewServer creates new Server.
func NewServer(h Handler, sec SecurityHandler, opts ...ServerOption) (*Server, error) {
	s, err := newServerConfig(opts...).baseServer()
	if err != nil {
		return nil, err
	}
	return &Server{
		h:          h,
		sec:        sec,
		baseServer: s,
	}, nil
}
