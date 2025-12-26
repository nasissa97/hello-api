package translation_test

import (
	"errors"
	"testing"

	"hello-api/translation"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

func TestRemoteServiceTestSuite(t *testing.T) {
	suite.Run(t, new(RemoteServiceTestSuite))
}

type RemoteServiceTestSuite struct {
	suite.Suite
	client    *MockHelloClient
	underTest *translation.RemoteService
}

func (suite *RemoteServiceTestSuite) SetupTest() {
	suite.client = new(MockHelloClient)
	suite.underTest = translation.NewRemoteService(suite.client)
}

type MockHelloClient struct {
	mock.Mock
}

func (m *MockHelloClient) Translate(language, word string) (string, error) {
	args := m.Called(language, word)
	return args.String(0), args.Error(1)
}

func (suite *RemoteServiceTestSuite) TestTranslate() {
	suite.client.On("Translate", "bar", "foo").Return("baz", nil)

	res := suite.underTest.Translate("bar", "foo")

	suite.Equal(res, "baz")
	suite.client.AssertExpectations(suite.T())
}

func (suite *RemoteServiceTestSuite) TestTranslate_CaseSensitive() {
	suite.client.On("Translate", "bar", "foo").Return("baz", nil)

	res := suite.underTest.Translate("bar", "Foo")

	suite.Equal(res, "baz")
	suite.client.AssertExpectations(suite.T())
}

func (suite *RemoteServiceTestSuite) TestTranslate_Error() {
	suite.client.On("Translate", "bar", "foo").Return("baz", errors.New("failure"))

	res := suite.underTest.Translate("bar", "foo")

	suite.Equal(res, "")
	suite.client.AssertExpectations(suite.T())
}

func (suite *RemoteServiceTestSuite) TestTranslate_Cache() {
	suite.client.On("Translate", "bar", "foo").Return("baz", nil).Times(1)

	res1 := suite.underTest.Translate("bar", "foo")
	res2 := suite.underTest.Translate("bar", "Foo")

	suite.Equal(res1, "baz")
	suite.Equal(res2, "baz")
	suite.client.AssertExpectations(suite.T())
}
