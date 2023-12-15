package mongo

import (
	clog "github.com/cherry-game/cherry/logger"
	"go.mongodb.org/mongo-driver/mongo"
)

type Invoker interface {
	GetDb(id string) (*mongo.Database, bool)
	GetHashDb(groupId string, hashFn HashDb) (*mongo.Database, bool)
	GetDbMap(groupId string) (map[string]*mongo.Database, bool)
}

func (s *Component) GetDb(id string) (*mongo.Database, bool) {
	for _, group := range s.dbMap {
		for k, v := range group {
			if k == id {
				return v, true
			}
		}
	}
	return nil, false
}

func (s *Component) GetHashDb(groupId string, hashFn HashDb) (*mongo.Database, bool) {
	dbGroup, found := s.GetDbMap(groupId)
	if !found {
		clog.Warnf("groupId = %s not found.", groupId)
		return nil, false
	}

	dbId := hashFn(dbGroup)
	db, found := dbGroup[dbId]
	return db, found
}

func (s *Component) GetDbMap(groupId string) (map[string]*mongo.Database, bool) {
	dbGroup, found := s.dbMap[groupId]
	return dbGroup, found
}
