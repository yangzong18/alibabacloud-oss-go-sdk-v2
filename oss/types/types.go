package types

type ApiMetadata struct {
	values map[interface{}]interface{}
}

func (m ApiMetadata) Get(key interface{}) interface{} {
	return m.values[key]
}

func (m ApiMetadata) Clone() ApiMetadata {
	vs := make(map[interface{}]interface{}, len(m.values))
	for k, v := range m.values {
		vs[k] = v
	}

	return ApiMetadata{
		values: vs,
	}
}

func (m *ApiMetadata) Set(key, value interface{}) {
	if m.values == nil {
		m.values = map[interface{}]interface{}{}
	}
	m.values[key] = value
}

func (m ApiMetadata) Has(key interface{}) bool {
	if m.values == nil {
		return false
	}
	_, ok := m.values[key]
	return ok
}
