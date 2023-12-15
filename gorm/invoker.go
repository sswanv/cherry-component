package gorm

import (
	clog "github.com/cherry-game/cherry/logger"
	"gorm.io/gorm"
)

type Invoker interface {
	GetDb(id string) (*gorm.DB, bool)
	GetHashDb(groupId string, hashFn HashDb) (*gorm.DB, bool)
	GetDbMap(groupId string) (map[string]*gorm.DB, bool)
}

func (s *Component) GetDb(id string) (*gorm.DB, bool) {
	for _, group := range s.ormMap {
		for k, v := range group {
			if k == id {
				return v, true
			}
		}
	}
	return nil, false
}

func (s *Component) GetHashDb(groupId string, hashFn HashDb) (*gorm.DB, bool) {
	dbGroup, found := s.GetDbMap(groupId)
	if !found {
		clog.Warnf("groupId = %s not found.", groupId)
		return nil, false
	}

	dbId := hashFn(dbGroup)
	db, found := dbGroup[dbId]
	return db, found
}

func (s *Component) GetDbMap(groupId string) (map[string]*gorm.DB, bool) {
	dbGroup, found := s.ormMap[groupId]
	return dbGroup, found
}
