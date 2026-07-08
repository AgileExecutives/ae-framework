// Package pdf provides pre-written swagger documentation for the PDF generation
// endpoints.  The spec is registered with the server's DocRegistry during
// module Initialize so that swagger.MergeAndRegister picks it up.
package pdf

// PDFSwaggerJSON is the Swagger 2.0 spec for the pdf module endpoints.
const PDFSwaggerJSON = `{
  "swagger": "2.0",
  "info": {
    "title": "PDF API",
    "description": "PDF document generation endpoints.",
    "version": "1.0.0"
  },
  "basePath": "/api/v1",
  "paths": {
    "/pdf/create": {
      "post": {
        "security": [{"BearerAuth": []}],
        "summary": "Generate a PDF document",
        "description": "Renders a named template with the provided data and returns the generated file name.",
        "consumes": ["application/json"],
        "produces": ["application/json"],
        "tags": ["pdf"],
        "operationId": "createPDF",
        "parameters": [
          {
            "name": "body",
            "in": "body",
            "required": true,
            "schema": {
              "type": "object",
              "required": ["data", "templateName", "fileName"],
              "properties": {
                "data":         {"type": "object",  "description": "Template data as key-value pairs"},
                "templateName": {"type": "string",  "example": "invoice"},
                "fileName":     {"type": "string",  "example": "invoice-2024-001.pdf"}
              }
            }
          }
        ],
        "responses": {
          "200": {
            "description": "PDF generated successfully",
            "schema": {
              "type": "object",
              "properties": {
                "success":  {"type": "boolean"},
                "message":  {"type": "string"},
                "filename": {"type": "string"}
              }
            }
          },
          "400": {"description": "Invalid request body or missing required fields"},
          "401": {"description": "Unauthorized"},
          "500": {"description": "PDF generation failed"}
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
