package system

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"ops-server/global"
	"ops-server/initialize"
	sysModel "ops-server/model/system"
	"os"
	"path/filepath"
)

type initApi struct{}

const initOrderApi = initialize.InitOrderSystem + 1

// auto run
func init() {
	initialize.RegisterInit(initOrderApi, &initApi{})
}

func (i initApi) InitializerName() string {
	return sysModel.SysApi{}.TableName()
}

func (i *initApi) MigrateTable(ctx context.Context) (context.Context, error) {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return ctx, initialize.ErrMissingDBContext
	}
	return ctx, db.AutoMigrate(&sysModel.SysApi{})
}

func (i *initApi) TableCreated(ctx context.Context) bool {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return false
	}
	return db.Migrator().HasTable(&sysModel.SysApi{})
}

func (i *initApi) DataInserted(ctx context.Context) bool {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return false
	}
	if errors.Is(db.Where("path = ? AND method = ?", "/user/getUserinfo/", "GET").
		First(&sysModel.SysApi{}).Error, gorm.ErrRecordNotFound) {
		return false
	}
	return true
	//return false
}

func (i *initApi) InitializeData(ctx context.Context) (context.Context, error) {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return ctx, initialize.ErrMissingDBContext
	}

	apiJsonFile, err := os.Open(filepath.Join(
		global.OPS_CONFIG.Local.JsonDir,
		fmt.Sprintf("%s.json", sysModel.SysApi{}.TableName()),
	))
	if err != nil {
		return ctx, errors.Wrap(err, sysModel.SysApi{}.TableName()+"打开json文件失败!")
	}
	defer apiJsonFile.Close()

	// 读取文件内容
	apiJsonData := json.NewDecoder(apiJsonFile)

	var entities []sysModel.SysApi
	err = apiJsonData.Decode(&entities)
	if err != nil {
		return nil, errors.Wrap(err, sysModel.SysApi{}.TableName()+".json解析失败!")
	}

	for _, entity := range entities {

		//if err := db.Where("path = ?", entity.Path).FirstOrCreate(&entity).Error; err != nil {
		//	return ctx, errors.Wrap(err, sysModel.SysApi{}.TableName()+"更新数据失败!")
		//}
		err = db.Where("path = ? and method = ?", entity.Path, entity.Method).First(&sysModel.SysApi{}).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = db.Create(&entity).Error
			if err != nil {
				return ctx, errors.Wrap(err, sysModel.SysApi{}.TableName()+"表数据初始化失败!")
			}
		} else if err != nil {
			return ctx, errors.Wrap(err, sysModel.SysApi{}.TableName()+"表数据初始化失败!")
		} else {
			err = db.Where("path = ? and method = ?", entity.Path, entity.Method).Updates(&entity).Error
		}
	}

	next := context.WithValue(ctx, i.InitializerName(), entities)
	return next, nil
}
