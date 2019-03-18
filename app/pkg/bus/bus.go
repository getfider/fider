package bus

import (
	"context"
	"fmt"
	"reflect"
)

type HandlerFunc interface{}
type Msg interface{}

type Service interface {
	Enabled() bool
	Init()
}

var services = make([]Service, 0)

func Register(s Service) {
	services = append(services, s)
}

func Reset() {
	services = make([]Service, 0)
}

func Init() {
	for _, s := range services {
		if s.Enabled() {
			s.Init()
		}
	}
}

var handlers = make(map[string]HandlerFunc)

func AddHandler(s Service, handler HandlerFunc) {
	handlerType := reflect.TypeOf(handler)
	elem := handlerType.In(1).Elem()
	handlers[keyForElement(elem)] = handler
}

func Dispatch(ctx context.Context, msg Msg) error {
	typeof := reflect.TypeOf(msg)
	if typeof.Kind() != reflect.Ptr {
		panic(fmt.Errorf("'%s' is not a pointer", keyForElement(typeof)))
	}
	elem := typeof.Elem()
	key := keyForElement(elem)
	handler := handlers[key]
	if handler == nil {
		panic(fmt.Errorf("could not find handler for '%s'", key))
	}

	var params = []reflect.Value{
		reflect.ValueOf(ctx),
		reflect.ValueOf(msg),
	}

	ret := reflect.ValueOf(handler).Call(params)
	err := ret[0].Interface()
	if err == nil {
		return nil
	}
	return err.(error)
}

func keyForElement(t reflect.Type) string {
	msgTypeName := t.Name()
	pkgPath := t.PkgPath()
	return pkgPath + "." + msgTypeName
}
