package postgres_test

import (
	"log"

	"github.com/nvtphong200401/store-management/pkg/db"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type PostgresSuite struct {
	suite.Suite
	TxStore db.TxStore
}

func (s *PostgresSuite) SetupSuite() {
	log.Println("Run set up suite")
	gormDB, err := db.ConnectPostgresDB("D:\\Documents\\store-management-backend\\develop.env")
	rd := db.ConnectRedis("D:\\Documents\\store-management-backend\\develop.env")
	require.NoError(s.T(), err)
	s.TxStore = db.NewTXStore(gormDB, rd)

	s.TxStore.MigrateUp()
}

// TearDownSuite teardown at the end of test
func (s *PostgresSuite) TearDownSuite() error {
	log.Println("Run tear down suite")
	s.TxStore.MigrateDown()
	// return s.TxStore.CloseStorage()
	return nil
}
