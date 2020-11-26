package webpack

import (
	"fmt"
	"reflect"
	"testing"
)

const (
	checkMark = "\u2713" //勾号
	ballotX   = "\u2717" // 叉号
)

func newRouterTester() *router {
	r := newRouter()
	r.addRoute("GET", "/", nil)
	r.addRoute("GET", "/hello/:name", nil)
	r.addRoute("GET", "/hello/b/c", nil)
	r.addRoute("GET", "/world/:name", nil)
	r.addRoute("GET", "/assets/*filepath", nil)
	return r
}

func TestParsePattern(t *testing.T) {
	t.Log("Testing parsePattern method for router.go:")
	{
		res := parsePattern("/p/:name")
		if reflect.DeepEqual(res, []string{"p", ":name"}) {
			t.Log("\tParsing \"/p/:name\" : ", checkMark)
		} else {
			t.Error("\tParsing \"/p/:name\" : ", ballotX)
			t.Errorf("\t\tYour Parsing result is %v\n", res)
			t.Errorf("\t\tIt should be %v\n", []string{"p", ":name"})
		}

		res = parsePattern("/p/*")
		if reflect.DeepEqual(res, []string{"p", "*"}) {
			t.Log("\tParsing \"/p/*\" : ", checkMark)
		} else {
			t.Error("\tParsing \"/p/*\" : ", ballotX)
			t.Errorf("\t\tYour Parsing result is %v\n", res)
			t.Errorf("\t\tIt should be %v\n", []string{"p", "*"})
		}

		res = parsePattern("/p/*name/*")
		if reflect.DeepEqual(res, []string{"p", "*name"}) {
			t.Log("\tParsing \"/p/*\" : ", checkMark)
		} else {
			t.Error("\tParsing \"/p/*\" : ", ballotX)
			t.Errorf("\t\tYour Parsing result is %v\n", res)
			t.Errorf("\t\tIt should be %v\n", []string{"p", "*name"})
		}
	}
}

func TestGetRoute(t *testing.T) {
	t.Log("Testing getRoute method for router.go:")
	{
		r := newRouterTester()
		t.Log("Testing getRoute(\"GET\",\"/hello/Phoenix\")(*node, map[string]string)")
		n, ps := r.getRoute("GET", "/hello/Phoenix")
		if n == nil {
			t.Fatal("\tnode should not be nil:", ballotX)
		}
		t.Log("\tnode should not be nil:", checkMark)

		if n.pattern != "/hello/:name" || ps["name"] != "Phoenix" {
			t.Fatal("\tnode.pattern == \"/hello/:name\", ps[\"name\"] == \"Phoenix\"", ballotX)
		}

		t.Log("\tnode.pattern == \"/hello/:name\", ps[\"name\"] == \"Phoenix\"", checkMark)

		t.Log("Testing getRoute(\"GET\",\"/assets/picture1.jpg\")")
		n1, ps1 := r.getRoute("GET", "/assets/picture1.jpg")
		if n1.pattern != "/assets/*filepath" || ps1["filepath"] != "picture1.jpg" {
			t.Fatal("\tnode.pattern == \"/assets/*filepath\", ps[\"filepath\"] == \"picture1.jpg\"", ballotX)
		}
		t.Log("\tnode.pattern == \"/assets/*filepath\", ps[\"filepath\"] == \"picture1.jpg\"", checkMark)
	}
}

func TestGetRoutes(t *testing.T) {
	r := newRouterTester()
	nodes := r.getRoutes("GET")
	for i, n := range nodes {
		fmt.Println(i+1, n)
	}
	if len(nodes) != 5 {
		t.Fatal("the number of routes shoule be 4")
	}
}
