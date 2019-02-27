package factory

import (
	"context"
	"errors"
	"fmt"
	"path"
	"reflect"
	"strings"
	"sync"

	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/conversion"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/apiserver/pkg/storage"
	"k8s.io/apiserver/pkg/storage/etcd"
	"k8s.io/apiserver/pkg/storage/storagebackend"
	"k8s.io/apiserver/pkg/storage/value"
)

type authenticatedDataString string

// AuthenticatedData implements the value.Context interface.
func (d authenticatedDataString) AuthenticatedData() []byte {
	return []byte(string(d))
}

var _ value.Context = authenticatedDataString("")

type pelotonStorage struct {
	sync.Mutex

	versioner   storage.Versioner
	codec       runtime.Codec
	transformer value.Transformer
	pathPrefix  string

	objs map[string][]byte
}

type pelotonWatchObj struct {
	results chan watch.Event
}

func (obj *pelotonWatchObj) Stop() {}

func (obj *pelotonWatchObj) ResultChan() <-chan watch.Event {
	return obj.results
}

type storageValue struct {
	content []byte
	version int64
}

func fakeStorage(c storagebackend.Config) storage.Interface {
	transformer := c.Transformer
	if transformer == nil {
		transformer = value.IdentityTransformer
	}
	return &pelotonStorage{
		versioner:   etcd.APIObjectVersioner{},
		codec:       c.Codec,
		transformer: transformer,
		pathPrefix:  c.Prefix,
		objs:        map[string][]byte{},
	}
}

func (s *pelotonStorage) Versioner() storage.Versioner {
	return s.versioner
}

// Get unmarshals json found at key into objPtr. On a not found error, will either
// return a zero object of the requested type, or an error, depending on ignoreNotFound.
// Treats empty responses and nil response nodes exactly like a not found error.
// The returned contents may be delayed, but it is guaranteed that they will
// be have at least 'resourceVersion'.
func (s *pelotonStorage) Get(ctx context.Context, key string, resourceVersion string, objPtr runtime.Object, ignoreNotFound bool) error {
	fmt.Printf("STORAGE_IN: Get: %+v, %+v, %+v, %+v\n", key, resourceVersion, objPtr, ignoreNotFound)
	key = path.Join(s.pathPrefix, key)
	s.Lock()
	defer s.Unlock()
	content, found := s.objs[key]
	if !found {
		if ignoreNotFound {
			return runtime.SetZeroValue(objPtr)
		}
		return storage.NewKeyNotFoundError(key, 0)
	}

	data, _, err := s.transformer.TransformFromStorage(content, authenticatedDataString(key))
	if err != nil {
		return storage.NewInternalError(err.Error())
	}

	return runtime.DecodeInto(s.codec, data, objPtr)
}

// Create adds a new object at a key unless it already exists. 'ttl' is time-to-live
// in seconds (0 means forever). If no error is returned and out is not nil, out will be
// set to the read value from database.
func (s *pelotonStorage) Create(ctx context.Context, key string, obj, out runtime.Object, ttl uint64) error {
	fmt.Printf("STORAGE_IN: Create: %+v, %+v, %+v, %+v\n", key, obj, out, ttl)
	if version, err := s.versioner.ObjectResourceVersion(obj); err == nil && version != 0 {
		return errors.New("resourceVersion should not be set on objects to be created")
	}
	if err := s.versioner.PrepareObjectForStorage(obj); err != nil {
		return fmt.Errorf("PrepareObjectForStorage failed: %v", err)
	}
	data, err := runtime.Encode(s.codec, obj)
	if err != nil {
		return err
	}
	key = path.Join(s.pathPrefix, key)

	newData, err := s.transformer.TransformToStorage(data, authenticatedDataString(key))
	if err != nil {
		return storage.NewInternalError(err.Error())
	}

	s.Lock()
	defer s.Unlock()
	if _, found := s.objs[key]; found {
		return storage.NewKeyExistsError(key, 0)
	}
	fmt.Println("STORAGE_DEBUG:", key)
	s.objs[key] = newData

	if out != nil {
		return runtime.DecodeInto(s.codec, newData, out)
	}
	return nil
}

