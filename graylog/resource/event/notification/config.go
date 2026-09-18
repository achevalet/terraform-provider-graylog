package notification

import (
	"encoding/json"
	"sort"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/suzuki-shunsuke/go-dataeq/dataeq"
)

// Graylog fills in config keys the configuration never set -- 7 adds
// include_event_procedure -- and returns the address lists in its own order
// without adopting the order it is sent. Only lists of strings are sorted.

// stripServerOnlyKeys drops what the configuration does not declare.
func stripServerOnlyKeys(server, config map[string]interface{}) {
	for k := range server {
		if _, ok := config[k]; !ok {
			delete(server, k)
		}
	}
}
func sortStringLists(m map[string]interface{}) {
	for _, v := range m {
		list, ok := v.([]interface{})
		if !ok {
			continue
		}
		strs := make([]string, 0, len(list))
		for _, e := range list {
			s, ok := e.(string)
			if !ok {
				break
			}
			strs = append(strs, s)
		}
		if len(strs) != len(list) {
			continue
		}
		sort.Strings(strs)
		for i, s := range strs {
			list[i] = s
		}
	}
}

func SchemaDiffSuppressConfig(k, oldV, newV string, d *schema.ResourceData) bool {
	var server, config map[string]interface{}
	if json.Unmarshal([]byte(oldV), &server) != nil || json.Unmarshal([]byte(newV), &config) != nil {
		b, err := dataeq.JSON.Equal([]byte(oldV), []byte(newV))
		return err == nil && b
	}
	stripServerOnlyKeys(server, config)
	sortStringLists(server)
	sortStringLists(config)

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
