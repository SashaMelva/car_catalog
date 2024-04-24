package client

import (
	"encoding/json"
	"io"
	"net/http"

	model "github.com/SashaMelva/car_catalog/internal/storage/models"
	"go.uber.org/zap"
)

func GetInfoCarByRegNum(regNum, clientHost string, log *zap.SugaredLogger) (*model.RequestBody, error) {
	client := http.Client{}
	req, err := http.NewRequest("get", clientHost+"/info?regNum="+regNum, nil)

	if err != nil {
		return nil, err
	}

	response, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	buf, err := io.ReadAll(response.Body)

	if err != nil {
		return nil, err
	}

	log.Debug(string(buf))

	requestBody := model.RequestBody{}
	err = json.Unmarshal(buf, &requestBody)

	if err != nil {
		return nil, err
	}

	return &requestBody, err

}
