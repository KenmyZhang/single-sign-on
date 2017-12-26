package sqlStore

import (
	"context"
)

type LayeredStore struct {
	TmpContext    context.Context
	DatabaseLayer *SqlSupplier
}

func NewLayeredStore() Store {
	return &LayeredStore{
		TmpContext:    context.TODO(),
		DatabaseLayer: NewSqlSupplier(),
	}
}

func (s *LayeredStore) User() UserStore {
	return s.DatabaseLayer.User()
}

func (s *LayeredStore) Token() TokenStore {
	return s.DatabaseLayer.Token()
}

func (s *LayeredStore) System() SystemStore {
	return s.DatabaseLayer.System()
}

func (s *LayeredStore) MarkSystemRanUnitTests() {
	s.DatabaseLayer.MarkSystemRanUnitTests()
}

func (s *LayeredStore) Close() {
	s.DatabaseLayer.Close()
}

func (s *LayeredStore) DropAllTables() {
	s.DatabaseLayer.DropAllTables()
}

func (s *LayeredStore) TotalMasterDbConnections() int {
	return s.DatabaseLayer.TotalMasterDbConnections()
}

func (s *LayeredStore) TotalReadDbConnections() int {
	return s.DatabaseLayer.TotalReadDbConnections()
}

func (s *LayeredStore) TotalSearchDbConnections() int {
	return s.DatabaseLayer.TotalSearchDbConnections()
}
