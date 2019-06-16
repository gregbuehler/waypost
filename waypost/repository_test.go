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

	lFile := `# list file
	foo.com	# inline note
	bar.com
	# baz.com # commented line
	`

	listTest := InitRepository(0)
	listCount, err := listTest.loadFromString(lFile)
	list1v, list1e := listTest.Get("foo.com")
	list2v, list2e := listTest.Get("bar.com")
	list3v, list3e := listTest.Get("baz.com")
	if err != nil {
		t.Errorf("Getting `err`, Expected `nil`, found `%v`", err)
	}
	if listCount != 2 {
		t.Errorf("Getting `listCount`, Expected `2`, found `%d`", listCount)
	}
	if list1e == false || list2e == false {
		t.Errorf("Getting 'list1e', Expected `true`, found `%t`", list1e)
		t.Errorf("Getting 'list2e', Expected `true`, found `%t`", list2e)
	}
	if list3e == true {
		t.Errorf("Getting 'list3e', Expected `false`, found `%t`", list3e)
	}
	if list1v != nil || list2v != nil || list3v != nil {
		t.Errorf("Getting list values, Expected `nil`, found [%t %t %t]", list1v, list2v, list3v)
	}

	hFile := `# hosts file
	127.0.0.1 foo.com
	1.2.3.4		bar.com # inline note
	# baz.com
	# biz.com #commented inline
	`

	hostTest := InitRepository(0)
	hostCount, err := hostTest.loadFromString(hFile)
	host1v, host1e := hostTest.Get("foo.com")
	host2v, host2e := hostTest.Get("bar.com")
	host3v, host3e := hostTest.Get("baz.com")
	host4v, host4e := hostTest.Get("biz.com")
	if err != nil {
		t.Errorf("Getting `err`, Expected `nil`, found `%v`", err)
	}
	if hostCount != 2 {
		t.Errorf("Getting `hostCount`, Expected `2`, found `%d`", hostCount)
	}
	if host1e == false || host2e == false {
		t.Errorf("Getting 'host1e', Expected `true`, found `%t`", host1e)
		t.Errorf("Getting 'host2e', Expected `true`, found `%t`", host2e)
	}
	if host3e != false || host4e != false {
		t.Errorf("Getting 'host3e', Expected `false`, found `%t`", host3e)
		t.Errorf("Getting 'host4e', Expected `false`, found `%t`", host4e)
	}
	if host1v != "127.0.0.1" || host2v != "1.2.3.4" {
		t.Errorf("Getting 'host1v', Expected `127.0.0.1`, found `%v`", host1v)
		t.Errorf("Getting 'host2v', Expected `1.2.3.4`, found `%v`", host2v)
	}
	if host3v != nil || host4v != nil {
		t.Errorf("Getting `host3v`, Expected `nil`, found `%v`", host3v)
		t.Errorf("Getting `host4v`, Expected `nil`, found `%v`", host4v)
	}

	// TODO: more tests
	// Fetch
}
