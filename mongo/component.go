package mongo

import (
	"context"
	"fmt"
	"time"

	cfacade "github.com/cherry-game/cherry/facade"
	clog "github.com/cherry-game/cherry/logger"
	cprofile "github.com/cherry-game/cherry/profile"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

const (
	Name = "mongo_component"
)

type config struct {
	Enable  bool   `json:"enable"`
	DbId    string `json:"dbId"`
	DbName  string `json:"dbName"`
	Uri     string `json:"uri"`
	Timeout int    `json:"timeout"`
}

type (
	Component struct {
		cfacade.Component
		dbMap map[string]map[string]*mongo.Database
	}

	// HashDb hash by group id
	HashDb func(dbMaps map[string]*mongo.Database) string
)

func NewComponent() *Component {
	return &Component{
		dbMap: make(map[string]map[string]*mongo.Database),
	}
}

func (*Component) Name() string {
	return Name
}

func (s *Component) Init() {
	mongoListConfig := s.App().Settings().Get("mongo")
	if mongoListConfig.LastError() != nil || mongoListConfig.Size() < 1 {
		clog.Warnf("[nodeId = %s] `mongo_id_list` property not exists.", s.App().NodeID())
		return
	}

	mongoConfig := cprofile.GetConfig("mongo")
	if mongoConfig.LastError() != nil {
		panic("`mongo` property not exists in profile file.")
	}

	for _, groupId := range mongoConfig.Keys() {
		s.dbMap[groupId] = make(map[string]*mongo.Database)

		dbGroup := mongoConfig.GetConfig(groupId)
		for i := 0; i < dbGroup.Size(); i++ {
			item := dbGroup.GetConfig(i)
			var conf config
			err := item.Unmarshal(&conf)
			if err != nil {
				clog.Fatalf("parse mongo config err: %v", err)
			}

			for j := 0; j < mongoListConfig.Size(); j++ {
				dbId := mongoListConfig.Get(j).ToString()
				if conf.DbId != dbId {
					continue
				}

				if !conf.Enable {
					panic(fmt.Sprintf("[dbName = %s] is disabled!", conf.DbName))
				}

				db, err := CreateDatabase(conf)
				if err != nil {
					panic(fmt.Sprintf("[dbName = %s] create mongodb fail. error = %s", conf.DbName, err))
				}

				s.dbMap[groupId][conf.DbId] = db
				clog.Infof("[dbGroup =%s, dbName = %s] is connected.", groupId, conf.DbName)
			}
		}
	}
}

func CreateDatabase(conf config) (*mongo.Database, error) {
	tt := 3 * time.Second
	if conf.Timeout > 3 {
		tt = time.Duration(conf.Timeout) * time.Second
	}

	o := options.Client().ApplyURI(conf.Uri)
	if err := o.Validate(); err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), tt)
	defer cancel()

	client, err := mongo.Connect(ctx, o)
	if err != nil {
		return nil, err
	}

	err = client.Ping(context.Background(), readpref.Primary())
	if err != nil {
		return nil, err
	}

	return client.Database(conf.DbName), nil
}
