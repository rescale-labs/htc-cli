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
        security: securityScheme,
      },
    },

    // Use named component for response and require token.
    '/auth/whoami'+: {
      get+: {
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
    '/htc/projects'+: {
      get+: {
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
    },
    '/htc/projects/{projectId}/dimensions'+: {
      get+: {
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
    },
    '/htc/projects/{projectId}/tasks/{taskId}/jobs/batch'+: {
      post+: {
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
  },
  components+: {
    schemas+: {
      HTCJobSubmitRequests: {
        items: {
          '$ref': '#/components/schemas/HTCJobSubmitRequest',
        },
        type: 'array',
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
        items: {
          '$ref': '#/components/schemas/HTCComputeEnvironment',
        },
        type: 'array',
      },
      HTCProjectLimits: {
        items: {
          '$ref': '#/components/schemas/HTCProjectLimit',
        },
        type: 'array',
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
          message: { type: 'string' },
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
    },
  },
};

// Render out the original swagger + our patches.
swagger + patches
