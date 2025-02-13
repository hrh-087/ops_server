package system

import (
	"errors"
	"gorm.io/gorm"
	"ops-server/global"
	"ops-server/model/common/request"
	"ops-server/model/system"
)

type ProjectService struct {
}

var ProjectServiceApp = new(ProjectService)

// CreateProject
// @author:rh
// @description: 新增项目
// @param: project model.SysProject
// @ return: err error
func (projectService *ProjectService) CreateProject(project system.SysProject) (err error) {
	if !errors.Is(global.OPS_DB.Where("project_name = ?", project.ProjectName).First(&system.SysProject{}).Error, gorm.ErrRecordNotFound) {
		return errors.New("项目名称重复")
	}
	return global.OPS_DB.Create(&project).Error
}

// UpdateProject
// @author:rh
// @description: 更新项目
// @param: project model.SysProject
// @ return: err error
func (projectService *ProjectService) UpdateProject(project system.SysProject) (err error) {
	var oldProject system.SysProject

	err = global.OPS_DB.First(&oldProject, "id = ?", project.ID).Error
	if err != nil {
		return err
	}

	return global.OPS_DB.Save(&project).Error
}

// GetProjectList
// @author:rh
// @description: 获取项目列表
// @param: project model.SysProject info request.PageInfo
// @ return: err error
func (projectService *ProjectService) GetProjectList(project system.SysProject, info request.PageInfo) (list interface{}, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	db := global.OPS_DB.Model(&system.SysProject{})
	var projectList []system.SysProject

	if project.ProjectName != "" {
		db = db.Where("project_name LIKE ?", "%"+project.ProjectName+"%")
	}

	err = db.Count(&total).Error

	if err != nil {
		return projectList, total, err
	}

	db = db.Limit(limit).Offset(offset)
	OrderStr := "id desc"
	err = db.Preload("Authorities").Order(OrderStr).Find(&projectList).Error
	return projectList, total, err
}

// GetApiById
// @author:
// @function: GetApiById
// @description: 根据id获取project
// @param: id float64
// @return: project model.SysProject, err error
func (projectService *ProjectService) GetProjectById(id int) (project system.SysProject, err error) {
	err = global.OPS_DB.Preload("Authorities").First(&project, "id = ?", id).Error
	return
}

// DeleteProject
// @author:
// @function: DeleteProject
// @description: 删除project
// @param: project system.SysProject
// @return: err error
func (projectService *ProjectService) DeleteProject(project system.SysProject) (err error) {
	var entity system.SysProject
	err = global.OPS_DB.Preload("Authorities").Where("id = ?", project.ID).First(&entity).Error
	if err != nil {
		return err
	}
	return global.OPS_DB.Transaction(func(tx *gorm.DB) error {
		if len(entity.Authorities) > 0 {
			return errors.New("此角色有用户正在使用此项目,请解除绑定后再试")
		}
		return tx.Delete(&system.SysProject{}, "id = ?", entity.ID).Error
	})
}

// GetAllProject
// @author:rh
// @description: 获取全部项目
// @param: project model.SysProject info request.PageInfo
// @ return: err error
func (projectService *ProjectService) GetAllProject() (projects []system.SysProject, err error) {

	db := global.OPS_DB.Model(&system.SysProject{})
	var projectList []system.SysProject

	if err != nil {
		return projectList, err
	}

	err = db.Preload("Authorities").Find(&projectList).Error
	return projectList, err
}

// GetAuthorityProject
// @author:rh
// @description: 获取角色项目
// @param:
// @ return:
func (projectService *ProjectService) GetAuthorityProject(info *request.GetAuthorityId) (projects []system.SysProject, err error) {

	var projectList []system.SysProject
	var projectAuthority []system.SysProjectAuthority
	err = global.OPS_DB.Where("sys_authority_authority_id = ?", info.AuthorityId).Find(&projectAuthority).Error
	if err != nil {
		return
	}

	var projectIds []uint

	for i := range projectAuthority {
		projectIds = append(projectIds, projectAuthority[i].SysProjectId)
	}

	err = global.OPS_DB.Where("id in (?)", projectIds).Find(&projectList).Error

	return projectList, err
}

func (projectService *ProjectService) SetAuthorityProject(authorityId uint, projects []system.SysProject) (err error) {
	var authority system.SysAuthority
	authority.AuthorityId = authorityId
	authority.Projects = projects
	err = AuthorityServiceApp.SetProjectAuthority(&authority)
	return err
}

func (projectService *ProjectService) CheckProject(roleId, projectId string) (err error) {
	if err = global.OPS_DB.Model(&system.SysProjectAuthority{}).Where("sys_authority_authority_id = ? AND sys_project_id = ?", roleId, projectId).First(&system.SysProjectAuthority{}).Error; err != nil {
		return err
	}
	return nil
}
