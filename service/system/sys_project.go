package system

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid/v5"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"ops-server/global"
	"ops-server/job/task"
	"ops-server/model/common/request"
	"ops-server/model/system"
	"ops-server/utils"
	"time"
)

type ProjectService struct {
}

//var ProjectServiceApp = new(ProjectService)

// CreateProject
// @author:rh
// @description: 新增项目
// @param: project model.SysProject
// @ return: err error
func (p *ProjectService) CreateProject(project system.SysProject) (err error) {
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
func (p *ProjectService) UpdateProject(project system.SysProject) (err error) {
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
func (p *ProjectService) GetProjectList(project system.SysProject, info request.PageInfo) (list interface{}, total int64, err error) {
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
func (p *ProjectService) GetProjectById(id int) (project system.SysProject, err error) {
	err = global.OPS_DB.Preload("Authorities").First(&project, "id = ?", id).Error
	return
}

// DeleteProject
// @author:
// @function: DeleteProject
// @description: 删除project
// @param: project system.SysProject
// @return: err error
func (p *ProjectService) DeleteProject(project system.SysProject) (err error) {
	var entity system.SysProject
	err = global.OPS_DB.Preload("Authorities").Where("id = ?", project.ID).First(&entity).Error
	if err != nil {
		return err
	}
	return global.OPS_DB.Transaction(func(tx *gorm.DB) error {
		if len(entity.Authorities) > 0 {
			return errors.New("此角色有用户正在使用此项目,请解除绑定后再试")
		}

		// 删除项目关联的gameServer
		if err = tx.Delete(&system.SysGameServer{}, "project_id = ?", entity.ID).Error; err != nil {
			return err
		}
		if err = tx.Delete(&system.SysGamePlatform{}, "project_id = ?", entity.ID).Error; err != nil {
			return err
		}
		if err = tx.Delete(&system.SysGameType{}, "project_id = ?", entity.ID).Error; err != nil {
			return err
		}
		if err = tx.Delete(&system.SysAssetsServer{}, "project_id = ?", entity.ID).Error; err != nil {
			return err
		}
		if err = tx.Delete(&system.SysAssetsListener{}, "project_id = ?", entity.ID).Error; err != nil {
			return err
		}
		if err = tx.Delete(&system.SysAssetsLb{}, "project_id = ?", entity.ID).Error; err != nil {
			return err
		}
		if err = tx.Delete(&system.SysAssetsRedis{}, "project_id = ?", entity.ID).Error; err != nil {
			return err
		}
		if err = tx.Delete(&system.SysAssetsMysql{}, "project_id = ?", entity.ID).Error; err != nil {
			return err
		}
		if err = tx.Delete(&system.SysAssetsMongo{}, "project_id = ?", entity.ID).Error; err != nil {
			return err
		}
		if err = tx.Delete(&system.SysAssetsKafka{}, "project_id = ?", entity.ID).Error; err != nil {
			return err
		}

		return tx.Delete(&system.SysProject{}, "id = ?", entity.ID).Error
	})
}

// GetAllProject
// @author:rh
// @description: 获取全部项目
// @param: project model.SysProject info request.PageInfo
// @ return: err error
func (p *ProjectService) GetAllProject() (projects []system.SysProject, err error) {

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
func (p *ProjectService) GetAuthorityProject(info *request.GetAuthorityId) (projects []system.SysProject, err error) {

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

func (p *ProjectService) SetAuthorityProject(authorityId uint, projects []system.SysProject) (err error) {
	var authority system.SysAuthority
	authority.AuthorityId = authorityId
	authority.Projects = projects
	err = AuthorityServiceApp.SetProjectAuthority(&authority)
	return err
}

func (p *ProjectService) CheckProject(roleId, projectId string) (err error) {
	if err = global.OPS_DB.Model(&system.SysProjectAuthority{}).Where("sys_authority_authority_id = ? AND sys_project_id = ?", roleId, projectId).First(&system.SysProjectAuthority{}).Error; err != nil {
		return err
	}
	return nil
}

// InitProject
// 初始化项目
func (p ProjectService) InitProject(ctx *gin.Context, project system.SysProject) (jobId uuid.UUID, err error) {

	var job system.Job
	var t system.JobTask
	var taskList []system.JobTask

	if err := global.OPS_DB.First(&project, "id = ?", project.ID).Error; err != nil {
		return uuid.UUID{}, err
	}

	jobId = uuid.Must(uuid.NewV4())
	taskId := uuid.Must(uuid.NewV4())

	taskInfo, err := task.NewInitProjectTask(task.InitProjectParams{
		TaskId:  taskId,
		Project: project,
	})

	if err != nil {
		global.OPS_LOG.Error("添加任务到队列失败", zap.String("jobId", jobId.String()), zap.Error(err))
	}

	t.JobId = jobId
	t.AsynqId = taskInfo.ID
	t.TaskId = taskId
	t.Status = taskInfo.State.String()
	t.HostName = global.OPS_CONFIG.Ops.Name
	t.HostIp = global.OPS_CONFIG.Ops.Host
	t.CreateAt = time.Now()

	if err := global.OPS_DB.Create(&t).Error; err != nil {
		global.OPS_LOG.Error("创建任务失败", zap.String("jobId", jobId.String()), zap.String("taskId", taskId.String()), zap.Error(err))
	}

	claims, _ := utils.GetClaims(ctx)

	job.JobId = jobId
	job.Name = "初始化项目"
	job.Status = 1
	job.Type = task.InitProjectTypeName
	job.Creator = claims.Username
	job.Tasks = taskList
	job.ProjectId = project.ID
	job.CreateAt = time.Now()

	// 创建作业任务
	err = global.OPS_DB.Create(&job).Error

	if err != nil {
		global.OPS_LOG.Error("创建作业任务失败", zap.String("jobId", jobId.String()), zap.Error(err))
		return
	}
	return
}

func (p ProjectService) SetProjectAuthorities(id uint, authorityIds []uint) (err error) {
	return global.OPS_DB.Transaction(func(tx *gorm.DB) error {
		var project system.SysProject

		//if len(authorityIds) == 0 {
		//	return errors.New("用户组不能为空")
		//}

		TxErr := tx.Where("id = ?", id).First(&project).Error
		if TxErr != nil {
			global.OPS_LOG.Debug(TxErr.Error())
			return errors.New("查询用户数据失败")
		}
		TxErr = tx.Delete(&[]system.SysProjectAuthority{}, "sys_project_id = ?", id).Error
		if TxErr != nil {
			return TxErr
		}
		var useAuthority []system.SysProjectAuthority
		for _, v := range authorityIds {
			useAuthority = append(useAuthority, system.SysProjectAuthority{
				SysProjectId:            id,
				SysAuthorityAuthorityId: v,
			})
		}
		TxErr = tx.Create(&useAuthority).Error
		if TxErr != nil {
			return TxErr
		}

		// 返回 nil 提交事务
		return nil
	})
}
