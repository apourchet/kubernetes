package factory

import (
	"context"
	"fmt"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/apiserver/pkg/storage"
)

type printedStorage struct {
	iface storage.Interface
}

func (s *printedStorage) Versioner() storage.Versioner {
	return s.iface.Versioner()
}

func (s *printedStorage) Get(ctx context.Context, key string, resourceVersion string, out runtime.Object, ignoreNotFound bool) error {
	fmt.Printf("STORAGE_IN: Get: %+v, %+v, %+v, %+v\n", key, resourceVersion, out, ignoreNotFound)
	err := s.iface.Get(ctx, key, resourceVersion, out, ignoreNotFound)
	fmt.Printf("STORAGE_OUT: Get: %+v, %+v\n", out, err)
	return err
}

func (s *printedStorage) Create(ctx context.Context, key string, obj, out runtime.Object, ttl uint64) error {
	fmt.Printf("STORAGE_IN: Create: %+v, %+v, %+v, %+v\n", key, obj, out, ttl)
	err := s.iface.Create(ctx, key, obj, out, ttl)
	fmt.Printf("STORAGE_OUT: Create: %+v, %+v\n", out, err)
	return err
}

func (s *printedStorage) Delete(ctx context.Context, key string, out runtime.Object, preconditions *storage.Preconditions) error {
	fmt.Printf("STORAGE_IN: Delete: %+v, %+v, %+v\n", key, out, preconditions)
	err := s.iface.Delete(ctx, key, out, preconditions)
	fmt.Printf("STORAGE_OUT: Delete: %+v, %+v\n", out, err)
	return err
}

func (s *printedStorage) Watch(ctx context.Context, key string, resourceVersion string, p storage.SelectionPredicate) (watch.Interface, error) {
	fmt.Printf("STORAGE_IN: Watch: %+v, %+v, %+v\n", key, resourceVersion, p)
	watcher, err := s.iface.Watch(ctx, key, resourceVersion, p)
	return watcher, err
}

func (s *printedStorage) WatchList(ctx context.Context, key string, resourceVersion string, p storage.SelectionPredicate) (watch.Interface, error) {
	fmt.Printf("STORAGE_IN: WatchList: %+v, %+v, %+v\n", key, resourceVersion, p)
	watcher, err := s.iface.WatchList(ctx, key, resourceVersion, p)
	return watcher, err
}

func (s *printedStorage) GetToList(ctx context.Context, key string, resourceVersion string, pred storage.SelectionPredicate, out runtime.Object) error {
	fmt.Printf("STORAGE_IN: GetToList: %+v, %+v, %+v, %+v\n", key, resourceVersion, pred, out)
	err := s.iface.GetToList(ctx, key, resourceVersion, pred, out)
	fmt.Printf("STORAGE_OUT: GetToList: %+v, %+v\n", out, err)
	return err
}

func (s *printedStorage) List(ctx context.Context, key string, resourceVersion string, pred storage.SelectionPredicate, out runtime.Object) error {
	fmt.Printf("STORAGE_IN: List: %+v, %+v, %+v, %+v\n", key, resourceVersion, pred, out)
	err := s.iface.List(ctx, key, resourceVersion, pred, out)
	fmt.Printf("STORAGE_OUT: List: %+v, %+v\n", out, err)
	return err
}

func (s *printedStorage) GuaranteedUpdate(ctx context.Context, key string, out runtime.Object, ignoreNotFound bool,
	preconditions *storage.Preconditions, tryUpdate storage.UpdateFunc, suggestion ...runtime.Object) error {
	fmt.Printf("STORAGE_IN: Update: %+v, %+v, %+v, %+v, %+v\n", key, out, ignoreNotFound, preconditions, suggestion)
	err := s.iface.GuaranteedUpdate(ctx, key, out, ignoreNotFound, preconditions, tryUpdate, suggestion...)
	fmt.Printf("STORAGE_OUT: Update: %+v, %+v\n", out, err)
	return err
}

func (s *printedStorage) Count(key string) (int64, error) {
	fmt.Printf("STORAGE_IN: Count: %s", key)
	cnt, err := s.iface.Count(key)
	fmt.Printf("STORAGE_OUT: Count: %+v, %+v", cnt, err)
	return cnt, err
}
