package request

import (
	"github.com/gofrs/uuid/v5"
	"github.com/golang-jwt/jwt/v4"
)

// 定义基础jwt认证结构体
type BaseClaims struct {
	UUID        uuid.UUID
	ID          uint
	Username    string
	NickName    string
	AuthorityId uint
}

// 自定义jwt结构体 继承BaseClaims
type CustomClaims struct {
	BaseClaims
	BufferTime int64
	jwt.RegisteredClaims
}
