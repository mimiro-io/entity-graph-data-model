
# Entity Graph Data Model

Data structures, parser and utilities for the entity graph data model

To use the module import the following:

`github.com/mimiro-io/entity-graph-data-model`

# Core Entity Graph Data Model structures

The following example shows how to create an Entity, add a property, reference and serialise it to JSON. 

``` go
package main

import ( 
    "fmt"
    egdm "github.com/mimiro-io/entity-graph-data-model"
    "encoding/json"
)

func main() {

    entity := egdm.NewEntity()
    entity.Properties["http://data.mimiro.io/schema/name"] = "homer"
    entity.References["http://data.mimiro.io/schema/worksFor"] = "http://data.mimiro.io/people/mrburns"

    entityJson, _ := json.Marshall(entity)
    fmt.Print(entityJson)
}
```

# Parsing Entity Graph Data Model JSON

The module can be used to parse entity graph data model JSON data. 

``` go
package main

import ( 
    "fmt"
    "github.com/mimiro-io/entity-graph-data-model"
)

func main() {
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
}
```
