package httd

import (
	"net/http"
	"strconv"
	"sync"
	"time"
)

const GlobalHash = "global"

func relationsByBucketID(relations map[string]string) map[string][]string {
	byHash := make(map[string][]string)
	for id, hash := range relations {
		if _, ok := byHash[hash]; !ok {
			byHash[hash] = []string{id}
		} else {
			byHash[hash] = append(byHash[hash], id)
		}
	}

	return byHash
}

func NewManager(defaultRelations map[string]string) *Manager {
	global := newBucket(nil)
	global.hash = GlobalHash

	m := &Manager{
		relations: make(map[string]*bucket),
		global:    global,
	}

	hashRelations := relationsByBucketID(defaultRelations)
	for hash, ids := range hashRelations {
		var bucket *bucket
		if hash == GlobalHash {
			bucket = m.global
		} else {
			bucket = newBucket(m.global)
		}

		for i := range ids {
			m.relations[ids[i]] = bucket
		}
	}

	return m
}

type BucketTransactioner interface {
	GetSet(atomicTransaction func(bucket string) (updated string))
}

type Manager struct {
	mu        sync.RWMutex
	relations map[string]*bucket

	buckets BucketTransactioner

	global *bucket
}

func (r *Manager) Relations() (relations map[string]string) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	relations = make(map[string]string)
	for k, v := range r.relations {
		relations[k] = v.hash
	}
	return relations
}

func (r *Manager) Bucket(id string) *bucket {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.relations[id]; !ok {
		r.relations[id] = newBucket(r.global)
	}
	return r.relations[id]
}

func (r *Manager) UpdateBucket(hash string, header http.Header) {
	// entry should always exist, as this is called after the bucket is ensured..
	r.mu.RLock()
	b := r.relations[hash]
	r.mu.RUnlock()

	// to synchronize the timestamp between the bot and the discord server
	// we assume the current time is equal the header date
	discordTime, err := HeaderToTime(header)
	if err != nil {
		discordTime = time.Now()
	}

	localTime := time.Now()
	diff := localTime.Sub(discordTime)

	var bu *bucket
	if global := header.Get(XRateLimitGlobal); global == "true" {
		bu = b.global
	} else {
		bu = b
	}

	bu.AcquireLock()
	defer bu.Unlock()

	if resetStr := header.Get(XRateLimitReset); resetStr != "" {
		epoch, err := strconv.ParseInt(resetStr, 10, 64)
		if err != nil {
			return
		}

		old := b.resetTime
		bu.resetTime = time.Unix(0, epoch+diff.Nanoseconds())

		oldNewDiffMs := uint(bu.resetTime.Sub(old).Nanoseconds() / int64(time.Millisecond))
		if !bu.resetTime.Equal(old) && bu.longestTimeout < oldNewDiffMs {
			bu.longestTimeout = oldNewDiffMs
		}
	}

	if remainingStr := header.Get(XRateLimitRemaining); remainingStr != "" {
		remaining, err := strconv.ParseInt(remainingStr, 10, 64)
		if err != nil {
			return
		}

		bu.remaining = uint(remaining)
	}

	if limitStr := header.Get(XRateLimitLimit); limitStr != "" {
		limit, err := strconv.ParseInt(limitStr, 10, 64)
		if err != nil {
			return
		}

		bu.limit = uint(limit)
	}

	if _, ok := header[XRateLimitBucket]; ok && header.Get(XRateLimitBucket) == "" {
		bu.hash = GlobalHash
	}

	if key := header.Get(XRateLimitBucket); key != "" {
		bu.hash = key
	}
}

func (r *Manager) Consolidate() {

}
