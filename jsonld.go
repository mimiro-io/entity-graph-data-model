package egdm

import (
	"encoding/json"
	"io"
)

type JsonLdRef struct {
	ID string `json:"@id"`
}

type JsonLDWriter struct {
}

func (jsonLDWriter *JsonLDWriter) Write(ec *EntityCollection, writer io.Writer) error {
	var err error

	// write [
	_, err = writer.Write([]byte("[\n"))
	if err != nil {
		return err
	}

	// write context
	mappings := ec.NamespaceManager.GetNamespaceMappings()
	context := jsonLDWriter.makeContext(mappings)
	contextJson, _ := json.Marshal(context)
	_, err = writer.Write(contextJson)
	if err != nil {
		return err
	}

	// write entities
	for _, entity := range ec.Entities {
		_, err = writer.Write([]byte(",\n"))
		if err != nil {
			return err
		}
		entityLD := jsonLDWriter.toJSONLD(entity)
		entityJson, err := json.Marshal(entityLD)
		if err != nil {
			return err
		}
		_, err = writer.Write(entityJson)
		if err != nil {
			return err
		}
	}

	// write continuation if not nil
	if ec.Continuation != nil {
		_, err = writer.Write([]byte(",\n"))
		if err != nil {
			return err
		}
		contToken := jsonLDWriter.makeContinuationToken(ec.Continuation.Token)
		_, err = writer.Write([]byte(contToken))
		if err != nil {
			return err
		}
	}

	// write ]
	_, err = writer.Write([]byte("\n]"))
	if err != nil {
		return err
	}

	return nil
}

func (jsonLDWriter *JsonLDWriter) makeContext(namespaceMappings map[string]string) map[string]interface{} {
	namespaces := make(map[string]string, len(namespaceMappings)+2)
	for k, v := range namespaceMappings {
		namespaces[k] = v
	}
	namespaces["core"] = "http://data.mimiro.io/core/uda/"
	namespaces["rdf"] = "http://www.w3.org/1999/02/22-rdf-syntax-ns#"

	jsonLdContext := make(map[string]interface{})
	jsonLdContext["@context"] = namespaces
	return jsonLdContext
}

func (jsonLDWriter *JsonLDWriter) makeContinuationToken(token string) string {
	contToken := make(map[string]interface{})
	contToken["rdf:type"] = map[string]string{"@id": "core:continuation"}
	contToken["core:token"] = token
	jsonData, _ := json.Marshal(contToken)
	return string(jsonData)
}

// Entity to JSON-LD representation
func (jsonLDWriter *JsonLDWriter) toJSONLD(entity *Entity) map[string]interface{} {
	jsonLd := make(map[string]interface{})

	// get the id and add that if not empty
	if entity.ID != "" {
		jsonLd["@id"] = entity.ID
	}

	// get props
	for key, value := range entity.Properties {
		// check the type of value
		switch v := value.(type) {
		case []interface{}:
			// array of values
			jsonLd[key] = jsonLDWriter.toJSONLDFromArray(v)
		case *Entity:
			// entity
			jsonLd[key] = jsonLDWriter.toJSONLD(v)
		default:
			// assume we can just put out the value
			jsonLd[key] = v
		}
	}

	// get the refs
	for key, value := range entity.References {
		// check the type of value
		switch v := value.(type) {
		case []string:
			refs := make([]JsonLdRef, len(v))
			for _, ref := range v {
				refs = append(refs, JsonLdRef{ID: ref})
			}
			jsonLd[key] = refs
		case string:
			jsonLd[key] = JsonLdRef{ID: v}
		}
	}

	return jsonLd
}

func (jsonLDWriter *JsonLDWriter) toJSONLDFromArray(entityArray []interface{}) []interface{} {
	jsonLd := make([]interface{}, len(entityArray))

	for i, value := range entityArray {
		switch v := value.(type) {
		case []interface{}:
			jsonLd[i] = jsonLDWriter.toJSONLDFromArray(v)
		case *Entity:
			jsonLd[i] = jsonLDWriter.toJSONLD(v)
		default:
			jsonLd[i] = value
		}
	}

	return jsonLd
}
