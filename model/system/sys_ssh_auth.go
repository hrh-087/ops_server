package system

import "ops-server/global"

type SysSshAuth struct {
	global.OpsModel
	ProjectId            uint   `json:"projectId" gorm:"uniqueIndex;comment:项目id"`
	User                 string `json:"user" gorm:"comment:用户名"`
	Password             string `json:"password" gorm:"comment:密码"`
	PublicKey            string `json:"publicKey" gorm:"type:text;comment:公钥"`
	PrivateKey           string `json:"privateKey" gorm:"type:text;comment:私钥"`
	PrivateKeyPassphrase string `json:"privateKeyPassphrase" gorm:"comment:私钥密码"`
	UsePass              bool   `json:"usePass" gorm:"comment:是否使用密码登录"`
}

func (s *SysSshAuth) TableName() string {
	return "sys_ssh_auth"
}
