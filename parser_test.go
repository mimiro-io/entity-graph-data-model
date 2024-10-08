package egdm

// basic go test for the parser using standard test setup
import (
	"bytes"
	"testing"
)

func TestParseValidSimpleEntity(t *testing.T) {

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

	nsManager := NewNamespaceContext()
	parser := NewEntityParser(nsManager) //.WithExpandURIs()
	entityCollection, err := parser.LoadEntityCollection(byteReader)

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

func TestParseValidSimpleEntityNoContext(t *testing.T) {

	byteReader := bytes.NewReader([]byte(`
		[
			{
				"id": "http://example.com/1",
				"props": {
					"http://example.com/name": "John Smith"
				}
			}
		]`))

	nsManager := NewNamespaceContext()
	parser := NewEntityParser(nsManager).WithNoContext().WithExpandURIs()
	entityCollection, err := parser.LoadEntityCollection(byteReader)

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

	byteReader := bytes.NewReader([]byte(`
		[
			{
				"id": "http://example.com/1",
				"props": {
					"http://example.com/name": "John Smith"
				}
			}
		]`))

	nsManager := NewNamespaceContext()
	parser := NewEntityParser(nsManager).WithExpandURIs()
	_, err := parser.LoadEntityCollection(byteReader)

	if err == nil {
		t.Errorf("Expected error with missing context")
	}

	if err != nil {
		t.Logf("Got expected error: %s", err)
	}
}

func TestParseMissingNamespaceMappings(t *testing.T) {

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

	nsManager := NewNamespaceContext()
	parser := NewEntityParser(nsManager).WithExpandURIs()
	_, err := parser.LoadEntityCollection(byteReader)

	if err != nil {
		t.Errorf("Error parsing entity collection: %s", err)
	}
}

func TestParseBadExpansionWithMissingHashOrSlash(t *testing.T) {

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

	nsManager := NewNamespaceContext()
	parser := NewEntityParser(nsManager).WithExpandURIs()
	_, err := parser.LoadEntityCollection(byteReader)

	if err == nil {
		t.Errorf("Expected error due to bad context definition")
	}

	if err != nil {
		t.Logf("Got expected error: %s", err)
	}
}

func TestAllowsBadExpansionWithMissingHashOrSlashAndLenientModeOn(t *testing.T) {

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

	nsManager := NewNamespaceContext()
	parser := NewEntityParser(nsManager).WithExpandURIs().WithLenientNamespaceChecks()
	_, err := parser.LoadEntityCollection(byteReader)

	if err != nil {
		t.Errorf("Expected no error due to bad context definition")
	}

	// lookup the namespace
	exp, err := nsManager.GetNamespaceExpansionForPrefix("ex")
	if err != nil {
		t.Errorf("Expected no error looking up namespace")
	}

	if exp != "http://example.com" {
		t.Errorf("Expected namespace expansion to be http://example.com, got %s", exp)
	}

	// check that when that namespace is used, it is expanded
	pfx, err := nsManager.AssertPrefixedIdentifierFromURI("http://example.com/1")
	if err != nil {
		t.Errorf("Expected no error looking up namespace")
	}

	if pfx != "ns1:1" {
		t.Errorf("Expected namespace expansion to be ex:1, got %s", pfx)
	}
}

func TestParseBadJSONForContextDefinition(t *testing.T) {

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

	nsManager := NewNamespaceContext()
	parser := NewEntityParser(nsManager).WithExpandURIs()
	_, err := parser.LoadEntityCollection(byteReader)

	if err == nil {
		t.Errorf("Expected error due to bad context definition")
	}

	if err != nil {
		t.Logf("Got expected error: %s", err)
	}
}

func TestParseInvalidJSONMissingComma(t *testing.T) {

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

	nsManager := NewNamespaceContext()
	parser := NewEntityParser(nsManager).WithExpandURIs()
	_, err := parser.LoadEntityCollection(byteReader)

	if err == nil {
		t.Errorf("Expected error due to invalid json")
	}

	if err != nil {
		t.Logf("Got expected error: %s", err)
	}
}

func TestParseWithNamespaceExpansionInPropName(t *testing.T) {

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

	nsManager := NewNamespaceContext()
	parser := NewEntityParser(nsManager).WithExpandURIs()
	entityCollection, err := parser.LoadEntityCollection(byteReader)

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

	nsManager := NewNamespaceContext()
	parser := NewEntityParser(nsManager).WithExpandURIs()
	entityCollection, err := parser.LoadEntityCollection(byteReader)

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

func TestParseWithNamespaceExpansionInIdAndExtraColonInId(t *testing.T) {

	byteReader := bytes.NewReader([]byte(`
		[
			{
				"id" : "@context", 
				"namespaces": {
					"ex": "http://example.com/"
				}
			},
			{
				"id" : "ex:my:1",
				"props": {
					"http://example.com/name": "John Smith"
				}
			}
		]`))

	nsManager := NewNamespaceContext()
	parser := NewEntityParser(nsManager).WithExpandURIs()
	entityCollection, err := parser.LoadEntityCollection(byteReader)

	if err != nil {
		t.Errorf("Error parsing entity collection: %s", err)
	}
	if len(entityCollection.Entities) != 1 {
		t.Errorf("Expected 1 entity, got %d", len(entityCollection.Entities))
	}
	if entityCollection.Entities[0].ID != "http://example.com/my:1" {
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

	nsManager := NewNamespaceContext()
	parser := NewEntityParser(nsManager).WithExpandURIs()
	entityCollection, err := parser.LoadEntityCollection(byteReader)

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
					"rdf:type": [ "ex:Person", "ex:Employee", "ex:my:valid_char-s" ]
				}
			}
		]`))

	nsManager := NewNamespaceContext()
	parser := NewEntityParser(nsManager).WithExpandURIs()
	entityCollection, err := parser.LoadEntityCollection(byteReader)

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

	if len(refTypes) != 3 {
		t.Errorf("Expected entity reference type to be array of 3, got %d", len(refTypes))
	}
	// check elements of array
	if refTypes[0] != "http://example.com/Person" {
		t.Errorf("Expected entity reference type to be http://example.com/Person, got %s", refTypes[0])
	}
	if refTypes[1] != "http://example.com/Employee" {
		t.Errorf("Expected entity reference type to be http://example.com/Employee, got %s", refTypes[1])
	}
	if refTypes[2] != "http://example.com/my:valid_char-s" {
		t.Errorf("Expected entity reference type to be http://example.com/my:valid_char-s, got %s", refTypes[1])
	}
}

func TestParseWithEmbeddedAnonymousEntity(t *testing.T) {

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

	nsManager := NewNamespaceContext()
	parser := NewEntityParser(nsManager).WithExpandURIs()
	entityCollection, err := parser.LoadEntityCollection(byteReader)

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

	nsManager := NewNamespaceContext()
	parser := NewEntityParser(nsManager).WithExpandURIs()
	entityCollection, err := parser.LoadEntityCollection(byteReader)

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

	nsManager := NewNamespaceContext()
	parser := NewEntityParser(nsManager).WithExpandURIs()
	entityCollection, err := parser.LoadEntityCollection(byteReader)

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

	nsManager := NewNamespaceContext()
	parser := NewEntityParser(nsManager).WithExpandURIs()
	entityCollection, err := parser.LoadEntityCollection(byteReader)
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
	parser = NewEntityParser(nsManager).WithExpandURIs()

	byteReader = bytes.NewReader(bytesBuffer.Bytes())
	entityCollection, err = parser.LoadEntityCollection(byteReader)

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

func TestParseRoundTripEntityCollectionWithNoContextWritten(t *testing.T) {

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

	nsManager := NewNamespaceContext()
	parser := NewEntityParser(nsManager).WithExpandURIs()
	entityCollection, err := parser.LoadEntityCollection(byteReader)
	if err != nil {
		t.Errorf("Error parsing entity collection: %s", err)
	}

	// write it out again
	bytesBuffer := bytes.Buffer{}
	entityCollection.SetOmitContextOnWrite(true)
	err = entityCollection.WriteEntityGraphJSON(&bytesBuffer)
	if err != nil {
		t.Errorf("Error writing entity collection: %s", err)
	}

	// and parse it back into a new entity collection
	parser = NewEntityParser(nsManager).WithExpandURIs().WithNoContext()

	byteReader = bytes.NewReader(bytesBuffer.Bytes())
	entityCollection, err = parser.LoadEntityCollection(byteReader)

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

	nsManager := NewNamespaceContext()
	parser := NewEntityParser(nsManager).WithExpandURIs()
	entityCollection, err := parser.LoadEntityCollection(byteReader)

	if err != nil {
		t.Errorf("Error parsing entity collection: %s", err)
	}

	// write it out again
	bytesBuffer := bytes.Buffer{}
	entityCollection.WriteJSON_LD(&bytesBuffer)
}

func TestParseWithDefaultNamespaceExpansionInPropName(t *testing.T) {

	byteReader := bytes.NewReader([]byte(`
		[
			{
				"id" : "@context", 
				"namespaces": {
					"_": "http://example.com/"
				}
			},
			{
				"id" : "http://example.com/1",
				"props": {
					"name": "John Smith"
				}
			}
		]`))

	nsManager := NewNamespaceContext()
	parser := NewEntityParser(nsManager).WithExpandURIs()
	entityCollection, err := parser.LoadEntityCollection(byteReader)

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

func TestParseWithDefaultNamespaceExpansionInRefs(t *testing.T) {

	byteReader := bytes.NewReader([]byte(`
		[
			{
				"id" : "@context", 
				"namespaces": {
					"_": "http://example.com/",
					"rdf": "http://www.w3.org/1999/02/22-rdf-syntax-ns#"
				}
			},
			{
				"id" : "1",
				"props": {
					"http://example.com/name": "John Smith"
				},
				"refs": {
					"parent": "2",
					"rdf:type": "Person"
				}
			}
		]`))

	nsManager := NewNamespaceContext()
	parser := NewEntityParser(nsManager).WithExpandURIs()
	entityCollection, err := parser.LoadEntityCollection(byteReader)

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

func TestParseWithoutContext(t *testing.T) {

	byteReader := bytes.NewReader([]byte(`
		[
			{
				"id" : "http://example.com/1",
				"props": {
					"http://example.com/name": "John Smith"
				},
				"refs": {
					"http://example.com/parent": "http://example.com/2"
				}
			}
		]`))

	nsManager := NewNamespaceContext()
	parser := NewEntityParser(nsManager).WithNoContext().WithExpandURIs()
	entityCollection, err := parser.LoadEntityCollection(byteReader)

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
	if len(entityCollection.Entities[0].References) != 1 {
		t.Errorf("Expected entity references to have 2 properties, got %d", len(entityCollection.Entities[0].References))
	}
	if entityCollection.Entities[0].References["http://example.com/parent"] != "http://example.com/2" {
		t.Errorf("Expected entity reference parent to be http://example.com/2, got %s", entityCollection.Entities[0].References["parent"])
	}
}

func TestParseCompressURIs(t *testing.T) {

	byteReader := bytes.NewReader([]byte(`
		[
			{
				"id" : "http://example.com/1",
				"props": {
					"http://example.com/name": "John Smith"
				},
				"refs": {
					"http://example.com/parent": "http://example.com/2",
					"http://www.w3.org/1999/02/22-rdf-syntax-ns#type" : "http://example.com/Person" 
				}
			}
		]`))

	nsManager := NewNamespaceContext()
	parser := NewEntityParser(nsManager).WithNoContext().WithCompressURIs()
	entityCollection, err := parser.LoadEntityCollection(byteReader)

	if err != nil {
		t.Errorf("Error parsing entity collection: %s", err)
	}
	if len(entityCollection.Entities) != 1 {
		t.Errorf("Expected 1 entity, got %d", len(entityCollection.Entities))
	}
	if entityCollection.Entities[0].ID != "ns0:1" {
		t.Errorf("Expected entity id to be ns0:1, got %s", entityCollection.Entities[0].ID)
	}
	if len(entityCollection.Entities[0].Properties) != 1 {
		t.Errorf("Expected entity properties to have 1 property, got %d", len(entityCollection.Entities[0].Properties))
	}
	if entityCollection.Entities[0].Properties["ns0:name"] != "John Smith" {
		t.Errorf("Expected entity property name to be John Smith, got %s", entityCollection.Entities[0].Properties["name"])
	}
	if len(entityCollection.Entities[0].References) != 2 {
		t.Errorf("Expected entity references to have 2 properties, got %d", len(entityCollection.Entities[0].References))
	}
	if entityCollection.Entities[0].References["ns0:parent"] != "ns0:2" {
		t.Errorf("Expected entity reference parent to be http://example.com/2, got %s", entityCollection.Entities[0].References["parent"])
	}

	// check that the context is correct
	context := entityCollection.NamespaceManager.AsContext()
	if len(context.Namespaces) != 2 {
		t.Errorf("Expected context to have 1 namespace, got %d", len(context.Namespaces))
	}
	if context.Namespaces["ns0"] != "http://example.com/" {
		t.Errorf("Expected context namespace ns0 to be http://example.com/, got %s", context.Namespaces["ns0"])
	}
	if context.Namespaces["ns1"] != "http://www.w3.org/1999/02/22-rdf-syntax-ns#" {
		t.Errorf("Expected context namespace ns1 to be http://www.w3.org/1999/02/22-rdf-syntax-ns#, got %s", context.Namespaces["ns1"])
	}

}

func TestParseNoSpecialInstruction(t *testing.T) {

	byteReader := bytes.NewReader([]byte(`
		[
{"id":"@context","namespaces":{}},
{"id":"http://data.sample.org/things/1","refs":{},"props":{"http://data.sample.org/Name":"John"}},
{"id":"http://data.sample.org/things/2","refs":{},"props":{"http://data.sample.org/Name":"Jane"}},
{"id":"http://data.sample.org/things/3","refs":{},"props":{"http://data.sample.org/Name":"Jim"}},
{"id":"@continuation","token":"1725182073988287"}]`))

	nsManager := NewNamespaceContext()
	parser := NewEntityParser(nsManager)
	entityCollection, err := parser.LoadEntityCollection(byteReader)

	if err != nil {
		t.Errorf("Error parsing entity collection: %s", err)
	}
	if len(entityCollection.Entities) != 3 {
		t.Errorf("Expected 1 entity, got %d", len(entityCollection.Entities))
	}
	if entityCollection.Entities[0].ID != "http://data.sample.org/things/1" {
		t.Errorf("Expected entity id to be ns0:1, got %s", entityCollection.Entities[0].ID)
	}
}

func TestParserDetectsMissingNamespace(t *testing.T) {

	byteReader := bytes.NewReader([]byte(`
		[ {"id":"@context","namespaces":{}},
		  {"id":"ns0:1","refs":{},"props":{"http://data.sample.org/Name":"John"}},
		  {"id":"@continuation","token":"1725182073988287"}
		]`))

	nsManager := NewNamespaceContext()
	parser := NewEntityParser(nsManager)
	_, err := parser.LoadEntityCollection(byteReader)

	if err == nil {
		t.Error("Failed to detect missing namespace expansion")
	}
}

func TestParserDetectsNamespace(t *testing.T) {

	byteReader := bytes.NewReader([]byte(`
		[ {"id":"@context","namespaces":{ "ns0" : "http://data.sample.org/"}},
		  {"id":"ns0:1","refs":{},"props":{"http://data.sample.org/Name":"John"}},
		  {"id":"@continuation","token":"1725182073988287"}
		]`))

	nsManager := NewNamespaceContext()
	parser := NewEntityParser(nsManager)
	_, err := parser.LoadEntityCollection(byteReader)

	if err != nil {
		t.Error("Error validating namespace expansion")
	}
}

func TestParserValidatesDefaultNamespace(t *testing.T) {

	byteReader := bytes.NewReader([]byte(`
		[ {"id":"@context","namespaces":{ "_" : "http://data.sample.org/"}},
		  {"id":"1","refs":{},"props":{"http://data.sample.org/Name":"John"}},
		  {"id":"@continuation","token":"1725182073988287"}
		]`))

	nsManager := NewNamespaceContext()
	parser := NewEntityParser(nsManager)
	_, err := parser.LoadEntityCollection(byteReader)

	if err != nil {
		t.Error("Error validating namespace expansion")
	}
}

func TestParserFailsWhenNoDefaultNamespace(t *testing.T) {

	byteReader := bytes.NewReader([]byte(`
		[ {"id":"@context","namespaces":{ "ns0" : "http://data.sample.org/"}},
		  {"id":"1","refs":{},"props":{"http://data.sample.org/Name":"John"}}
		]`))

	nsManager := NewNamespaceContext()
	parser := NewEntityParser(nsManager)
	_, err := parser.LoadEntityCollection(byteReader)

	if err == nil {
		t.Error("Error validating namespace expansion")
	}
}

func TestParserDealsWithNullRefsAndProps(t *testing.T) {

	byteReader := bytes.NewReader([]byte(`
		[ {"id":"@context","namespaces":{}},
		  {"id":"http://data.sample.org/1","refs":null,"props":{"http://data.sample.org/Name":"John"}},
		  {"id":"http://data.sample.org/2","refs":{},"props":null}
		]`))

	nsManager := NewNamespaceContext()
	parser := NewEntityParser(nsManager)
	_, err := parser.LoadEntityCollection(byteReader)

	if err != nil {
		t.Errorf("Failed to parse collection with null refs %s", err)
	}
}
