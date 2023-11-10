package egdm

import "errors"

type Entity struct {
	ID         string         `json:"id,omitempty"`
	InternalID uint64         `json:"internalId,omitempty"`
	Recorded   uint64         `json:"recorded,omitempty"`
	IsDeleted  bool           `json:"deleted,omitempty"`
	References map[string]any `json:"refs"`
	Properties map[string]any `json:"props"`
}

func NewEntity() *Entity {
	e := &Entity{}
	e.References = make(map[string]any)
	e.Properties = make(map[string]any)
	return e
}

func (anEntity *Entity) SetID(id string) *Entity {
	anEntity.ID = id
	return anEntity
}

func (anEntity *Entity) SetProperty(property string, value any) *Entity {
	anEntity.Properties[property] = value
	return anEntity
}

func (anEntity *Entity) SetReference(reference string, value any) *Entity {
	anEntity.References[reference] = value
	return anEntity
}

func (anEntity *Entity) GetFirstReferenceValue(typeURI string) (string, error) {
	if values, found := anEntity.GetReferenceValues(typeURI); found == nil {
		if len(values) == 0 {
			return "", errors.New("no reference for type")
		}
		return values[0], nil
	}
	return "", errors.New("no reference for type")
}

func (anEntity *Entity) GetReferenceValues(typeURI string) ([]string, error) {
	if values, found := anEntity.References[typeURI]; found {
		switch v := values.(type) {
		case []string:
			return v, nil
		case string:
			result := make([]string, 1)
			result[0] = v
			return result, nil
		}
	}
	return nil, errors.New("no reference for type")
}

func (anEntity *Entity) GetFirstStringPropertyValue(typeURI string) (string, error) {
	if values, found := anEntity.GetStringPropertyValues(typeURI); found == nil {
		if len(values) == 0 {
			return "", errors.New("no reference for type")
		}
		return values[0], nil
	}
	return "", errors.New("no reference for type")
}

func (anEntity *Entity) GetStringPropertyValues(typeURI string) ([]string, error) {
	if values, found := anEntity.Properties[typeURI]; found {
		switch v := values.(type) {
		case []string:
			return v, nil
		case string:
			result := make([]string, 1)
			result[0] = v
			return result, nil
		}
		return nil, errors.New("property key exists but type is not string")
	}
	return nil, errors.New("no property string literal")
}

func (anEntity *Entity) GetFirstBooleanPropertyValue(typeURI string) (bool, error) {
	if values, found := anEntity.GetBooleanPropertyValues(typeURI); found == nil {
		if len(values) == 0 {
			return false, errors.New("no reference for type")
		}
		return values[0], nil
	}
	return false, errors.New("no reference for type")
}

func (anEntity *Entity) GetBooleanPropertyValues(typeURI string) ([]bool, error) {
	if values, found := anEntity.Properties[typeURI]; found {
		switch v := values.(type) {
		case []bool:
			return v, nil
		case bool:
			result := make([]bool, 1)
			result[0] = v
			return result, nil
		}
	}
	return nil, errors.New("no property boolean literal")
}

func (anEntity *Entity) GetFirstIntPropertyValue(typeURI string) (int, error) {
	if values, found := anEntity.GetIntPropertyValues(typeURI); found == nil {
		if len(values) == 0 {
			return 0, errors.New("no reference for type")
		}
		return values[0], nil
	}
	return 0, errors.New("no reference for type")
}

func (anEntity *Entity) GetIntPropertyValues(typeURI string) ([]int, error) {
	if values, found := anEntity.Properties[typeURI]; found {
		switch v := values.(type) {
		case []int:
			return v, nil
		case int:
			result := make([]int, 1)
			result[0] = v
			return result, nil
		case []float64:
			result := make([]int, len(v))
			for i, val := range v {
				result[i] = int(val)
			}
			return result, nil
		case float64:
			result := make([]int, 1)
			result[0] = int(v)
			return result, nil
		}
	}
	return nil, errors.New("no property int32 literal")
}

func (anEntity *Entity) GetFirstFloatPropertyValue(typeURI string) (float64, error) {
	if values, found := anEntity.GetFloatPropertyValues(typeURI); found == nil {
		if len(values) == 0 {
			return 0, errors.New("no reference for type")
		}
		return values[0], nil
	}
	return 0, errors.New("no reference for type")
}

func (anEntity *Entity) GetFloatPropertyValues(typeURI string) ([]float64, error) {
	if values, found := anEntity.Properties[typeURI]; found {
		switch v := values.(type) {
		case []float64:
			return v, nil
		case float64:
			result := make([]float64, 1)
			result[0] = v
			return result, nil
		}
	}
	return nil, errors.New("no property int32 literal")
}
