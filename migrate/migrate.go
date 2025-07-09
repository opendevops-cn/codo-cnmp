package migrate

import (
	"codo-cnmp/internal/biz"
	"context"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	"os"
	"path/filepath"
)

var ProviderSet = wire.NewSet(NewMigration)

type Migration struct {
	role      *biz.RoleUseCase
	Directory string
	log       *log.Helper
}

func NewMigration(roleUseCase *biz.RoleUseCase, logger log.Logger) *Migration {
	// 获取当前目录
	dir, err := os.Getwd()
	if err != nil {
		return nil
	}
	migrationDir := filepath.Join(dir, "migrate", "yaml", "rbac")
	return &Migration{role: roleUseCase, Directory: migrationDir, log: log.NewHelper(log.With(logger, "module", "migrate"))}
}

// ReadYaml 读取yaml文件
func (x *Migration) readYAMLFiles(ctx context.Context) error {
	// 读取目录下的所有文件
	files, err := os.ReadDir(x.Directory)
	if err != nil {
		x.log.WithContext(ctx).Errorf("读取目录失败: %v", err)
		return err
	}

	for _, file := range files {
		// 检查文件是否是 YAML 文件
		if filepath.Ext(file.Name()) == ".yaml" || filepath.Ext(file.Name()) == ".yml" {
			filePath := filepath.Join(x.Directory, file.Name())

			// 读取文件内容
			data, err := os.ReadFile(filePath)
			if err != nil {
				continue
			}
			// 解析 YAML 文件
			if err := x.migrateData(ctx, file.Name(), data); err != nil {
				x.log.WithContext(ctx).Errorf("迁移数据失败: %v", err)
			}
		}
	}

	return nil
}

func (x *Migration) migrateData(ctx context.Context, fileName string, data []byte) error {
	var role biz.RoleItem

	// 根据文件名判断角色
	switch fileName {
	case "cluster-admin.yaml":
		role.Name = "集群管理员"
		role.Description = "集群管理员"
		role.YamlStr = string(data)
		role.ISDefault = true
		role.UpdateBy = "system"
	case "edit.yaml":
		role.Name = "运维管理员"
		role.Description = "运维管理员"
		role.YamlStr = string(data)
		role.ISDefault = true
		role.UpdateBy = "system"
	case "view.yaml":
		role.Name = "只读角色"
		role.Description = "只读角色"
		role.YamlStr = string(data)
		role.ISDefault = true
		role.UpdateBy = "system"
	default:
		x.log.WithContext(ctx).Errorf("未知的角色文件: %s", fileName)
		return fmt.Errorf("未知的角色文件: %s", fileName)
	}
	// 检查角色是否已存在
	exist, err := x.role.ExistByRoleName(ctx, role.Name)
	if err != nil {
		x.log.WithContext(ctx).Errorf("检查角色是否存在失败: %v", err)
		return err
	}
	if exist {
		x.log.WithContext(ctx).Infof("角色已存在: %s", role.Name)
		return nil
	}

	// 调用 roleUseCase 创建角色
	err = x.role.CreateRole(ctx, &role)
	if err != nil {
		x.log.WithContext(ctx).Errorf("创建角色失败: %v", err)
		return fmt.Errorf("创建角色失败: %v", err)
	}
	x.log.WithContext(ctx).Infof("成功创建角色: %s", role.Name)
	return nil
}

func (x *Migration) Run(ctx context.Context) error {
	x.log.WithContext(ctx).Infof("开始迁移数据")
	if err := x.readYAMLFiles(ctx); err != nil {
		x.log.Error(err)
		return err
	}
	x.log.WithContext(ctx).Infof("迁移数据完成")
	return nil
}
