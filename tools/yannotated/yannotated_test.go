package main

import (
	"io/ioutil"
	"os"
	"testing"
)

// Config is a config.
type Config struct {
	// Foo is foo.
	Foo string `yaml:"foo"`
	// Bar is bar.
	// So useful.
	Bar Bar `yaml:"bar"`
	// Rebar is another bar.
	Rebar Bar `yaml:"rebar"`
	Quz   map[string]string
}

// Bar is a struct.
type Bar struct {
	// Baz is baz.
	Baz int `yaml:"baz"`
}

func TestMain(t *testing.T) {
	tmp, err := ioutil.TempFile("", "")
	if err != nil {
		t.Fatal(err)
	}
	tmp.Close()
	defer os.RemoveAll(tmp.Name())

	err = mainE(Flags{
		Dir:     ".",
		Package: "main",
		Type:    "Config",
		Output:  tmp.Name(),
	})
	if err != nil {
		t.Fatal(err)
	}

	want := `# Foo is foo.
foo: ""
# Bar is bar.
# So useful.
bar:
  # Baz is baz.
  baz: 0
# Rebar is another bar.
rebar:
  # Baz is baz.
  baz: 0
quz: {}
`
	b, err := ioutil.ReadFile(tmp.Name())
	if err != nil {
		t.Fatal(err)
	}

	if got := string(b); got != want {
		t.Fatalf("got:\n%s\nwant:\n%s", got, want)
	}
}

func TestGo(t *testing.T) {
	tmp, err := ioutil.TempFile("", "")
	if err != nil {
		t.Fatal(err)
	}
	tmp.Close()
	defer os.RemoveAll(tmp.Name())

	err = mainE(Flags{
		Dir:     ".",
		Package: "main",
		Type:    "Bar",
		Output:  tmp.Name(),
		Format:  GoFormat,
	})
	if err != nil {
		t.Fatal(err)
	}

	want := `package main

var yannotated = ` + "`" + `# Baz is baz.
baz: 0
` + "`\n"

	b, err := ioutil.ReadFile(tmp.Name())
	if err != nil {
		t.Fatal(err)
	}

	if got := string(b); got != want {
		t.Fatalf("got:\n%s\nwant:\n%s", got, want)
	}
}
