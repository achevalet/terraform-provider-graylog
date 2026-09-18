package request

// CreateEntity wraps data the way Graylog 7 expects on the endpoints that take
// a CreateEntityRequest; earlier versions take the entity itself.
func CreateEntity(serverMajor int, data map[string]interface{}) interface{} {
	if serverMajor >= 7 {
		return map[string]interface{}{"entity": data}
	}
	return data
}
