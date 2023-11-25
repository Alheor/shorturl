// Package repository
// map implementation
package repository

import (
	"context"
	"errors"
	"sync"
)

// ShortNameMap struct
type ShortNameMap struct {
	URLMap map[string]string
	sync.RWMutex
}

func (sn *ShortNameMap) Init(ctx context.Context) error {

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	sn.URLMap = make(map[string]string)
	return nil
}

func (sn *ShortNameMap) Add(ctx context.Context, id string, value string) error {

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	sn.Lock()
	defer sn.Unlock()

	_, exists := sn.URLMap[id]
	if exists {
		return NewUniqueError(id, nil)
	}

	for _, mapValue := range sn.URLMap {
		if mapValue == value {
			return NewUniqueError(id, nil)
		}
	}

	sn.URLMap[id] = value

	return nil
}

func (sn *ShortNameMap) Get(ctx context.Context, id string) (value string, error error) {

	select {
	case <-ctx.Done():
		return ``, errors.New(ErrIDNotFound)
	default:
	}

	sn.RLock()
	defer sn.RUnlock()

	url, exists := sn.URLMap[id]
	if !exists {
		return ``, errors.New(ErrIDNotFound)
	}

	return url, nil
}

func (sn *ShortNameMap) Remove(ctx context.Context, id string) {

	select {
	case <-ctx.Done():
		return
	default:
	}

	sn.Lock()
	defer sn.Unlock()

	delete(sn.URLMap, id)
}

func (sn *ShortNameMap) IsReady(ctx context.Context) bool {

	select {
	case <-ctx.Done():
		return false
	default:
	}

	return sn.URLMap != nil
}

func (sn *ShortNameMap) AddBatch(ctx context.Context, in []BatchEl) error {

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	sn.Lock()
	defer sn.Unlock()

	for _, v := range in {

		_, exists := sn.URLMap[v.ShortURL]
		if exists {
			return errors.New(ErrValueAlreadyExist)
		}

		for _, mapValue := range sn.URLMap {
			if mapValue == v.OriginalURL {
				return errors.New(ErrValueAlreadyExist)
			}
		}

		sn.URLMap[v.ShortURL] = v.OriginalURL
	}

	return nil
}