// Delete removes the specified key and returns the value that existed at that spot.
// If key didn't exist, it will return NotFound storage error.
func (s *pelotonStorage) Delete(ctx context.Context, key string, out runtime.Object, preconditions *storage.Preconditions) error {
	fmt.Printf("STORAGE_IN: Delete: %+v, %+v, %+v\n", key, out, preconditions)
	_, err := conversion.EnforcePtr(out)
	if err != nil {
		panic("unable to convert output object to pointer")
	}
	key = path.Join(s.pathPrefix, key)

	s.Lock()
	defer s.Unlock()
	delete(s.objs, key)
	return nil
}

// Watch begins watching the specified key. Events are decoded into API objects,
// and any items selected by 'p' are sent down to returned watch.Interface.
// resourceVersion may be used to specify what version to begin watching,
// which should be the current resourceVersion, and no longer rv+1
// (e.g. reconnecting without missing any updates).
// If resource version is "0", this interface will get current object at given key
// and send it in an "ADDED" event, before watch starts.
func (s *pelotonStorage) Watch(ctx context.Context, key string, resourceVersion string, p storage.SelectionPredicate) (watch.Interface, error) {
	fmt.Printf("STORAGE_IN: Watch: %+v, %+v, %+v\n", key, resourceVersion, p)
	return &pelotonWatchObj{make(chan watch.Event)}, nil
}

// WatchList begins watching the specified key's items. Items are decoded into API
// objects and any item selected by 'p' are sent down to returned watch.Interface.
// resourceVersion may be used to specify what version to begin watching,
// which should be the current resourceVersion, and no longer rv+1
// (e.g. reconnecting without missing any updates).
// If resource version is "0", this interface will list current objects directory defined by key
// and send them in "ADDED" events, before watch starts.
func (s *pelotonStorage) WatchList(ctx context.Context, key string, resourceVersion string, p storage.SelectionPredicate) (watch.Interface, error) {
	fmt.Printf("STORAGE_IN: WatchList: %+v, %+v, %+v\n", key, resourceVersion, p)
	return &pelotonWatchObj{make(chan watch.Event)}, nil
}

// GetToList unmarshals json found at key and opaque it into *List api object
// (an object that satisfies the runtime.IsList definition).
// The returned contents may be delayed, but it is guaranteed that they will
// be have at least 'resourceVersion'.
func (s *pelotonStorage) GetToList(ctx context.Context, key string, resourceVersion string, pred storage.SelectionPredicate, listObj runtime.Object) error {
	fmt.Printf("STORAGE_IN: GetToList: %+v, %+v, %+v, %+v\n", key, resourceVersion, pred, listObj)
	listPtr, err := meta.GetItemsPtr(listObj)
	if err != nil {
		return err
	}
	v, err := conversion.EnforcePtr(listPtr)
	if err != nil || v.Kind() != reflect.Slice {
		panic("need ptr to slice")
	}

	key = path.Join(s.pathPrefix, key)

	s.Lock()
	defer s.Unlock()
	content, found := s.objs[key]
	if !found {
		return nil
	}

	data, _, err := s.transformer.TransformFromStorage(content, authenticatedDataString(key))
	if err != nil {
		return storage.NewInternalError(err.Error())
	}

	obj, _, err := s.codec.Decode(data, nil, reflect.New(v.Type().Elem()).Interface().(runtime.Object))
	if err != nil {
		return err
	}
	if matched, err := pred.Matches(obj); err == nil && matched {
		v.Set(reflect.Append(v, reflect.ValueOf(obj).Elem()))
	}
	return nil
}

