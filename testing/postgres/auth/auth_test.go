package postgres_test

import (
	"testing"

	"github.com/nvtphong200401/store-management/pkg/handlers/respository"
	"github.com/nvtphong200401/store-management/pkg/helpers"
	postgres_test "github.com/nvtphong200401/store-management/testing/postgres"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type authSuiteTest struct {
	postgres_test.PostgresSuite
}

func TestAuthSuite(t *testing.T) {

	authSuite := &authSuiteTest{
		postgres_test.PostgresSuite{},
	}
	suite.Run(t, authSuite)
}

func (s *authSuiteTest) TearDownTest() {

	err := s.TearDownSuite()
	require.NoError(s.T(), err)

}

func (s *authSuiteTest) TestSignUp() {
	repo := respository.NewAuthRepositopry(&s.TxStore)
	err := repo.SignUp("daylarapviet", "wtf123456")
	require.NoError(s.T(), err)
}

func (s *authSuiteTest) TestLogin() {
	repo := respository.NewAuthRepositopry(&s.TxStore)
	err := repo.SignUp("daylarapviet", "wtf123456")
	require.NoError(s.T(), err)
	tokenMap, err := repo.Login("daylarapviet", "wtf123456")
	require.NoError(s.T(), err)
	require.NotEmpty(s.T(), tokenMap[helpers.ACCESS_TOKEN])
	require.NotEmpty(s.T(), tokenMap[helpers.REFRESH_TOKEN])
}
