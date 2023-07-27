package egdm

// basic go test for the parser using standard test setup
import (
	"bytes"
	"testing"
)

func TestParseValidSimpleEntity(t *testing.T) {

	nsManager := NewNamespaceContext()
	entityCollection := NewEntityCollection(nsManager)
	parser := NewEntityParser(nsManager, true)

	byteReader := bytes.NewReader([]byte(`
		[
			{
				"id" : "@context",
				"namespaces": {
				}
			},
			{
				"id": "http://example.com/1",
				"props": {
					"http://example.com/name": "John Smith"
				}
			}
		]`))

	err := parser.Parse(byteReader, entityCollection.AddEntity, entityCollection.SetContinuationToken)

	if err != nil {
		t.Errorf("Error parsing entity collection: %s", err)
	}
	if len(entityCollection.Entities) != 1 {
		t.Errorf("Expected 1 entity, got %d", len(entityCollection.Entities))
	}
	if entityCollection.Entities[0].ID != "http://example.com/1" {
		t.Errorf("Expected entity id to be http://example.com/1, got %s", entityCollection.Entities[0].ID)
	}
	if len(entityCollection.Entities[0].Properties) != 1 {
		t.Errorf("Expected entity properties to have 1 property, got %d", len(entityCollection.Entities[0].Properties))
	}
	if entityCollection.Entities[0].Properties["http://example.com/name"] != "John Smith" {
		t.Errorf("Expected entity property name to be John Smith, got %s", entityCollection.Entities[0].Properties["name"])
	}
}

func TestParseMissingNamespaceContext(t *testing.T) {

	nsManager := NewNamespaceContext()
	entityCollection := NewEntityCollection(nsManager)
	parser := NewEntityParser(nsManager, true)

	byteReader := bytes.NewReader([]byte(`
		[
			{
				"id": "http://example.com/1",
				"props": {
					"http://example.com/name": "John Smith"
				}
			}
		]`))

	err := parser.Parse(byteReader, entityCollection.AddEntity, entityCollection.SetContinuationToken)

	if err == nil {
		t.Errorf("Expected error with missing context")
	}

	if err != nil {
		t.Logf("Got expected error: %s", err)
	}
}

func TestParseMissingNamespaceMappings(t *testing.T) {

	nsManager := NewNamespaceContext()
	entityCollection := NewEntityCollection(nsManager)
	parser := NewEntityParser(nsManager, true)

	byteReader := bytes.NewReader([]byte(`
		[
			{
				"id" : "@context"
			},
			{
				"id": "http://example.com/1",
				"props": {
					"http://example.com/name": "John Smith"
				}
			}
		]`))

	err := parser.Parse(byteReader, entityCollection.AddEntity, entityCollection.SetContinuationToken)

	if err != nil {
		t.Errorf("Error parsing entity collection: %s", err)
	}
}

func TestParseBadExpansionWithMissingHashOrSlash(t *testing.T) {

	nsManager := NewNamespaceContext()
	entityCollection := NewEntityCollection(nsManager)
	parser := NewEntityParser(nsManager, true)

	byteReader := bytes.NewReader([]byte(`
		[
			{
				"id" : "@context",
				"namespaces": {
					"ex": "http://example.com"
				}
			},
			{
				"id": "ex:1",
				"props": {
					"http://example.com/name": "John Smith"
				}
			}
		]`))

	err := parser.Parse(byteReader, entityCollection.AddEntity, entityCollection.SetContinuationToken)

	if err == nil {
		t.Errorf("Expected error due to bad context definition")
	}

	if err != nil {
		t.Logf("Got expected error: %s", err)
	}
}

func TestParseBadJSONForContextDefinition(t *testing.T) {

	nsManager := NewNamespaceContext()
	entityCollection := NewEntityCollection(nsManager)
	parser := NewEntityParser(nsManager, true)

	byteReader := bytes.NewReader([]byte(`
		[
			[
				"id", "@context"
			],
			{
				"id": "http://example.com/1",
				"props": {
					"http://example.com/name": "John Smith"
				}
			}
		]`))

	err := parser.Parse(byteReader, entityCollection.AddEntity, entityCollection.SetContinuationToken)

	if err == nil {
		t.Errorf("Expected error due to bad context definition")
	}

	if err != nil {
		t.Logf("Got expected error: %s", err)
	}
}

