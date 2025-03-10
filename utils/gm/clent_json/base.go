package clent_json

import (
	"context"
	"errors"
	"ops-server/global"
	"ops-server/model/system"
)

func GetConfigDir(ctx context.Context) (configDir string, err error) {
	projectId := ctx.Value("projectId").(string)

	if projectId == "" {
		return configDir, errors.New("projectId is empty")
	}

	var project system.SysProject

	err = global.OPS_DB.First(&project, "id = ?", projectId).Error
	if err != nil {
		return
	}

	if project.ClientConfigDir == "" {
		return configDir, errors.New("clientConfigDir is empty")
	}

	return project.ClientConfigDir, nil
}
