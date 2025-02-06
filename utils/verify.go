package utils

var (
	LoginVerify     = Rules{"Username": {NotEmpty()}, "Password": {NotEmpty()}}
	RegisterVerify  = Rules{"Username": {NotEmpty()}, "NickName": {NotEmpty()}, "Password": {NotEmpty()}, "AuthorityId": {NotEmpty()}}
	AuthorityVerify = Rules{"AuthorityId": {NotEmpty()}, "AuthorityName": {NotEmpty()}}
	ApiVerify       = Rules{"Path": {NotEmpty()}, "Description": {NotEmpty()}, "ApiGroup": {NotEmpty()}, "Method": {NotEmpty()}}

	PageInfoVerify     = Rules{"Page": {NotEmpty()}, "PageSize": {NotEmpty()}}
	SearchServerVerify = Rules{"PlatformId": {NotEmpty()}}

	IdVerify           = Rules{"ID": []string{NotEmpty()}}
	MenuVerify         = Rules{"Path": {NotEmpty()}, "Name": {NotEmpty()}, "Component": {NotEmpty()}, "Sort": {Ge("0")}}
	MenuMetaVerify     = Rules{"Title": {NotEmpty()}}
	AuthorityIdVerify  = Rules{"AuthorityId": {NotEmpty()}}
	OldAuthorityVerify = Rules{"OldAuthorityId": {NotEmpty()}}

	SetUserAuthorityVerify = Rules{"AuthorityId": {NotEmpty()}}
	SetUserProjectVerify   = Rules{"ProjectId": {NotEmpty()}}
	ChangePasswordVerify   = Rules{"Password": {NotEmpty()}, "NewPassword": {NotEmpty()}}

	ProjectVerify  = Rules{"ProjectName": {NotEmpty()}}
	SshAuthVerify  = Rules{"User": {NotEmpty()}}
	CloudVerify    = Rules{"CloudName": {NotEmpty()}, "RegionId": {NotEmpty()}, "RegionName": {NotEmpty()}, "SecretId": {NotEmpty()}, "SecretKey": {NotEmpty()}}
	PlatformVerify = Rules{"PlatformCode": {NotEmpty()}, "PlatformName": {NotEmpty()}}

	AssetsServerVerify = Rules{"PrivateIp": {NotEmpty()}, "PubIp": {NotEmpty()}, "SSHPort": {NotEmpty()}, "ServerName": {NotEmpty()}, "PlatformId": {NotEmpty()}}
	AssetsMysqlVerify  = Rules{"Host": {NotEmpty()}, "Name": {NotEmpty()}, "Pass": {NotEmpty()}, "Port": {NotEmpty()}, "PlatformId": {NotEmpty()}}
	AssetsRedisVerify  = Rules{"Host": {NotEmpty()}, "Name": {NotEmpty()}, "PlatformId": {NotEmpty()}, "Port": {NotEmpty()}}
	AssetsMongoVerify  = Rules{"Host": {NotEmpty()}, "Name": {NotEmpty()}, "PlatformId": {NotEmpty()}, "Auth": {NotEmpty()}}
	AssetsKafkaVerify  = Rules{"Host": {NotEmpty()}, "Name": {NotEmpty()}, "PlatformId": {NotEmpty()}}
	AssetsLbVerify     = Rules{"PlatformId": {NotEmpty()}, "CloudProduceId": {NotEmpty()}}

	GameTypeVerify   = Rules{"Code": {NotEmpty()}, "Name": {NotEmpty()}, "VmidRule": {NotEmpty()}}
	GameServerVerify = Rules{
		"PlatformId": {NotEmpty()},
		"Name":       {NotEmpty()},
		"GameTypeId": {NotEmpty()},
		"RedisId":    {NotEmpty()},
		"KafkaId":    {NotEmpty()},
		"MongoId":    {NotEmpty()},
	}

	JobVerify          = Rules{"JobId": {NotEmpty()}}
	CommandVerify      = Rules{"Command": {NotEmpty()}, "Name": {NotEmpty()}}
	BatchCommandVerify = Rules{"BatchType": {NotEmpty()}, "CommandId": {NotEmpty()}, "ServerList": []string{NotEmpty()}}

	GameUpdateVerify = Rules{"Name": {NotEmpty()}, "UpdateType": {NotEmpty()}}
)