func TestParseInvalidJSONMissingComma(t *testing.T) {

	nsManager := NewNamespaceContext()
	entityCollection := NewEntityCollection(nsManager)
	parser := NewEntityParser(nsManager, true)

	byteReader := bytes.NewReader([]byte(`
		[
			{
				"id" : "@context"
			},
			{
				"id" : "http://example.com/1"
				"props": {
					"http://example.com/name": "John Smith"
				}
			}
		]`))

	err := parser.Parse(byteReader, entityCollection.AddEntity, entityCollection.SetContinuationToken)

	if err == nil {
		t.Errorf("Expected error due to invalid json")
	}

	if err != nil {
		t.Logf("Got expected error: %s", err)
	}
}

func TestParseWithNamespaceExpansionInPropName(t *testing.T) {

	nsManager := NewNamespaceContext()
	entityCollection := NewEntityCollection(nsManager)
	parser := NewEntityParser(nsManager, true)

	byteReader := bytes.NewReader([]byte(`
		[
			{
				"id" : "@context", 
				"namespaces": {
					"ex": "http://example.com/"
				}
			},
			{
				"id" : "http://example.com/1",
				"props": {
					"ex:name": "John Smith"
				}
			}
		]`))

	err := parser.Parse(byteReader, entityCollection.AddEntity, entityCollection.SetContinuationToken)

	if err != nil {
		t.Errorf("Error parsing entity collection: %s", err)
	}
	if len(entityCollection.Entities) != 1 {
		t.Errorf("Expected 1 entity, got %d", len(entityCollection.Entities))
	}
	if entityCollection.Entities[0].ID != "http://example.com/1" {
		t.Errorf("Expected entity id to be http://example.com/1, got %s", entityCollection.Entities[0].ID)
	}
	if len(entityCollection.Entities[0].Properties) != 1 {
		t.Errorf("Expected entity properties to have 1 property, got %d", len(entityCollection.Entities[0].Properties))
	}
	if entityCollection.Entities[0].Properties["http://example.com/name"] != "John Smith" {
		t.Errorf("Expected entity property name to be John Smith, got %s", entityCollection.Entities[0].Properties["name"])
	}
}

func TestParseWithNamespaceExpansionInId(t *testing.T) {

	nsManager := NewNamespaceContext()
	entityCollection := NewEntityCollection(nsManager)
	parser := NewEntityParser(nsManager, true)

	byteReader := bytes.NewReader([]byte(`
		[
			{
				"id" : "@context", 
				"namespaces": {
					"ex": "http://example.com/"
				}
			},
			{
				"id" : "ex:1",
				"props": {
					"http://example.com/name": "John Smith"
				}
			}
		]`))

	err := parser.Parse(byteReader, entityCollection.AddEntity, entityCollection.SetContinuationToken)

	if err != nil {
		t.Errorf("Error parsing entity collection: %s", err)
	}
	if len(entityCollection.Entities) != 1 {
		t.Errorf("Expected 1 entity, got %d", len(entityCollection.Entities))
	}
	if entityCollection.Entities[0].ID != "http://example.com/1" {
		t.Errorf("Expected entity id to be http://example.com/1, got %s", entityCollection.Entities[0].ID)
	}
	if len(entityCollection.Entities[0].Properties) != 1 {
		t.Errorf("Expected entity properties to have 1 property, got %d", len(entityCollection.Entities[0].Properties))
	}
	if entityCollection.Entities[0].Properties["http://example.com/name"] != "John Smith" {
		t.Errorf("Expected entity property name to be John Smith, got %s", entityCollection.Entities[0].Properties["name"])
	}
}

