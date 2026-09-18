package input

import "testing"

func TestSchemaDiffSuppressAttributes_unsetKeys(t *testing.T) {
	t.Parallel()
	data := []struct {
		title  string
		server string
		config string
		exp    bool
	}{
		{"null the server drops", `{"port":5044}`, `{"port":5044,"override_source":null}`, true},
		{"encrypted the server drops", `{"port":5044}`, `{"port":5044,"tls_key_password":{"is_set":false}}`, true},
		{"both at once", `{"port":5044}`, `{"port":5044,"override_source":null,"tls_key_password":{"is_set":false}}`, true},
		{"null the server does report", `{"override_source":"x"}`, `{"override_source":null}`, false},
		{"a real change alongside", `{"port":5044}`, `{"port":5045,"override_source":null}`, false},
	}
	for _, d := range data {
		d := d
		t.Run(d.title, func(t *testing.T) {
			t.Parallel()
			if got := SchemaDiffSuppressAttributes("attributes", d.server, d.config, nil); got != d.exp {
				t.Fatalf("got %v, expected %v", got, d.exp)
			}
		})
	}
}
