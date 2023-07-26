package egdm

func toJsonLdFromMap(entityMap map[string]interface{}) map[string]interface{} {
	jsonLd := make(map[string]interface{})

	// add id
	if entityMap["id"] != nil {
		jsonLd["@id"] = entityMap["id"]
	}

	if entityMap["props"] != nil {
		for key, value := range entityMap["props"].(map[string]interface{}) {
			// check the type of value
			switch v := value.(type) {
			case []interface{}:
				// array of entities
				jsonLd[key] = toJsonLdFromArray(v)
			case map[string]interface{}:
				// entity as json
				jsonLd[key] = toJsonLdFromMap(v)
			default:
				// assume we can just put out the value
				jsonLd[key] = v
			}
		}
	}

	// if references
	if entityMap["refs"] != nil {
		for key, value := range entityMap["refs"].(map[string]interface{}) {
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
	}

	return jsonLd
}

func toJsonLdFromArray(entityArray []interface{}) []interface{} {
	jsonLd := make([]interface{}, len(entityArray))

	for i, value := range entityArray {
		switch value.(type) {
		case []interface{}:
			jsonLd[i] = toJsonLdFromArray(value.([]interface{}))
		case map[string]interface{}:
			jsonLd[i] = toJsonLdFromMap(value.(map[string]interface{}))
		default:
			jsonLd[i] = value
		}
	}

	return jsonLd
}

// Convert Entity JSON-LD representation
func toJSONLD(entity *Entity) map[string]interface{} {
	jsonLd := make(map[string]interface{})

	// get the id and add that
	jsonLd["@id"] = entity.ID

	// get props
	for key, value := range entity.Properties {
		// check the type of value
		switch v := value.(type) {
		case []interface{}:
			// array of values
			jsonLd[key] = toJsonLdFromArray(v)
		case map[string]interface{}:
			// entity as json
			jsonLd[key] = toJsonLdFromMap(v)
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

type JsonLdRef struct {
	ID string `json:"@id"`
}
