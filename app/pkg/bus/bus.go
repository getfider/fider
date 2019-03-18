package bus

import (
	"context"
	"fmt"
	"reflect"
)

type HandlerFunc interface{}
type Msg interface{}

type Service interface {
	Category() string
	Enabled() bool
	Init()
}

var services = make(map[string]Service)

func Register(svc Service) {
	services[svc.Category()] = svc
}

func Reset() {
	services = make(map[string]Service)
	handlers = make(map[string]HandlerFunc)
}

// Initializes the bus services that have been registered via bus.Register
// Services that set via Init(...services) are always registered (regardless of Enabled() function)
/// and have preference over services registered from bus.Register
func Init(forcedServices ...Service) {
	var initializedServices = make(map[string]bool)
	for _, svc := range forcedServices {
		svc.Init()
		initializedServices[svc.Category()] = true
	}

	for category, svc := range services {
		_, found := initializedServices[category]
		if !found && svc.Enabled() {
			svc.Init()
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
