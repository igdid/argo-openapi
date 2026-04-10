package openapiutil

import (
	"fmt"

	"github.com/getkin/kin-openapi/openapi3"
)

type ParamInfo struct {
	Name     string
	Type     string
	Required bool
}

type OperationParams struct {
	Path  []ParamInfo
	Other []ParamInfo
}

// FindOperation ищет operation по operationID во всём документе.
func FindOperation(doc *openapi3.T, operationID string) (*openapi3.Operation, string, string, error) {
	if doc == nil {
		return nil, "", "", fmt.Errorf("nil document")
	}

	for path, pathItem := range doc.Paths.Map() {
		for method, op := range operationsFromPathItem(pathItem) {
			if op != nil && op.OperationID == operationID {
				return op, method, path, nil
			}
		}
	}

	return nil, "", "", fmt.Errorf("operation %q not found", operationID)
}

func ExtractOperationParams(op *openapi3.Operation) *OperationParams {
	result := &OperationParams{}

	if op == nil {
		return result
	}

	for _, paramRef := range op.Parameters {
		if paramRef == nil || paramRef.Value == nil {
			continue
		}

		param := paramRef.Value

		info := ParamInfo{
			Name:     param.Name,
			Type:     goTypeFromSchema(param.Schema),
			Required: param.Required,
		}

		if param.In == "path" {
			result.Path = append(result.Path, info)
		} else {
			result.Other = append(result.Other, info)
		}
	}

	return result
}

func operationsFromPathItem(pathItem *openapi3.PathItem) []*openapi3.Operation {
	if pathItem == nil {
		return nil
	}

	return []*openapi3.Operation{
		pathItem.Get,
		pathItem.Post,
		pathItem.Put,
		pathItem.Delete,
		pathItem.Patch,
		pathItem.Head,
		pathItem.Options,
		pathItem.Trace,
	}
}

func goTypeFromSchema(schemaRef *openapi3.SchemaRef) string {
	if schemaRef == nil || schemaRef.Value == nil {
		return "interface{}"
	}

	schema := schemaRef.Value

	switch schema.Type {
	case "string":
		switch schema.Format {
		case "date-time":
			return "time.Time"
		case "date":
			return "civil.Date" // либо string, если civil.Date не используется
		case "byte":
			return "[]byte"
		default:
			return "string"
		}

	case "integer":
		switch schema.Format {
		case "int32":
			return "int32"
		case "int64":
			return "int64"
		default:
			return "int"
		}

	case "number":
		switch schema.Format {
		case "float":
			return "float32"
		case "double":
			return "float64"
		default:
			return "float64"
		}

	case "boolean":
		return "bool"

	case "array":
		if schema.Items == nil {
			return "[]interface{}"
		}
		return "[]" + goTypeFromSchema(schema.Items)

	case "object":
		// Можно дополнительно учитывать AdditionalProperties / свойства.
		return "map[string]interface{}"

	default:
		return "interface{}"
	}
}

