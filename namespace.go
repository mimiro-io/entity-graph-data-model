package egdm

type NamespaceManager interface {
	GetNamespaceExpansionForPrefix(prefix string) (string, error)
	GetPrefixForExpansion(expansion string) (string, error)
	StorePrefixExpansionMapping(prefix string, expansion string)
	IsFullUri(value string) bool
	GetFullURI(value string) (string, error)
	GetPrefixedIdentifier(value string) (string, error)
	GetNamespaceMappings() map[string]string
	AssertPrefixFromURI(URI string) (string, error)
	AsContext() *Context
}
