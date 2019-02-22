package factory

import (
	"context"
	"fmt"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/apiserver/pkg/storage"
	"k8s.io/apiserver/pkg/storage/etcd"
)

type pelotonStorage struct{}

type pelotonWatchObj struct{}

func (obj pelotonWatchObj) Stop() {}

func (obj pelotonWatchObj) ResultChan() <-chan watch.Event {
	return make(chan watch.Event, 0)
}

func fakeStorage() storage.Interface {
	return pelotonStorage{}
}

func (s pelotonStorage) Versioner() storage.Versioner {
	return etcd.APIObjectVersioner{}
}

// Create adds a new object at a key unless it already exists. 'ttl' is time-to-live
// in seconds (0 means forever). If no error is returned and out is not nil, out will be
// set to the read value from database.
func (s pelotonStorage) Create(ctx context.Context, key string, obj, out runtime.Object, ttl uint64) error {
	fmt.Printf("STORAGE: Create: %+v, %+v, %+v, %+v\n", key, obj, out, ttl)
	return nil
}

// Delete removes the specified key and returns the value that existed at that spot.
// If key didn't exist, it will return NotFound storage error.
func (s pelotonStorage) Delete(ctx context.Context, key string, out runtime.Object, preconditions *storage.Preconditions) error {
	fmt.Printf("STORAGE: Delete: %+v, %+v, %+v\n", key, out, preconditions)
	return nil
}

// Watch begins watching the specified key. Events are decoded into API objects,
// and any items selected by 'p' are sent down to returned watch.Interface.
// resourceVersion may be used to specify what version to begin watching,
// which should be the current resourceVersion, and no longer rv+1
// (e.g. reconnecting without missing any updates).
// If resource version is "0", this interface will get current object at given key
// and send it in an "ADDED" event, before watch starts.
func (s pelotonStorage) Watch(ctx context.Context, key string, resourceVersion string, p storage.SelectionPredicate) (watch.Interface, error) {
	fmt.Printf("STORAGE: Watch: %+v, %+v, %+v\n", key, resourceVersion, p)
	return pelotonWatchObj{}, nil
}

// WatchList begins watching the specified key's items. Items are decoded into API
// objects and any item selected by 'p' are sent down to returned watch.Interface.
// resourceVersion may be used to specify what version to begin watching,
// which should be the current resourceVersion, and no longer rv+1
// (e.g. reconnecting without missing any updates).
// If resource version is "0", this interface will list current objects directory defined by key
// and send them in "ADDED" events, before watch starts.
func (s pelotonStorage) WatchList(ctx context.Context, key string, resourceVersion string, p storage.SelectionPredicate) (watch.Interface, error) {
	fmt.Printf("STORAGE: WatchList: %+v, %+v, %+v\n", key, resourceVersion, p)
	return pelotonWatchObj{}, nil
}

// Get unmarshals json found at key into objPtr. On a not found error, will either
// return a zero object of the requested type, or an error, depending on ignoreNotFound.
// Treats empty responses and nil response nodes exactly like a not found error.
// The returned contents may be delayed, but it is guaranteed that they will
// be have at least 'resourceVersion'.
func (s pelotonStorage) Get(ctx context.Context, key string, resourceVersion string, objPtr runtime.Object, ignoreNotFound bool) error {
	fmt.Printf("STORAGE: Get: %+v, %+v, %+v, %+v\n", key, resourceVersion, objPtr, ignoreNotFound)
	return nil
}

// GetToList unmarshals json found at key and opaque it into *List api object
// (an object that satisfies the runtime.IsList definition).
// The returned contents may be delayed, but it is guaranteed that they will
// be have at least 'resourceVersion'.
func (s pelotonStorage) GetToList(ctx context.Context, key string, resourceVersion string, p storage.SelectionPredicate, listObj runtime.Object) error {
	fmt.Printf("STORAGE: GetToList: %+v, %+v, %+v, %+v\n", key, resourceVersion, p, listObj)
	return nil
}

// List unmarshalls jsons found at directory defined by key and opaque them
// into *List api object (an object that satisfies runtime.IsList definition).
// The returned contents may be delayed, but it is guaranteed that they will
// be have at least 'resourceVersion'.
func (s pelotonStorage) List(ctx context.Context, key string, resourceVersion string, p storage.SelectionPredicate, listObj runtime.Object) error {
	fmt.Printf("STORAGE: List: %+v, %+v, %+v, %+v\n", key, resourceVersion, p, listObj)
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
func (s pelotonStorage) GuaranteedUpdate(
	ctx context.Context, key string, ptrToType runtime.Object, ignoreNotFound bool,
	precondtions *storage.Preconditions, tryUpdate storage.UpdateFunc, suggestion ...runtime.Object) error {
	return nil
}

// Count returns number of different entries under the key (generally being path prefix).
func (s pelotonStorage) Count(key string) (int64, error) {
	fmt.Printf("STORAGE: Count: %+v\n", key)
	return 0, nil
}
