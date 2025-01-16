package system

import (
	"github.com/gofrs/uuid/v5"
	"ops-server/global"
)

type SysUser struct {
	global.OpsModel
	UUID        uuid.UUID      `json:"uuid" gorm:"index;comment:用户UUID"`
	Username    string         `json:"userName" gorm:"index;comment:用户名称"`
	Password    string         `json:"-" gorm:"comment:用户密码"`
	NickName    string         `json:"nickName" gorm:"comment:用户昵称"`
	HeaderImg   string         `json:"headerImg" gorm:"comment:用户头像"`
	Phone       string         `json:"phone" gorm:"index;comment:用户手机号"`
	Email       string         `json:"email" gorm:"index;comment:用户邮箱"`
	Enable      bool           `json:"enable" gorm:"default:true;comment:用户状态 1 正常 0 禁用"`
	AuthorityId uint           `json:"authorityId" gorm:"comment:用户角色ID"` // 用户角色ID
	ProjectId   uint           `json:"projectId" gorm:"comment:项目ID"`     // 用户当前所选择的项目ID
	Authority   SysAuthority   `json:"authority" gorm:"foreignKey:AuthorityId;references:AuthorityId;comment:用户角色"`
	Authorities []SysAuthority `json:"authorities" gorm:"many2many:sys_user_authority;"`
}

func (SysUser) TableName() string {
	return "sys_user"
}
