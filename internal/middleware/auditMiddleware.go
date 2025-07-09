package middleware

import (
	"codo-cnmp/common/result"
	"codo-cnmp/common/utils"
	"codo-cnmp/common/xerr"
	"codo-cnmp/internal/biz"
	"codo-cnmp/internal/dep"
	"codo-cnmp/pb"
	"context"
	"encoding/json"
	"github.com/Ccheers/protoc-gen-go-kratos-http/audit"
	"github.com/Ccheers/protoc-gen-go-kratos-http/kcontext"
	"github.com/go-kratos/kratos/v2/log"
	kmiddleware "github.com/go-kratos/kratos/v2/middleware"
	"go.opentelemetry.io/otel/trace"
	"net"
	"strconv"
	"strings"
	"time"
)

type ResourceType string

var (
	RoleResourceKind       = "角色"
	GameServerResourceKind = "进程"
	UserGroupResourceKind  = "用户组授权"
)

const (
	ClusterResourceType    ResourceType = "cluster"
	UserGroupResourceType  ResourceType = "userGroup"
	RoleResourceType       ResourceType = "role"
	GameServerResourceType ResourceType = "gameServer"
)

type AuditMiddleware struct {
	audit   *biz.AuditLogUseCase
	cluster *biz.ClusterUseCase
	logger  *log.Helper
	kafka   dep.IKafka
	ug      *biz.UserGroupUseCase
	role    *biz.RoleUseCase
}

func NewAuditMiddleware(audit *biz.AuditLogUseCase, cluster *biz.ClusterUseCase, logger log.Logger,
	kafka dep.IKafka, ug *biz.UserGroupUseCase, role *biz.RoleUseCase) *AuditMiddleware {
	return &AuditMiddleware{
		audit:   audit,
		cluster: cluster,
		logger:  log.NewHelper(log.With(logger, "module", "middleware/audit")),
		kafka:   kafka,
		ug:      ug,
		role:    role,
	}
}

func GetClientIP(ctx context.Context) string {
	httpTr, err := extraKratosHTTPTransport(ctx)
	if err != nil {
		return ""
	}
	// 检查 X-Forwarded-For
	if forwardedFor := httpTr.RequestHeader().Get("X-Forwarded-For"); forwardedFor != "" {
		if ip := getForwardedIP(forwardedFor); ip != "" {
			return ip
		}
	}

	// 检查 X-Real-IP
	if clientIP := httpTr.RequestHeader().Get("X-Real-IP"); clientIP != "" {
		if ip := cleanIP(clientIP); ip != "" {
			return ip
		}
	}

	// 获取 RemoteAddr
	if remoteAddr := httpTr.Request().RemoteAddr; remoteAddr != "" {
		if ip := getRemoteIP(remoteAddr); ip != "" {
			return ip
		}
	}

	return ""
}

// cleanIP 清理和验证 IP 地址
func cleanIP(ip string) string {
	ip = strings.TrimSpace(ip)
	if parsed := net.ParseIP(ip); parsed != nil {
		return ip
	}
	return ""
}

// getForwardedIP 从 X-Forwarded-For 获取第一个有效 IP
func getForwardedIP(forwarded string) string {
	ips := strings.Split(forwarded, ",")
	for _, ip := range ips {
		if cleaned := cleanIP(ip); cleaned != "" {
			return cleaned
		}
	}
	return ""
}

// getRemoteIP 从 RemoteAddr (IP:Port) 获取 IP
func getRemoteIP(remoteAddr string) string {
	// 尝试分离 host 和 port
	host, _, err := net.SplitHostPort(remoteAddr)
	if err != nil {
		// 如果分离失败，可能没有端口，直接尝试作为 IP
		return cleanIP(remoteAddr)
	}
	return cleanIP(host)
}

// extractMetaData 提取元数据
func extractMetaData(metas []audit.ExtractedMeta) map[string]string {
	data := make(map[string]string, len(metas))
	for _, meta := range metas {
		data[meta.Key] = meta.Value
	}
	return data
}

// formatJSON 格式化 JSON
func formatJSON(v interface{}) string {
	if v == nil {
		return ""
	}
	bs, err := json.Marshal(v)
	if err != nil {
		return ""
	}
	return string(bs)
}

// formatDuration 格式化持续时间
func formatDuration(d time.Duration) string {
	return strconv.FormatInt(d.Milliseconds(), 10)
}

// getStatus 获取状态码
func getStatus(err error) int {
	if err != nil {
		return int(pb.OperationStatus_Fail)
	}
	return int(pb.OperationStatus_Success)
}

type AuditRule struct {
	Module       string
	Action       string
	ExtractMetas map[string]string
}

func (x *AuditMiddleware) getClusterName(ctx context.Context, cluster, kind string) string {
	if ResourceType(kind) == ClusterResourceType {
		// 集群参数为数字, 则为集群 ID, 需要转换为集群名称
		clusterId, err := strconv.Atoi(cluster)
		if err != nil {
			// 处理转换错误
			x.logger.WithContext(ctx).Errorf("集群ID转换失败: %v", err)
			return cluster // 返回原始的 cluster ID 作为 fallback
		}
		clusterItem, err := x.cluster.GetClusterByID(ctx, uint32(clusterId))
		if err != nil {
			x.logger.WithContext(ctx).Errorf("获取集群信息失败: %v", err)
			return cluster // 如果获取失败，返回原始 cluster
		}
		return clusterItem.Name
	}
	// 如果不是集群资源类型，直接返回 cluster
	return cluster
}

