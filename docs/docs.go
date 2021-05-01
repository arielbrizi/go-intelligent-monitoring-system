// Package docs GENERATED BY THE COMMAND ABOVE; DO NOT EDIT
// This file was generated by swaggo/swag
package docs

import (
	"bytes"
	"encoding/json"
	"strings"
	"text/template"

	"github.com/swaggo/swag"
)

var doc = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{.Description}}",
        "title": "{{.Title}}",
        "termsOfService": "http://swagger.io/terms/",
        "contact": {
            "name": "API Support",
            "url": "https://arielbrizi.github.io/go-intelligent-monitoring-system/",
            "email": "arielbrizi@gmail.com"
        },
        "license": {
            "name": "Apache 2.0",
            "url": "http://www.apache.org/licenses/LICENSE-2.0.html"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/configuration-core/authorized-face": {
            "post": {
                "description": "add the person in the image (parameter) as an authorized face",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "configuration-core"
                ],
                "summary": "add authorized face",
                "parameters": [
                    {
                        "description": "name, bucket and collection are required",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/domain.Image"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": ""
                    }
                }
            },
            "delete": {
                "description": "delete the person in the image (parameter) as an authorized face",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "configuration-core"
                ],
                "summary": "delete authorized face",
                "parameters": [
                    {
                        "description": "Authorized FaceId (id), and Collection ID (collection) are required",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/domain.AuthorizedFace"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": ""
                    }
                }
            }
        },
        "/configuration-core/authorized-face/{collectionName}": {
            "get": {
                "description": "get the authorized faces for collection Id",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "configuration-core"
                ],
                "summary": "get authorized faces",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Collection ID",
                        "name": "collectionName",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": ""
                    }
                }
            }
        }
    },
    "definitions": {
        "domain.AuthorizedFace": {
            "type": "object",
            "required": [
                "collection",
                "id"
            ],
            "properties": {
                "bucket": {
                    "type": "string"
                },
                "bytes": {
                    "type": "array",
                    "items": {
                        "type": "integer"
                    }
                },
                "collection": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                }
            }
        },
        "domain.Image": {
            "type": "object",
            "required": [
                "bucket",
                "collection",
                "name"
            ],
            "properties": {
                "bucket": {
                    "type": "string"
                },
                "bytes": {
                    "type": "array",
                    "items": {
                        "type": "integer"
                    }
                },
                "collection": {
                    "type": "string"
                },
                "day": {
                    "type": "string"
                },
                "hour": {
                    "description": "time",
                    "type": "string"
                },
                "month": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "url": {
                    "description": "To get the image file",
                    "type": "string"
                },
                "year": {
                    "type": "string"
                }
            }
        }
    }
}`

type swaggerInfo struct {
	Version     string
	Host        string
	BasePath    string
	Schemes     []string
	Title       string
	Description string
}

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = swaggerInfo{
	Version:     "1.0",
	Host:        "",
	BasePath:    "",
	Schemes:     []string{},
	Title:       "Intelligent Monitoring System",
	Description: "Configuration and Analize Services.",
}

type s struct{}

func (s *s) ReadDoc() string {
	sInfo := SwaggerInfo
	sInfo.Description = strings.Replace(sInfo.Description, "\n", "\\n", -1)

	t, err := template.New("swagger_info").Funcs(template.FuncMap{
		"marshal": func(v interface{}) string {
			a, _ := json.Marshal(v)
			return string(a)
		},
	}).Parse(doc)
	if err != nil {
		return doc
	}

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, sInfo); err != nil {
		return doc
	}

	return tpl.String()
}

func init() {
	swag.Register(swag.Name, &s{})
}
