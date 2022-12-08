package utils

import (
	"testing"

	"github.com/ArkeoNetwork/directory/pkg/types"
)

func TestParseCoordinates(t *testing.T) {
	epsilon := .0001
	coordinateString := "67.3523,-47.6878"
	coordinates, err := ParseCoordinates(coordinateString)
	if err != nil {
		t.FailNow()
	}
	if !IsNearEqual(coordinates.Latitude, 67.35234, epsilon) ||
		!IsNearEqual(coordinates.Longitude, -47.6878, epsilon) {
		t.FailNow()
	}

	coordinateString = "67.3523,-x"
	coordinates, err = ParseCoordinates(coordinateString)
	if err == nil {
		t.FailNow()
	}

	coordinateString = "yy,-47.6878"
	coordinates, err = ParseCoordinates(coordinateString)
	if err == nil {
		t.FailNow()
	}

	coordinateString = "67.3523,-47.6878,666"
	coordinates, err = ParseCoordinates(coordinateString)
	if err == nil {
		t.FailNow()
	}
	coordinateString = "67.3523"
	coordinates, err = ParseCoordinates(coordinateString)
	if err == nil {
		t.FailNow()
	}
}

func TestParseContractType(t *testing.T) {
	contract := "paygo"
	_, err := ParseContractType(contract)
	if err == nil {
		t.FailNow()
	}

	contract = "PayAsYouGo"
	contractType, err := ParseContractType(contract)
	if err != nil {
		t.FailNow()
	}

	if contractType != types.ContractTypePayAsYouGo {
		t.FailNow()
	}

	contract = "Subscription"
	contractType, err = ParseContractType(contract)
	if err != nil {
		t.FailNow()
	}
	if contractType != types.ContractTypeSubscription {
		t.FailNow()
	}
}

func TestDownloadProviderMetadata(t *testing.T) {

	data, err := DownloadProviderMetadata("https://petstore.swagger.io/v2/swagger.json", 5, 1e6)
	if data == nil {
		t.FailNow()
	}

	if err != nil {
		t.FailNow()
	}

	if _, exists := (*data)["host"]; !exists {
		t.FailNow()
	}

	_, err = DownloadProviderMetadata("https://petstore.swagger.io/v2/swagger.json", 5, 1)
	if err == nil {
		t.FailNow()
	}

}
