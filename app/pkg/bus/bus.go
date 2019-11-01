package bus

import (
	"context"
	"fmt"
	"reflect"
	"sync"

	"github.com/getfider/fider/app/models/cmd"
	"github.com/getfider/fider/app/pkg/errors"
)

type HandlerFunc interface{}
type Msg interface{}

type Service interface {
	Name() string
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
func Init(forcedServices ...Service) []Service {
	initializedServices := make([]Service, 0)
	for _, svc := range forcedServices {
		initializedServices = append(initializedServices, svc)
		svc.Init()
	}

	for _, svc := range services {
		if svc.Enabled() {
			initializedServices = append(initializedServices, svc)
			svc.Init()
		}
	}
	return initializedServices
}

func AddHandler(handler HandlerFunc) {
	busLock.Lock()
	defer busLock.Unlock()

	handlerType := reflect.TypeOf(handler)
	elem := handlerType.In(1).Elem()
	handlers[keyForElement(elem)] = handler
}

func AddListener(handler HandlerFunc) {
	busLock.Lock()
	defer busLock.Unlock()

	handlerType := reflect.TypeOf(handler)
	elem := handlerType.In(1).Elem()
	eventName := keyForElement(elem)
	_, exists := listeners[eventName]
	if !exists {
		listeners[eventName] = make([]HandlerFunc, 0)
	}
	listeners[eventName] = append(listeners[eventName], handler)
}

func MustDispatch(ctx context.Context, msgs ...Msg) {
	err := Dispatch(ctx, msgs...)
	if err != nil {
		panic(err)
	}
}

func Dispatch(ctx context.Context, msgs ...Msg) error {
	if len(msgs) == 0 {
		return nil
	}

	busLock.RLock()
	defer busLock.RUnlock()

	for _, msg := range msgs {
		typeof := reflect.TypeOf(msg)
		if typeof.Kind() != reflect.Ptr {
			panic(fmt.Errorf("'%s' is not a pointer", keyForElement(typeof)))
		}

		elem := typeof.Elem()
		key := keyForElement(elem)
		handler := handlers[key]
		if handler == nil {
			panic(fmt.Errorf("could not find handler for '%s'.", key))
		}

		var params = []reflect.Value{
			reflect.ValueOf(ctx),
			reflect.ValueOf(msg),
		}

		ret := reflect.ValueOf(handler).Call(params)
		if err := ret[0].Interface(); err != nil {
			return err.(error)
		}
	}

	return nil
}

func Publish(ctx context.Context, msgs ...Msg) {
	if len(msgs) == 0 {
		return
	}

	busLock.RLock()
	defer busLock.RUnlock()

	for _, msg := range msgs {
		typeof := reflect.TypeOf(msg)
		if typeof.Kind() != reflect.Ptr {
			panic(fmt.Errorf("'%s' is not a pointer", keyForElement(typeof)))
		}

		elem := typeof.Elem()
		key := keyForElement(elem)
		msgListeners := listeners[key]

		var params = []reflect.Value{
			reflect.ValueOf(ctx),
			reflect.ValueOf(msg),
		}

		for _, msgListener := range msgListeners {
			ret := reflect.ValueOf(msgListener).Call(params)
			if len(ret) > 0 {
				if err, isErr := ret[0].Interface().(error); isErr {
					Publish(ctx, &cmd.LogError{
						Err: errors.Wrap(err, "failed to execute msg '%s'", key),
					})
				}
			}
		}
	}
}

func keyForElement(t reflect.Type) string {
	msgTypeName := t.Name()
	pkgPath := t.PkgPath()
	return pkgPath + "." + msgTypeName
}
