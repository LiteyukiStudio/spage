package agent

import (
	"context"
	"fmt"
	agentUtils "github.com/LiteyukiStudio/spage/agent/utils"
	"github.com/LiteyukiStudio/spage/pkg/config"
	"github.com/LiteyukiStudio/spage/pkg/utils"
	pb "github.com/LiteyukiStudio/spage/protos/result/protos/source"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io"
	"os"
	"path/filepath"
)

type ServerVisit struct {
	pb.UnimplementedAgentServiceServer
}

// RegisterGrpcApps 注册grpc服务
func RegisterGrpcApps() (*grpc.Server, error) {
	// 创建gRPC服务器并添加拦截器
	s := grpc.NewServer(
		grpc.UnaryInterceptor(UnaryInterceptor),
	)
	pb.RegisterAgentServiceServer(s, &ServerVisit{})
	return s, nil
}

/*
站点保存路径规则：
站点保存路径：<static_path>/<OwnerName>/<ProjectName>/<Name>
*/
func getSitePath(ownerName, projectName, siteName string) (string, error) {
	sitePath := filepath.Join(config.AgentConfig.Service.Static, ownerName, projectName, siteName)
	// 判断路径是否存在
	if _, err := os.Stat(sitePath); os.IsNotExist(err) {
		// 创建目录
		if err = os.MkdirAll(sitePath, os.ModePerm); err != nil {
			return "", err
		}
	}
	return sitePath, nil
}

func (ServerVisit) CreateSite(ctx context.Context, request *pb.CreateSiteRequest) (response *pb.CreateSiteResponse, err error) {
	sitePath, err := getSitePath(request.OwnerName, request.ProjectName, request.Name)
	if err != nil {
		return
	}
	// TODO 在caddy创建站点
	response.Success = true
	response.Message = fmt.Sprintf("站点创建成功，保存路径: %s", sitePath)
	return
}

func (ServerVisit) UpdateSite(ctx context.Context, request *pb.UpdateSiteRequest) (response *pb.UpdateSiteResponse, err error) {
	sitePath, err := getSitePath(request.OwnerName, request.ProjectName, request.Name)
	if err != nil {
		return
	}
	// TODO 在caddy更新站点
	response.Success = true
	response.Message = sitePath
	return
}

func (ServerVisit) DeleteSite(ctx context.Context, request *pb.DeleteSiteRequest) (response *pb.DeleteSiteResponse, err error) {
	sitePath, err := getSitePath(request.OwnerName, request.ProjectName, request.Name)
	if err != nil {
		return
	}
	err = os.RemoveAll(sitePath)
	// TODO 在caddy删除站点
	response.Success = true
	return
}

func (ServerVisit) GetSite(ctx context.Context, request *pb.GetSiteRequest) (response *pb.GetSiteResponse, err error) {
	sitePath, err := getSitePath(request.OwnerName, request.ProjectName, request.Name)
	if err != nil {
		return
	}
	hash, updateAt, err := agentUtils.GetHash(sitePath)
	if err != nil {
		return
	}
	response.ReleaseHash = hash
	response.ReleaseUpdateAt = updateAt.String()
	response.Success = true
	response.SitePath = sitePath
	/*
		TODO 获取caddy中的站点信息
		response.Domains = []string{"example.com"}
		response.SubDomain = "sub.example.com"
	*/
	return
}

func (ServerVisit) UploadRelease(stream grpc.ClientStreamingServer[pb.UploadReleaseRequest, pb.UploadReleaseResponse]) error {
	var sitePath string
	var contentBytes []byte
	var outputFile *os.File
	var outputPath string
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			// 流结束，关闭文件并返回响应
			if outputFile != nil {
				outputFile.Close()
			}
			// 更新站点 hash
			hash, err := agentUtils.UpdateSiteHash(sitePath)
			if err != nil {
				return err
			}
			// 对比hash
			if hash != req.ReleaseHash {
				return status.Errorf(codes.Internal, "站点 hash 不一致")
			}
			// 解压 zip
			err = utils.UnzipFromBytes(contentBytes, sitePath)
			if err != nil {
				return err
			}
			return stream.SendAndClose(&pb.UploadReleaseResponse{
				Success: true,
				Message: fmt.Sprintf("文件上传成功，保存路径: %s", outputPath),
			})
		}
		if err != nil {
			if outputFile != nil {
				outputFile.Close()
				// 尝试删除可能已部分写入的文件
				os.Remove(outputPath)
			}
			return status.Errorf(codes.Internal, "接收流数据错误: %v", err)
		}

		// 首个请求处理站点路径
		if sitePath == "" {
			sitePath, err = getSitePath(req.OwnerName, req.ProjectName, req.SiteName)
			if err != nil {
				return err
			}
			// 创建输出文件
			outputPath = filepath.Join(sitePath, agentUtils.ReleaseName)
			outputFile, err = os.Create(outputPath)
			if err != nil {
				return status.Errorf(codes.Internal, "创建文件失败: %v", err)
			}
		}

		// 写入文件内容
		content := req.GetContent()
		if len(content) > 0 {
			contentBytes = append(contentBytes, content...)
			if _, err := outputFile.Write(content); err != nil {
				outputFile.Close()
				os.Remove(outputPath)
				return status.Errorf(codes.Internal, "写入文件失败: %v", err)
			}
		}
	}
}
