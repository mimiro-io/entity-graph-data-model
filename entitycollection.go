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

// Set Continuation	Token
func (ec *EntityCollection) SetContinuationToken(continuation *Continuation) {
	ec.Continuation = continuation
}

// add entity to collection
func (ec *EntityCollection) AddEntity(entity *Entity) error {
	ec.Entities = append(ec.Entities, entity)
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
		err = ec.expandEntityNamespaces(entity, ec.NamespaceManager)
		if err != nil {
			return err
		}
	}

	return err
}

func (ec *EntityCollection) expandEntityNamespaces(entity *Entity, namespaceManager NamespaceManager) error {
	// expand id
	fullID, err := namespaceManager.GetFullURI(entity.ID)
	if err != nil {
		return err
	}
	entity.ID = fullID

	// expand property types
	for typeURI, propertyValue := range entity.Properties {
		fullType, err := namespaceManager.GetFullURI(typeURI)
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
		fullType, err := namespaceManager.GetFullURI(typeURI)
		if err != nil {
			return err
		}

		// get updated values
		values, err := ec.expandRefValues(refValues, namespaceManager)
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

func (ec *EntityCollection) expandRefValues(values any, namespaceManage NamespaceManager) (any, error) {
	// switch if string or []string
	switch values.(type) {
	case string:
		// expand ref value
		fullRefValue, err := namespaceManage.GetFullURI(values.(string))
		if err != nil {
			return nil, err
		}
		return fullRefValue, nil
	case []string:
		// expand ref values
		for i, refValue := range values.([]string) {
			fullRefValue, err := namespaceManage.GetFullURI(refValue)
			if err != nil {
				return nil, err
			}
			values.([]string)[i] = fullRefValue
		}
		return values, nil
	}

	return nil, errors.New("unexpected type in refs")
}
