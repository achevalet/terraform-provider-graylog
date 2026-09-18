package notification

import "testing"

func TestSchemaDiffSuppressConfig(t *testing.T) {
	t.Parallel()
	data := []struct {
		title  string
		server string
		config string
		exp    bool
	}{
		{"key the config never set", `{"sender":"a","include_event_procedure":false}`, `{"sender":"a"}`, true},
		{"recipients in another order", `{"email_recipients":["b","a"]}`, `{"email_recipients":["a","b"]}`, true},
		{"both at once", `{"email_recipients":["b","a"],"include_event_procedure":false}`, `{"email_recipients":["a","b"]}`, true},
		{"a declared value that differs", `{"sender":"a"}`, `{"sender":"b"}`, false},
		{"a recipient added", `{"email_recipients":["a"]}`, `{"email_recipients":["a","b"]}`, false},
		{"a key only the config sets", `{"sender":"a"}`, `{"sender":"a","time_zone":"UTC"}`, false},
	}
	for _, d := range data {
		d := d
		t.Run(d.title, func(t *testing.T) {
			t.Parallel()
			if got := SchemaDiffSuppressConfig("config", d.server, d.config, nil); got != d.exp {
				t.Fatalf("got %v, expected %v", got, d.exp)
			}
		})
	}
}
