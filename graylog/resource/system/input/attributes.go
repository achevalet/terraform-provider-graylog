package input

import (
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/suzuki-shunsuke/go-dataeq/dataeq"
)

// An encrypted attribute is reported as {"is_set": bool} on 6.x and
// {"encrypted_value": "", "salt": ""} on 7.x, and refused back on write.

func isEncryptedValueRead(v interface{}) bool {
	m, ok := v.(map[string]interface{})
	if !ok {
		return false
	}
	if _, ok := m["is_set"]; ok {
		return true
	}
	_, enc := m["encrypted_value"]
	_, salt := m["salt"]
	return enc && salt
}

func isEncryptedValueWrite(v interface{}) bool {
	m, ok := v.(map[string]interface{})
	if !ok {
		return false
	}
	for _, k := range []string{"set_value", "keep_value", "delete_value"} {
		if _, ok := m[k]; ok {
			return true
		}
	}
	return false
}

func stripEncryptedValueRead(attributes interface{}) {
	m, ok := attributes.(map[string]interface{})
	if !ok {
		return
	}
	for k, v := range m {
		if isEncryptedValueRead(v) {
			delete(m, k)
		}
	}
}

func SchemaDiffSuppressAttributes(k, oldV, newV string, d *schema.ResourceData) bool {
	var server, config map[string]interface{}
	if json.Unmarshal([]byte(oldV), &server) != nil || json.Unmarshal([]byte(newV), &config) != nil {
		b, err := dataeq.JSON.Equal([]byte(oldV), []byte(newV))
		return err == nil && b
	}

	for key, sv := range server {
		if !isEncryptedValueRead(sv) {
			continue
		}
		if cv, ok := config[key]; ok && isEncryptedValueWrite(cv) {
			continue
		}
		delete(server, key)
		delete(config, key)
	}
	// a fresh create reports nothing for the attribute at all, and a null the
	// server does not report back is simply unset
	for key, cv := range config {
		if _, ok := server[key]; ok {
			continue
		}
		if cv == nil || isEncryptedValueRead(cv) {
			delete(config, key)
		}
	}

	serverB, err := json.Marshal(server)
	if err != nil {
		return false
	}
	configB, err := json.Marshal(config)
	if err != nil {
		return false
	}
	b, err := dataeq.JSON.Equal(serverB, configB)
	return err == nil && b
}
