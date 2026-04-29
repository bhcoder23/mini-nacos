package data

import (
	"context"
	"mini-nacos/internal/biz"
	"sync"

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

type configWatchHub struct {
	data    *Data
	log     *log.Helper
	mu      sync.RWMutex
	changes []biz.ConfigChange
}

func NewConfigWatchHub(data *Data, logger log.Logger) biz.ConfigWatchHub {

	return &configWatchHub{
		data:    data,
		log:     log.NewHelper(logger),
		changes: make([]biz.ConfigChange, 0),
	}

}

func (h *configWatchHub) Notify(_ context.Context, change *biz.ConfigChange) {

	if change == nil {
		return
	}
	h.mu.Lock()
	defer h.mu.Unlock()
	h.changes = append(h.changes, *change)

}
