package handler

import (
	"context"

	"github.com/jinmukeji/jiujiantang-services/analysis/config"
	proto "github.com/jinmukeji/proto/v3/gen/micro/idl/partner/xima/analysis/v1"
)

// GetVersion 获取服务版本信息
func (j *AnalysisManagerService) GetVersion(ctx context.Context, req *proto.GetVersionRequest, resp *proto.GetVersionResponse) error {
	resp.ServiceName = config.FullServiceName()
	resp.Version = config.ProductVersion
	return nil
}
