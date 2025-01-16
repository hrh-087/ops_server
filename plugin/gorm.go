package plugin

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// 定义一个gorm插件, 用来过滤项目的数据
type ProjectFilterPlugin struct{}

func (p ProjectFilterPlugin) Name() string {
	return "project_filter"
}

func (p ProjectFilterPlugin) Initialize(db *gorm.DB) error {
	db.Callback().Query().Before("gorm:query").Register("project_filter:query", p.queryFilter)
	db.Callback().Create().Before("gorm:create").Register("project_filter:create", p.createFilter)
	db.Callback().Update().Before("gorm:update").Register("project_filter:update", p.queryFilter)
	db.Callback().Delete().Before("gorm:delete").Register("project_filter:delete", p.queryFilter)
	return nil
}

func (p *ProjectFilterPlugin) queryFilter(db *gorm.DB) {
	// 从 Settings 中读取 'skip_project_filter' 标志
	if skip, ok := db.Statement.Settings.Load("skip_project_filter"); ok && skip == true {
		// 如果设置为 true，则跳过该过滤逻辑
		return
	}

	projectId, ok := db.Statement.Context.Value("projectId").(string)
	if !ok {
		return
	}

	// 动态检查是否有 project_id 字段
	if !hasProjectIDField(db.Statement) {
		return // 如果模型中没有 project_id 字段，不添加过滤条件
	}

	if ok && projectId != "" {
		db.Statement.AddClause(clause.Where{
			Exprs: []clause.Expression{
				clause.Eq{
					Column: "project_id",
					Value:  projectId,
				},
			},
		})
	}
}

func (p *ProjectFilterPlugin) createFilter(db *gorm.DB) {
	projectId, ok := db.Statement.Context.Value("projectId").(string)

	if !ok {
		return
	}

	// 动态检查是否有 project_id 字段
	if !hasProjectIDField(db.Statement) {
		return // 如果模型中没有 project_id 字段，不添加过滤条件
	}

	if ok && projectId != "" {
		db.Statement.SetColumn("project_id", projectId)
	}
}

func hasProjectIDField(stmt *gorm.Statement) bool {
	// 遍历模型字段，检查是否有 `project_id` 字段
	if stmt.Schema == nil {
		return false
	}
	_, ok := stmt.Schema.FieldsByDBName["project_id"]
	return ok
}
