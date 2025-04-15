package system

import (
	"context"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"ops-server/initialize"
	sysModel "ops-server/model/system"
)

const initOrderProjectAuthority = initOrderProject + initOrderAuthority

type initProjectAuthority struct{}

// auto run
func init() {
	initialize.RegisterInit(initOrderProjectAuthority, &initProjectAuthority{})
}

func (i *initProjectAuthority) MigrateTable(ctx context.Context) (context.Context, error) {
	return ctx, nil // do nothing
}

func (i *initProjectAuthority) TableCreated(ctx context.Context) bool {
	return false // always replace
}

func (i initProjectAuthority) InitializerName() string {
	return "sys_project_authority"
}

func (i *initProjectAuthority) DataInserted(ctx context.Context) bool {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return false
	}
	auth := &sysModel.SysAuthority{}
	if ret := db.Model(auth).
		Where("authority_id = ?", 888).Preload("Projects").Find(auth); ret != nil {
		if ret.Error != nil {
			return false
		}
		return len(auth.Projects) > 0
	}
	return false
}

func (i *initProjectAuthority) InitializeData(ctx context.Context) (next context.Context, err error) {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return ctx, initialize.ErrMissingDBContext
	}
	authorities, ok := ctx.Value(initAuthority{}.InitializerName()).([]sysModel.SysAuthority)
	if !ok {
		return ctx, errors.Wrap(initialize.ErrMissingDependentContext, "创建 [菜单-权限] 关联失败, 未找到权限表初始化数据")
	}
	var projects []sysModel.SysProject
	if err := db.Model(&sysModel.SysProject{}).Find(&projects).Error; err != nil {
		return next, errors.Wrap(errors.New(""), "创建 [菜单-权限] 关联失败, 未找到菜单表初始化数据")
	}
	next = ctx
	// 888
	if err = db.Model(&authorities[0]).Association("Projects").Replace(projects); err != nil {
		return next, err
	}

	return next, nil
}
