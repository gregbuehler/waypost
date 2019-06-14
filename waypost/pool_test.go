package waypost

import "testing"

func TestPool(t *testing.T) {
	p := InitPool(nil)

	// a := p.Select()
	// if a != "" {
	// 	t.Errorf("Getting 'a'. Expected `nil`, found `%s`", a)
	// }

	p.Set("b")
	b := p.Select()
	if b != "b" {
		t.Errorf("Getting 'b'. Expected `b`, found `%s`", b)
	}

	p.Set("bar")
	p.Set("baz")
	c := p.List()
	if len(c) != 3 {
		t.Errorf("Getting 'c'. Expected 3 elements, found %d. result: %v", len(c), c)
	}

	p.Unset("baz")
	d := p.List()
	if len(d) != 2 {
		t.Errorf("Getting 'd'. Expected 2 elements, found %d. result: %v", len(d), d)
	}

	// TODO: more tests
}
