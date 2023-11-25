package egdm

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"
)

type Parser interface {
	Parse(data io.Reader, entity func(*Entity) error, continuation func(*Continuation)) error
	LoadEntityCollection(reader io.Reader) (*EntityCollection, error)
	GetNamespaceManager() NamespaceManager
}

type EntityParser struct {
	nsManager             NamespaceManager
	expandURIs            bool
	compressURIs          bool
	requireContext        bool
	contextParsedCallback func(*Context)
}

func NewEntityParser(nsmanager NamespaceManager) *EntityParser {
	ep := &EntityParser{}
	ep.nsManager = nsmanager
	ep.expandURIs = false
	ep.compressURIs = false
	ep.requireContext = true
	return ep
}

func (esp *EntityParser) WithNoContext() *EntityParser {
	esp.requireContext = false
	return esp
}

func (esp *EntityParser) WithExpandURIs() *EntityParser {
	esp.expandURIs = true
	return esp
}

func (esp *EntityParser) WithCompressURIs() *EntityParser {
	esp.compressURIs = true
	return esp
}

func (esp *EntityParser) WithParsedContextCallback(callback func(context *Context)) *EntityParser {
	esp.contextParsedCallback = callback
	return esp
}

func (esp *EntityParser) GetNamespaceManager() NamespaceManager {
	return esp.nsManager
}

func (esp *EntityParser) LoadEntityCollection(reader io.Reader) (*EntityCollection, error) {
	ec := NewEntityCollection(esp.nsManager)
	err := esp.Parse(reader, func(e *Entity) error {
		return ec.AddEntity(e)
	}, func(c *Continuation) {
		ec.SetContinuationToken(c)
	})
	if err != nil {
		return nil, err
	}
	return ec, nil
}

func (esp *EntityParser) GetIdentityValue(value string) (string, error) {

	identity := value
	var err error

	if esp.compressURIs {
		identity, err = esp.nsManager.AssertPrefixedIdentifierFromURI(value)
		if err != nil {
			return "", err
		}
	}

	if esp.expandURIs {
		return esp.nsManager.GetFullURI(identity)
	} else {
		return esp.nsManager.GetPrefixedIdentifier(identity)
	}
}

func (esp *EntityParser) Parse(reader io.Reader, emitEntity func(*Entity) error, emitContinuation func(*Continuation)) error {

	decoder := json.NewDecoder(reader)

	// expect start of array
	t, err := decoder.Token()
	if err != nil {
		return fmt.Errorf("parsing error: Bad token at start of stream: %w", err)
	}

	if delim, ok := t.(json.Delim); !ok || delim != '[' {
		return errors.New("parsing error: Expected [ at start of document")
	}

	// decode context object
	if esp.requireContext {
		if esp.nsManager == nil {
			return errors.New("parsing error: Namespace manager required when parsing with context")
		}
		context := make(map[string]any)
		err = decoder.Decode(&context)
		if err != nil {
			return fmt.Errorf("parsing error: Unable to decode context: %w", err)
		}

		if context["id"] == "@context" {
			if context["namespaces"] != nil {
				for k, v := range context["namespaces"].(map[string]any) {
					expansion := v.(string)
					if strings.HasSuffix(expansion, "/") || strings.HasSuffix(expansion, "#") {
						esp.nsManager.StorePrefixExpansionMapping(k, v.(string))
					} else {
						return fmt.Errorf("expansion %s for prefix %s must end with / or #", expansion, k)
					}
				}
			}
		} else {
			return errors.New("first object in array must be a context with id @context")
		}

		// if a callback func for the parsed context is registered we call it
		if esp.contextParsedCallback != nil {
			esp.contextParsedCallback(esp.nsManager.AsContext())
		}
	}

	for {
		t, err = decoder.Token()
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return fmt.Errorf("parsing error: Unable to read next token: %w", err)
			}
		}

		switch v := t.(type) {
		case json.Delim:
			if v == '{' {
				e, err := esp.parseEntity(decoder)
				if err != nil {
					return fmt.Errorf("parsing error: Unable to parse entity: %w", err)
				}
				if e.ID == "@continuation" {
					if emitContinuation != nil {
						continuation := NewContinuation()
						continuation.Token = e.Properties["token"].(string)
						emitContinuation(continuation)
					}
				} else {
					err = emitEntity(e)
					if err != nil {
						return err
					}
				}
			} else if v == ']' {
				// done
				break
			}
		default:
			return errors.New("parsing error: unexpected value in entity array")
		}
	}

	return nil
}

func (esp *EntityParser) parseEntity(decoder *json.Decoder) (*Entity, error) {
	e := &Entity{}
	e.Properties = make(map[string]any)
	e.References = make(map[string]any)
	isContinuation := false
	for {
		t, err := decoder.Token()
		if err != nil {
			return nil, fmt.Errorf("unable to read token: %w", err)
		}

		switch v := t.(type) {
		case json.Delim:
			if v == '}' {
				return e, nil
			}
		case string:
			if v == "id" {
				val, err := decoder.Token()
				if err != nil {
					return nil, fmt.Errorf("unable to read token of id value: %w", err)
				}

				if val.(string) == "@continuation" {
					e.ID = "@continuation"
					isContinuation = true
				} else if val.(string) == "@context" {
					return nil, errors.New("context object found when entity expected")
				} else {
					id, err := esp.GetIdentityValue(val.(string))
					if err != nil {
						return nil, err
					}
					e.ID = id
				}
			} else if v == "recorded" {
				val, err := decoder.Token()
				if err != nil {
					return nil, fmt.Errorf("unable to read token of recorded value: %w", err)
				}
				e.Recorded = uint64(val.(float64))
			} else if v == "deleted" {
				val, err := decoder.Token()
				if err != nil {
					return nil, fmt.Errorf("unable to read token of deleted value: %w", err)
				}
				e.IsDeleted = val.(bool)

			} else if v == "props" {
				e.Properties, err = esp.parseProperties(decoder)
				if err != nil {
					return nil, fmt.Errorf("unable to parse properties: %w", err)
				}
			} else if v == "refs" {
				e.References, err = esp.parseReferences(decoder)
				if err != nil {
					return nil, fmt.Errorf("unable to parse references %w", err)
				}
			} else if v == "token" {
				if !isContinuation {
					return nil, errors.New("token property found but not a continuation entity")
				}
				val, err := decoder.Token()
				if err != nil {
					return nil, fmt.Errorf("unable to read continuation token value: %w", err)
				}
				e.Properties = make(map[string]any)
				e.Properties["token"] = val
			} else {
				// log named property
				// read value
				_, err := decoder.Token()
				if err != nil {
					return nil, fmt.Errorf("unable to parse value of unknown key: %s %w", v, err)
				}
			}
		default:
			return nil, errors.New("unexpected value in entity")
		}
	}
}

