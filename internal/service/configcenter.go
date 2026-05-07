// Package service service
package service

import (
	"context"
	"time"

	pb "mini-nacos/api/configcenter/v1"
	"mini-nacos/internal/biz"
)

type ConfigCenterService struct {
	pb.UnimplementedConfigCenterServer

	uc *biz.ConfigUsecase
}

func NewConfigCenterService(uc *biz.ConfigUsecase) *ConfigCenterService {
	return &ConfigCenterService{uc: uc}
}

func (s *ConfigCenterService) PublishConfig(ctx context.Context, req *pb.PublishConfigRequest) (*pb.PublishConfigResponse, error) {
	item, err := s.uc.Publish(ctx, biz.ConfigKey{
		Namespace: req.Namespace,
		Group:     req.Group,
		DataID:    req.DataId,
	}, req.Content)
	if err != nil {
		return nil, err
	}

	return &pb.PublishConfigResponse{
		Namespace: item.Key.Namespace,
		Group:     item.Key.Group,
		DataId:    item.Key.DataID,
		Content:   item.Content,
		Md5:       item.MD5,
	}, nil
}

func (s *ConfigCenterService) GetConfig(ctx context.Context, req *pb.GetConfigRequest) (*pb.GetConfigResponse, error) {
	item, err := s.uc.Get(ctx, biz.ConfigKey{
		Namespace: req.Namespace,
		Group:     req.Group,
		DataID:    req.DataId,
	})

	if err != nil {
		return nil, err
	}

	return &pb.GetConfigResponse{
		Namespace: item.Key.Namespace,
		Group:     item.Key.Group,
		DataId:    item.Key.DataID,
		Content:   item.Content,
		Md5:       item.MD5,
	}, nil
}

func (s *ConfigCenterService) ListenConfig(ctx context.Context, req *pb.ListenConfigRequest) (*pb.ListenConfigResponse, error) {
	result, err := s.uc.Listen(ctx, biz.ConfigKey{
		Namespace: req.Namespace,
		Group:     req.Group,
		DataID:    req.DataId,
	}, req.Md5, time.Duration(req.TimeoutMs)*time.Millisecond)

	if err != nil {
		return nil, err
	}

	return &pb.ListenConfigResponse{
		Namespace: result.Key.Namespace,
		Group:     result.Key.Group,
		DataId:    result.Key.DataID,
		Md5:       result.MD5,
		Changed:   result.Changed,
	}, nil
}
