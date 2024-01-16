package egdm

import (
	"encoding/json"
	"errors"
	"io"
)

// EntityCollection is a utility structure for collecting together a set of entities, namespace mappings and a continuation token
type EntityCollection struct {
	Entities         []*Entity
	Continuation     *Continuation
	NamespaceManager NamespaceManager
}

func NewEntityCollection(nsManager NamespaceManager) *EntityCollection {
	ec := &EntityCollection{}
	// default to inbuilt namespace manager if not defined
	if nsManager == nil {
		nsManager = NewNamespaceContext()
	}
	ec.NamespaceManager = nsManager
	ec.Entities = make([]*Entity, 0)
	return ec
}

// SetContinuationToken sets the continuation token on the EntityCollection
func (ec *EntityCollection) SetContinuationToken(continuation *Continuation) {
	ec.Continuation = continuation
}

// AddEntity adds the given entity to the collection
func (ec *EntityCollection) AddEntity(entity *Entity) error {
	ec.Entities = append(ec.Entities, entity)
	return nil
}

// AddEntityFromMap adds an entity to the collection from a map
// The map should have the following structure (the keys are case sensitive):
//
//	{
//	  "id": "ns0:entity1",
//	  "deleted": false,
//	  "recorded": 1234567890,
//	  "props": {
//	    "ns0:property1": "value1"
//	  },
//	  "refs": {
//	    "ns0:reference1": "ns0:entity2"
//	  }
func (ec *EntityCollection) AddEntityFromMap(data map[string]any) error {
	entity := NewEntity()

	// get metadata
	if id, found := data["id"]; found {
		entity.ID = id.(string)
	}

	if isDeleted, found := data["deleted"]; found {
		// check type it bool
		if _, ok := isDeleted.(bool); ok {
			entity.IsDeleted = isDeleted.(bool)
		}
	}

	if recorded, found := data["recorded"]; found {
		if _, ok := recorded.(float64); ok {
			entity.Recorded = uint64(recorded.(float64))
		}

		if _, ok := recorded.(uint64); ok {
			entity.Recorded = recorded.(uint64)
		}
	}

	// get props
	if props, found := data["props"]; found {
		for key, value := range props.(map[string]any) {
			entity.Properties[key] = value
		}
	}

	// get refs
	if refs, found := data["refs"]; found {
		for key, value := range refs.(map[string]any) {
			entity.References[key] = value
		}
	}

	// add entity to collection
	err := ec.AddEntity(entity)
	if err != nil {
		return err
	}

	return nil
}

func (ec *EntityCollection) GetEntities() []*Entity {
	return ec.Entities
}

func (ec *EntityCollection) GetContinuationToken() *Continuation {
	return ec.Continuation
}

func (ec *EntityCollection) GetNamespaceManager() NamespaceManager {
	return ec.NamespaceManager
}

func (ec *EntityCollection) GetNamespaceMappings() map[string]string {
	return ec.NamespaceManager.GetNamespaceMappings()
}

type Context struct {
	ID         string            `json:"id"`
	Namespaces map[string]string `json:"namespaces"`
}

func NewContext() *Context {
	c := &Context{}
	c.ID = "@context"
	c.Namespaces = make(map[string]string)
	return c
}

func (ec *EntityCollection) WriteEntityGraphJSON(writer io.Writer) error {
	var err error

	// write [
	_, err = writer.Write([]byte("[\n"))
	if err != nil {
		return err
	}

	// write context
	context := NewContext()
	context.ID = "@context"
	context.Namespaces = ec.NamespaceManager.GetNamespaceMappings()
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
		entityJson, err := json.Marshal(entity)
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
		contJson, err := json.Marshal(ec.Continuation)
		if err != nil {
			return err
		}
		_, err = writer.Write([]byte(", "))
		if err != nil {
			return err
		}
		_, err = writer.Write([]byte(contJson))
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

func (ec *EntityCollection) WriteJSON_LD(writer io.Writer) error {
	jsonLDWriter := &JsonLDWriter{}
	return jsonLDWriter.Write(ec, writer)
}

func (ec *EntityCollection) ExpandNamespacePrefixes() error {
	var err error
	for _, entity := range ec.Entities {
		err = ec.expandEntityNamespaces(entity)
		if err != nil {
			return err
		}
	}

	return err
}

func (ec *EntityCollection) expandEntityNamespaces(entity *Entity) error {
	// expand id
	fullID, err := ec.NamespaceManager.GetFullURI(entity.ID)
	if err != nil {
		return err
	}
	entity.ID = fullID

	// expand property types
	for typeURI, propertyValue := range entity.Properties {
		fullType, err := ec.NamespaceManager.GetFullURI(typeURI)
		if err != nil {
			return err
		}
		// remove old key
		delete(entity.Properties, typeURI)

		// add new key and value
		entity.Properties[fullType] = propertyValue
	}

	// expand ref types and values
	for typeURI, refValues := range entity.References {
		fullType, err := ec.NamespaceManager.GetFullURI(typeURI)
		if err != nil {
			return err
		}

		// get updated values
		values, err := ec.expandRefValues(refValues)
		if err != nil {
			return err
		}

		// remove old key
		delete(entity.References, typeURI)

		// add new key and value
		entity.References[fullType] = values
	}

	return nil
}

func (ec *EntityCollection) expandRefValues(values any) (any, error) {
	// switch if string or []string
	switch values.(type) {
	case string:
		// expand ref value
		fullRefValue, err := ec.NamespaceManager.GetFullURI(values.(string))
		if err != nil {
			return nil, err
		}
		return fullRefValue, nil
	case []string:
		// expand ref values
		for i, refValue := range values.([]string) {
			fullRefValue, err := ec.NamespaceManager.GetFullURI(refValue)
			if err != nil {
				return nil, err
			}
			values.([]string)[i] = fullRefValue
		}
		return values, nil
	}

	return nil, errors.New("unexpected type in refs")
}
