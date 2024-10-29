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
  },
  components+: {
    schemas+: {
      HTCTokenPayload: { type: 'string' },
      // They said it was only ever a string! But this isn't always
      // true.
      HTCProject+: {
        properties+: {
          organizationCode+: {
            nullable: true,
          },
        },
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
    },
  },
};

// Render out the original swagger + our patches.
swagger + patches
