package advertisementSearch

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"api-tests-template/internal/assertions"
	advList "api-tests-template/internal/managers/advertisements"
	"api-tests-template/internal/managers/advertisements/models"
	"api-tests-template/internal/managers/auth"
	authModels "api-tests-template/internal/managers/auth/models"
	"api-tests-template/internal/utils"
	base "api-tests-template/tests"
)

const (
	existingAdTitle = "Муся в добрые руки"
	searchBeginning = "Муся"
	searchMiddle    = "добрые"
	searchEnd       = "руки"
)

type TestSuite struct {
	suite.Suite
	loginData authModels.LoginOkResponse
}

func TestSuiteRun(t *testing.T) {
	suite.Run(t, &TestSuite{})
}

func (s *TestSuite) SetupSuite() {
	base.SetupSuite()

	base.Precondition("Проверяем, что название объявления и подстроки заданы")
	require.NotEmpty(s.T(), existingAdTitle, "Задайте existingAdTitle")
	require.NotEmpty(s.T(), searchBeginning, "Задайте searchBeginning")
	require.NotEmpty(s.T(), searchMiddle, "Задайте searchMiddle")
	require.NotEmpty(s.T(), searchEnd, "Задайте searchEnd")

	base.Precondition("Авторизация пользователя с данными из переменных окружения")
	s.loginData = auth.Login(s.T(), os.Getenv("TEST_LOGIN"), os.Getenv("TEST_PASSWORD"))
}

func (s *TestSuite) TearDownSuite() {
	base.TearDownSuite()
}

func (s *TestSuite) TestSearchByTitle() {
	var list models.AdvertisementsList

	s.Run("Ищем объявление по подстроке в названии", func() {
		list = advList.GetAdvertisementsList(s.T(), s.loginData.Token, searchMiddle)
	})

	s.Run("В ответе есть хотя бы одно объявление", func() {
		require.NotEmpty(s.T(), list.Items)
	})

	s.Run("Среди найденных есть наше объявление", func() {
		assertions.ContainsAdvertisementByTitle(s.T(), list.Items, existingAdTitle)
	})
}

func (s *TestSuite) TestSearchByTitle_Substring() {
	cases := []struct {
		name   string
		search string
	}{
		{"начало", searchBeginning},
		{"середина", searchMiddle},
		{"конец", searchEnd},
		{"нижний регистр", strings.ToLower(searchMiddle)},
		{"верхний регистр", strings.ToUpper(searchMiddle)},
		{"смешанный регистр", utils.MixCase(searchMiddle)},
	}

	for _, tc := range cases {
		s.Run(tc.name, func() {
			list := advList.GetAdvertisementsList(s.T(), s.loginData.Token, tc.search)
			require.NotEmpty(s.T(), list.Items)
			assertions.ContainsAdvertisementByTitle(s.T(), list.Items, existingAdTitle)
		})
	}
}