func (x *AuditMiddleware) GetRoleName(ctx context.Context, roleId string) string {
	intRoleId, err := strconv.Atoi(roleId)
	if err != nil {
		x.logger.WithContext(ctx).Errorf("获取角色信息失败: %v", err) // 返回空字符串表示无法获取角色信息
	}
	role, err := x.role.GetRole(ctx, uint32(intRoleId))
	if err != nil {
		x.logger.WithContext(ctx).Errorf("获取角色信息失败: %v", err)
		return "" // 如果获取失败，返回空字符串
	}
	return role.Name
}

func (x *AuditMiddleware) getUserGroupName(ctx context.Context, userGroupId string) string {
	if userGroupId != "" {
		intUserGroupId, err := strconv.Atoi(userGroupId)
		if err != nil {
			x.logger.WithContext(ctx).Errorf("获取用户组信息失败: %v", err) // 返回空字符串表示无法获取用户组信息
		}
		input := strings.Trim(userGroupId, "[]")
		if input != "" {
			// 尝试将其按逗号分割
			parts := strings.Split(input, ",")
			if len(parts) > 0 {
				// 尝试将切片的第一个元素转换为整数
				if val, err := strconv.Atoi(parts[0]); err == nil {
					intUserGroupId = val
				}
			}
		}
		userGroups, err := x.ug.GetUserGroup(ctx, uint32(intUserGroupId))
		if err != nil {
			x.logger.WithContext(ctx).Errorf("获取用户组信息失败: %v", err)
			return "" // 如果获取失败，返回空字符串
		}
		return userGroups.RoleName
	}
	return ""
}

// Server 审计中间件
func (x *AuditMiddleware) Server() kmiddleware.Middleware {
	return func(handler kmiddleware.Handler) kmiddleware.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			startTime := time.Now()
			httpTr, err := extraKratosHTTPTransport(ctx)
			if err != nil {
				return nil, err
			}
			auditInfo, ok := kcontext.GetKHTTPAuditContextWithContext(ctx)
			if ok {
				// 请求路径
				path := httpTr.Request().RequestURI

				// 提取元数据
				metaData := extractMetaData(auditInfo.ExtractedMetas)

				// 执行请求
				reply, replyErr := handler(ctx, req)
				replyJsonData := formatJSON(reply)
				var responseBody string
				if replyErr != nil {
					errCode := result.KHttpParseErr(replyErr)
					responseBean := result.Error(errCode, xerr.MapErrMsg(errCode), replyErr.Error())
					responseBody = formatJSON(responseBean)
				} else {
					responseBean := result.Success(replyJsonData)
					responseBody = formatJSON(responseBean)
				}

				// 请求耗时
				duration := time.Since(startTime)

				// 操作人
				userName, err := utils.GetUserNameFromCtx(ctx)
				if err != nil {
					userName = ""
				}

				resourceType := metaData["kind"]
				switch ResourceType(resourceType) {
				case ClusterResourceType:
					// 处理集群资源类型
					clusterName := x.getClusterName(ctx, metaData["cluster"], resourceType)
					metaData["cluster"] = clusterName
				case UserGroupResourceType:
					// 处理用户组资源类型
					userGroupName := x.getUserGroupName(ctx, metaData["name"])
					if userGroupName != "" {
						metaData["name"] = userGroupName
					}
					metaData["kind"] = UserGroupResourceKind
				case RoleResourceType:
					// 处理角色资源类型
					roleName := x.GetRoleName(ctx, metaData["name"])
					if roleName != "" {
						metaData["name"] = roleName
					}
					metaData["kind"] = RoleResourceKind
				default:
				}

				//创建新的 context，避免使用可能被取消的原始 context
				auditCtx := context.Background()
				traceID := trace.SpanContextFromContext(ctx).TraceID()

				// 创建 AuditLogItem
				auditLogItem := biz.AuditLogItem{
					UserName:      userName,
					ClientIP:      GetClientIP(ctx),
					Module:        auditInfo.Module,
					Action:        auditInfo.Action,
					Cluster:       metaData["cluster"],
					Namespace:     metaData["namespace"],
					ResourceType:  metaData["kind"],
					ResourceName:  metaData["name"],
					RequestPath:   path,
					RequestBody:   formatJSON(req),
					ResponseBody:  responseBody,
					Status:        getStatus(replyErr),
					Duration:      formatDuration(duration),
					OperationTime: startTime.Format(time.DateTime),
					TraceID:       traceID.String(),
				}

				// 记录审计日志
				auditLog := &biz.CreateAuditLogRequest{
					AuditLogItem: auditLogItem,
				}
				go func(ctx context.Context) {
					if err := x.audit.CreateAuditLog(ctx, auditLog); err != nil {
						x.logger.WithContext(ctx).Errorf("记录审计日志失败: %v", err)
					}
				}(auditCtx)
				if x.kafka != nil {
					// 发送消息到 Kafka
					go func(ctx context.Context, log biz.AuditLogItem) {
						defer func(ctx context.Context) {
							if r := recover(); r != nil {
								x.logger.Errorf("发送消息到 Kafka 时发生 panic: %v", r)
								x.kafka.Close(ctx)
							}
						}(ctx)

						auditLogBytes, err := json.Marshal(auditLogItem)
						if err != nil {
							x.logger.Errorf("序列化审计日志失败	: %v", err)
						}
						err = x.kafka.SendMessage(ctx, auditLogBytes)
						if err != nil {
							x.logger.Errorf("发送审计日志到 Kafka 失败: %v", err)
						} else {
							x.logger.WithContext(ctx).Debug("发送审计日志到 Kafka成功")
						}
					}(context.Background(), auditLogItem)
				}

				return reply, replyErr
			}
			return handler(ctx, req)
		}
	}
}
