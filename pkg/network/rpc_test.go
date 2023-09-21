package network

import (
	"fmt"
	"testing"
)

type TestService struct {
}

type testRequest1 struct {
	field1 string
	field2 int
}

func (t *TestService) Test1(req *testRequest1) error {
	fmt.Println(fmt.Sprintf("test1 methods, req:%v", req))
	return nil
}

type testRequest2 struct {
	field1 string
	field2 int
	field3 map[string]int
}

func (t *TestService) test2(req *testRequest2) error {
	fmt.Println(fmt.Sprintf("test1 methods, req:%v", req))
	return nil
}

func TestNewService(t *testing.T) {
	var testService TestService
	s := newService(&testService)
	if s.name != "TestService" {
		t.Errorf("service name mismatch:%s", s.name)
	}
	if len(s.methods) != 1 {
		t.Errorf("service method mismatch")
	}
}

func TestServiceCall(t *testing.T) {

}
