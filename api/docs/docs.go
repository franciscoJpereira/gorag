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
        "/chat": {
            "get": {
                "description": "Returns the names of all existing chats",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "chat"
                ],
                "summary": "Retrieves all available Chats",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Chat ID",
                        "name": "chatID",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/store.ChatHistory"
                        }
                    }
                }
            },
            "post": {
                "description": "Sends a message to a chat and creates it if needed",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "chat"
                ],
                "summary": "Send a new message to a chat",
                "parameters": [
                    {
                        "description": "Message to send to Chat",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/pkg.ChatInstruct"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/pkg.MessageResponse"
                        }
                    },
                    "400": {
                        "description": "Error sending message",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/knowledge-base": {
            "get": {
                "description": "Returns a list with the names of all available knwoledge bases",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "knowledge-base"
                ],
                "summary": "Get available Knowledge Bases",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "type": "string"
                            }
                        }
                    }
                }
            },
            "post": {
                "description": "Add string data to a knwoledge base, it creates the KB if the flag is set",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "knowledge-base"
                ],
                "summary": "Add data to a knowledge base",
                "parameters": [
                    {
                        "description": "Data to add to the KB",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/pkg.KBAddDataInstruct"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK"
                    },
                    "400": {
                        "description": "KB does not exist",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/knowledge-base/{KBName}": {
            "post": {
                "description": "Create a new knowledge base",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "knowledge-base"
                ],
                "summary": "Create a new knowledge base",
                "parameters": [
                    {
                        "type": "string",
                        "description": "KBName",
                        "name": "KBName",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK"
                    },
                    "400": {
                        "description": "KB already exists",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/message": {
            "post": {
                "description": "Send a one-shot message to get a response",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "chat"
                ],
                "summary": "Send a one-shot message",
                "parameters": [
                    {
                        "description": "Message to send",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/pkg.MessageInstruct"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/pkg.MessageResponse"
                        }
                    },
                    "400": {
                        "description": "Error sending message",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "apiinterface.ChatMessage": {
            "type": "object",
            "properties": {
                "content": {
                    "type": "string"
                },
                "role": {
                    "type": "string"
                }
            }
        },
        "pkg.ChatInstruct": {
            "type": "object",
            "properties": {
                "ChatName": {
                    "type": "string"
                },
                "Message": {
                    "$ref": "#/definitions/pkg.MessageInstruct"
                },
                "NewChat": {
                    "type": "boolean"
                }
            }
        },
        "pkg.KBAddDataInstruct": {
            "type": "object",
            "properties": {
                "Create": {
                    "type": "boolean"
                },
                "KBName": {
                    "type": "string"
                },
                "data": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                }
            }
        },
        "pkg.MessageInstruct": {
            "type": "object",
            "properties": {
                "KB": {
                    "type": "boolean"
                },
                "KBName": {
                    "type": "string"
                },
                "Message": {
                    "type": "string"
                }
            }
        },
        "pkg.MessageResponse": {
            "type": "object",
            "properties": {
                "Context": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "Query": {
                    "type": "string"
                },
                "Response": {
                    "type": "string"
                }
            }
        },
        "store.ChatHistory": {
            "type": "object",
            "properties": {
                "history": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/apiinterface.ChatMessage"
                    }
                },
                "name": {
                    "type": "string"
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
