package system

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"ops-server/global"
	"ops-server/model/common/request"
	"ops-server/model/system"
	"strconv"
	"strings"
	"text/template"
)

type GameTypeService struct {
}

var GameTypeApp = new(GameTypeService)

func (g *GameTypeService) CreateGameType(ctx context.Context, gameType system.SysGameType) (err error) {
	if !errors.Is(global.OPS_DB.WithContext(ctx).Where("code = ?", gameType.Code).First(&system.SysGameType{}).Error, gorm.ErrRecordNotFound) {
		return errors.New("游戏类型已存在")
	}
	err = global.OPS_DB.WithContext(ctx).Create(&gameType).Error
	return
}

func (g *GameTypeService) UpdateGameType(ctx context.Context, gameType system.SysGameType) (err error) {
	var old system.SysGameType

	updateField := []string{
		"code",
		"name",
		"order",
		"compose_template",
		"config_template",
		"grpc_port",
		"http_port",
		"tcp_port",
		"vmid_rule",
		"is_fight",
	}

	if errors.Is(global.OPS_DB.WithContext(ctx).Where("id = ?", gameType.ID).First(&old).Error, gorm.ErrRecordNotFound) {
		return errors.New("记录不存在")
	}
	err = global.OPS_DB.WithContext(ctx).Model(&old).Select(updateField).Updates(gameType).Error
	return
}

func (g *GameTypeService) DeleteGameType(ctx context.Context, id int) (err error) {
	return global.OPS_DB.WithContext(ctx).Where("id = ?", id).Delete(&system.SysGameType{}).Error
}

func (g *GameTypeService) GetGameTypeById(ctx context.Context, id int) (result system.SysGameType, err error) {
	err = global.OPS_DB.WithContext(ctx).Preload("SysProject").First(&result, "id = ?", id).Error
	return
}

func (g *GameTypeService) GetGameTypeList(ctx context.Context, info request.PageInfo, server system.SysGameType) (list interface{}, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	db := global.OPS_DB.WithContext(ctx).Model(&system.SysGameType{})

	var resultList []system.SysGameType

	// 在count的时候已经执行了插件逻辑, 添加一个上下文使后续的db操作跳过插件逻辑
	err = db.Count(&total).Set("skip_project_filter", true).Error
	if err != nil {
		return resultList, total, err
	}

	if server.Name != "" {
		db = db.Where("name like ?", "%"+server.Name+"%")
	}

	db = db.Limit(limit).Offset(offset)
	OrderStr := "id desc"

	err = db.Preload("SysProject").Order(OrderStr).Find(&resultList).Error
	return resultList, total, err
}

func (g *GameTypeService) GetGameTypeAll(ctx context.Context) (result []system.SysGameType, err error) {
	err = global.OPS_DB.WithContext(ctx).Find(&result).Error
	return
}

func (g GameTypeService) GenerateConfigFile(game system.SysGameServer) (content string, err error) {
	var dbNamePrefix string
	var buf bytes.Buffer
	var fightType string

	if strings.TrimSpace(game.GameType.ConfigTemplate) == "" {
		return "", errors.New("配置模板为空")
	}

	if game.GameType.IsFight {
		dbNamePrefix = "fight"
	} else {
		dbNamePrefix = game.GameType.Code
	}

	if game.GameType.Code == "fight" {
		fightType = "default"
	} else {
		fightType = game.GameType.Code
	}

	templateData := request.GameConfigFile{
		PlatformCode:   game.Platform.PlatformCode,
		Vmid:           game.Vmid,
		Name:           game.Name,
		PubIp:          game.Host.PubIp,
		InnerIp:        game.Host.PrivateIp,
		HttpPort:       game.HttpPort,
		GrpcPort:       game.GrpcPort,
		TcpPort:        game.TcpPort,
		MongoUri:       game.Mongo.Host,
		MongoAuth:      game.Mongo.Auth,
		DbName:         fmt.Sprintf("%s_%s", dbNamePrefix, game.Platform.PlatformCode),
		KafkaUri:       game.Kafka.Host,
		RedisUri:       fmt.Sprintf("%s:%d", game.Redis.Host, game.Redis.Port),
		RedisPass:      game.Redis.Password,
		RedisMeshUri:   fmt.Sprintf("%s:%d", game.Redis.Host, game.Redis.Port),
		RedisMeshPass:  game.Redis.Password,
		GatewayUri:     game.Platform.GatewayUrl,
		FightType:      fightType,
		LtsGroupId:     game.Platform.LtsLogGroupId,
		LtsStreamId:    game.Platform.LtsLogStreamId,
		SecretKey:      game.Platform.CloudSecretKey,
		AccessKey:      game.Platform.CloudSecretId,
		CloudProjectId: game.Platform.CloudProjectId,
		CloudRegionId:  game.Platform.CloudRegionId,
	}

	tmpl, err := template.New("config").Parse(game.GameType.ConfigTemplate)
	if err != nil {
		return "", err
	}

	err = tmpl.Execute(&buf, templateData)
	if err != nil {
		global.OPS_LOG.Error("生成配置文件失败", zap.Error(err))
		return "", errors.New("生成配置文件失败")
	}

	return buf.String(), nil
}

func (g GameTypeService) GenerateComposeFile(game system.SysGameServer) (content string, err error) {
	var buf bytes.Buffer
	var imageName string

	if strings.TrimSpace(game.GameType.ComposeTemplate) == "" {
		return "", errors.New("compose模板为空")
	}

	if game.GameType.IsFight {
		imageName = "fight"
	} else {
		imageName = game.GameType.Code
	}

	templateData := request.DockerComposeFile{
		ImageService:     game.GameType.Code,
		ImageTag:         game.Platform.ImageTag,
		JsonConfigVolume: global.OPS_CONFIG.Game.GameConfigDir,
		ImageUri:         game.Platform.ImageUri,
		ImageName:        imageName,
	}

	tmpl, err := template.New("config").Parse(game.GameType.ComposeTemplate)
	if err != nil {
		return "", err
	}

	err = tmpl.Execute(&buf, templateData)
	if err != nil {
		global.OPS_LOG.Error("生成docker-compose文件失败", zap.Error(err))
		return "", errors.New("生成docker-compose文件失败")
	}

	return buf.String(), nil
}

func (g GameTypeService) CopyGameType(ctx context.Context, projectId uint, gameTypeIds []int) (err error) {
	var gameTypeList []system.SysGameType

	headerProjectId := ctx.Value("projectId").(string)
	if headerProjectId == strconv.Itoa(int(projectId)) {
		return errors.New("不可以复制到同一项目下")
	}

	err = global.OPS_DB.WithContext(ctx).Where("id in ?", gameTypeIds).Find(&gameTypeList).Error
	if err != nil {
		return
	}

	return global.OPS_DB.Transaction(func(tx *gorm.DB) error {
		for _, gameType := range gameTypeList {
			gameType.ID = 0
			gameType.ProjectId = projectId
			err = tx.Create(&gameType).Error
			if err != nil {
				return err
			}
		}
		return nil
	})
}
