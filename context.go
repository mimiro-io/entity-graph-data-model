package egdm

import (
	"errors"
	"strings"
)

func NewNamespaceContext() *NamespaceContext {
	context := &NamespaceContext{}
	context.prefixToExpansionMappings = make(map[string]string)
	context.expansionToPrefixMappings = make(map[string]string)
	return context
}

type NamespaceContext struct {
	prefixToExpansionMappings map[string]string
	expansionToPrefixMappings map[string]string
}

func (aContext *NamespaceContext) GetNamespaceExpansionForPrefix(prefix string) (string, error) {
	if expansion, found := aContext.prefixToExpansionMappings[prefix]; found {
		return expansion, nil
	} else {
		return "", errors.New("no expansion for prefix: " + prefix)
	}
}

func (aContext *NamespaceContext) GetPrefixForExpansion(expansion string) (string, error) {
	if prefix, found := aContext.expansionToPrefixMappings[expansion]; found {
		return prefix, nil
	} else {
		return "", errors.New("no expansion for prefix: " + expansion)
	}
}

func (aContext *NamespaceContext) StorePrefixExpansionMapping(prefix string, expansion string) {
	aContext.prefixToExpansionMappings[prefix] = expansion
	aContext.expansionToPrefixMappings[expansion] = prefix
}

func (aContext *NamespaceContext) IsFullUri(value string) bool {
	return strings.HasPrefix(value, "https:") || strings.HasPrefix(value, "http:")
}

func (aContext *NamespaceContext) isCURIE(value string) (bool, string, string) {
	if aContext.IsFullUri(value) {
		return false, "", ""
	} else {
		parts := strings.Split(value, ":")
		if len(parts) != 2 {
			return false, "", ""
		}
		return true, parts[0], parts[1]
	}
}

func (aContext *NamespaceContext) GetFullURI(value string) (string, error) {
	if aContext.IsFullUri(value) {
		return value, nil
	}

	isCURIE, prefix, postfix := aContext.isCURIE(value)
	if isCURIE {
		expansion, err := aContext.GetNamespaceExpansionForPrefix(prefix)
		if err != nil {
			return "", err
		}
		return expansion + postfix, nil
	} else {
		// lookup default namespace expansion
		// and append the original value
		expansion, err := aContext.GetNamespaceExpansionForPrefix("_")
		if err != nil {
			return "", err
		}
		return expansion + value, nil
	}
}

// implement get namespace mappings
func (aContext *NamespaceContext) GetNamespaceMappings() map[string]string {
	return aContext.prefixToExpansionMappings
}

// implement get prefixed identifier
func (aContext *NamespaceContext) GetPrefixedIdentifier(value string) (string, error) {
	if aContext.IsFullUri(value) {
		for prefix, expansion := range aContext.prefixToExpansionMappings {
			if strings.HasPrefix(value, expansion) {
				return prefix + ":" + strings.TrimPrefix(value, expansion), nil
			}
		}
		return "", errors.New("unable to find prefix for expansion: " + value)
	} else {
		return value, nil
	}
}
