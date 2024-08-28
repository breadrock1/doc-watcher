// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/hello/": {
            "get": {
                "description": "Check service is available",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "hello"
                ],
                "summary": "Hello",
                "operationId": "hello",
                "responses": {
                    "200": {
                        "description": "Ok",
                        "schema": {
                            "$ref": "#/definitions/server.ResponseForm"
                        }
                    },
                    "503": {
                        "description": "Server does not available",
                        "schema": {
                            "$ref": "#/definitions/server.ServerErrorForm"
                        }
                    }
                }
            }
        },
        "/storage/bucket": {
            "put": {
                "description": "Create new bucket into storage",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "storage"
                ],
                "summary": "Create new bucket into storage",
                "operationId": "create-bucket",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Bucket name to create",
                        "name": "bucket",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Bucket name to create",
                        "name": "jsonQuery",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/server.BucketNameForm"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Ok",
                        "schema": {
                            "$ref": "#/definitions/server.ResponseForm"
                        }
                    },
                    "400": {
                        "description": "Bad Request message",
                        "schema": {
                            "$ref": "#/definitions/server.BadRequestForm"
                        }
                    },
                    "503": {
                        "description": "Server does not available",
                        "schema": {
                            "$ref": "#/definitions/server.ServerErrorForm"
                        }
                    }
                }
            }
        },
        "/storage/buckets": {
            "get": {
                "description": "Get watched bucket list",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "storage"
                ],
                "summary": "Get watched bucket list",
                "operationId": "get-buckets",
                "responses": {
                    "200": {
                        "description": "Ok",
                        "schema": {
                            "type": "array",
                            "items": {
                                "type": "string"
                            }
                        }
                    },
                    "503": {
                        "description": "Server does not available",
                        "schema": {
                            "$ref": "#/definitions/server.ServerErrorForm"
                        }
                    }
                }
            }
        },
        "/storage/{bucket}": {
            "delete": {
                "description": "Remove bucket from storage",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "storage"
                ],
                "summary": "Remove bucket from storage",
                "operationId": "remove-bucket",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Bucket name to remove",
                        "name": "bucket",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Ok",
                        "schema": {
                            "$ref": "#/definitions/server.ResponseForm"
                        }
                    },
                    "400": {
                        "description": "Bad Request message",
                        "schema": {
                            "$ref": "#/definitions/server.BadRequestForm"
                        }
                    },
                    "503": {
                        "description": "Server does not available",
                        "schema": {
                            "$ref": "#/definitions/server.ServerErrorForm"
                        }
                    }
                }
            }
        },
        "/storage/{bucket}/file/copy": {
            "post": {
                "description": "Copy file to another location into bucket",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "storage"
                ],
                "summary": "Copy file to another location into bucket",
                "operationId": "copy-file",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Bucket name of src file",
                        "name": "bucket",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Params to copy file",
                        "name": "jsonQuery",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/server.CopyFileForm"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Ok",
                        "schema": {
                            "$ref": "#/definitions/server.ResponseForm"
                        }
                    },
                    "400": {
                        "description": "Bad Request message",
                        "schema": {
                            "$ref": "#/definitions/server.BadRequestForm"
                        }
                    },
                    "503": {
                        "description": "Server does not available",
                        "schema": {
                            "$ref": "#/definitions/server.ServerErrorForm"
                        }
                    }
                }
            }
        },
        "/storage/{bucket}/file/download": {
            "post": {
                "description": "Download file from storage",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "storage"
                ],
                "summary": "Download file from storage",
                "operationId": "download-file",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Bucket name to download file",
                        "name": "bucket",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Parameters to download file",
                        "name": "jsonQuery",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/server.DownloadFile"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Ok",
                        "schema": {
                            "type": "file"
                        }
                    },
                    "400": {
                        "description": "Bad Request message",
                        "schema": {
                            "$ref": "#/definitions/server.BadRequestForm"
                        }
                    },
                    "503": {
                        "description": "Server does not available",
                        "schema": {
                            "$ref": "#/definitions/server.ServerErrorForm"
                        }
                    }
                }
            }
        },
        "/storage/{bucket}/file/move": {
            "post": {
                "description": "Move file to another location into bucket",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "storage"
                ],
                "summary": "Move file to another location into bucket",
                "operationId": "move-file",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Bucket name of src file",
                        "name": "bucket",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Params to move file",
                        "name": "jsonQuery",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/server.CopyFileForm"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Ok",
                        "schema": {
                            "$ref": "#/definitions/server.ResponseForm"
                        }
                    },
                    "400": {
                        "description": "Bad Request message",
                        "schema": {
                            "$ref": "#/definitions/server.BadRequestForm"
                        }
                    },
                    "503": {
                        "description": "Server does not available",
                        "schema": {
                            "$ref": "#/definitions/server.ServerErrorForm"
                        }
                    }
                }
            }
        },
        "/storage/{bucket}/file/remove": {
            "post": {
                "description": "Remove file from storage",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "storage"
                ],
                "summary": "Remove file from storage",
                "operationId": "remove-file",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Bucket name to remove file",
                        "name": "bucket",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Parameters to remove file",
                        "name": "jsonQuery",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/server.RemoveFile"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Ok",
                        "schema": {
                            "$ref": "#/definitions/server.ResponseForm"
                        }
                    },
                    "400": {
                        "description": "Bad Request message",
                        "schema": {
                            "$ref": "#/definitions/server.BadRequestForm"
                        }
                    },
                    "503": {
                        "description": "Server does not available",
                        "schema": {
                            "$ref": "#/definitions/server.ServerErrorForm"
                        }
                    }
                }
            }
        },
        "/storage/{bucket}/file/upload": {
            "post": {
                "description": "Upload files to storage",
                "consumes": [
                    "multipart/form"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "storage"
                ],
                "summary": "Upload files to storage",
                "operationId": "upload-files",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Bucket name to upload files",
                        "name": "bucket",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "file",
                        "description": "Files multipart form",
                        "name": "files",
                        "in": "formData",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Ok",
                        "schema": {
                            "$ref": "#/definitions/server.ResponseForm"
                        }
                    },
                    "400": {
                        "description": "Bad Request message",
                        "schema": {
                            "$ref": "#/definitions/server.BadRequestForm"
                        }
                    },
                    "503": {
                        "description": "Server does not available",
                        "schema": {
                            "$ref": "#/definitions/server.ServerErrorForm"
                        }
                    }
                }
            }
        },
        "/storage/{bucket}/files": {
            "post": {
                "description": "Get files list into bucket",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "storage"
                ],
                "summary": "Get files list into bucket",
                "operationId": "get-list-files",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Bucket name to get list files",
                        "name": "bucket",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Parameters to get list files",
                        "name": "jsonQuery",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/server.ListFilesForm"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Ok",
                        "schema": {
                            "$ref": "#/definitions/server.ResponseForm"
                        }
                    },
                    "400": {
                        "description": "Bad Request message",
                        "schema": {
                            "$ref": "#/definitions/server.BadRequestForm"
                        }
                    },
                    "503": {
                        "description": "Server does not available",
                        "schema": {
                            "$ref": "#/definitions/server.ServerErrorForm"
                        }
                    }
                }
            }
        },
        "/watcher/attach": {
            "post": {
                "description": "Attach new directory to watcher",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "watcher"
                ],
                "summary": "Attach new directory to watcher",
                "operationId": "folders-attach",
                "parameters": [
                    {
                        "description": "File entity",
                        "name": "jsonQuery",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/server.WatcherDirectoriesForm"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Ok",
                        "schema": {
                            "$ref": "#/definitions/server.ResponseForm"
                        }
                    },
                    "400": {
                        "description": "Bad Request message",
                        "schema": {
                            "$ref": "#/definitions/server.BadRequestForm"
                        }
                    },
                    "503": {
                        "description": "Server does not available",
                        "schema": {
                            "$ref": "#/definitions/server.ServerErrorForm"
                        }
                    }
                }
            }
        },
        "/watcher/detach": {
            "post": {
                "description": "Attach new directory to watcher",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "watcher"
                ],
                "summary": "Attach new directory to watcher",
                "operationId": "folders-detach",
                "parameters": [
                    {
                        "description": "Folder ids",
                        "name": "jsonQuery",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/server.WatcherDirectoriesForm"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Ok",
                        "schema": {
                            "$ref": "#/definitions/server.ResponseForm"
                        }
                    },
                    "400": {
                        "description": "Bad Request message",
                        "schema": {
                            "$ref": "#/definitions/server.BadRequestForm"
                        }
                    },
                    "503": {
                        "description": "Server does not available",
                        "schema": {
                            "$ref": "#/definitions/server.ServerErrorForm"
                        }
                    }
                }
            }
        },
        "/watcher/is-pause": {
            "get": {
                "description": "Check does watcher has been paused",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "watcher"
                ],
                "summary": "Check does watcher has been paused",
                "operationId": "is-watcher-pause",
                "responses": {
                    "200": {
                        "description": "Ok",
                        "schema": {
                            "$ref": "#/definitions/server.ResponseForm"
                        }
                    },
                    "400": {
                        "description": "Bad Request message",
                        "schema": {
                            "$ref": "#/definitions/server.BadRequestForm"
                        }
                    },
                    "503": {
                        "description": "Server does not available",
                        "schema": {
                            "$ref": "#/definitions/server.ServerErrorForm"
                        }
                    }
                }
            }
        },
        "/watcher/pause": {
            "get": {
                "description": "Pause all watchers",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "watcher"
                ],
                "summary": "Pause all watchers",
                "operationId": "watcher-pause",
                "responses": {
                    "200": {
                        "description": "Ok",
                        "schema": {
                            "$ref": "#/definitions/server.ResponseForm"
                        }
                    },
                    "400": {
                        "description": "Bad Request message",
                        "schema": {
                            "$ref": "#/definitions/server.BadRequestForm"
                        }
                    },
                    "503": {
                        "description": "Server does not available",
                        "schema": {
                            "$ref": "#/definitions/server.ServerErrorForm"
                        }
                    }
                }
            }
        },
        "/watcher/run": {
            "get": {
                "description": "Run all watchers",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "watcher"
                ],
                "summary": "Run all watchers",
                "operationId": "watcher-run",
                "responses": {
                    "200": {
                        "description": "Ok",
                        "schema": {
                            "$ref": "#/definitions/server.ResponseForm"
                        }
                    },
                    "400": {
                        "description": "Bad Request message",
                        "schema": {
                            "$ref": "#/definitions/server.BadRequestForm"
                        }
                    },
                    "503": {
                        "description": "Server does not available",
                        "schema": {
                            "$ref": "#/definitions/server.ServerErrorForm"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "server.BadRequestForm": {
            "type": "object",
            "properties": {
                "message": {
                    "type": "string",
                    "example": "Bad Request message"
                },
                "status": {
                    "type": "integer",
                    "example": 400
                }
            }
        },
        "server.BucketNameForm": {
            "type": "object",
            "properties": {
                "bucket_name": {
                    "type": "string",
                    "example": "test-bucket"
                }
            }
        },
        "server.CopyFileForm": {
            "type": "object",
            "properties": {
                "dst_path": {
                    "type": "string",
                    "example": "test-document.docx"
                },
                "src_path": {
                    "type": "string",
                    "example": "old-test-document.docx"
                }
            }
        },
        "server.DownloadFile": {
            "type": "object",
            "properties": {
                "file_name": {
                    "type": "string",
                    "example": "test-file.docx"
                }
            }
        },
        "server.ListFilesForm": {
            "type": "object",
            "properties": {
                "directory": {
                    "type": "string",
                    "example": "test-folder/"
                }
            }
        },
        "server.RemoveFile": {
            "type": "object",
            "properties": {
                "file_name": {
                    "type": "string",
                    "example": "test-file.docx"
                }
            }
        },
        "server.ResponseForm": {
            "type": "object",
            "properties": {
                "message": {
                    "type": "string",
                    "example": "Done"
                },
                "status": {
                    "type": "integer",
                    "example": 200
                }
            }
        },
        "server.ServerErrorForm": {
            "type": "object",
            "properties": {
                "message": {
                    "type": "string",
                    "example": "Server Error message"
                },
                "status": {
                    "type": "integer",
                    "example": 503
                }
            }
        },
        "server.WatcherDirectoriesForm": {
            "type": "object",
            "properties": {
                "paths": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    },
                    "example": [
                        "./indexer/test_folder"
                    ]
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "",
	Host:             "",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "",
	Description:      "",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
