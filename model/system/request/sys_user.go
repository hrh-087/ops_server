package request

import "ops-server/model/system"

type Register struct {
	Username     string `json:"username"`
	Password     string `json:"password"`
	NickName     string `json:"nickName"`
	HeaderImg    string `json:"headerImg"`
	Phone        string `json:"phone"`
	Email        string `json:"email"`
	Enable       bool   `json:"enable"`
	AuthorityId  uint   `json:"authorityId" swaggertype:"string" example:"int 角色id"`
	AuthorityIds []uint `json:"authorityIds" swaggertype:"string" example:"[]uint 角色id"`
}

type Login struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type ChangePasswordReq struct {
	ID          uint   `json:"-"`           // 从 JWT 中提取 user id，避免越权
	Password    string `json:"password"`    // 密码
	NewPassword string `json:"newPassword"` // 新密码
}

// Modify  user's auth structure
type SetUserAuth struct {
	AuthorityId uint `json:"authorityId"` // 角色ID
}

type SetUserProject struct {
	ProjectId uint `json:"projectId"`
}

type ChangeUserInfo struct {
	ID           uint   `gorm:"primarykey"`                                                                           // 主键ID
	NickName     string `json:"nickName" gorm:"default:系统用户;comment:用户昵称"`                                            // 用户昵称
	Phone        string `json:"phone"  gorm:"comment:用户手机号"`                                                          // 用户手机号
	AuthorityIds []uint `json:"authorityIds" gorm:"-"`                                                                // 角色ID
	Email        string `json:"email"  gorm:"comment:用户邮箱"`                                                           // 用户邮箱
	HeaderImg    string `json:"headerImg" gorm:"default:https://qmplusimg.henrongyi.top/gva_header.jpg;comment:用户头像"` // 用户头像
	//SideMode     string                `json:"sideMode"  gorm:"comment:用户侧边主题"`                                                      // 用户侧边主题
	Enable      bool                  `json:"enable" gorm:"comment:冻结用户"` //冻结用户
	Authorities []system.SysAuthority `json:"-" gorm:"many2many:sys_user_authority;"`
}

// Modify  user's auth structure
type SetUserAuthorities struct {
	ID           uint
	AuthorityIds []uint `json:"authorityIds"` // 角色ID
}
