package webpack

import "testing"

func TestNestedGroup(t *testing.T) {
	t.Log("Testing Group method for webpack.go")
	r := New()
	v1 := r.Group("/v1")
	v2 := v1.Group("/v2")
	v3 := v2.Group("/v3")
	if v2.prefix != "/v1/v2" {
		t.Log("\tTestNestedGroup:", ballotX)
		t.Fatal("\t\tv2 prefix should be /v1/v2")
	}
	if v3.prefix != "/v1/v2/v3" {
		t.Log("\tTestNestedGroup:", ballotX)
		t.Fatal("\t\tv3 prefix should be /v1/v2/v3")
	}
	t.Log("\tTestNestedGroup:", checkMark)
}
