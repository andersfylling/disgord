package ratelimit

import (
	"sync"
	"time"

	"github.com/andersfylling/disgord/internal/httd"
)

type GroupID uint8

const (
	GroupOthers GroupID = iota
	GroupChannels
	GroupGuilds
	GroupWebhooks

	NonMajor httd.Snowflake = 0
)

func newBucketGroup() bucketGroup {
	return bucketGroup{
		buckets: map[httd.Snowflake][]*httd.bucket{},
	}
}

type bucketGroup struct {
	sync.RWMutex
	buckets map[httd.Snowflake][]*httd.bucket
}

func (r *bucketGroup) add(majorID httd.Snowflake, b *httd.bucket) {
	r.Lock()
	if _, exists := r.buckets[majorID]; !exists {
		r.buckets[majorID] = []*httd.bucket{b}
	}
	r.buckets[majorID] = append(r.buckets[majorID], b)
	r.Unlock()
}

func (r *bucketGroup) bucket(majorID httd.Snowflake, localBucketKey httd.LocalKey) (b *httd.bucket, ok bool) {
	r.RLock()
	defer r.RUnlock()
	if buckets, exists := r.buckets[majorID]; exists {
		for i := range buckets {
			if buckets[i].LinkedTo(localBucketKey) {
				b = buckets[i]
				ok = true
				break
			}
		}
		return b, ok
	}
	return nil, false
}

// consolidate makes sure that the given bucket has a unique hash. On matches the data of the oldest bucket is copied
// to the newest, and the oldest is marked as invalid. After 10 minutes, it is deleted.
func (r *bucketGroup) consolidate(majorID httd.Snowflake, b *httd.bucket) {
	b.AcquireLock()
	defer b.Unlock()
	if b.invalid {
		return
	}

	r.RLock()
	defer r.RUnlock()
	buckets, ok := r.buckets[majorID]
	if !ok { // if the major (channel, guild,...) was deleted for some reason...
		return
	}

	var invalid *httd.bucket
	for i := range buckets {
		if buckets[i] == b {
			continue
		}

		buckets[i].AcquireLock()
		if buckets[i].invalid {
			buckets[i].Unlock()
			continue
		}
		if buckets[i].hash != b.hash {
			continue
		}

		var oldest *httd.bucket
		var newest *httd.bucket
		if buckets[i].reset.Before(b.reset) {
			oldest = b
			newest = buckets[i]
		} else {
			oldest = buckets[i]
			newest = b
		}

		// mark invalid
		oldest.invalid = true
		invalid = oldest

		// copy data into the new one
		// local keys
		for j := range oldest.localKeys {
			newest.AddLocalKey(oldest.localKeys[j])
		}

		buckets[i].Unlock()
		break
	}

	// delete invalid bucket
	go func(b *httd.bucket) {
		<-time.After(10 * time.Minute)
		r.Lock()
		if buckets, ok := r.buckets[majorID]; ok {
			for i := range buckets {
				if buckets[i] != b {
					continue
				}

				buckets[i] = nil
				buckets = append(buckets[:i], buckets[i+1:]...)
				r.buckets[majorID] = buckets
				break
			}
		}
		r.Unlock()
	}(invalid)
}
