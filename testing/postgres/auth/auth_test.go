package postgres_test

import (
	"fmt"
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
	authUseCase := usecases.NewAuthUseCases(repo)
	status, response := authUseCase.SignUp("nvtphong20042001@gmail.com", "wtf123456")
	fmt.Println(response)
	require.Equal(s.T(), http.StatusOK, status)
}

// func (s *authSuiteTest) TestVerifyCode() {
// 	repo := respository.NewAuthRepositopry(&s.TxStore)
// 	authUseCase := usecases.NewAuthUseCases(repo)
// 	status, _ := authUseCase.SignUp("nvtphong20042001@gmail.com", "wtf123456")
// 	require.Equal(s.T(), http.StatusOK, status)
// 	statusCode, _ := authUseCase.Login("nvtphong20042001@gmail.com", "wtf123456")
// 	require.Equal(s.T(), http.StatusOK, statusCode)

// 	authUseCase.V
// }

func (s *authSuiteTest) TestLogin() {
	repo := respository.NewAuthRepositopry(&s.TxStore)
	authUseCase := usecases.NewAuthUseCases(repo)
	status, _ := authUseCase.SignUp("daylarapviet", "wtf123456")
	require.Equal(s.T(), http.StatusOK, status)
	statusCode, _ := authUseCase.Login("daylarapviet", "wtf123456")
	require.Equal(s.T(), http.StatusOK, statusCode)
}
