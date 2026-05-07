package data

import (
	"context"
	"mini-nacos/internal/biz"
	"sync"
	"time"

	"github.com/go-kratos/kratos/v2/log"
)

type configRepo struct {
	data  *Data
	log   *log.Helper
	mu    sync.RWMutex
	items map[biz.ConfigKey]biz.ConfigItem
}

func NewConfigRepo(data *Data, logger log.Logger) biz.ConfigRepo {
	return &configRepo{
		data:  data,
		log:   log.NewHelper(logger),
		items: make(map[biz.ConfigKey]biz.ConfigItem),
	}
}

func (r *configRepo) Save(_ context.Context, item *biz.ConfigItem) error {

	if item == nil {
		return biz.ErrInvalidConfigItem
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	r.items[item.Key] = *item
	return nil

}

func (r *configRepo) Get(_ context.Context, key biz.ConfigKey) (biz.ConfigItem, error) {

	r.mu.RLock()
	defer r.mu.RUnlock()
	item, ok := r.items[key]
	if !ok {
		return biz.ConfigItem{}, biz.ErrConfigNotFound
	}
	return item, nil

}

type waiter struct {
	ch chan biz.ConfigChange
}

type configWatchHub struct {
	data    *Data
	log     *log.Helper
	mu      sync.RWMutex
	changes []biz.ConfigChange

	waiters map[biz.ConfigKey]map[*waiter]struct{}
}

func NewConfigWatchHub(data *Data, logger log.Logger) biz.ConfigWatchHub {

	return &configWatchHub{
		data:    data,
		log:     log.NewHelper(logger),
		changes: make([]biz.ConfigChange, 0),
		waiters: make(map[biz.ConfigKey]map[*waiter]struct{}),
	}

}

func (h *configWatchHub) Notify(_ context.Context, change *biz.ConfigChange) {

	if change == nil {
		return
	}
	h.mu.Lock()
	defer h.mu.Unlock()

	ws, ok := h.waiters[change.Key]
	if !ok {
		return
	}

	for w := range ws {
		select {
		case w.ch <- *change:
		default:
		}
	}

	delete(h.waiters, change.Key)

}

func (h *configWatchHub) Wait(ctx context.Context, key biz.ConfigKey, duration time.Duration) (biz.ConfigChange, bool, error) {
	w := &waiter{
		ch: make(chan biz.ConfigChange, 1),
	}

	h.mu.Lock()
	if h.waiters[key] == nil {
		h.waiters[key] = make(map[*waiter]struct{})
	}
	h.waiters[key][w] = struct{}{}
	h.mu.Unlock()

	timer := time.NewTimer(duration)
	defer timer.Stop()

	defer func() {
		h.mu.Lock()
		defer h.mu.Unlock()
		delete(h.waiters[key], w)
		if len(h.waiters[key]) == 0 {
			delete(h.waiters, key)
		}
	}()

	select {
	case change := <-w.ch:
		return change, true, nil
	case <-timer.C:
		return biz.ConfigChange{}, false, nil
	case <-ctx.Done():
		return biz.ConfigChange{}, false, ctx.Err()
	}
}
