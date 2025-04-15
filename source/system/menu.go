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

const initOrderMenu = initOrderAuthority + 1

type initMenu struct{}

// auto run
func init() {
	initialize.RegisterInit(initOrderMenu, &initMenu{})
}

func (i initMenu) InitializerName() string {
	return sysModel.SysBaseMenu{}.TableName()
}

func (i *initMenu) MigrateTable(ctx context.Context) (context.Context, error) {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return ctx, initialize.ErrMissingDBContext
	}
	return ctx, db.AutoMigrate(
		&sysModel.SysBaseMenu{},
		&sysModel.SysBaseMenuParameter{},
		&sysModel.SysBaseMenuBtn{},
	)
}

func (i *initMenu) TableCreated(ctx context.Context) bool {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return false
	}
	m := db.Migrator()
	return m.HasTable(&sysModel.SysBaseMenu{}) &&
		m.HasTable(&sysModel.SysBaseMenuParameter{}) &&
		m.HasTable(&sysModel.SysBaseMenuBtn{})
}

func (i *initMenu) DataInserted(ctx context.Context) bool {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return false
	}
	if errors.Is(db.Where("path = ?", "system").First(&sysModel.SysBaseMenu{}).Error, gorm.ErrRecordNotFound) { // 判断是否存在数据
		return false
	}
	return true
	//return false
}

func (i *initMenu) InitializeData(ctx context.Context) (next context.Context, err error) {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return ctx, initialize.ErrMissingDBContext
	}

	menuJsonFile, err := os.Open(filepath.Join(
		global.OPS_CONFIG.Local.JsonDir,
		fmt.Sprintf("%s.json", sysModel.SysBaseMenu{}.TableName()),
	))

	if err != nil {
		return ctx, errors.Wrap(err, sysModel.SysBaseMenu{}.TableName()+"打开json文件失败!")
	}
	defer menuJsonFile.Close()

	// 读取文件内容
	jsonData := json.NewDecoder(menuJsonFile)

	var entities []sysModel.SysBaseMenu
	err = jsonData.Decode(&entities)
	if err != nil {
		return ctx, errors.Wrap(err, sysModel.SysBaseMenu{}.TableName()+".json解析失败!")
	}

	// 根据父菜单与子菜单的关系，进行递归创建菜单

	var menuIdMap = make(map[uint]bool)
	var menuList []sysModel.SysBaseMenu
	for _, entity := range entities {
		// 如果ID已经存在，则跳过当前循环
		if ok := menuIdMap[entity.ID]; ok {
			continue
		}
		entity.Children = getMenuChildrenList(entity, entities, menuIdMap)

		menuIdMap[entity.ID] = true
		menuList = append(menuList, entity)
	}

	//menuData, err := json.MarshalIndent(menuList, "", "  ")
	//if err != nil {
	//	return ctx, errors.Wrap(err, sysModel.SysBaseMenu{}.TableName()+"json序列化失败!")
	//}
	//fmt.Println(string(menuData))

	for menuIndex := range menuList {
		err = createMenu(db, &menuList[menuIndex], 0)
		if err != nil {
			return ctx, errors.Wrap(err, sysModel.SysBaseMenu{}.TableName()+"表数据初始化失败!")
		}

	}

	next = context.WithValue(ctx, i.InitializerName(), menuList)
	return next, nil
}

func getMenuChildrenList(menu sysModel.SysBaseMenu, allMenu []sysModel.SysBaseMenu, menuMap map[uint]bool) []sysModel.SysBaseMenu {
	var menuList []sysModel.SysBaseMenu

	for _, v := range allMenu {

		if v.ParentId == menu.ID {
			menuMap[v.ID] = true

			v.Children = getMenuChildrenList(v, allMenu, menuMap)
			menuList = append(menuList, v)
		}
	}
	return menuList
}

func createMenu(db *gorm.DB, menu *sysModel.SysBaseMenu, parentId uint) error {
	menu.ID = 0
	menu.ParentId = parentId
	err := db.Debug().Create(menu).Error
	if err != nil {
		return err
	}

	for i := range menu.Children {
		err = createMenu(db, &menu.Children[i], menu.ID)
		if err != nil {
			return err
		}
	}

	return err
}
