package biz

import (
	"context"
	"testing"
)

type fakeConfigRepo struct {
	getItem   ConfigItem
	getErr    error
	saved     []ConfigItem
	callOrder *[]string
}

func (r *fakeConfigRepo) Save(_ context.Context, item *ConfigItem) error {
	if r.callOrder != nil {
		*r.callOrder = append(*r.callOrder, "save")
	}
	r.saved = append(r.saved, *item)

	return nil
}

func (r *fakeConfigRepo) Get(_ context.Context, _ ConfigKey) (ConfigItem, error) {
	return r.getItem, r.getErr
}

type fakeConfigWatchHub struct {
	notified  []ConfigChange
	callOrder *[]string
}

func (h *fakeConfigWatchHub) Notify(_ context.Context, change *ConfigChange) {
	if h.callOrder != nil {
		*h.callOrder = append(*h.callOrder, "notify")
	}
	if change != nil {
		h.notified = append(h.notified, *change)
	}
}

func TestConfigUsecasePublishFirstPublish(t *testing.T) {
	key := ConfigKey{
		Namespace: "public",
		Group:     "DEFAULT_GROUP",
		DataID:    "app.yaml",
	}
	order := make([]string, 0, 2)
	repo := &fakeConfigRepo{
		getErr:    ErrConfigNotFound,
		callOrder: &order,
	}
	hub := &fakeConfigWatchHub{callOrder: &order}
	uc := NewConfigUseCase(repo, hub)

	item, err := uc.Publish(context.Background(), key, "content-v1")
	if err != nil {
		t.Fatalf("Publish() error = %v", err)
	}
	if len(repo.saved) != 1 {
		t.Fatalf("expected 1 saved item, got %d", len(repo.saved))
	}
	if len(hub.notified) != 1 {
		t.Fatalf("expected 1 notify event, got %d", len(hub.notified))
	}
	if len(order) != 2 || order[0] != "save" || order[1] != "notify" {
		t.Fatalf("expected call order [save notify], got %v", order)
	}
	if item.Key != key {
		t.Fatalf("expected key %#v, got %#v", key, item.Key)
	}
	if item.Content != "content-v1" {
		t.Fatalf("expected content-v1, got %q", item.Content)
	}
	if item.MD5 == "" {
		t.Fatal("expected MD5 to be calculated")
	}
}

func TestConfigUsecasePublishUnchangedStillSaveButNoNotify(t *testing.T) {
	key := ConfigKey{
		Namespace: "public",
		Group:     "DEFAULT_GROUP",
		DataID:    "app.yaml",
	}
	old := ConfigItem{
		Key:     key,
		Content: "content-v1",
		MD5:     calcMD5("content-v1"),
	}
	repo := &fakeConfigRepo{getItem: old}
	hub := &fakeConfigWatchHub{}
	uc := NewConfigUseCase(repo, hub)

	item, err := uc.Publish(context.Background(), key, "content-v1")
	if err != nil {
		t.Fatalf("Publish() error = %v", err)
	}
	if len(repo.saved) != 1 {
		t.Fatalf("expected 1 saved item, got %d", len(repo.saved))
	}
	if len(hub.notified) != 0 {
		t.Fatalf("expected 0 notify events, got %d", len(hub.notified))
	}
	if item.MD5 != old.MD5 {
		t.Fatalf("expected unchanged md5 %q, got %q", old.MD5, item.MD5)
	}
}
