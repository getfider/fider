package bus

import (
	"context"
	"fmt"
	"reflect"
	"sync"
)

type HandlerFunc interface{}
type Msg interface{}
type Event interface{}

type Service interface {
	Category() string
	Enabled() bool
	Init()
}

var handlers = make(map[string]HandlerFunc)
var listeners = make(map[string][]HandlerFunc)
var services = make([]Service, 0)
var busLock = &sync.RWMutex{}

func Register(svc Service) {
	busLock.Lock()
	defer busLock.Unlock()

	services = append(services, svc)
}

func Reset() {
	busLock.Lock()
	defer busLock.Unlock()
	services = make([]Service, 0)
	handlers = make(map[string]HandlerFunc)
	listeners = make(map[string][]HandlerFunc)
}

// Initializes the bus services that have been registered via bus.Register
// Services that set via Init(...services) are always registered (regardless of Enabled() function)
/// and have preference over services registered from bus.Register
func Init(forcedServices ...Service) {
	for _, svc := range forcedServices {
		svc.Init()
	}

	for _, svc := range services {
		if svc.Enabled() {
			svc.Init()
		}
	}
}

func AddHandler(handler HandlerFunc) {
	busLock.RLock()
	defer busLock.RUnlock()

	handlerType := reflect.TypeOf(handler)
	elem := handlerType.In(1).Elem()
	handlers[keyForElement(elem)] = handler
}

func AddEventListener(handler HandlerFunc) {
	busLock.RLock()
	defer busLock.RUnlock()

	handlerType := reflect.TypeOf(handler)
	elem := handlerType.In(1).Elem()
	eventName := keyForElement(elem)
	_, exists := listeners[eventName]
	if !exists {
		listeners[eventName] = make([]HandlerFunc, 0)
	}
	listeners[eventName] = append(listeners[eventName], handler)
}

func Dispatch(ctx context.Context, msg Msg) error {
	busLock.RLock()
	defer busLock.RUnlock()

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

func Publish(ctx context.Context, evt Event) {
	busLock.RLock()
	defer busLock.RUnlock()

	typeof := reflect.TypeOf(evt)
	if typeof.Kind() != reflect.Ptr {
		panic(fmt.Errorf("'%s' is not a pointer", keyForElement(typeof)))
	}
	elem := typeof.Elem()
	key := keyForElement(elem)
	eventListeners := listeners[key]

	var params = []reflect.Value{
		reflect.ValueOf(ctx),
		reflect.ValueOf(evt),
	}

	for _, evtListener := range eventListeners {
		reflect.ValueOf(evtListener).Call(params)
	}
}

func keyForElement(t reflect.Type) string {
	msgTypeName := t.Name()
	pkgPath := t.PkgPath()
	return pkgPath + "." + msgTypeName
}