// List unmarshalls jsons found at directory defined by key and opaque them
// into *List api object (an object that satisfies runtime.IsList definition).
// The returned contents may be delayed, but it is guaranteed that they will
// be have at least 'resourceVersion'.
func (s *pelotonStorage) List(ctx context.Context, key string, resourceVersion string, pred storage.SelectionPredicate, listObj runtime.Object) error {
	fmt.Printf("STORAGE_IN: List: %+v, %+v, %+v, %+v\n", key, resourceVersion, pred, listObj)
	listPtr, err := meta.GetItemsPtr(listObj)
	if err != nil {
		return storage.NewInternalErrorf(err.Error())
	}
	v, err := conversion.EnforcePtr(listPtr)
	if err != nil || v.Kind() != reflect.Slice {
		panic("need ptr to slice")
	}

	if s.pathPrefix != "" {
		key = path.Join(s.pathPrefix, key)
	}
	// We need to make sure the key ended with "/" so that we only get children "directories".
	// e.g. if we have key "/a", "/a/b", "/ab", getting keys with prefix "/a" will return all three,
	// while with prefix "/a/" will return only "/a/b" which is the correct answer.
	if !strings.HasSuffix(key, "/") {
		key += "/"
	}

	s.Lock()
	defer s.Unlock()
	for k, content := range s.objs {
		fmt.Println("STORAGE_DEBUG:", k, key)
		if !strings.HasPrefix(k, key) {
			continue
		}
		fmt.Printf("STORAGE_DEBUG: ADDING %s\n", k)
		obj, _, err := s.codec.Decode(content, nil, reflect.New(v.Type().Elem()).Interface().(runtime.Object))
		if err != nil {
			fmt.Printf("STORAGE_ERROR: %+v, %+v\n", k, err)
			return storage.NewInternalErrorf(err.Error())
		}

		fmt.Printf("STORAGE_DEBUG: %+v, %+v", k, obj)
		if matched, err := pred.Matches(obj); err == nil && matched {
			v.Set(reflect.Append(v, reflect.ValueOf(obj).Elem()))
		} else if err != nil {
			fmt.Println("STORAGE_ERR:", err)
		} else if !matched {
			fmt.Println("STORAGE: NOT MATCHED")
		}
	}

	fmt.Printf("STORAGE_OUT: List: %+v, %+v\n", key, listObj)
	return nil
}

// GuaranteedUpdate keeps calling 'tryUpdate()' to update key 'key' (of type 'ptrToType')
// retrying the update until success if there is index conflict.
// Note that object passed to tryUpdate may change across invocations of tryUpdate() if
// other writers are simultaneously updating it, so tryUpdate() needs to take into account
// the current contents of the object when deciding how the update object should look.
// If the key doesn't exist, it will return NotFound storage error if ignoreNotFound=false
// or zero value in 'ptrToType' parameter otherwise.
// If the object to update has the same value as previous, it won't do any update
// but will return the object in 'ptrToType' parameter.
// If 'suggestion' can contain zero or one element - in such case this can be used as
// a suggestion about the current version of the object to avoid read operation from
// storage to get it.
//
// Example:
//
// s := /* implementation of Interface */
// err := s.GuaranteedUpdate(
//     "myKey", &MyType{}, true,
//     func(input runtime.Object, res ResponseMeta) (runtime.Object, *uint64, error) {
//       // Before each incovation of the user defined function, "input" is reset to
//       // current contents for "myKey" in database.
//       curr := input.(*MyType)  // Guaranteed to succeed.
//
//       // Make the modification
//       curr.Counter++
//
//       // Return the modified object - return an error to stop iterating. Return
//       // a uint64 to alter the TTL on the object, or nil to keep it the same value.
//       return cur, nil, nil
//    }
// })
func (s *pelotonStorage) GuaranteedUpdate(
	ctx context.Context, key string, ptrToType runtime.Object, ignoreNotFound bool,
	precondtions *storage.Preconditions, tryUpdate storage.UpdateFunc, suggestion ...runtime.Object) error {
	return nil
}

// Count returns number of different entries under the key (generally being path prefix).
func (s *pelotonStorage) Count(key string) (int64, error) {
	s.Lock()
	defer s.Unlock()
	count := int64(0)
	for k := range s.objs {
		if strings.HasPrefix(k, key) {
			count++
		}
	}
	fmt.Printf("STORAGE_OUT: Count: %+v, %+v\n", key, count)
	return count, nil
}