func TestParseWithNamespaceExpansionInIRefs(t *testing.T) {

	nsManager := NewNamespaceContext()
	entityCollection := NewEntityCollection(nsManager)
	parser := NewEntityParser(nsManager, true)

	byteReader := bytes.NewReader([]byte(`
		[
			{
				"id" : "@context", 
				"namespaces": {
					"ex": "http://example.com/",
					"rdf": "http://www.w3.org/1999/02/22-rdf-syntax-ns#"
				}
			},
			{
				"id" : "ex:1",
				"props": {
					"http://example.com/name": "John Smith"
				},
				"refs": {
					"http://example.com/parent": "ex:2",
					"rdf:type": "ex:Person"
				}
			}
		]`))

	err := parser.Parse(byteReader, entityCollection.AddEntity, entityCollection.SetContinuationToken)

	if err != nil {
		t.Errorf("Error parsing entity collection: %s", err)
	}
	if len(entityCollection.Entities) != 1 {
		t.Errorf("Expected 1 entity, got %d", len(entityCollection.Entities))
	}
	if entityCollection.Entities[0].ID != "http://example.com/1" {
		t.Errorf("Expected entity id to be http://example.com/1, got %s", entityCollection.Entities[0].ID)
	}
	if len(entityCollection.Entities[0].Properties) != 1 {
		t.Errorf("Expected entity properties to have 1 property, got %d", len(entityCollection.Entities[0].Properties))
	}
	if entityCollection.Entities[0].Properties["http://example.com/name"] != "John Smith" {
		t.Errorf("Expected entity property name to be John Smith, got %s", entityCollection.Entities[0].Properties["name"])
	}
	if len(entityCollection.Entities[0].References) != 2 {
		t.Errorf("Expected entity references to have 2 properties, got %d", len(entityCollection.Entities[0].References))
	}
	if entityCollection.Entities[0].References["http://example.com/parent"] != "http://example.com/2" {
		t.Errorf("Expected entity reference parent to be http://example.com/2, got %s", entityCollection.Entities[0].References["parent"])
	}
	if entityCollection.Entities[0].References["http://www.w3.org/1999/02/22-rdf-syntax-ns#type"] != "http://example.com/Person" {
		t.Errorf("Expected entity reference type to be http://example.com/Person, got %s", entityCollection.Entities[0].References["type"])
	}

}

func TestParseWithRefArray(t *testing.T) {

	nsManager := NewNamespaceContext()
	entityCollection := NewEntityCollection(nsManager)
	parser := NewEntityParser(nsManager, true)

	byteReader := bytes.NewReader([]byte(`
		[
			{
				"id" : "@context", 
				"namespaces": {
					"ex": "http://example.com/",
					"rdf": "http://www.w3.org/1999/02/22-rdf-syntax-ns#"
				}
			},
			{
				"id" : "ex:1",
				"props": {
					"http://example.com/name": "John Smith"
				},
				"refs": {
					"http://example.com/parent": "ex:2",
					"rdf:type": [ "ex:Person", "ex:Employee" ]
				}
			}
		]`))

	err := parser.Parse(byteReader, entityCollection.AddEntity, entityCollection.SetContinuationToken)

	if err != nil {
		t.Errorf("Error parsing entity collection: %s", err)
	}
	if len(entityCollection.Entities) != 1 {
		t.Errorf("Expected 1 entity, got %d", len(entityCollection.Entities))
	}
	if entityCollection.Entities[0].ID != "http://example.com/1" {
		t.Errorf("Expected entity id to be http://example.com/1, got %s", entityCollection.Entities[0].ID)
	}
	if len(entityCollection.Entities[0].Properties) != 1 {
		t.Errorf("Expected entity properties to have 1 property, got %d", len(entityCollection.Entities[0].Properties))
	}
	if entityCollection.Entities[0].Properties["http://example.com/name"] != "John Smith" {
		t.Errorf("Expected entity property name to be John Smith, got %s", entityCollection.Entities[0].Properties["name"])
	}
	if len(entityCollection.Entities[0].References) != 2 {
		t.Errorf("Expected entity references to have 2 properties, got %d", len(entityCollection.Entities[0].References))
	}
	if entityCollection.Entities[0].References["http://example.com/parent"] != "http://example.com/2" {
		t.Errorf("Expected entity reference parent to be http://example.com/2, got %s", entityCollection.Entities[0].References["parent"])
	}

	refTypes := entityCollection.Entities[0].References["http://www.w3.org/1999/02/22-rdf-syntax-ns#type"].([]string)

	if len(refTypes) != 2 {
		t.Errorf("Expected entity reference type to be array of 2, got %d", len(refTypes))
	}
	// check elements of array
	if refTypes[0] != "http://example.com/Person" {
		t.Errorf("Expected entity reference type to be http://example.com/Person, got %s", refTypes[0])
	}
	if refTypes[1] != "http://example.com/Employee" {
		t.Errorf("Expected entity reference type to be http://example.com/Employee, got %s", refTypes[1])
	}
}

