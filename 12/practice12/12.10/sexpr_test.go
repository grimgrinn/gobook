package sexpr

import (
	"reflect"
	"testing"
)

func TestBoolean(t *testing.T) {
	var b bool

	data, _ := Marshal(true)

	if err := Unmarshal(data, &b); err != nil {
		t.Fatal(err)
	}
	if !b {
		t.Errorf("expected true, got %v", b)
	}

	data, _ = Marshal(false)
	if err := Unmarshal(data, &b); err != nil {
		t.Fatal(err)
	}
	if b {
		t.Errorf("expected false, got %v", b)
	}
}

func TestFloat(t *testing.T) {
	var f float64
	data, _ := Marshal(3.14)
	if err := Unmarshal(data, &f); err != nil {
		t.Fatal(err)
	}
	if f != 3.14 {
		t.Errorf("exptected 3.14, got %v", f)
	}
}

type MyType struct {
	Value int
}

func TestInterface(t *testing.T) {
	RegisterType("main.MyType", reflect.TypeOf(MyType{}))

	var result MyType

	data, err := Marshal(MyType{Value: 42})
	if err != nil {
		t.Fatal(err)
	}

	if err := Unmarshal(data, &result); err != nil {
		t.Fatal(err)
	}

	if result.Value != 42 {
		t.Errorf("expected 42, got %d", result.Value)
	}
}
