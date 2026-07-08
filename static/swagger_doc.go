// Package static provides pre-written swagger documentation for the static JSON
// file serving endpoints.  The spec is registered with the server's DocRegistry
// during module Initialize so that swagger.MergeAndRegister picks it up.
package static

// StaticSwaggerJSON is the Swagger 2.0 spec for the static module endpoints.
const StaticSwaggerJSON = `{
  "swagger": "2.0",
  "info": {
    "title": "Static Assets API",
    "description": "Securely serve static JSON files from the server.",
    "version": "1.0.0"
  },
  "basePath": "/api/v1",
  "paths": {
    "/static": {
      "get": {
        "security": [{"BearerAuth": []}],
        "summary": "List available static JSON files",
        "description": "Returns the names of all JSON files available in the statics/json/ directory.",
        "produces": ["application/json"],
        "tags": ["static"],
        "operationId": "listStaticFiles",
        "responses": {
          "200": {
            "description": "List of available file names (without .json extension)",
            "schema": {
              "type": "object",
              "properties": {
                "available_files": {"type": "array", "items": {"type": "string"}},
                "base_url":        {"type": "string"},
                "example_usage":   {"type": "string"},
                "security_note":   {"type": "string"},
                "restrictions":    {"type": "string"},
                "note":            {"type": "string"}
              }
            }
          },
          "401": {"description": "Unauthorized"},
          "500": {"description": "Internal server error"}
        }
      }
    },
    "/static/{filename}": {
      "get": {
        "security": [{"BearerAuth": []}],
        "summary": "Serve a static JSON file",
        "description": "Returns the raw JSON content of the requested file. Filename must be alphanumeric with hyphens/underscores only (no extension).",
        "produces": ["application/json"],
        "tags": ["static"],
        "operationId": "serveStaticFile",
        "parameters": [
          {
            "name": "filename",
            "in": "path",
            "required": true,
            "type": "string",
            "description": "File name without .json extension (e.g. 'diagnostic_std')"
          }
        ],
        "responses": {
          "200": {"description": "Raw JSON file content"},
          "400": {"description": "Invalid file name (contains illegal characters)"},
          "401": {"description": "Unauthorized"},
          "404": {"description": "File not found"}
        }
      }
    }
  },
  "securityDefinitions": {
    "BearerAuth": {
      "type": "apiKey",
      "name": "Authorization",
      "in": "header"
    }
  },
  "definitions": {}
}`