func (esp *EntityParser) parseReferences(decoder *json.Decoder) (map[string]any, error) {
	refs := make(map[string]any)

	_, err := decoder.Token()
	if err != nil {
		return nil, fmt.Errorf("unable to read token of at start of references: %w", err)
	}

	for {
		t, err := decoder.Token()
		if err != nil {
			return nil, fmt.Errorf("unable to read token in parse references: %w", err)
		}

		switch v := t.(type) {
		case json.Delim:
			if v == '}' {
				return refs, nil
			}
		case string:
			val, err := esp.parseRefValue(decoder)
			if err != nil {
				return nil, fmt.Errorf("unable to parse value of reference key %s", v)
			}

			id, err := esp.GetIdentityValue(v)
			if err != nil {
				return nil, err
			}
			refs[id] = val
		default:
			return nil, errors.New("unknown type")
		}
	}
}

func (esp *EntityParser) parseProperties(decoder *json.Decoder) (map[string]any, error) {
	props := make(map[string]any)

	_, err := decoder.Token()
	if err != nil {
		return nil, fmt.Errorf("unable to read token of at start of properties: %w ", err)
	}

	for {
		t, err := decoder.Token()
		if err != nil {
			return nil, fmt.Errorf("unable to read token in parse properties: %w", err)
		}

		switch v := t.(type) {
		case json.Delim:
			if v == '}' {
				return props, nil
			}
		case string:
			val, err := esp.parseValue(decoder)
			if err != nil {
				return nil, fmt.Errorf("unable to parse property value of key %s err: %w", v, err)
			}

			if val != nil {
				id, err := esp.GetIdentityValue(v)
				if err != nil {
					return nil, err
				}
				props[id] = val
			}
		default:
			return nil, errors.New("unknown type")
		}
	}
}

func (esp *EntityParser) parseRefValue(decoder *json.Decoder) (any, error) {
	for {
		t, err := decoder.Token()
		if err != nil {
			return nil, fmt.Errorf("unable to read token in parse value: %w", err)
		}

		switch v := t.(type) {
		case json.Delim:
			if v == '[' {
				return esp.parseRefArray(decoder)
			}
		case string:
			id, err := esp.GetIdentityValue(v)
			if err != nil {
				return nil, err
			}
			return id, nil
		default:
			return nil, errors.New("unknown token in parse ref value")
		}
	}
}

func (esp *EntityParser) parseRefArray(decoder *json.Decoder) ([]string, error) {
	array := make([]string, 0)
	for {
		t, err := decoder.Token()
		if err != nil {
			return nil, fmt.Errorf("unable to read token in parse ref array: %w", err)
		}

		switch v := t.(type) {
		case json.Delim:
			if v == ']' {
				return array, nil
			}
		case string:
			id, err := esp.GetIdentityValue(v)
			if err != nil {
				return nil, err
			}
			array = append(array, id)
		default:
			return nil, errors.New("unknown type")
		}
	}
}

func (esp *EntityParser) parseArray(decoder *json.Decoder) ([]any, error) {
	array := make([]any, 0)
	for {
		t, err := decoder.Token()
		if err != nil {
			return nil, fmt.Errorf("unable to read token in parse array: %w", err)
		}

		switch v := t.(type) {
		case json.Delim:
			if v == '{' {
				r, err := esp.parseEntity(decoder)
				if err != nil {
					return nil, fmt.Errorf("unable to parse array: %w", err)
				}
				array = append(array, r)
			} else if v == ']' {
				return array, nil
			} else if v == '[' {
				r, err := esp.parseArray(decoder)
				if err != nil {
					return nil, fmt.Errorf("unable to parse array: %w", err)
				}
				array = append(array, r)
			}
		case string:
			array = append(array, v)
		case int:
			array = append(array, v)
		case float64:
			array = append(array, v)
		case bool:
			array = append(array, v)
		default:
			return nil, errors.New("unknown type")
		}
	}
}

func (esp *EntityParser) parseValue(decoder *json.Decoder) (any, error) {
	for {
		t, err := decoder.Token()
		if err != nil {
			return nil, fmt.Errorf("unable to read token in parse value: %w", err)
		}

		if t == nil {
			// there is a good chance that we got a null value, and we need to handle that
			return nil, nil
		}

		switch v := t.(type) {
		case json.Delim:
			if v == '{' {
				return esp.parseEntity(decoder)
			} else if v == '[' {
				return esp.parseArray(decoder)
			}
		case string:
			return v, nil
		case int:
			return v, nil
		case float64:
			return v, nil
		case bool:
			return v, nil
		default:
			return nil, errors.New("unknown token in parse value")
		}
	}
}
