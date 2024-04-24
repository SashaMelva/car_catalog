package hendler

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/SashaMelva/car_catalog/internal/app"
	"github.com/SashaMelva/car_catalog/internal/config"
	"github.com/SashaMelva/car_catalog/internal/logger"
	"github.com/SashaMelva/car_catalog/internal/storage/connection"
	"github.com/SashaMelva/car_catalog/internal/storage/memory"
	"go.uber.org/zap/zapcore"
)

func TestCarCatalogHendler(t *testing.T) {
	serever := testService()
	testCase := []struct {
		name       string
		method     string
		path       string
		body       []byte
		want       string
		statusCode int
	}{
		// {
		// 	name:   "Update car",
		// 	method: http.MethodPut,
		// 	path:   "/car-catalog",
		// 	body: []byte(`{
		// 		"cars":[
		// 		{
		// 		"regNum":"A777AA129",
		// 		"mark":"Lada",
		// 		"model":"Vesta",
		// 		"year":2002,
		// 		"owner":{
		// 			"name":"q",
		// 			"patronymic":"q"
		// 		}
		// 		}
		// 		]
		// 		}
		// 		`),
		// 	want:       ``,
		// 	statusCode: http.StatusOK,
		// },
		// {
		// 	name:   "Update cars",
		// 	method: http.MethodPut,
		// 	path:   "/car-catalog",
		// 	body: []byte(`{
		// 		"cars":[
		// 		{
		// 		"regNum":"A777AA123",
		// 		"mark":"Lada",
		// 		"model":"Vesta",
		// 		"year":2002,
		// 		"owner":{
		// 			"name":"q",
		// 			"surname":"q",
		// 			"patronymic":"q"}
		// 		},
		// 		{
		// 		"regNum":"A777AA124",
		// 		"mark":"BMW",
		// 		"model":" M1 HOMMAGE",
		// 		"year":2002,
		// 		"owner":{
		// 			"name":"q",
		// 			"surname":"q",
		// 			"patronymic":"q"}
		// 		},
		// 		{
		// 		"regNum":"X123XX150",
		// 		"mark":"Lada",
		// 		"model":"Vesta",
		// 		"year":2002,
		// 		"owner":{
		// 			"name":"q",
		// 			"surname":"q",
		// 			"patronymic":"q"}
		// 		}
		// 		]
		// 		}`),
		// 	want:       ``,
		// 	statusCode: http.StatusOK,
		// },

		// {
		// 	name:       "Delete car regNum empty",
		// 	method:     http.MethodDelete,
		// 	path:       "/car-catalog",
		// 	body:       []byte(``),
		// 	want:       "Для удаления машины необходим регистрационный номер",
		// 	statusCode: http.StatusBadRequest,
		// },
		// {
		// 	name:       "Delete cars",
		// 	method:     http.MethodDelete,
		// 	path:       "/car-catalog?regNums=A777AA200",
		// 	body:       []byte(``),
		// 	want:       ``,
		// 	statusCode: http.StatusOK,
		// },
		// {
		// 	name:       "Delete cars",
		// 	method:     http.MethodDelete,
		// 	path:       "/car-catalog?regNums=A777AA129,A777AA201,A777AA209",
		// 	body:       []byte(``),
		// 	want:       ``,
		// 	statusCode: http.StatusBadRequest,
		// },
		// {
		// 	name:       "Get all cars",
		// 	method:     http.MethodGet,
		// 	path:       "/car-catalog",
		// 	body:       []byte(``),
		// 	want:       ``,
		// 	statusCode: http.StatusOK,
		// },
		// {
		// 	name:       "Get cars by regNums",
		// 	method:     http.MethodGet,
		// 	path:       "/car-catalog?regNum=A777AA111",
		// 	body:       []byte(``),
		// 	want:       ``,
		// 	statusCode: http.StatusOK,
		// },
		// {
		// 	name:       "Get all cars by limit and offset",
		// 	method:     http.MethodGet,
		// 	path:       "/car-catalog?limit=2&offset=2",
		// 	body:       []byte(``),
		// 	want:       ``,
		// 	statusCode: http.StatusOK,
		// },
		// {
		// 	name:       "Get cars by model and mark",
		// 	method:     http.MethodGet,
		// 	path:       "/car-catalog?mark=BMW",
		// 	body:       []byte(``),
		// 	want:       ``,
		// 	statusCode: http.StatusOK,
		// },
		// {
		// 	name:       "Get cars by model and mark",
		// 	method:     http.MethodGet,
		// 	path:       "/car-catalog?periodYear=2003:2024",
		// 	body:       []byte(``),
		// 	want:       ``,
		// 	statusCode: http.StatusOK,
		// },
		// {
		// 	name:       "Get cars by year lower",
		// 	method:     http.MethodGet,
		// 	path:       "/car-catalog?periodYear=:2024",
		// 	body:       []byte(``),
		// 	want:       ``,
		// 	statusCode: http.StatusOK,
		// },
		// {
		// 	name:       "Get cars by year upper",
		// 	method:     http.MethodGet,
		// 	path:       "/car-catalog?periodYear=2000:&limit=1",
		// 	body:       []byte(``),
		// 	want:       ``,
		// 	statusCode: http.StatusOK,
		// },
		{
			name:   "add cars",
			method: http.MethodPost,
			path:   "/car-catalog?periodYear=2003:2024",
			body: []byte(`{
				"regNums":[
				"X123XX150",
				"X123XX151"
				]
				}`),
			want:       ``,
			statusCode: http.StatusOK,
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			bodyReader := bytes.NewReader(tc.body)
			request := httptest.NewRequest(tc.method, tc.path, bodyReader)
			responseRecorder := httptest.NewRecorder()

			serever.CarCatalogHendler(responseRecorder, request)

			if responseRecorder.Code != tc.statusCode {
				t.Errorf("Want status '%d', got '%d'", tc.statusCode, responseRecorder.Code)
			}

			if strings.TrimSpace(responseRecorder.Body.String()) != tc.want {
				t.Errorf("Want '%s', got '%s'", tc.want, responseRecorder.Body)
			}
		})
	}
}

func testService() *Service {

	log := logger.New(&config.ConfigLogger{
		Level:       zapcore.InfoLevel,
		LogEncoding: "console",
	}, "../../../logFiles")

	connection := connection.New(&config.ConfigDB{
		NameDB:   "car_catalog_db",
		Host:     "127.0.0.1",
		Port:     "5432",
		User:     "postgres",
		Password: "123456",
	}, log)

	memstorage := memory.New(connection.StorageDb, log)
	calendar := app.New(log, memstorage)

	return &Service{
		Logger: *log,
		app:    *calendar,
	}
}
