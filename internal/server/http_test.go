package server

import (
	"context"
	"io"
	"net/http/httptest"
	"testing"

	pb "mini-nacos/api/configcenter/v1"
	"mini-nacos/internal/biz"
	"mini-nacos/internal/conf"
	"mini-nacos/internal/data"
	"mini-nacos/internal/service"

	"github.com/go-kratos/kratos/v2/log"
	kratoshttp "github.com/go-kratos/kratos/v2/transport/http"
)

func newTestConfigCenterService(t *testing.T, item biz.ConfigItem) *service.ConfigCenterService {
	t.Helper()

	logger := log.NewStdLogger(io.Discard)
	dataStore, cleanup, err := data.NewData(&conf.Data{})
	if err != nil {
		t.Fatalf("NewData() error = %v", err)
	}
	t.Cleanup(cleanup)

	repo := data.NewConfigRepo(dataStore, logger)
	hub := data.NewConfigWatchHub(dataStore, logger)
	uc := biz.NewConfigUseCase(repo, hub)

	if err := repo.Save(context.Background(), &item); err != nil {
		t.Fatalf("repo.Save() error = %v", err)
	}

	return service.NewConfigCenterService(uc)
}

func newTestConfigCenterHTTPClient(t *testing.T, configSvc *service.ConfigCenterService) pb.ConfigCenterHTTPClient {
	t.Helper()

	logger := log.NewStdLogger(io.Discard)
	srv := NewHTTPServer(&conf.Server{Http: &conf.Server_HTTP{}}, &service.GreeterService{}, configSvc, logger)
	ts := httptest.NewServer(srv)
	t.Cleanup(ts.Close)

	client, err := kratoshttp.NewClient(context.Background(), kratoshttp.WithEndpoint(ts.URL))
	if err != nil {
		t.Fatalf("http.NewClient() error = %v", err)
	}
	t.Cleanup(func() {
		_ = client.Close()
	})

	return pb.NewConfigCenterHTTPClient(client)
}

func TestHTTPServerGetConfigRoute(t *testing.T) {
	configSvc := newTestConfigCenterService(t, biz.ConfigItem{
		Key: biz.ConfigKey{
			Namespace: "public",
			Group:     "DEFAULT_GROUP",
			DataID:    "app.yaml",
		},
		Content: "content-v1",
		MD5:     "md5-v1",
	})
	client := newTestConfigCenterHTTPClient(t, configSvc)

	resp, err := client.GetConfig(context.Background(), &pb.GetConfigRequest{
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
}

func TestHTTPServerListenConfigRoute(t *testing.T) {
	configSvc := newTestConfigCenterService(t, biz.ConfigItem{
		Key: biz.ConfigKey{
			Namespace: "public",
			Group:     "DEFAULT_GROUP",
			DataID:    "app.yaml",
		},
		Content: "content-v1",
		MD5:     "server-md5",
	})
	client := newTestConfigCenterHTTPClient(t, configSvc)

	resp, err := client.ListenConfig(context.Background(), &pb.ListenConfigRequest{
		Namespace: "public",
		Group:     "DEFAULT_GROUP",
		DataId:    "app.yaml",
		Md5:       "client-md5",
		TimeoutMs: 1000,
	})
	if err != nil {
		t.Fatalf("ListenConfig() error = %v", err)
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
	if resp.Md5 != "server-md5" {
		t.Fatalf("expected md5 server-md5, got %q", resp.Md5)
	}
	if !resp.Changed {
		t.Fatal("expected changed to be true")
	}
}
