package excache

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/boltdb/bolt"
	"github.com/vmihailenco/msgpack"
)

var (
	// ErrItemExpired if item is got from local storage
	ErrItemExpired = errors.New("item was expired")
)

var _DefaultWriteQueueSize = 10000
var _DefaultCacheBucket = "excache"
var _DefaultBatchSize = 1 // no batch

type asyncSet struct {
	key interface{}
	val interface{}
	err chan<- error
}

type serializable interface {
	Load([]byte)
	Dump() []byte
}

type PLRUStat struct {
	LRUItem  uint64
	LRUQuery uint64
	LRUHit   uint64
	DBStat   bolt.Stats
}

// PersistentLRUCache LRUCache with local storage
type PersistentLRUCache struct {
	*LRUCache
	batchSize  int
	db         *bolt.DB
	rotateSize int

	currBuckID uint8
	setQueue   chan asyncSet
	closeCh    chan struct{}
}

// NewPersistentLRUCache return new persistent lru cache
func NewPersistentLRUCache(db string, rotateSize int, size int, age int, expires time.Duration) (
	p *PersistentLRUCache, err error) {
	// 临时将逻辑修改为不持久化看程序表现
	return &PersistentLRUCache{
		LRUCache: NewLRUCache(size, age, expires),
	}, nil

	if rotateSize <= 0 {
		panic("rotateSize should greater than zero")
	}
	p = &PersistentLRUCache{
		LRUCache:   NewLRUCache(size, age, expires),
		batchSize:  _DefaultBatchSize,
		rotateSize: rotateSize,
		currBuckID: 0,
		setQueue:   make(chan asyncSet, _DefaultWriteQueueSize),
		closeCh:    make(chan struct{}),
	}

	p.db, err = bolt.Open(db, os.ModePerm, &bolt.Options{})

	go p.backupLoop()
	return
}

// SetBatchSize modify batch size
func (p *PersistentLRUCache) SetBatchSize(n int) {
	p.batchSize = n
}

// Stat statistics
func (p *PersistentLRUCache) Stat() *PLRUStat {
	ret := &PLRUStat{}
	ret.LRUItem, ret.LRUQuery, ret.LRUHit = p.LRUCache.Count()
	ret.DBStat = p.db.Stats()

	return ret
}

// Set value with key, and write to local storage.
// Set is non-block api for write local storage
func (p *PersistentLRUCache) Set(key, val interface{}, errCh chan<- error) {
	p.LRUCache.Set(key, val)
	// 临时将逻辑修改为不持久化看程序表现
	return

	p.setQueue <- asyncSet{
		key: key,
		val: val,
		err: errCh,
	}
}

// SyncSet is same as Set() except block and return error
func (p *PersistentLRUCache) SyncSet(key, val interface{}) error {
	errCh := make(chan error)
	p.Set(key, val, errCh)
	return <-errCh
}

// Get try get from memory or disk
func (p *PersistentLRUCache) Get(key interface{}, allocFunc func() interface{}) (
	val interface{}, ok bool, err error) {
	val, ok = p.LRUCache.Get(key)
	// 临时将逻辑修改为不持久化看程序表现
	return val, ok, nil

	if ok {
		return val, ok, nil
	}
	err = p.db.View(func(tx *bolt.Tx) error {
		bkey, err := msgpack.Marshal(key)
		if err != nil {
			return fmt.Errorf("Pack key error: %s", err.Error())
		}
		var bucketSeek = func(bn []byte) []byte {
			bucket := tx.Bucket(bn)
			if bucket == nil {
				return nil
			}
			bval := bucket.Get(bkey)
			if bval == nil {
				return nil
			}
			return bval
		}
		bval := bucketSeek(p.firstBucket())
		if bval == nil {
			bval = bucketSeek(p.sencondBucket())
		}
		if bval == nil {
			return nil
		}
		val = allocFunc()
		if serialier, ok := val.(serializable); ok {
			serialier.Load(bval)
		} else {
			err = msgpack.Unmarshal(bval, val)
			if err != nil {
				fmt.Printf("bval: %+v, val: %+v, err: %+v\n", bval, val, err)
				return err
			}
		}
		ok = true
		return ErrItemExpired
	})
	return val, ok, err
}

// PurgeBackup purge backup
func (p *PersistentLRUCache) PurgeBackup() error {
	err := p.db.Update(func(tx *bolt.Tx) error {
		return tx.DeleteBucket(p.firstBucket())
	})
	if err == bolt.ErrBucketNotFound {
		return nil
	}
	return err
}

func (p *PersistentLRUCache) firstBucket() []byte {
	return []byte(fmt.Sprintf("%s_%d", _DefaultCacheBucket, p.currBuckID))
}

func (p *PersistentLRUCache) sencondBucket() []byte {
	return []byte(fmt.Sprintf("%s_%d", _DefaultCacheBucket, 1-p.currBuckID))
}

func (p *PersistentLRUCache) backupLoop() {
	batch := make([]asyncSet, 0)
	for {
		select {
		case <-p.closeCh:
			p.db.Close()
			break
		case msg := <-p.setQueue:
			batch = append(batch, msg)
			if len(batch) < p.batchSize {
				if msg.err != nil {
					msg.err <- nil
				}
				continue
			}
			// t0 := time.Now()
			err := p.db.Update(func(tx *bolt.Tx) error {
				bucket, err := tx.CreateBucketIfNotExists(p.firstBucket())
				if err != nil {
					return fmt.Errorf("Create bucket: %s", err.Error())
				}
				for _, record := range batch {
					// fmt.Printf("key: %v, begin: %s\n", msg.key, time.Since(t0).String())
					bkey, err := msgpack.Marshal(record.key)
					if err != nil {
						return fmt.Errorf("Pack Key error: %s", err.Error())
					}
					var bval []byte
					if serializer, ok := record.val.(serializable); ok {
						bval = serializer.Dump()
					} else {
						bval, err = msgpack.Marshal(record.val)
						if err != nil {
							return fmt.Errorf("Pack value error: %s", err.Error())
						}
					}
					// fmt.Printf("key: %v, serialize: %s\n", msg.key, time.Since(t0).String())
					if err := bucket.Put(bkey, bval); err != nil {
						return err
					}
					if bucket.Stats().KeyN >= p.rotateSize {
						// Ignore rotate failed
						err := tx.DeleteBucket(p.sencondBucket())
						if err == nil {
							p.currBuckID = 1 - p.currBuckID
						}
					}
				}
				return nil
			})
			// fmt.Printf("key: %+v, batch: %s\n", msg.key, time.Since(t0).String())
			batch = make([]asyncSet, 0)
			if msg.err != nil {
				msg.err <- err
			}
		}
	}
}
