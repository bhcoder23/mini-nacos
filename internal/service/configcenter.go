// Package service service
package service

import (
	"context"

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
