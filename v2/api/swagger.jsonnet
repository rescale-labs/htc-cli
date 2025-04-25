local swagger = import 'swagger.json';

// Annotates that a resource requires a token. Used in two places where
// we are expecting a bearer token (rescale API key).
local securityScheme = [{ SecurityScheme: [] }];

// Patches we apply to the stock swagger JSON.
//
// Long term, all of these should move into the API itself. But this
// works fairly well as a stopgap.
//
// Note that `+:` informs jsonnet that, when adding this object to
// another object, we shouldn't clobber the earlier (left-hand) object's
// other keys.
local patches = {
  paths+: {
    '/auth/token'+: {
      get+: {
        'x-ogen-operation-group': 'Auth',
        operationId: 'getToken',
        security: securityScheme,
      },
    },

    // Use named component for response and require token.
    '/auth/whoami'+: {
      get+: {
        'x-ogen-operation-group': 'Auth',
        operationId: 'whoAmI',
        responses+: {
          '200'+: {
            content+: {
              'application/json'+: {
                schema: { '$ref': '#/components/schemas/WhoAmI' },
              },
            },
          },
        },
        security: securityScheme,
      },
    },
    '/htc/metrics'+: {
      get+: {
        'x-ogen-operation-group': 'Metrics',
        operationId: 'getMetrics',
      },
    },
    '/htc/projects'+: {
      get+: {
        'x-ogen-operation-group': 'Project',
        operationId: 'getProjects',
        responses+: {
          '200'+: {
            content+: {
              'application/json'+: {
                schema: { '$ref':
                  '#/components/schemas/HTCProjectsResponse' },

              },
            },
          },
        },
      },
      post+: {
        'x-ogen-operation-group': 'Project',
        operationId: 'createProject',
        responses+: {
            '400'+: {
                content+: {
                    'application/json'+: {
                        schema: { '$ref': '#/components/schemas/HTCRequestError' },
                    }
                }
            }
        }
      },
    },
    '/htc/projects/{projectId}'+: {
      get+: {
        'x-ogen-operation-group': 'Project',
        operationId: 'getProject',
      },
    },
    '/htc/projects/{projectId}/container-registry/images'+: {
      get+: {
        'x-ogen-operation-group': 'Image',
        operationId: 'getImages',
      },
    },
    '/htc/projects/{projectId}/container-registry/images/{imageName}'+: {
      get+: {
        'x-ogen-operation-group': 'Image',
        operationId: 'getImage',
      },
    },
    '/htc/projects/{projectId}/container-registry/repo/{repoName}'+: {
      post+: {
        'x-ogen-operation-group': 'Image',
        operationId: 'createRepo',
      },
    },
    '/htc/projects/{projectId}/container-registry/token'+: {
      get+: {
        'x-ogen-operation-group': 'Image',
        operationId: 'getRegistryToken',
        responses+: {
          '404': {
            content: {
              'application/json': {
                schema: { '$ref': '#/components/schemas/OAuth2ErrorResponse'}
              }
            }
          }
        }
      },
    },
    '/htc/projects/{projectId}/dimensions'+: {
      get+: {
        'x-ogen-operation-group': 'Project',
        operationId: 'getDimensions',
        responses+: {
          '200'+: {
            content+: {
              'application/json'+: {
                schema: { '$ref': '#/components/schemas/HTCProjectDimensions' },
              },
            },
          },
        },
      },
    },
    '/htc/projects/{projectId}/limits'+: {
      get+: {
        'x-ogen-operation-group': 'Project',
        operationId: 'getLimits',
        responses+: {
          '200'+: {
            content+: {
              'application/json'+: {
                schema: { '$ref': '#/components/schemas/HTCProjectLimits' },
              },
            },
          },
        },
      },
    },
    '/htc/projects/{projectId}/tasks'+: {
      get+: {
        'x-ogen-operation-group': 'Task',
        operationId: 'getTasks',
        responses+: {
          '200'+: {
            content+: {
              'application/json'+: {
                schema: { '$ref':
                  '#/components/schemas/HTCTasksResponse' },

              },
            },
          },
        },
      },
      post+: {
        'x-ogen-operation-group': 'Task',
        operationId: 'createTask',
      },
    },
  '/htc/projects/{projectId}/tasks/{taskId}/summary-statistics'+: {
    get+: {
      'x-ogen-operation-group': 'Task',
      operationId: 'GetTaskStats',
    }
  },
    '/htc/projects/{projectId}/tasks/{taskId}/jobs'+: {
      get+: {
        'x-ogen-operation-group': 'Job',
        operationId: 'getJobs',
        responses+: {
          '200'+: {
            content+: {
              'application/json'+: {
                schema: { '$ref':
                  '#/components/schemas/HTCJobs' },
              },
            },
          },
          '404'+: {
            content+: {
              'application/json'+: {
                schema: { '$ref': '#/components/schemas/HTCRequestError' },
              },
            },
          }
        },
      },
    },
    '/htc/projects/{projectId}/tasks/{taskId}/jobs/{jobId}'+: {
      get+: {
        'x-ogen-operation-group': 'Job',
        operationId: 'getJob',
        responses+: {
          '404': {
            content: {
              'application/json': {
                schema: { '$ref': '#/components/schemas/HTCRequestError' },
              },
            },
          },
        },
      },
    },
    '/htc/projects/{projectId}/tasks/{taskId}/jobs/batch'+: {
      post+: {
        'x-ogen-operation-group': 'Job',
        operationId: 'submitJobs',
        responses+: {
          '200'+: {
            content+: {
              'application/json'+: {
                schema: { '$ref':
                  '#/components/schemas/HTCJobSubmitRequests' },
              },
            },
          },
          '400': {
            content: {
              'application/json': {
                schema: { '$ref': '#/components/schemas/HTCRequestError' },
              },
            },
          },
        },
      },
    },
    '/htc/projects/{projectId}/tasks/{taskId}/jobs/cancel'+: {
      post+: {
        'x-ogen-operation-group': 'Job',
        operationId: 'cancelJobs',
      },
    },
    '/htc/projects/{projectId}/tasks/{taskId}/jobs/{jobId}/logs'+: {
      get+: {
        'x-ogen-operation-group': 'Job',
        operationId: 'getLogs',
        responses+: {
          '200'+: {
              content+: {
                'application/json'+: {
                  schema: { '$ref':
                    '#/components/schemas/HTCJobLogs' },
                },
              },
            },
          '404': {
            content: {
              'application/json': {
                schema: { '$ref': '#/components/schemas/HTCRequestError' },
              },
            },
          },
        },
      }
    },
    '/htc/projects/{projectId}/tasks/{taskId}/jobs/{jobId}/events'+: {
      get+: {
        'x-ogen-operation-group': 'Job',
        operationId: 'getEvents',
        responses+: {
          '200'+: {
            content+: {
              'application/json'+: {
                schema: { '$ref':
                  '#/components/schemas/HTCJobStatusEvents' },
              }
            }
          },
          '404': {
            content: {
              'application/json': {
                schema: { '$ref': '#/components/schemas/HTCRequestError' },
              },
            },
          },
        }
      }
    }
  },
  components+: {
    schemas+: {
      ContainerDetails+: {
        nullable: true,
      },
      // Need to remove the enum to avoid `non-primitive enum not implemented`
      // In this schema the status enum string is used as a key
      JobStatusSummary+: {
        properties+: {
          jobStatuses: {
            type: 'object',
            properties: {
              FAILED: {
                example: 5,
                type: 'integer'
              },
              RUNNABLE: {
                example: 5,
                type: 'integer'
              },
              RUNNING: {
                example: 5,
                type: 'integer'
              },
              STARTING: {
                example: 5,
                type: 'integer'
              },
              SUBMITTED_TO_PROVIDER: {
                example: 5,
                type: 'integer'
              },
              SUBMITTED_TO_RESCALE: {
                example: 5,
                type: 'integer'
              },
              SUCCEEDED: {
                example: 5,
                type: 'integer'
              }
            }
          }
        }
      },
      HTCJob+: {
        properties+: {
          completedAt: {
            '$ref': '#/components/schemas/NullableInstant',
          },
          statusReason+: {
            nullable: true,
          },
          startedAt: {
            '$ref': '#/components/schemas/NullableInstant',
          },
        },
      },
      HTCJobErrorCodeName+: {
        nullable: true,
      },
      HTCJobs: {
        type: 'object',
        properties: {
          items: {
            items: {
              '$ref': '#/components/schemas/HTCJob',
            },
            type: 'array',
          },
          next: {
            example: 'https://page2.com',
            format: 'uri',
            type: 'string',
          },
        },
      },
      HTCJobStatusEvent+: {
        properties+: {
          statusReason+: {
            nullable: true,
          },
        },
      },
      HTCJobStatusEvents+: {
        type: 'object',
        properties+: {
          items: {
            items: {
              '$ref': '#/components/schemas/RescaleJobStatusEvent',
            },
            type: 'array',
          },
          next: {
            example: 'https://page2.com',
            format: 'uri',
            type: 'string',
          },
        },
      },
      HTCJobSubmitRequest+: {
        properties+: {
          rescaleProjectId+: {
            nullable: true,
          },

        },
      },
      HTCJobDefinition+: {
        properties+: {
            workingDir+: {
                nullable: true,
            },
        },
      },
      HTCJobSubmitRequests: {
        type: 'array',
        items: {
          '$ref': '#/components/schemas/HTCJobSubmitRequest',
        },
      },
      HTCJobLogs: {
        type: 'object',
        properties: {
          items: {
            type: 'array',
            items: {
              '$ref': '#/components/schemas/HTCLogEvent',
            }
          },
          next: {
            format: 'uri',
            type: 'string',
            example: 'https://page2.com',
          }
        }
      },
      HTCProject+: {
        properties+: {
          organizationCode+: {
            // The docs said it was only ever a string! But sometimes
            // it's null and then we fail at runtime :-(
            nullable: true,
          },
        },
      },
      HTCProjectDimensions: {
        type: 'array',
        items: {
          '$ref': '#/components/schemas/HTCComputeEnvironment',
        },
      },
      HTCProjectLimits: {
        type: 'array',
        items: {
          '$ref': '#/components/schemas/HTCProjectLimit',
        },
      },
      HTCProjectsResponse: {
        type: 'object',
        properties: {
          items: {
            type: 'array',
            items: {
              '$ref': '#/components/schemas/HTCProject',
            },
          },
          next: {
            format: 'uri',
            type: 'string',
            example: 'https://page2.com',
          },
        },
      },
      HTCRequestError: {  // Catchall for many HTC error responses
        type: 'object',
        properties: {
          errorType: { type: 'string' },
          errorDescription: { type: 'string' },
          message: {
            type: 'string',
            nullable: true
          },
        },
      },
      HTCTask+: {
        properties+: {
          archivedAt+: {
            nullable: true,
          },
          deletedAt+: {
            nullable: true,
          },
        },
      },
      HTCTasksResponse: {
        type: 'object',
        properties: {
          items: {
            type: 'array',
            items: {
              '$ref': '#/components/schemas/HTCTask',
            },
          },
          next: {
            format: 'uri',
            type: 'string',
            example: 'https://page2.com',
          },
        },
      },
      HTCTokenPayload: { type: 'string' },
      NullableInstant: {
        example: '2022-03-10T16:15:50Z',
        format: 'date-time',
        type: 'string',
        nullable: true,
      },
      InstanceLabels+: {
        nullable: true,
      },
      RescaleJobStatusEvent+: {
        properties+: {
          statusReason+: {
            nullable: true,
          },
          instanceId+: {
            nullable: true,
          }
        }
      },
    },
  },
};

// Render out the original swagger + our patches.
swagger + patches
