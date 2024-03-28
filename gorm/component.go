package gorm

import (
	"time"

	cfacade "github.com/cherry-game/cherry/facade"
	clog "github.com/cherry-game/cherry/logger"
	cprofile "github.com/cherry-game/cherry/profile"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const (
	Name = "gorm_component"
)

type (
	Component struct {
		cfacade.Component
		// key:groupId,value:{key:id,value:*gorm.Db}
		ormMap map[string]map[string]*gorm.DB
	}

	generalDB struct {
		DbName         string `json:"db_name"`          // 数据库名
		Host           string `json:"host"`             // 连接地址
		UserName       string `json:"user_name"`        // 数据库用户名
		Password       string `json:"password"`         // 数据库密码
		Config         string `json:"config"`           // 高级配置
		MaxIdleConnect int    `json:"max_idle_connect"` // 空闲中的最大连接数
		MaxOpenConnect int    `json:"max_open_connect"` // 打开到数据库的最大连接数
		LogMode        bool   `json:"log_mode"`
	}

	specializedDB struct {
		generalDB
		Enable  bool   `json:"enable"`
		GroupId string `json:"group_id"`
		DbId    string `json:"db_id"`
		Type    string `json:"type"`
	}

	// HashDb hash by group id
	HashDb func(dbMaps map[string]*gorm.DB) string
)

func NewComponent() *Component {
	return &Component{
		ormMap: make(map[string]map[string]*gorm.DB),
	}
}

func (s *Component) Name() string {
	return Name
}

func parseGormConfig(groupId string, item cfacade.ProfileJSON) *specializedDB {
	config := new(specializedDB)
	err := item.Unmarshal(config)
	if err != nil {
		clog.Fatalf("failed to parse configuration: %v", err)
	}
	config.GroupId = groupId
	return config
}

func (s *Component) Init() {
	// load only the database contained in the `db_id_list`
	dbIdList := s.App().Settings().Get("db_id_list")
	if dbIdList.LastError() != nil || dbIdList.Size() < 1 {
		clog.Warnf("[nodeId = %s] `db_id_list` property not exists.", s.App().NodeId())
		return
	}

	dbConfig := cprofile.GetConfig("db")
	if dbConfig.LastError() != nil {
		clog.Panic("`db` property not exists in profile file.")
	}

	for _, groupId := range dbConfig.Keys() {
		s.ormMap[groupId] = make(map[string]*gorm.DB)

		dbGroup := dbConfig.GetConfig(groupId)
		for i := 0; i < dbGroup.Size(); i++ {
			item := dbGroup.GetConfig(i)
			gormConfig := parseGormConfig(groupId, item)

			for j := 0; j < dbIdList.Size(); j++ {
				if dbIdList.Get(j).ToString() != gormConfig.DbId {
					continue
				}

				if !gormConfig.Enable {
					clog.Fatalf("[dbName = %s] is disabled!", gormConfig.DbName)
				}

				db, err := s.createORM(gormConfig)
				if err != nil {
					clog.Fatalf("[dbName = %s] create orm fail. error = %s", gormConfig.DbName, err)
				}

				s.ormMap[groupId][gormConfig.DbId] = db
				clog.Infof("[dbGroup =%s, dbName = %s] is connected.", gormConfig.GroupId, gormConfig.DbId)
			}
		}
	}
}

func (s *Component) createORM(cfg *specializedDB) (*gorm.DB, error) {
	switch cfg.Type {
	case "mysql":
		return NewMysql(mysqlConfig{generalDB: cfg.generalDB})
	case "clickhouse":
		return NewClickhouse(clickhouseConfig{generalDB: cfg.generalDB})
	default:
		return nil, errors.Errorf("not support gorm type: [%v]", cfg.Type)
	}
}

func getLogger() logger.Interface {
	return logger.New(
		gormLogger{log: clog.DefaultLogger},
		logger.Config{
			SlowThreshold: time.Second,
			LogLevel:      logger.Silent,
			Colorful:      true,
		},
	)
}
