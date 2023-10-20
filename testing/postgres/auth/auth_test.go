package postgres_test

import (
	"net/http"
	"testing"

	"github.com/nvtphong200401/store-management/pkg/handlers/respository"
	"github.com/nvtphong200401/store-management/pkg/handlers/usecases"
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
	err := repo.SignUp("daylarapviet", []byte("wtf123456"))
	require.NoError(s.T(), err)
}

func (s *authSuiteTest) TestLogin() {
	repo := respository.NewAuthRepositopry(&s.TxStore)
	authUseCase := usecases.NewAuthUseCases(repo)
	err := repo.SignUp("daylarapviet", []byte("wtf123456"))
	require.NoError(s.T(), err)
	statusCode, _ := authUseCase.Login("daylarapviet", "wtf123456")
	require.Equal(s.T(), http.StatusOK, statusCode)

}
