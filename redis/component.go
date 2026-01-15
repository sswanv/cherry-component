package redis

import (
	"context"
	"time"

	cfacade "github.com/cherry-game/cherry/facade"
	clog "github.com/cherry-game/cherry/logger"
	cprofile "github.com/cherry-game/cherry/profile"
	"github.com/redis/go-redis/v9"
)

const (
	Name = "redis_component"
)

func NewComponent() *Component {
	return &Component{
		rdsMap: make(map[string]*redis.Client),
	}
}

type (
	config struct {
		Addr       string
		Password   string
		DB         int
		ClientName string
	}
	Component struct {
		cfacade.Component

		rdsMap map[string]*redis.Client
	}
)

func (c *Component) Name() string {
	return Name
}

func (c *Component) Init() {
	redisConfig := c.App().Settings().GetConfig("redis")
	if redisConfig.LastError() != nil || redisConfig.Size() == 0 {
		clog.Fatalf("[nodeId = %s] `redis` property not exists.", c.App().NodeID())
		return
	}

	redisListConfig := cprofile.GetConfig("redis")
	if redisListConfig.LastError() != nil {
		clog.Panic("`redis` property not exists in profile file.")
	}

	for i := 0; i < redisConfig.Size(); i++ {
		id1 := redisConfig.GetString(i)
		for _, id2 := range redisListConfig.Keys() {
			if id2 != id1 {
				continue
			}
		}

		item := redisListConfig.GetConfig(id1)
		var conf config
		err := item.Unmarshal(&conf)
		if err != nil {
			clog.Fatalf("parse redis config err: %v", err)
		}

		rds, err := c.createRedisClient(conf)
		if err != nil {
			clog.Fatalf("[id = %s] create redis client fail. error = %s", id1, err)
		}
		c.rdsMap[id1] = rds
	}
}

func (c *Component) createRedisClient(conf config) (*redis.Client, error) {
	rds := redis.NewClient(&redis.Options{
		Addr:       conf.Addr,
		Password:   conf.Password,
		DB:         conf.DB,
		ClientName: conf.ClientName,
	})
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := rds.Ping(ctx).Err(); err != nil {
		return nil, err
	}
	return rds, nil
}
