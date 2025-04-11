package system

import (
	"context"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"ops-server/initialize"
	sysModel "ops-server/model/system"
)

const initOrderProject = initOrderAuthority + 1

type initProject struct{}

// auto run
func init() {
	initialize.RegisterInit(initOrderProject, &initProject{})
}

func (i initProject) InitializerName() string {
	return sysModel.SysProject{}.TableName()
}

func (i *initProject) MigrateTable(ctx context.Context) (context.Context, error) {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return ctx, initialize.ErrMissingDBContext
	}
	return ctx, db.AutoMigrate(
		&sysModel.SysProject{},
	)
}

func (i *initProject) TableCreated(ctx context.Context) bool {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return false
	}
	m := db.Migrator()
	return m.HasTable(&sysModel.SysProject{})
}

func (i *initProject) DataInserted(ctx context.Context) bool {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return false
	}
	if errors.Is(db.Where("project_name = ?", "初始项目").First(&sysModel.SysProject{}).Error, gorm.ErrRecordNotFound) { // 判断是否存在数据
		return false
	}
	return true
	//return false
}

func (i *initProject) InitializeData(ctx context.Context) (next context.Context, err error) {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return ctx, initialize.ErrMissingDBContext
	}

	entities := []sysModel.SysProject{
		{
			ProjectName:     "初始项目",
			SvnUrl:          "svn://127.0.0.1/ops/",
			GmUrl:           "http://127.0.0.1:8080/",
			GatewayUrl:      "http://127.0.0.1:8080/",
			ClientSvnUrl:    "svn://127.0.0.1/ops/",
			ClientConfigDir: "config/",
			ConfigDir:       "config/",
			IsTest:          false,
		},
	}
	if err = db.Create(&entities).Error; err != nil {
		return ctx, errors.Wrap(err, sysModel.SysProject{}.TableName()+"表数据初始化失败!")
	}
	next = context.WithValue(ctx, i.InitializerName(), entities)
	return next, err
}
