package system

import (
	"context"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"ops-server/initialize"
	sysModel "ops-server/model/system"
)

type initApiIgnore struct{}

const initOrderApiIgnore = initOrderApi + 1

// auto run
func init() {
	initialize.RegisterInit(initOrderApiIgnore, &initApiIgnore{})
}

func (i initApiIgnore) InitializerName() string {
	return sysModel.SysIgnoreApi{}.TableName()
}

func (i *initApiIgnore) MigrateTable(ctx context.Context) (context.Context, error) {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return ctx, initialize.ErrMissingDBContext
	}
	return ctx, db.AutoMigrate(&sysModel.SysIgnoreApi{})
}

func (i *initApiIgnore) TableCreated(ctx context.Context) bool {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return false
	}
	return db.Migrator().HasTable(&sysModel.SysIgnoreApi{})
}

func (i *initApiIgnore) DataInserted(ctx context.Context) bool {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return false
	}
	if errors.Is(db.Where("path = ? AND method = ?", "/base/login", "POST").
		First(&sysModel.SysIgnoreApi{}).Error, gorm.ErrRecordNotFound) {
		return false
	}
	return true
}

func (i *initApiIgnore) InitializeData(ctx context.Context) (context.Context, error) {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return ctx, initialize.ErrMissingDBContext
	}
	entities := []sysModel.SysIgnoreApi{
		{Method: "POST", Path: "/base/login"},
		{Method: "POST", Path: "/base/captcha"},
		{Method: "POST", Path: "/base/uploadFile/"},
		{Method: "POST", Path: "/base/generateExcel/"},
	}
	if err := db.Create(&entities).Error; err != nil {
		return ctx, errors.Wrap(err, sysModel.SysIgnoreApi{}.TableName()+"表数据初始化失败!")
	}
	next := context.WithValue(ctx, i.InitializerName(), entities)
	return next, nil
}
