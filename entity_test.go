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

func TestExpandPrefixes(t *testing.T) {
	// namespace manager
	nsManager := NewNamespaceContext()
	nsManager.StorePrefixExpansionMapping("ns0", "http://data.example.com/things/")
	// create entity using short form
	entity := NewEntity().SetID("ns0:entity1")

	// add some properties and references
	entity.SetProperty("ns0:property1", "value1")
	entity.SetReference("ns0:reference1", "ns0:entity2")
	// and some refs in a list
	entity.SetReference("ns0:reference2", []string{"ns0:entity3", "ns0:entity4"})

	// create entity collection and add entity
	ec := NewEntityCollection(nsManager)
	err := ec.AddEntity(entity)
	if err != nil {
		t.Error(err)
	}

	// expand prefixes
	err = ec.ExpandNamespacePrefixes()
	if err != nil {
		t.Error(err)
	}

	// check that the entity id has been expanded
	if entity.ID != "http://data.example.com/things/entity1" {
		t.Errorf("expected entity id to be 'http://data.example.com/things/entity1', got '%s'", entity.ID)
	}

	// check that the property has been expanded
	if entity.Properties["http://data.example.com/things/property1"] != "value1" {
		t.Errorf("expected entity property to be 'value1', got '%s'", entity.Properties["http://data.example.com/things/property1"])
	}

	// check that the reference has been expanded
	if entity.References["http://data.example.com/things/reference1"] != "http://data.example.com/things/entity2" {
		t.Errorf("expected entity reference to be 'http://data.example.com/things/entity2', got '%s'", entity.References["http://data.example.com/things/reference1"])
	}

	// check that the reference list has been expanded
	if entity.References["http://data.example.com/things/reference2"].([]string)[0] != "http://data.example.com/things/entity3" {
		t.Errorf("expected entity reference to be 'http://data.example.com/things/entity3', got '%s'", entity.References["http://data.example.com/things/reference2"].([]string)[0])
	}
}

func TestExpandPrefixesWithMissingExpansion(t *testing.T) {
	// namespace manager
	nsManager := NewNamespaceContext()
	nsManager.StorePrefixExpansionMapping("ns0", "http://data.example.com/things/")
	// create entity using short form
	entity := NewEntity().SetID("ns0:entity1")

	// add some properties and references
	entity.SetProperty("ns0:property1", "value1")
	entity.SetReference("ns0:reference1", "ns0:entity2")
	// and some refs in a list
	entity.SetReference("ns0:reference2", []string{"ns0:entity3", "ns0:entity4"})
	// add a reference with a missing expansion
	entity.SetReference("ns0:reference3", "ns1:entity5")

	// create entity collection and add entity
	ec := NewEntityCollection(nsManager)
	err := ec.AddEntity(entity)
	if err != nil {
		t.Error(err)
	}

	// expand prefixes
	err = ec.ExpandNamespacePrefixes()

	// expecting an error
	if err == nil {
		t.Error("expected error")
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
