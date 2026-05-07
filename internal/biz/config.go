// Package biz usecase
package biz

import (
	"context"
	// #nosec G501 -- MD5 used for config content comparison, not cryptographic security
	"crypto/md5"
	"fmt"
	v1 "mini-nacos/api/configcenter/v1"
	"time"

	"github.com/go-kratos/kratos/v2/errors"
)

var (
	ErrConfigNotFound    = errors.NotFound(v1.ErrorReason_CONFIG_NOT_FOUND.String(), "config not found")
	ErrInvalidConfigItem = errors.BadRequest("CONFIG_ITEM_INVALID", "config item is nil")
)

type ConfigKey struct {
	Namespace string
	Group     string
	DataID    string
}

type ConfigItem struct {
	Key       ConfigKey
	Content   string
	MD5       string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type ConfigChange struct {
	Key       ConfigKey
	MD5       string
	ChangedAt time.Time
}

type ListenResult struct {
	Key     ConfigKey
	MD5     string
	Changed bool
}

type ConfigRepo interface {
	Save(context.Context, *ConfigItem) error
	Get(context.Context, ConfigKey) (ConfigItem, error)
}

type ConfigWatchHub interface {
	Notify(context.Context, *ConfigChange)
	Wait(context.Context, ConfigKey, time.Duration) (ConfigChange, bool, error)
}

type ConfigUsecase struct {
	repo ConfigRepo
	hub  ConfigWatchHub
}

func NewConfigUseCase(repo ConfigRepo, hub ConfigWatchHub) *ConfigUsecase {
	return &ConfigUsecase{
		repo: repo,
		hub:  hub,
	}
}

func (uc *ConfigUsecase) Publish(ctx context.Context, key ConfigKey, content string) (ConfigItem, error) {
	now := time.Now().UTC()
	newMD5 := calcMD5(content)

	currentItem, err := uc.repo.Get(ctx, key)
	if err != nil && !errors.Is(err, ErrConfigNotFound) {
		return ConfigItem{}, fmt.Errorf("ConfigCenterUseCase - Publish - uc.repo.Get: %w", err)
	}

	isFirstPublish := errors.Is(err, ErrConfigNotFound)

	item := ConfigItem{
		Key:       key,
		Content:   content,
		MD5:       newMD5,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if !isFirstPublish {
		item.CreatedAt = currentItem.CreatedAt
	}

	if err := uc.repo.Save(ctx, &item); err != nil {
		return ConfigItem{}, fmt.Errorf("ConfigCenterUseCase - Publish - uc.repo.Save: %w", err)
	}

	if isFirstPublish || currentItem.MD5 != item.MD5 {
		change := &ConfigChange{
			Key:       key,
			MD5:       newMD5,
			ChangedAt: now,
		}
		uc.hub.Notify(ctx, change)
	}

	return item, nil
}

func (uc *ConfigUsecase) Get(ctx context.Context, key ConfigKey) (ConfigItem, error) {
	item, err := uc.repo.Get(ctx, key)
	if err != nil {
		return ConfigItem{}, err
	}
	return item, nil
}

func (uc *ConfigUsecase) Listen(ctx context.Context, key ConfigKey, clientMD5 string, timeout time.Duration) (ListenResult, error) {
	item, err := uc.repo.Get(ctx, key)
	if err != nil {
		return ListenResult{}, err
	}

	if item.MD5 != clientMD5 {
		return ListenResult{
			Key:     key,
			MD5:     item.MD5,
			Changed: true,
		}, nil
	}

	change, ok, err := uc.hub.Wait(ctx, key, timeout)
	if err != nil {
		return ListenResult{}, err
	}

	if !ok {
		return ListenResult{
			Key:     key,
			MD5:     item.MD5,
			Changed: false,
		}, nil
	}

	return ListenResult{
		Key:     change.Key,
		MD5:     change.MD5,
		Changed: true,
	}, nil
}

func calcMD5(content string) string {
	// #nosec G401 -- MD5 used for config content comparison, not cryptographic security
	return fmt.Sprintf("%x", md5.Sum([]byte(content)))
}
