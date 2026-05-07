package service

import (
	"context"
	"testing"
	"time"

	pb "mini-nacos/api/configcenter/v1"
	"mini-nacos/internal/biz"
)

type serviceTestRepo struct {
	item biz.ConfigItem
	err  error
}

func (r *serviceTestRepo) Save(_ context.Context, item *biz.ConfigItem) error {
	r.item = *item
	return nil
}

func (r *serviceTestRepo) Get(_ context.Context, _ biz.ConfigKey) (biz.ConfigItem, error) {
	return r.item, r.err
}

type serviceTestHub struct {
	waitChange  biz.ConfigChange
	waitOK      bool
	waitErr     error
	waitCalled  bool
	waitKey     biz.ConfigKey
	waitTimeout time.Duration
}

func (h *serviceTestHub) Notify(context.Context, *biz.ConfigChange) {}

func (h *serviceTestHub) Wait(_ context.Context, key biz.ConfigKey, timeout time.Duration) (biz.ConfigChange, bool, error) {
	h.waitCalled = true
	h.waitKey = key
	h.waitTimeout = timeout
	return h.waitChange, h.waitOK, h.waitErr
}

func TestConfigCenterServicePublishConfig(t *testing.T) {
	repo := &serviceTestRepo{err: biz.ErrConfigNotFound}
	hub := &serviceTestHub{}
	uc := biz.NewConfigUseCase(repo, hub)
	svc := NewConfigCenterService(uc)

	resp, err := svc.PublishConfig(context.Background(), &pb.PublishConfigRequest{
		Namespace: "public",
		Group:     "DEFAULT_GROUP",
		DataId:    "app.yaml",
		Content:   "content-v1",
	})
	if err != nil {
		t.Fatalf("PublishConfig() error = %v", err)
	}
	if resp.Namespace != "public" {
		t.Fatalf("expected namespace public, got %q", resp.Namespace)
	}
	if resp.Group != "DEFAULT_GROUP" {
		t.Fatalf("expected group DEFAULT_GROUP, got %q", resp.Group)
	}
	if resp.DataId != "app.yaml" {
		t.Fatalf("expected dataId app.yaml, got %q", resp.DataId)
	}
	if resp.Content != "content-v1" {
		t.Fatalf("expected content-v1, got %q", resp.Content)
	}
	if resp.Md5 == "" {
		t.Fatal("expected md5 in response")
	}
}

func TestConfigCenterServiceGetConfig(t *testing.T) {
	repo := &serviceTestRepo{
		item: biz.ConfigItem{
			Key: biz.ConfigKey{
				Namespace: "public",
				Group:     "DEFAULT_GROUP",
				DataID:    "app.yaml",
			},
			Content: "content-v1",
			MD5:     "md5-v1",
		},
	}
	hub := &serviceTestHub{}
	uc := biz.NewConfigUseCase(repo, hub)
	svc := NewConfigCenterService(uc)

	resp, err := svc.GetConfig(context.Background(), &pb.GetConfigRequest{
		Namespace: "public",
		Group:     "DEFAULT_GROUP",
		DataId:    "app.yaml",
	})
	if err != nil {
		t.Fatalf("GetConfig() error = %v", err)
	}
	if resp.Namespace != "public" {
		t.Fatalf("expected namespace public, got %q", resp.Namespace)
	}
	if resp.Group != "DEFAULT_GROUP" {
		t.Fatalf("expected group DEFAULT_GROUP, got %q", resp.Group)
	}
	if resp.DataId != "app.yaml" {
		t.Fatalf("expected dataId app.yaml, got %q", resp.DataId)
	}
	if resp.Content != "content-v1" {
		t.Fatalf("expected content-v1, got %q", resp.Content)
	}
	if resp.Md5 != "md5-v1" {
		t.Fatalf("expected md5 md5-v1, got %q", resp.Md5)
	}
	if hub.waitCalled {
		t.Fatal("expected GetConfig not to call Wait")
	}
}

func TestConfigCenterServiceListenConfig(t *testing.T) {
	key := biz.ConfigKey{
		Namespace: "public",
		Group:     "DEFAULT_GROUP",
		DataID:    "app.yaml",
	}
	repo := &serviceTestRepo{
		item: biz.ConfigItem{
			Key: key,
			MD5: "same-md5",
		},
	}
	hub := &serviceTestHub{
		waitChange: biz.ConfigChange{
			Key: key,
			MD5: "new-md5",
		},
		waitOK: true,
	}
	uc := biz.NewConfigUseCase(repo, hub)
	svc := NewConfigCenterService(uc)

	resp, err := svc.ListenConfig(context.Background(), &pb.ListenConfigRequest{
		Namespace: "public",
		Group:     "DEFAULT_GROUP",
		DataId:    "app.yaml",
		Md5:       "same-md5",
		TimeoutMs: 1500,
	})
	if err != nil {
		t.Fatalf("ListenConfig() error = %v", err)
	}
	if !hub.waitCalled {
		t.Fatal("expected ListenConfig to call Wait")
	}
	if hub.waitKey != key {
		t.Fatalf("expected wait key %#v, got %#v", key, hub.waitKey)
	}
	if hub.waitTimeout != 1500*time.Millisecond {
		t.Fatalf("expected wait timeout 1500ms, got %v", hub.waitTimeout)
	}
	if resp.Namespace != "public" {
		t.Fatalf("expected namespace public, got %q", resp.Namespace)
	}
	if resp.Group != "DEFAULT_GROUP" {
		t.Fatalf("expected group DEFAULT_GROUP, got %q", resp.Group)
	}
	if resp.DataId != "app.yaml" {
		t.Fatalf("expected dataId app.yaml, got %q", resp.DataId)
	}
	if resp.Md5 != "new-md5" {
		t.Fatalf("expected md5 new-md5, got %q", resp.Md5)
	}
	if !resp.Changed {
		t.Fatal("expected changed to be true")
	}
}