func TestParseWithEmbeddedAnonymousEntity(t *testing.T) {

	nsManager := NewNamespaceContext()
	entityCollection := NewEntityCollection(nsManager)
	parser := NewEntityParser(nsManager, true)

	byteReader := bytes.NewReader([]byte(`
		[
			{
				"id" : "@context", 
				"namespaces": {
					"ex": "http://example.com/",
					"rdf": "http://www.w3.org/1999/02/22-rdf-syntax-ns#"
				}
			},
			{
				"id" : "ex:1",
				"props": {
					"http://example.com/address": {
						"props": {
							"http://example.com/street": "123 Main Street"
						}
					}
				}
			}
		]`))

	err := parser.Parse(byteReader, entityCollection.AddEntity, entityCollection.SetContinuationToken)

	if err != nil {
		t.Errorf("Error parsing entity collection: %s", err)
	}
	if len(entityCollection.Entities) != 1 {
		t.Errorf("Expected 1 entity, got %d", len(entityCollection.Entities))
	}
	if entityCollection.Entities[0].ID != "http://example.com/1" {
		t.Errorf("Expected entity id to be http://example.com/1, got %s", entityCollection.Entities[0].ID)
	}
	if len(entityCollection.Entities[0].Properties) != 1 {
		t.Errorf("Expected entity properties to have 1 property, got %d", len(entityCollection.Entities[0].Properties))
	}

	embeddedEntity := entityCollection.Entities[0].Properties["http://example.com/address"].(*Entity)

	if len(embeddedEntity.Properties) != 1 {
		t.Errorf("Expected embedded entity properties to have 1 property, got %d", len(embeddedEntity.Properties))
	}

	if embeddedEntity.Properties["http://example.com/street"] != "123 Main Street" {
		t.Errorf("Expected embedded entity property street to be 123 Main Street, got %s", embeddedEntity.Properties["http://example.com/street"])
	}

}

func TestParseWithEmbeddedEntityWithIdentity(t *testing.T) {

	nsManager := NewNamespaceContext()
	entityCollection := NewEntityCollection(nsManager)
	parser := NewEntityParser(nsManager, true)

	byteReader := bytes.NewReader([]byte(`
		[
			{
				"id" : "@context", 
				"namespaces": {
					"ex": "http://example.com/",
					"rdf": "http://www.w3.org/1999/02/22-rdf-syntax-ns#"
				}
			},
			{
				"id" : "ex:1",
				"props": {
					"http://example.com/address": {
						"id": "ex:2",
						"props": {
							"http://example.com/street": "123 Main Street"
						}
					}
				}
			}
		]`))

	err := parser.Parse(byteReader, entityCollection.AddEntity, entityCollection.SetContinuationToken)

	if err != nil {
		t.Errorf("Error parsing entity collection: %s", err)
	}
	if len(entityCollection.Entities) != 1 {
		t.Errorf("Expected 1 entity, got %d", len(entityCollection.Entities))
	}
	if entityCollection.Entities[0].ID != "http://example.com/1" {
		t.Errorf("Expected entity id to be http://example.com/1, got %s", entityCollection.Entities[0].ID)
	}
	if len(entityCollection.Entities[0].Properties) != 1 {
		t.Errorf("Expected entity properties to have 1 property, got %d", len(entityCollection.Entities[0].Properties))
	}

	embeddedEntity := entityCollection.Entities[0].Properties["http://example.com/address"].(*Entity)

	if len(embeddedEntity.Properties) != 1 {
		t.Errorf("Expected embedded entity properties to have 1 property, got %d", len(embeddedEntity.Properties))
	}
	if embeddedEntity.Properties["http://example.com/street"] != "123 Main Street" {
		t.Errorf("Expected embedded entity property street to be 123 Main Street, got %s", embeddedEntity.Properties["http://example.com/street"])
	}
	if embeddedEntity.ID != "http://example.com/2" {
		t.Errorf("Expected embedded entity id to be http://example.com/2, got %s", embeddedEntity.ID)
	}

}

