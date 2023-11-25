package egdm

// basic go test for the parser using standard test setup
import (
	"testing"
)

func TestAddEntityToEntityCollection(t *testing.T) {
	ec := NewEntityCollection(nil)
	entity := NewEntity().SetID("ns0:entity1")
	err := ec.AddEntity(entity)
	if err != nil {
		t.Error(err)
	}
	if len(ec.Entities) != 1 {
		t.Errorf("expected entity collection to have 1 entity, got %d", len(ec.Entities))
	}
}

func TestCreateEntity(t *testing.T) {
	// create a new entity
	entity := NewEntity().SetID("ns0:entity1")
	if entity.ID != "ns0:entity1" {
		t.Errorf("expected entity id to be 'ns0:entity1', got '%s'", entity.ID)
	}

	// add a property
	entity.SetProperty("ns0:property1", "value1")
	if entity.Properties["ns0:property1"] != "value1" {
		t.Errorf("expected entity property to be 'value1', got '%s'", entity.Properties["ns0:property1"])
	}

	// use get property function
	if value, err := entity.GetFirstStringPropertyValue("ns0:property1"); err != nil {
		t.Errorf("expected entity property to be 'value1', got error '%s'", err)
	} else {
		if value != "value1" {
			t.Errorf("expected entity property to be 'value1', got '%s'", value)
		}
	}

	// add a reference
	entity.SetReference("ns0:reference1", "ns0:entity2")
	if entity.References["ns0:reference1"] != "ns0:entity2" {
		t.Errorf("expected entity reference to be 'ns0:entity2', got '%s'", entity.References["ns0:reference1"])
	}
}

func TestAssertIdentifierReturnsErrorWhenMissingPostfix(t *testing.T) {
	// namespace manager
	nsManager := NewNamespaceContext()
	_, err := nsManager.AssertPrefixedIdentifierFromURI("http://data.example.com/things/")
	if err == nil {
		t.Error(err)
	}
}

func TestCreateEntityUsingContextManager(t *testing.T) {
	// namespace manager
	nsManager := NewNamespaceContext()
	entityId, err := nsManager.AssertPrefixedIdentifierFromURI("http://data.example.com/things/entity1")
	if err != nil {
		t.Error(err)
	}

	entity := NewEntity().SetID(entityId)
	if entity.ID != "ns0:entity1" {
		t.Errorf("expected entity id to be 'ns0:entity1', got '%s'", entity.ID)
	}

	// add a property
	entity.SetProperty("ns0:property1", "value1")
	if entity.Properties["ns0:property1"] != "value1" {
		t.Errorf("expected entity property to be 'value1', got '%s'", entity.Properties["ns0:property1"])
	}

	// use get property function
	if value, err := entity.GetFirstStringPropertyValue("ns0:property1"); err != nil {
		t.Errorf("expected entity property to be 'value1', got error '%s'", err)
	} else {
		if value != "value1" {
			t.Errorf("expected entity property to be 'value1', got '%s'", value)
		}
	}

	// add a reference
	entity.SetReference("ns0:reference1", "ns0:entity2")
	if entity.References["ns0:reference1"] != "ns0:entity2" {
		t.Errorf("expected entity reference to be 'ns0:entity2', got '%s'", entity.References["ns0:reference1"])
	}
}
