package network

import (
	"go/ast"
	"reflect"
)

type method struct {
	methodFunc reflect.Method
	ReqType    reflect.Type
}

func (m *method) newReq() reflect.Value {
	var req reflect.Value
	if m.ReqType.Kind() == reflect.Ptr {
		req = reflect.New(m.ReqType.Elem())
	} else {
		req = reflect.New(m.ReqType).Elem()
	}

	return req
}

type service struct {
	name     string
	typ      reflect.Type
	receiver reflect.Value
	methods  map[string]*method
}

func newService(receiver interface{}) *service {
	s := service{}
	s.receiver = reflect.ValueOf(receiver)
	s.name = reflect.Indirect(s.receiver).Type().Name()
	s.typ = reflect.TypeOf(receiver)
	s.registerMethods()

	return &s
}

func (s *service) registerMethods() {
	s.methods = make(map[string]*method)
	for i := 0; i < s.typ.NumMethod(); i++ {
		m := s.typ.Method(i)
		methodType := m.Type
		//TODO check methods format

		reqType := methodType.In(1)
		if !(ast.IsExported(reqType.Name()) || reqType.PkgPath() == "") {
			continue
		}
		//logutil.Logger.Debugf("register methods:%s of service %s", methods.Name, s.name)

		s.methods[m.Name] = &method{
			methodFunc: m,
			ReqType:    reqType,
		}
	}
}

func (s *service) call(m *method, req reflect.Value) error {
	f := m.methodFunc.Func
	retValues := f.Call([]reflect.Value{s.receiver, req})
	if err := retValues[0].Interface(); err != nil {
		return err.(error)
	}

	return nil
}
