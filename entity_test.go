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
	entity.SetReference("ns0:reference2", []any{"ns0:entity3", "ns0:entity4"})

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
	if entity.References["http://data.example.com/things/reference2"].([]any)[0] != "http://data.example.com/things/entity3" {
		t.Errorf("expected entity reference to be 'http://data.example.com/things/entity3', got '%s'", entity.References["http://data.example.com/things/reference2"].([]string)[0])
	}
}

// expandEntityNamespaces: []any of scalar values (e.g. string URIs from JSON unmarshal)
func TestExpandPrefixesWithAnySliceOfScalars(t *testing.T) {
	nsManager := NewNamespaceContext()
	nsManager.StorePrefixExpansionMapping("ns0", "http://data.example.com/things/")
	entity := NewEntity().SetID("ns0:entity1")
	entity.SetProperty("ns0:tags", []any{"value1", "value2", "value3"})

	ec := NewEntityCollection(nsManager)
	if err := ec.AddEntity(entity); err != nil {
		t.Fatal(err)
	}
	if err := ec.ExpandNamespacePrefixes(); err != nil {
		t.Fatal(err)
	}

	vals, ok := entity.Properties["http://data.example.com/things/tags"].([]any)
	if !ok || len(vals) != 3 {
		t.Fatalf("expected []any with 3 elements, got %T %v", entity.Properties["http://data.example.com/things/tags"], entity.Properties["http://data.example.com/things/tags"])
	}
	if vals[0] != "value1" || vals[1] != "value2" || vals[2] != "value3" {
		t.Errorf("unexpected values: %v", vals)
	}
}

// expandEntityNamespaces: []any of map[string]any (sub-entities from JSON unmarshal)
func TestExpandPrefixesWithAnySliceOfSubEntityMaps(t *testing.T) {
	nsManager := NewNamespaceContext()
	nsManager.StorePrefixExpansionMapping("ns0", "http://data.example.com/things/")
	entity := NewEntity().SetID("ns0:entity1")

	subEntities := []any{
		map[string]any{"id": "ns0:sub1", "props": map[string]any{"ns0:subprop": "val1"}},
		map[string]any{"id": "ns0:sub2", "props": map[string]any{"ns0:subprop": "val2"}},
	}
	entity.SetProperty("ns0:children", subEntities)

	ec := NewEntityCollection(nsManager)
	if err := ec.AddEntity(entity); err != nil {
		t.Fatal(err)
	}
	if err := ec.ExpandNamespacePrefixes(); err != nil {
		t.Fatal(err)
	}

	expanded, ok := entity.Properties["http://data.example.com/things/children"].([]*Entity)
	if !ok || len(expanded) != 2 {
		t.Fatalf("expected []*Entity with 2 elements, got %T", entity.Properties["http://data.example.com/things/children"])
	}
	if expanded[0].ID != "http://data.example.com/things/sub1" {
		t.Errorf("expected sub1 id to be expanded, got %s", expanded[0].ID)
	}
	if expanded[1].Properties["http://data.example.com/things/subprop"] != "val2" {
		t.Errorf("expected sub2 subprop to be expanded, got %v", expanded[1].Properties)
	}
}

// expandEntityNamespaces: map[string]any sub-entity (single, from JSON unmarshal)
func TestExpandPrefixesWithMapSubEntity(t *testing.T) {
	nsManager := NewNamespaceContext()
	nsManager.StorePrefixExpansionMapping("ns0", "http://data.example.com/things/")
	entity := NewEntity().SetID("ns0:entity1")
	entity.SetProperty("ns0:child", map[string]any{
		"id":    "ns0:sub1",
		"props": map[string]any{"ns0:subprop": "val1"},
	})

	ec := NewEntityCollection(nsManager)
	if err := ec.AddEntity(entity); err != nil {
		t.Fatal(err)
	}
	if err := ec.ExpandNamespacePrefixes(); err != nil {
		t.Fatal(err)
	}

	sub, ok := entity.Properties["http://data.example.com/things/child"].(*Entity)
	if !ok {
		t.Fatalf("expected *Entity, got %T", entity.Properties["http://data.example.com/things/child"])
	}
	if sub.ID != "http://data.example.com/things/sub1" {
		t.Errorf("expected sub id to be expanded, got %s", sub.ID)
	}
	if sub.Properties["http://data.example.com/things/subprop"] != "val1" {
		t.Errorf("expected subprop to be expanded, got %v", sub.Properties)
	}
}

