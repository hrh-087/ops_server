package system

import (
	"context"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"ops-server/initialize"
	sysModel "ops-server/model/system"
	"ops-server/utils"
)

const initOrderAuthority = initOrderCasbin + 1

type initAuthority struct{}

func init() {
	initialize.RegisterInit(initOrderAuthority, &initAuthority{})
}

func (i *initAuthority) MigrateTable(ctx context.Context) (context.Context, error) {

	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return ctx, initialize.ErrMissingDBContext
	}
	return ctx, db.AutoMigrate(&sysModel.SysAuthority{})
}

func (i *initAuthority) TableCreated(ctx context.Context) bool {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return false
	}
	return db.Migrator().HasTable(&sysModel.SysAuthority{})
}

func (i initAuthority) InitializerName() string {
	return sysModel.SysAuthority{}.TableName()
}

func (i *initAuthority) DataInserted(ctx context.Context) bool {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return false
	}
	if errors.Is(db.Where("authority_id = ?", "888").
		First(&sysModel.SysAuthority{}).Error, gorm.ErrRecordNotFound) { // 判断是否存在数据
		return false
	}
	return true
}

func (i *initAuthority) InitializeData(ctx context.Context) (context.Context, error) {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return ctx, initialize.ErrMissingDBContext
	}
	entities := []sysModel.SysAuthority{
		{AuthorityId: 888, AuthorityName: "普通用户", ParentId: utils.Pointer[uint](0), DefaultRouter: "dashboard"},
	}

	if err := db.Create(&entities).Error; err != nil {
		return ctx, errors.Wrapf(err, "%s表数据初始化失败!", sysModel.SysAuthority{}.TableName())
	}
	// data authority
	//if err := db.Model(&entities[0]).Association("DataAuthorityId").Replace(
	//	[]*sysModel.SysAuthority{
	//		{AuthorityId: 888},
	//	}); err != nil {
	//	return ctx, errors.Wrapf(err, "%s表数据初始化失败!",
	//		db.Model(&entities[0]).Association("DataAuthorityId").Relationship.JoinTable.Name)
	//}
	//if err := db.Model(&entities[1]).Association("DataAuthorityId").Replace(
	//	[]*sysModel.SysAuthority{
	//		{AuthorityId: 9528},
	//		{AuthorityId: 8881},
	//	}); err != nil {
	//	return ctx, errors.Wrapf(err, "%s表数据初始化失败!",
	//		db.Model(&entities[1]).Association("DataAuthorityId").Relationship.JoinTable.Name)
	//}

	next := context.WithValue(ctx, i.InitializerName(), entities)
	return next, nil
}
