package translation_test

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"hello-api/translation"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

func TestHelloClientSuite(t *testing.T) {
	suite.Run(t, new(HelloClientSuite))
}

type HelloClientSuite struct {
	suite.Suite
	mockServerService *MockService
	server            *httptest.Server
	underTest         translation.HelloClient
}

type MockService struct {
	mock.Mock
}

func (m *MockService) Translate(language, word string) (string, error) {
	args := m.Called(language, word)
	return args.String(0), args.Error(1)
}

func (suite *HelloClientSuite) SetupSuite() {
	suite.mockServerService = new(MockService)
	handler := func(w http.ResponseWriter, r *http.Request) {
		b, _ := io.ReadAll(r.Body)
		defer func(r *http.Request) {
			_ = r.Body.Close()
		}(r)

		var m map[string]any
		_ = json.Unmarshal(b, &m)

		word := m["word"].(string)
		language := m["language"].(string)

		resp, err := suite.mockServerService.Translate(language, word)
		if err != nil {
			http.Error(w, "error", http.StatusInternalServerError)
			return
		}
		if resp == "" {
			http.Error(w, "missing", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = io.WriteString(w, resp)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", handler)
	suite.server = httptest.NewServer(mux)
	suite.underTest = translation.NewHelloClient(suite.server.URL)
}

func (suite *HelloClientSuite) SetupTest() {
	suite.mockServerService = new(MockService)
}

func (suite *HelloClientSuite) TearDownSuite() {
	suite.server.Close()
}

func (suite *HelloClientSuite) TestCall() {
	suite.mockServerService.On("Translate", "bar", "foo").Return(`{"translation":"baz"}`, nil)

	resp, err := suite.underTest.Translate("bar", "foo")

	suite.NoError(err)
	suite.Equal(resp, "baz")
}

func (suite *HelloClientSuite) TestCall_APIError() {
	suite.mockServerService.On("Translate", "bar", "foo").Return("", errors.New("this is a test"))

	resp, err := suite.underTest.Translate("bar", "foo")

	suite.EqualError(err, "error in api")
	suite.Equal(resp, "")
}

func (suite *HelloClientSuite) TestCall_InvalidJSON() {
	suite.mockServerService.On("Translate", "bar", "foo").Return(`invalid json`, nil)

	resp, err := suite.underTest.Translate("bar", "foo")

	suite.EqualError(err, "unable to decode message")
	suite.Equal(resp, "")
}
