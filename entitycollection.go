package egdm

// EntityCollection is a utility structure for collecting together a set of entities, namespace mappings and a continuation token
type EntityCollection struct {
	Entities         []*Entity
	Continuation     *Continuation
	NamespaceManager NamespaceManager
}

func NewEntityCollection(nsManager NamespaceManager) *EntityCollection {
	ec := &EntityCollection{}
	ec.NamespaceManager = nsManager
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
