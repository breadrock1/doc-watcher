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
                            "$ref": "#/definitions/endpoints.ResponseForm"
                        }
                    },
                    "503": {
                        "description": "Server does not available",
                        "schema": {
                            "$ref": "#/definitions/endpoints.ServerErrorForm"
                        }
                    }
                }
            }
        },
        "/watcher/files/analyse": {
            "post": {
                "description": "Analyse uploaded files by ids",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "files"
                ],
                "summary": "Analyse uploaded files by ids",
                "operationId": "files-analyse",
                "parameters": [
                    {
                        "description": "Document ids to analyse",
                        "name": "jsonQuery",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/endpoints.AnalyseFilesForm"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Ok",
                        "schema": {
                            "$ref": "#/definitions/endpoints.ResponseForm"
                        }
                    },
                    "400": {
                        "description": "Bad Request message",
                        "schema": {
                            "$ref": "#/definitions/endpoints.BadRequestForm"
                        }
                    },
                    "503": {
                        "description": "Server does not available",
                        "schema": {
                            "$ref": "#/definitions/endpoints.ServerErrorForm"
                        }
                    }
                }
            }
        },
        "/watcher/files/download": {
            "post": {
                "description": "Download file by path",
                "produces": [
                    "multipart/form"
                ],
                "tags": [
                    "files"
                ],
                "summary": "Download file by path",
                "operationId": "files-download",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Path to file",
                        "name": "file_path",
                        "in": "formData",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Ok",
                        "schema": {
                            "$ref": "#/definitions/endpoints.ResponseForm"
                        }
                    },
                    "400": {
                        "description": "Bad Request message",
                        "schema": {
                            "$ref": "#/definitions/endpoints.BadRequestForm"
                        }
                    },
                    "503": {
                        "description": "Server does not available",
                        "schema": {
                            "$ref": "#/definitions/endpoints.ServerErrorForm"
                        }
                    }
                }
            }
        },
        "/watcher/files/move": {
            "post": {
                "description": "Moving files to target directory",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "files"
                ],
                "summary": "Moving files to target directory",
                "operationId": "moving",
                "parameters": [
                    {
                        "description": "Document ids to move",
                        "name": "jsonQuery",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/endpoints.MoveFilesForm"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Ok",
                        "schema": {
                            "$ref": "#/definitions/endpoints.ResponseForm"
                        }
                    },
                    "400": {
                        "description": "Bad Request message",
                        "schema": {
                            "$ref": "#/definitions/endpoints.BadRequestForm"
                        }
                    },
                    "503": {
                        "description": "Server does not available",
                        "schema": {
                            "$ref": "#/definitions/endpoints.ServerErrorForm"
                        }
                    }
                }
            }
        },
        "/watcher/files/remove": {
            "post": {
                "description": "Remove files from directory",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "files"
                ],
                "summary": "Remove files from directory",
                "operationId": "files-remove",
                "parameters": [
                    {
                        "description": "Document paths to remove",
                        "name": "jsonQuery",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/endpoints.RemoveFilesForm"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Ok",
                        "schema": {
                            "$ref": "#/definitions/endpoints.ResponseForm"
                        }
                    },
                    "400": {
                        "description": "Bad Request message",
                        "schema": {
                            "$ref": "#/definitions/endpoints.RemoveFilesError"
                        }
                    },
                    "503": {
                        "description": "Server does not available",
                        "schema": {
                            "$ref": "#/definitions/endpoints.ServerErrorForm"
                        }
                    }
                }
            }
        },
        "/watcher/files/unrecognized": {
            "get": {
                "description": "Get unrecognized documents",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "files"
                ],
                "summary": "Get unrecognized documents",
                "operationId": "files-unrecognized",
                "responses": {
                    "200": {
                        "description": "Ok",
                        "schema": {
                            "$ref": "#/definitions/endpoints.UnrecognizedDocuments"
                        }
                    },
                    "400": {
                        "description": "Bad Request message",
                        "schema": {
                            "$ref": "#/definitions/endpoints.BadRequestForm"
                        }
                    },
                    "503": {
                        "description": "Server does not available",
                        "schema": {
                            "$ref": "#/definitions/endpoints.ServerErrorForm"
                        }
                    }
                }
            }
        },
        "/watcher/files/upload": {
            "post": {
                "description": "Upload files to analyse",
                "consumes": [
                    "multipart/form"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "files"
                ],
                "summary": "Upload files to analyse",
                "operationId": "files-upload",
                "parameters": [
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
                            "$ref": "#/definitions/endpoints.ResponseForm"
                        }
                    },
                    "400": {
                        "description": "Bad Request message",
                        "schema": {
                            "$ref": "#/definitions/endpoints.BadRequestForm"
                        }
                    },
                    "503": {
                        "description": "Server does not available",
                        "schema": {
                            "$ref": "#/definitions/endpoints.ServerErrorForm"
                        }
                    }
                }
            }
        },
        "/watcher/folders/": {
            "get": {
                "description": "Get watched directories list",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "watcher"
                ],
                "summary": "Get watched directories list",
                "operationId": "folders-all",
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
                            "$ref": "#/definitions/endpoints.ServerErrorForm"
                        }
                    }
                }
            }
        },
        "/watcher/folders/attach": {
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
                            "$ref": "#/definitions/endpoints.WatcherDirectoriesForm"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Ok",
                        "schema": {
                            "$ref": "#/definitions/endpoints.ResponseForm"
                        }
                    },
                    "400": {
                        "description": "Bad Request message",
                        "schema": {
                            "$ref": "#/definitions/endpoints.BadRequestForm"
                        }
                    },
                    "503": {
                        "description": "Server does not available",
                        "schema": {
                            "$ref": "#/definitions/endpoints.ServerErrorForm"
                        }
                    }
                }
            }
        },
        "/watcher/folders/create": {
            "post": {
                "description": "Create folder to store documents",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "files"
                ],
                "summary": "Create folder to store documents",
                "operationId": "folder-create",
                "parameters": [
                    {
                        "description": "Folder name to create",
                        "name": "jsonQuery",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/endpoints.FolderNameForm"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Ok",
                        "schema": {
                            "$ref": "#/definitions/endpoints.ResponseForm"
                        }
                    },
                    "400": {
                        "description": "Bad Request message",
                        "schema": {
                            "$ref": "#/definitions/endpoints.BadRequestForm"
                        }
                    },
                    "503": {
                        "description": "Server does not available",
                        "schema": {
                            "$ref": "#/definitions/endpoints.ServerErrorForm"
                        }
                    }
                }
            }
        },
        "/watcher/folders/detach": {
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
                            "$ref": "#/definitions/endpoints.WatcherDirectoriesForm"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Ok",
                        "schema": {
                            "$ref": "#/definitions/endpoints.ResponseForm"
                        }
                    },
                    "400": {
                        "description": "Bad Request message",
                        "schema": {
                            "$ref": "#/definitions/endpoints.BadRequestForm"
                        }
                    },
                    "503": {
                        "description": "Server does not available",
                        "schema": {
                            "$ref": "#/definitions/endpoints.ServerErrorForm"
                        }
                    }
                }
            }
        },
        "/watcher/folders/remove": {
            "post": {
                "description": "Remove folder",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "files"
                ],
                "summary": "Remove folder",
                "operationId": "folder-remove",
                "parameters": [
                    {
                        "description": "Folder name to remove",
                        "name": "jsonQuery",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/endpoints.FolderNameForm"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Ok",
                        "schema": {
                            "$ref": "#/definitions/endpoints.ResponseForm"
                        }
                    },
                    "400": {
                        "description": "Bad Request message",
                        "schema": {
                            "$ref": "#/definitions/endpoints.BadRequestForm"
                        }
                    },
                    "503": {
                        "description": "Server does not available",
                        "schema": {
                            "$ref": "#/definitions/endpoints.ServerErrorForm"
                        }
                    }
                }
            }
        },
        "/watcher/folders/upload": {
            "post": {
                "description": "Upload file to watcher directory",
                "consumes": [
                    "multipart/form"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "watcher"
                ],
                "summary": "Upload file to watcher directory",
                "operationId": "watcher-upload",
                "parameters": [
                    {
                        "type": "file",
                        "description": "Files multipart form",
                        "name": "files",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Directory to upload",
                        "name": "directory",
                        "in": "formData",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Ok",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/reader.DocumentPreview"
                            }
                        }
                    },
                    "400": {
                        "description": "Bad Request message",
                        "schema": {
                            "$ref": "#/definitions/endpoints.BadRequestForm"
                        }
                    },
                    "503": {
                        "description": "Server does not available",
                        "schema": {
                            "$ref": "#/definitions/endpoints.ServerErrorForm"
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
                            "$ref": "#/definitions/endpoints.ResponseForm"
                        }
                    },
                    "400": {
                        "description": "Bad Request message",
                        "schema": {
                            "$ref": "#/definitions/endpoints.BadRequestForm"
                        }
                    },
                    "503": {
                        "description": "Server does not available",
                        "schema": {
                            "$ref": "#/definitions/endpoints.ServerErrorForm"
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
                            "$ref": "#/definitions/endpoints.ResponseForm"
                        }
                    },
                    "400": {
                        "description": "Bad Request message",
                        "schema": {
                            "$ref": "#/definitions/endpoints.BadRequestForm"
                        }
                    },
                    "503": {
                        "description": "Server does not available",
                        "schema": {
                            "$ref": "#/definitions/endpoints.ServerErrorForm"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "endpoints.AnalyseFilesForm": {
            "type": "object",
            "properties": {
                "document_ids": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    },
                    "example": [
                        "886f7e11874040ca8b8461fb4cd1aa2c"
                    ]
                }
            }
        },
        "endpoints.BadRequestForm": {
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
        "endpoints.FolderNameForm": {
            "type": "object",
            "properties": {
                "folder_name": {
                    "type": "string",
                    "example": "test_folder"
                }
            }
        },
        "endpoints.MoveFilesForm": {
            "type": "object",
            "properties": {
                "document_paths": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    },
                    "example": [
                        "./indexer/upload/test.txt"
                    ]
                },
                "source_directory": {
                    "type": "string",
                    "example": "unrecognized"
                },
                "target_directory": {
                    "type": "string",
                    "example": "common_folder"
                }
            }
        },
        "endpoints.RemoveFilesError": {
            "type": "object",
            "properties": {
                "code": {
                    "type": "integer",
                    "example": 403
                },
                "file_paths": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    },
                    "example": [
                        "./indexer/upload/test.txt"
                    ]
                },
                "message": {
                    "type": "string",
                    "example": "File not found"
                }
            }
        },
        "endpoints.RemoveFilesForm": {
            "type": "object",
            "properties": {
                "document_paths": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    },
                    "example": [
                        "./indexer/upload/test.txt"
                    ]
                }
            }
        },
        "endpoints.ResponseForm": {
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
        "endpoints.ServerErrorForm": {
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
        "endpoints.UnrecognizedDocuments": {
            "type": "object",
            "properties": {
                "unrecognized": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/reader.DocumentPreview"
                    }
                }
            }
        },
        "endpoints.WatcherDirectoriesForm": {
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
        },
        "reader.Artifacts": {
            "type": "object",
            "properties": {
                "group_json_name": {
                    "type": "string"
                },
                "group_name": {
                    "type": "string"
                },
                "group_values": {
                    "type": "array",
                    "items": {
                        "type": "object",
                        "properties": {
                            "json_name": {
                                "type": "string"
                            },
                            "name": {
                                "type": "string"
                            },
                            "type": {
                                "type": "string"
                            },
                            "value": {
                                "type": "string"
                            }
                        }
                    }
                }
            }
        },
        "reader.DocumentPreview": {
            "type": "object",
            "properties": {
                "created_at": {
                    "type": "string",
                    "example": "2024-05-04T22:53:57Z"
                },
                "file_size": {
                    "type": "integer",
                    "example": 311652
                },
                "id": {
                    "type": "string",
                    "example": "886f7e11874040ca8b8461fb4cd1aa2c"
                },
                "location": {
                    "type": "string",
                    "example": "unrecognized"
                },
                "name": {
                    "type": "string",
                    "example": "document_name.pdf"
                },
                "preview_properties": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/reader.Artifacts"
                    }
                },
                "quality_recognition": {
                    "type": "integer",
                    "example": 10000
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
