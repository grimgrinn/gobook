package sexpr

import (
	"bytes"
	"testing"
)

type Person struct {
	Name string
	Age  int
	Tags []string
}

func TestMarshalUnmarshal(t *testing.T) {
	p := Person{
		Name: "Alice",
		Age:  30,
		Tags: []string{"go", "programming"},
	}

	data, err := Marshal(p)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("S-выражение: %s", data)

	var decoded Person
	if err := Unmarshal(data, &decoded); err != nil {
		t.Fatal(err)
	}

	if decoded.Name != p.Name {
		t.Errorf("Name: got %q, want %q", decoded.Name, p.Name)
	}

	if decoded.Age != p.Age {
		t.Errorf("Age: got %d, want %d", decoded.Age, p.Age)
	}

	if len(decoded.Tags) != len(p.Tags) {
		t.Errorf("Tags length: got %d, want %d", len(decoded.Tags), len(p.Tags))
	}
}

func TestDecoder(t *testing.T) {
	s := `((Name "Bob") (Age 25) (Tags ("rust" "haskell")))`
	r := bytes.NewReader([]byte(s))

	var p Person
	dec := NewDecoder(r)
	if err := dec.Decode(&p); err != nil {
		t.Fatal(err)
	}

	if p.Name != "Bob" {
		t.Errorf("Name: got %q, want %q", p.Name, "Bob")
	}
	if p.Age != 25 {
		t.Errorf("Age: got %d, want %d", p.Age, 25)
	}
}
