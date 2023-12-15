package jwt

import (
	cfacade "github.com/cherry-game/cherry/facade"
	clog "github.com/cherry-game/cherry/logger"
	"github.com/golang-jwt/jwt/v5"
)

const (
	Name = "jwt_component"
)

type (
	Claims struct {
		Pid        string `json:"pid"`
		OpenId     string `json:"openId"`
		DeviceType string `json:"deviceType"`
		BundleName string `json:"bundleName"`
		jwt.RegisteredClaims
	}

	config struct {
		SecretKey      []byte // key
		ExpireDuration int    // 过期时长(单位: 小时)
	}
)

func NewComponent() *Component {
	return &Component{}
}

type Component struct {
	cfacade.Component

	config config
}

func (c *Component) Name() string {
	return Name
}

func (c *Component) Init() {
	jwtConfig := c.App().Settings().GetConfig("jwt")
	if jwtConfig.LastError() != nil {
		clog.Panic("`jwt` property not exists in profile file")
	}
	c.config.SecretKey = []byte(jwtConfig.GetString("secretKey"))
	c.config.ExpireDuration = jwtConfig.GetInt("expireDuration")
}
