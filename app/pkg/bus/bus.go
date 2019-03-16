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
	msgTypeName := handlerType.In(1).Elem().Name()
	handlers[msgTypeName] = handler
}

func Dispatch(ctx context.Context, msg Msg) error {
	msgName := reflect.TypeOf(msg).Elem().Name()
	handler := handlers[msgName]
	if handler == nil {
		panic(fmt.Errorf("could not find handler for '%s'", msgName))
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
