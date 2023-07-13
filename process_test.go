package confo

import "testing"

func TestEnvFormat(t *testing.T) {
	cases := map[string]string{
		"Hello":        "HELLO",
		"HelloWorld":   "HELLO_WORLD",
		"helloWorld":   "HELLO_WORLD",
		"hello_world":  "HELLO_WORLD",
		"Hello-World":  "HELLO_WORLD",
		"-World":       "WORLD",
		"World__hello": "WORLD_HELLO",
		"__world__":    "WORLD",
		"HTTP":         "HTTP",
		"MySQL":        "MY_SQL",
	}

	for k, target := range cases {
		if s := envNameFormat(k); s != target {
			t.Fatalf("%s return %s, should be %s", k, s, target)
			t.Fail()
		}
	}
}

func TestArgFormat(t *testing.T) {
	cases := map[string]string{
		"Hello":        "hello",
		"HelloWorld":   "hello-world",
		"helloWorld":   "hello-world",
		"hello_world":  "hello-world",
		"Hello-World":  "hello-world",
		"-World":       "world",
		"World__hello": "world-hello",
		"__world__":    "world",
		"HTTP":         "http",
		"MySQL":        "my-sql",
	}

	for k, target := range cases {
		if s := argNameFormat(k); s != target {
			t.Fatalf("%s return %s, should be %s", k, s, target)
			t.Fail()
		}
	}
}