func TestParseWithEmbeddedEntityArray(t *testing.T) {

	nsManager := NewNamespaceContext()
	entityCollection := NewEntityCollection(nsManager)
	parser := NewEntityParser(nsManager, true)

	byteReader := bytes.NewReader([]byte(`
		[
			{
				"id" : "@context", 
				"namespaces": {
					"ex": "http://example.com/",
					"rdf": "http://www.w3.org/1999/02/22-rdf-syntax-ns#"
				}
			},
			{
				"id" : "ex:1",
				"props": {
					"http://example.com/addresses": 
						[
							{
								"id": "ex:2",
								"props": {
									"http://example.com/street": "123 Main Street"
								},
								"refs": {
									"http://example.com/country": "ex:5"
								}
							},
							{
								"id": "ex:3",
								"props": {
									"http://example.com/street": "125 Main Street"
								},
								"refs": {
									"http://example.com/country": "ex:6"
								}
							}
						]	
				}
			}
		]`))

	err := parser.Parse(byteReader, entityCollection.AddEntity, entityCollection.SetContinuationToken)

	if err != nil {
		t.Errorf("Error parsing entity collection: %s", err)
	}
	if len(entityCollection.Entities) != 1 {
		t.Errorf("Expected 1 entity, got %d", len(entityCollection.Entities))
	}
	if entityCollection.Entities[0].ID != "http://example.com/1" {
		t.Errorf("Expected entity id to be http://example.com/1, got %s", entityCollection.Entities[0].ID)
	}
	if len(entityCollection.Entities[0].Properties) != 1 {
		t.Errorf("Expected entity properties to have 1 property, got %d", len(entityCollection.Entities[0].Properties))
	}
	if len(entityCollection.Entities[0].References) != 0 {
		t.Errorf("Expected entity references to have 0 properties, got %d", len(entityCollection.Entities[0].References))
	}

	embeddedEntityArrayAny := entityCollection.Entities[0].Properties["http://example.com/addresses"].([]any)
	embeddedEntityArray := make([]*Entity, len(embeddedEntityArrayAny))
	for i, v := range embeddedEntityArrayAny {
		embeddedEntityArray[i] = v.(*Entity)
	}

	if len(embeddedEntityArray) != 2 {
		t.Errorf("Expected embedded entity array to have 2 elements, got %d", len(embeddedEntityArray))
	}

	embeddedEntity := embeddedEntityArray[0]

	if len(embeddedEntity.Properties) != 1 {
		t.Errorf("Expected embedded entity properties to have 1 property, got %d", len(embeddedEntity.Properties))
	}
	if embeddedEntity.Properties["http://example.com/street"] != "123 Main Street" {
		t.Errorf("Expected embedded entity property street to be 123 Main Street, got %s", embeddedEntity.Properties["http://example.com/street"])
	}

	if embeddedEntity.ID != "http://example.com/2" {
		t.Errorf("Expected embedded entity id to be http://example.com/2, got %s", embeddedEntity.ID)
	}

	if len(embeddedEntity.References) != 1 {
		t.Errorf("Expected embedded entity references to have 1 property, got %d", len(embeddedEntity.References))
	}
	if embeddedEntity.References["http://example.com/country"] != "http://example.com/5" {
		t.Errorf("Expected embedded entity reference country to be http://example.com/5, got %s", embeddedEntity.References["http://example.com/country"])
	}

	embeddedEntity = embeddedEntityArray[1]

	if len(embeddedEntity.Properties) != 1 {
		t.Errorf("Expected embedded entity properties to have 1 property, got %d", len(embeddedEntity.Properties))
	}
	if embeddedEntity.Properties["http://example.com/street"] != "125 Main Street" {
		t.Errorf("Expected embedded entity property street to be 125 Main Street, got %s", embeddedEntity.Properties["http://example.com/street"])
	}
	if embeddedEntity.ID != "http://example.com/3" {
		t.Errorf("Expected embedded entity id to be http://example.com/3, got %s", embeddedEntity.ID)
	}
	if len(embeddedEntity.References) != 1 {
		t.Errorf("Expected embedded entity references to have 1 property, got %d", len(embeddedEntity.References))
	}
	if embeddedEntity.References["http://example.com/country"] != "http://example.com/6" {
		t.Errorf("Expected embedded entity reference country to be http://example.com/6, got %s", embeddedEntity.References["http://example.com/country"])
	}

}

