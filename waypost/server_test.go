package waypost

import "testing"

func TestServer(t *testing.T) {

	a1, a2, authoritative := examineName("foo.com.")
	if a1 != "foo.com." {
		t.Errorf("Getting authoritative name. Expected `foo.com.`, got `%s`", a1)
	}
	if a2 != "foo.com" {
		t.Errorf("Getting simple name. Expected `foo.com`, got `%s`", a2)
	}
	if !authoritative {
		t.Errorf("Getting authoritative quality. Expected `true`, got `%t`", authoritative)
	}

	b1, b2, authoritative := examineName("foo.com")
	if b1 != "foo.com" {
		t.Errorf("Getting authoritative name. Expected `foo.com.`, got `%s`", b1)
	}
	if b2 != "foo.com" {
		t.Errorf("Getting simple name. Expected `foo.com`, got `%s`", b2)
	}
	if authoritative {
		t.Errorf("Getting authoritative quality. Expected `false`, got `%t`", authoritative)
	}

	// TODO: more tests
}
