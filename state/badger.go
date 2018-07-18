// Copyright 2017, 2018 Tensigma Ltd. All rights reserved.
// Use of this source code is governed by Microsoft Reference Source
// License (MS-RSL) that can be found in the LICENSE file.

package state

import (
	"fmt"
	"time"

	"github.com/dgraph-io/badger"
)

// badgerStore implements IndexedStore.
type badgerStore struct {
	opts *storeOptions
	db   *badger.DB
}

func newBadgerStore(prefix string, opts ...storeOpt) (*badgerStore, error) {
	s := &badgerStore{
		opts: defaultStoreOptions(),
	}
	for _, o := range opts {
		if o != nil {
			o(s.opts)
		}
	}
	badgerOpts := badger.DefaultOptions
	badgerOpts.Dir = prefix
	badgerOpts.ValueDir = prefix
	badgerOpts.SyncWrites = s.opts.SyncWrites
	db, err := badger.Open(badgerOpts)
	if err != nil {
		return nil, err
	}
	s.db = db
	go s.runGC(s.opts.GCRatio, s.opts.GCInterval)
	return s, nil
}

func (s *badgerStore) runGC(ratio float64, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for range ticker.C {
	again:
		err := s.db.RunValueLogGC(ratio)
		if err == nil {
			goto again
		}
	}
}

func (s *badgerStore) View(k *Key, fn PeekFunc) error {
	return s.db.View(func(tx *badger.Txn) error {
		v, err := tx.Get(k.Bytes())
		if err == badger.ErrKeyNotFound {
			return ErrNotFound
		} else if err != nil {
			err = fmt.Errorf("item get error: %v", err)
			return err
		}
		vv, err := v.Value()
		if err != nil {
			err = fmt.Errorf("value read error: %v", err)
			return err
		}
		return fn(k, vv)
	})
}

func (s *badgerStore) Update(k *Key, fn ModifyFunc) error {
	return s.db.Update(func(tx *badger.Txn) error {
		if fn == nil {
			return nil
		}
		key := k.Bytes()
		v, err := tx.Get(key)
		if err == badger.ErrKeyNotFound {
			vv, err := fn(k, nil)
			if err == ErrNoUpdate {
				return nil
			} else if err != nil {
				return err
			}
			if k.TTL > 0 {
				return tx.SetWithTTL(key, vv, k.TTL)
			}
			return tx.Set(key, vv)
		} else if err != nil {
			err = fmt.Errorf("item set error: %v", err)
			return err
		}
		vv, err := v.ValueCopy(nil)
		if err != nil {
			return err
		}
		vv, err = fn(k, vv)
		if err == ErrNoUpdate {
			return nil
		} else if err != nil {
			return err
		}
		if k.TTL > 0 {
			return tx.SetWithTTL(key, vv, k.TTL)
		}
		return tx.Set(key, vv)
	})
}

func (s *badgerStore) RangeKeys(b Bucket, fn KeyFunc) (*RangeOptions, error) {
	var opt *RangeOptions
	err := s.db.View(func(tx *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchSize = 10
		opts.PrefetchValues = false
		it := tx.NewIterator(opts)
		defer it.Close()

		for it.Seek(b.NewKey(nil).Bytes()); it.Valid(); it.Next() {
			item := it.Item()
			k := (&Key{}).Unmarshal(item.Key())
			if k.Bucket.ID != b.ID {
				return nil
			}
			fn(k)
		}
		return nil
	})
	return opt, err
}

func (s *badgerStore) RangePeek(b Bucket, fn PeekFunc) (*RangeOptions, error) {
	var opt *RangeOptions
	err := s.db.View(func(tx *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchSize = 10
		if b.RangeOptions.Prefetch > 0 {
			opts.PrefetchSize = b.RangeOptions.Prefetch
		}
		it := tx.NewIterator(opts)
		defer it.Close()

		for it.Seek(b.NewKey(nil).Bytes()); it.Valid(); it.Next() {
			item := it.Item()
			k := (&Key{}).Unmarshal(item.Key())
			if k.Bucket.ID != b.ID {
				return nil
			}
			v, err := it.Item().Value()
			if err != nil {
				return err
			}
			if err := fn(k, v); err == ErrRangeStop {
				return nil
			} else if err != nil {
				return err
			}
		}
		return nil
	})
	return opt, err
}

func (s *badgerStore) RangeModify(b Bucket, fn ModifyFunc) (*RangeOptions, error) {
	var opt *RangeOptions
	err := s.db.View(func(tx *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchSize = 10
		if b.RangeOptions.Prefetch > 0 {
			opts.PrefetchSize = b.RangeOptions.Prefetch
		}
		it := tx.NewIterator(opts)
		defer it.Close()

		for it.Seek(b.NewKey(nil).Bytes()); it.Valid(); it.Next() {
			item := it.Item()
			k := (&Key{}).Unmarshal(item.Key())
			if k.Bucket.ID != b.ID {
				return nil
			}
			v, err := it.Item().Value()
			if err != nil {
				return err
			}
			vv, err := fn(k, v)
			if err == ErrNoUpdate {
				continue
			} else if err != nil && err != ErrRangeStop {
				return err
			}
			if setErr := s.db.Update(func(tx *badger.Txn) error {
				return tx.Set(item.Key(), vv)
			}); setErr != nil {
				return setErr
			}
			if err == ErrRangeStop {
				return nil
			}
		}
		return nil
	})
	return opt, err
}

func (s *badgerStore) Delete(k *Key) error {
	if k == nil {
		return nil
	}
	return s.db.View(func(tx *badger.Txn) error {
		if err := tx.Delete(k.Bytes()); err == badger.ErrKeyNotFound {
			return nil
		} else if err != nil {
			return err
		}
		return nil
	})
}

func (s *badgerStore) Close() error {
	return s.db.Close()
}