// expandRefValues: []interface{} refs (from JSON unmarshal)
func TestExpandPrefixesWithInterfaceSliceRefs(t *testing.T) {
	nsManager := NewNamespaceContext()
	nsManager.StorePrefixExpansionMapping("ns0", "http://data.example.com/things/")
	entity := NewEntity().SetID("ns0:entity1")
	entity.SetReference("ns0:links", []interface{}{"ns0:entity2", "ns0:entity3"})

	ec := NewEntityCollection(nsManager)
	if err := ec.AddEntity(entity); err != nil {
		t.Fatal(err)
	}
	if err := ec.ExpandNamespacePrefixes(); err != nil {
		t.Fatal(err)
	}

	refs, ok := entity.References["http://data.example.com/things/links"].([]interface{})
	if !ok || len(refs) != 2 {
		t.Fatalf("expected []interface{} with 2 elements, got %T", entity.References["http://data.example.com/things/links"])
	}
	if refs[0] != "http://data.example.com/things/entity2" {
		t.Errorf("expected entity2 ref to be expanded, got %s", refs[0])
	}
	if refs[1] != "http://data.example.com/things/entity3" {
		t.Errorf("expected entity3 ref to be expanded, got %s", refs[1])
	}
}

// expandEntityNamespaces: scalar default (non-string, non-entity property value)
func TestExpandPrefixesWithScalarProperty(t *testing.T) {
	nsManager := NewNamespaceContext()
	nsManager.StorePrefixExpansionMapping("ns0", "http://data.example.com/things/")
	entity := NewEntity().SetID("ns0:entity1")
	entity.SetProperty("ns0:count", 42)
	entity.SetProperty("ns0:flag", false)

	ec := NewEntityCollection(nsManager)
	if err := ec.AddEntity(entity); err != nil {
		t.Fatal(err)
	}
	if err := ec.ExpandNamespacePrefixes(); err != nil {
		t.Fatal(err)
	}

	if entity.Properties["http://data.example.com/things/count"] != 42 {
		t.Errorf("expected count=42, got %v", entity.Properties["http://data.example.com/things/count"])
	}
	if entity.Properties["http://data.example.com/things/flag"] != false {
		t.Errorf("expected flag=false, got %v", entity.Properties["http://data.example.com/things/flag"])
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

func TestExpandPrefixesWithValueArray(t *testing.T) {
	// namespace manager
	nsManager := NewNamespaceContext()
	nsManager.StorePrefixExpansionMapping("ns0", "http://data.example.com/things/")
	// create entity using short form
	entity := NewEntity().SetID("ns0:entity1")

	// add some properties and references
	propArray := []any{"value1"}
	entity.SetProperty("ns0:property1", propArray)

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

	// check that the property has been expanded
	if entity.Properties["http://data.example.com/things/property1"].([]any)[0] != propArray[0] {
		t.Errorf("expected entity property to be array, got '%s'", entity.Properties["http://data.example.com/things/property1"])
	}

}

func TestExpandPrefixesWithSubEntity(t *testing.T) {
	// namespace manager
	nsManager := NewNamespaceContext()
	nsManager.StorePrefixExpansionMapping("ns0", "http://data.example.com/things/")
	// create entity using short form
	entity := NewEntity().SetID("ns0:entity1")

	// add some properties and references
	subEntity := NewEntity().SetID("ns0:entity2")
	subEntity.Properties["ns0:subproperty1"] = "value2"

	entity.SetProperty("ns0:subEntity", subEntity)

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

	// check that the property has been expanded
	val, exist := entity.Properties["http://data.example.com/things/subEntity"]
	if exist {
		sub := val.(*Entity)
		if sub.Properties["http://data.example.com/things/subproperty1"] != "value2" {
			t.Errorf("expected sub entity property to be 'value2', got '%s'", sub.Properties["http://data.example.com/things/subproperty1"])
		}
	} else {
		t.Error("expected resolved sub entity property to exist")
	}
}

func TestExpandPrefixesWithSubEntityArray(t *testing.T) {
	// namespace manager
	nsManager := NewNamespaceContext()
	nsManager.StorePrefixExpansionMapping("ns0", "http://data.example.com/things/")
	// create entity using short form
	entity := NewEntity().SetID("ns0:entity1")

	// add some properties and references
	var subEntities []*Entity
	subEntity := NewEntity().SetID("ns0:entity2")
	subEntity.Properties["ns0:subproperty1"] = "value2"
	subEntities = append(subEntities, subEntity)

	entity.SetProperty("ns0:subEntities", subEntities)

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

	// check that the property has been expanded
	val, exist := entity.Properties["http://data.example.com/things/subEntities"]
	if exist {
		sub := val.([]*Entity)[0]
		if sub.Properties["http://data.example.com/things/subproperty1"] != "value2" {
			t.Errorf("expected sub entity property to be 'value2', got '%s'", sub.Properties["http://data.example.com/things/subproperty1"])
		}
	} else {
		t.Error("expected resolved sub entity property to exist")
	}
}

func TestExpandPrefixesWithSubEntityAsMap(t *testing.T) {
	// namespace manager
	nsManager := NewNamespaceContext()
	nsManager.StorePrefixExpansionMapping("ns0", "http://data.example.com/things/")
	// create entity using short form
	entity := NewEntity().SetID("ns0:entity1")

	// add some properties and references
	subEntity := map[string]any{"id": "ns0:entity2", "props": map[string]any{"ns0:subproperty1": "value2"}}

	entity.SetProperty("ns0:subEntity", subEntity)

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

	// check that the property has been expanded
	val, exist := entity.Properties["http://data.example.com/things/subEntity"]
	if exist {
		sub := val.(*Entity)
		if sub.Properties["http://data.example.com/things/subproperty1"] != "value2" {
			t.Errorf("expected sub entity property to be 'value2', got '%s'", sub.Properties["http://data.example.com/things/subproperty1"])
		}
	} else {
		t.Error("expected resolved sub entity property to exist")
	}
}

func TestExpandPrefixesWithSubEntityAsMapArray(t *testing.T) {
	// namespace manager
	nsManager := NewNamespaceContext()
	nsManager.StorePrefixExpansionMapping("ns0", "http://data.example.com/things/")
	// create entity using short form
	entity := NewEntity().SetID("ns0:entity1")

	// add some properties and references
	var subEntities []any
	subEntity := map[string]any{"id": "ns0:entity2", "props": map[string]any{"ns0:subproperty1": "value2"}}
	subEntities = append(subEntities, subEntity)

	entity.SetProperty("ns0:subEntities", subEntities)

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

	// check that the property has been expanded
	val, exist := entity.Properties["http://data.example.com/things/subEntities"]
	if exist {
		sub := val.([]*Entity)[0]
		if sub.Properties["http://data.example.com/things/subproperty1"] != "value2" {
			t.Errorf("expected sub entity property to be 'value2', got '%s'", sub.Properties["http://data.example.com/things/subproperty1"])
		}
	} else {
		t.Error("expected resolved sub entity property to exist")
	}
}

func TestGetStringPropertyValuesFromParsedJSON(t *testing.T) {
	// This is the real-world case that exposed the bug: JSON arrays unmarshal as []interface{},
	// not []string, so GetStringPropertyValues must handle []any.
	data := map[string]any{
		"id":      "http://data.mimiro.io/e360/organizations/018eccea-0040-78ae-b066-e8bd8dfd9b2d",
		"deleted": false,
		"refs": map[string]any{
			"http://www.w3.org/1999/02/22-rdf-syntax-ns/type": "http://data.mimiro.io/e360/Organization",
		},
		"props": map[string]any{
			"http://data.mimiro.io/e360/props.name": "Test Johan",
			"http://data.mimiro.io/e360/props.externalIdentifiers.val": []any{
				"http://data.mimiro.io/prodreg/producer-code/3425011396",
				"http://data.mimiro.io/prodreg/farm-code/34250113",
				"http://data.mimiro.io/prodreg/herdidentifier/2504433",
				"http://data.mimiro.io/prodreg/pid/10003427241",
				"http://data.mimiro.io/bronnoysund/organization-number/992666862",
			},
		},
	}

	ec := NewEntityCollection(nil)
	if err := ec.AddEntityFromMap(data); err != nil {
		t.Fatalf("unexpected error adding entity from map: %s", err)
	}

	entity := ec.Entities[0]
	values, err := entity.GetStringPropertyValues("http://data.mimiro.io/e360/props.externalIdentifiers.val")
	if err != nil {
		t.Fatalf("GetStringPropertyValues returned error: %s", err)
	}

	expected := []string{
		"http://data.mimiro.io/prodreg/producer-code/3425011396",
		"http://data.mimiro.io/prodreg/farm-code/34250113",
		"http://data.mimiro.io/prodreg/herdidentifier/2504433",
		"http://data.mimiro.io/prodreg/pid/10003427241",
		"http://data.mimiro.io/bronnoysund/organization-number/992666862",
	}
	if len(values) != len(expected) {
		t.Fatalf("expected %d values, got %d: %v", len(expected), len(values), values)
	}
	for i, v := range values {
		if v != expected[i] {
			t.Errorf("values[%d]: expected %q, got %q", i, expected[i], v)
		}
	}
}

func TestGetStringPropertyValues(t *testing.T) {
	entity := NewEntity().SetID("ns0:entity1")

	// []string case
	entity.SetProperty("ns0:strSlice", []string{"a", "b", "c"})
	if values, err := entity.GetStringPropertyValues("ns0:strSlice"); err != nil {
		t.Errorf("unexpected error: %s", err)
	} else if len(values) != 3 || values[0] != "a" || values[1] != "b" || values[2] != "c" {
		t.Errorf("unexpected values: %v", values)
	}

	// []any with all strings
	entity.SetProperty("ns0:anySlice", []any{"x", "y"})
	if values, err := entity.GetStringPropertyValues("ns0:anySlice"); err != nil {
		t.Errorf("unexpected error: %s", err)
	} else if len(values) != 2 || values[0] != "x" || values[1] != "y" {
		t.Errorf("unexpected values: %v", values)
	}

	// []any with non-string elements should skip them, not insert ""
	entity.SetProperty("ns0:mixedSlice", []any{"hello", 42, "world"})
	if values, err := entity.GetStringPropertyValues("ns0:mixedSlice"); err != nil {
		t.Errorf("unexpected error: %s", err)
	} else if len(values) != 2 || values[0] != "hello" || values[1] != "world" {
		t.Errorf("expected non-string elements to be skipped, got: %v", values)
	}

	// single string case
	entity.SetProperty("ns0:singleStr", "solo")
	if values, err := entity.GetStringPropertyValues("ns0:singleStr"); err != nil {
		t.Errorf("unexpected error: %s", err)
	} else if len(values) != 1 || values[0] != "solo" {
		t.Errorf("unexpected values: %v", values)
	}

	// missing key
	if _, err := entity.GetStringPropertyValues("ns0:missing"); err == nil {
		t.Error("expected error for missing key")
	}

	// wrong type
	entity.SetProperty("ns0:badType", 123)
	if _, err := entity.GetStringPropertyValues("ns0:badType"); err == nil {
		t.Error("expected error for non-string property type")
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

func TestCreateEntityFromMap(t *testing.T) {
	// define map with id, props and refs
	data := make(map[string]any)
	data["id"] = "ns0:entity1"
	data["props"] = make(map[string]any)
	data["props"].(map[string]any)["ns0:property1"] = "value1"
	data["refs"] = make(map[string]any)
	data["refs"].(map[string]any)["ns0:reference1"] = "ns0:entity2"

	// create a new entity collection
	ec := NewEntityCollection(nil)
	err := ec.AddEntityFromMap(data)

	// check for error
	if err != nil {
		t.Error(err)
	}

	// check the entity id
	if ec.Entities[0].ID != "ns0:entity1" {
		t.Errorf("expected entity id to be 'ns0:entity1', got '%s'", ec.Entities[0].ID)
	}

	// check the entity property
	if ec.Entities[0].Properties["ns0:property1"] != "value1" {
		t.Errorf("expected entity property to be 'value1', got '%s'", ec.Entities[0].Properties["ns0:property1"])
	}

	// check the entity reference
	if ec.Entities[0].References["ns0:reference1"] != "ns0:entity2" {
		t.Errorf("expected entity reference to be 'ns0:entity2', got '%s'", ec.Entities[0].References["ns0:reference1"])
	}
}

func TestCreateEntityFromMapWithWrongDataTypeForDeleted(t *testing.T) {
	// define map with id, props and refs
	data := make(map[string]any)
	data["id"] = "ns0:entity1"
	data["deleted"] = "true"

	// create a new entity collection
	ec := NewEntityCollection(nil)
	err := ec.AddEntityFromMap(data)
	if err != nil {
		t.Error("unexpected error")
	}

	// check that deleted is false
	if ec.Entities[0].IsDeleted != false {
		t.Errorf("expected entity deleted to be false, got '%t'", ec.Entities[0].IsDeleted)
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