func TestParseRoundTripEntityCollection(t *testing.T) {

	nsManager := NewNamespaceContext()
	entityCollection := NewEntityCollection(nsManager)
	parser := NewEntityParser(nsManager, true)

	byteReader := bytes.NewReader([]byte(`
		[
			{
				"id" : "@context", 
				"namespaces": {
					"ex": "http://example.com/",
					"rdf": "http://www.w3.org/1999/02/22-rdf-syntax-ns#"
				}
			},
			{
				"id" : "ex:1",
				"props": {
					"http://example.com/addresses": 
						[
							{
								"id": "ex:2",
								"props": {
									"http://example.com/street": "123 Main Street"
								},
								"refs": {
									"http://example.com/country": "ex:5"
								}
							},
							{
								"id": "ex:3",
								"props": {
									"http://example.com/street": "125 Main Street"
								},
								"refs": {
									"http://example.com/country": "ex:6"
								}
							}
						]	
				}
			},
			{
				"id" : "@continuation",
				"token" : "1234567890"	
			}
		]`))

	err := parser.Parse(byteReader, entityCollection.AddEntity, entityCollection.SetContinuationToken)
	if err != nil {
		t.Errorf("Error parsing entity collection: %s", err)
	}

	// write it out again
	bytesBuffer := bytes.Buffer{}
	err = entityCollection.WriteEntityGraphJSON(&bytesBuffer)
	if err != nil {
		t.Errorf("Error writing entity collection: %s", err)
	}

	// and parse it back into a new entity collection
	entityCollection = NewEntityCollection(nsManager)
	parser = NewEntityParser(nsManager, true)
	byteReader = bytes.NewReader(bytesBuffer.Bytes())
	err = parser.Parse(byteReader, entityCollection.AddEntity, entityCollection.SetContinuationToken)
	if err != nil {
		t.Errorf("Error parsing entity collection: %s", err)
	}

	if len(entityCollection.Entities) != 1 {
		t.Errorf("Expected 1 entity, got %d", len(entityCollection.Entities))
	}
	if entityCollection.Entities[0].ID != "http://example.com/1" {
		t.Errorf("Expected entity id to be http://example.com/1, got %s", entityCollection.Entities[0].ID)
	}
	if len(entityCollection.Entities[0].Properties) != 1 {
		t.Errorf("Expected entity properties to have 1 property, got %d", len(entityCollection.Entities[0].Properties))
	}
	if len(entityCollection.Entities[0].References) != 0 {
		t.Errorf("Expected entity references to have 0 properties, got %d", len(entityCollection.Entities[0].References))
	}

	embeddedEntityArrayAny := entityCollection.Entities[0].Properties["http://example.com/addresses"].([]any)
	embeddedEntityArray := make([]*Entity, len(embeddedEntityArrayAny))
	for i, v := range embeddedEntityArrayAny {
		embeddedEntityArray[i] = v.(*Entity)
	}

	if len(embeddedEntityArray) != 2 {
		t.Errorf("Expected embedded entity array to have 2 elements, got %d", len(embeddedEntityArray))
	}

	embeddedEntity := embeddedEntityArray[0]

	if len(embeddedEntity.Properties) != 1 {
		t.Errorf("Expected embedded entity properties to have 1 property, got %d", len(embeddedEntity.Properties))
	}
	if embeddedEntity.Properties["http://example.com/street"] != "123 Main Street" {
		t.Errorf("Expected embedded entity property street to be 123 Main Street, got %s", embeddedEntity.Properties["http://example.com/street"])
	}

	if embeddedEntity.ID != "http://example.com/2" {
		t.Errorf("Expected embedded entity id to be http://example.com/2, got %s", embeddedEntity.ID)
	}

	if len(embeddedEntity.References) != 1 {
		t.Errorf("Expected embedded entity references to have 1 property, got %d", len(embeddedEntity.References))
	}
	if embeddedEntity.References["http://example.com/country"] != "http://example.com/5" {
		t.Errorf("Expected embedded entity reference country to be http://example.com/5, got %s", embeddedEntity.References["http://example.com/country"])
	}

	embeddedEntity = embeddedEntityArray[1]

	if len(embeddedEntity.Properties) != 1 {
		t.Errorf("Expected embedded entity properties to have 1 property, got %d", len(embeddedEntity.Properties))
	}
	if embeddedEntity.Properties["http://example.com/street"] != "125 Main Street" {
		t.Errorf("Expected embedded entity property street to be 125 Main Street, got %s", embeddedEntity.Properties["http://example.com/street"])
	}
	if embeddedEntity.ID != "http://example.com/3" {
		t.Errorf("Expected embedded entity id to be http://example.com/3, got %s", embeddedEntity.ID)
	}
	if len(embeddedEntity.References) != 1 {
		t.Errorf("Expected embedded entity references to have 1 property, got %d", len(embeddedEntity.References))
	}
	if embeddedEntity.References["http://example.com/country"] != "http://example.com/6" {
		t.Errorf("Expected embedded entity reference country to be http://example.com/6, got %s", embeddedEntity.References["http://example.com/country"])
	}

	if entityCollection.Continuation.Token != "1234567890" {
		t.Errorf("Expected continuation token to be 1234567890, got %s", entityCollection.Continuation.Token)
	}
}

func TestToJSONLDWithEmbeddedEntityArray(t *testing.T) {

	nsManager := NewNamespaceContext()
	entityCollection := NewEntityCollection(nsManager)
	parser := NewEntityParser(nsManager, true)

	byteReader := bytes.NewReader([]byte(`
		[
			{
				"id" : "@context", 
				"namespaces": {
					"ex": "http://example.com/",
					"rdf": "http://www.w3.org/1999/02/22-rdf-syntax-ns#"
				}
			},
			{
				"id" : "ex:1",
				"props": {
					"http://example.com/addresses": 
						[
							{
								"props": {
									"http://example.com/street": "123 Main Street"
								},
								"refs": {
									"http://example.com/country": "ex:5"
								}
							},
							{
								"props": {
									"http://example.com/street": "125 Main Street"
								},
								"refs": {
									"http://example.com/country": "ex:6"
								}
							}
						]	
				}
			}
		]`))

	err := parser.Parse(byteReader, entityCollection.AddEntity, entityCollection.SetContinuationToken)

	if err != nil {
		t.Errorf("Error parsing entity collection: %s", err)
	}

	// write it out again
	bytesBuffer := bytes.Buffer{}
	entityCollection.WriteJSON_LD(&bytesBuffer)

}
