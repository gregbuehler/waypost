package waypost

import (
	"testing"
)

func TestRepository(t *testing.T) {
	r := InitRepository(0)

	a, found := r.Get("a")
	if found || a != nil {
		t.Errorf("Getting 'a'. Expected `nil`, found `%s`", a)
	}

	r.Set("b", "bar")
	b, found := r.Get("b")
	if !found || b != "bar" {
		t.Errorf("Getting 'b'. Expected `bar`, unexpected value `%s` returned", b)
	}

	r.Set("c", "cool")
	r.Unset("c")
	c, found := r.Get("c")
	if found || c != nil {
		t.Errorf("Getting 'c'. Expected `nil`, found `%s`", c)
	}

	r.Set("foo", nil)
	r.Set("bar", nil)
	r.Set("baz", nil)
	d := r.List()
	if len(d) != 4 {
		t.Errorf("Getting 'd'. Expected 4 elements, found %d. result: %v", len(d), d)
	}

	r.Set("Fizz-Buzz", "asdf")
	e1, e2, e3 := r.Search("Fizz")
	if e1 != "Fizz-Buzz" {
		t.Errorf("Getting 'e1', Expected `Fizz-Buzz`, found `%s`", e1)
	}
	if e2.(string) != "asdf" {
		t.Errorf("Getting 'e2', Expected `asdf`, found `%s`", e2)
	}
	if e3 == false {
		t.Errorf("Getting 'e3', Expected `true`, found `%t`", e3)
	}
	_, _, f3 := r.Search("Fuzz")
	if f3 {
		t.Errorf("Getting 'e3', Expected `false`, found `%t`", f3)
	}
	// TODO: more tests
	// Fetch
}
