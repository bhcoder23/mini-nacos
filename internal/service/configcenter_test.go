package service

import (
	"context"
	"testing"

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

type serviceTestHub struct{}

func (h *serviceTestHub) Notify(context.Context, *biz.ConfigChange) {}

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
