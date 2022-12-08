package utils

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/ArkeoNetwork/directory/pkg/types"
	resty "github.com/go-resty/resty/v2"
)

func ParseCoordinates(coordinates string) (types.Coordinates, error) {
	if coordinates == "" {
		return types.Coordinates{}, errors.New("empty string cannot be parsed into coordinates")
	}
	coordinatesSplit := strings.Split(coordinates, ",")
	if len(coordinatesSplit) != 2 {
		return types.Coordinates{}, errors.New("too many parameters passed to coordinates")
	}
	latitude, err := strconv.ParseFloat(coordinatesSplit[0], 32)
	if err != nil {
		return types.Coordinates{}, errors.New("latitude cannot be parsed")
	}
	longitude, err := strconv.ParseFloat(coordinatesSplit[1], 32)
	if err != nil {
		return types.Coordinates{}, errors.New("longitude cannot be parsed")
	}
	return types.Coordinates{Latitude: latitude, Longitude: longitude}, nil
}

func ParseContractType(contractType string) (types.ContractType, error) {
	if types.ContractType(contractType) == types.ContractTypePayAsYouGo {
		return types.ContractType(contractType), nil
	} else if types.ContractType(contractType) == types.ContractTypeSubscription {
		return types.ContractType(contractType), nil
	} else {
		return types.ContractTypePayAsYouGo, fmt.Errorf("unexpected contract type %s", contractType)
	}
}

func IsNearEqual(a float64, b float64, epsilon float64) bool {
	return math.Abs(a-b) <= epsilon
}

// see arkeo-protocol/common/chain.go
var validChains = map[string]struct{}{"arkeo-mainnet-fullnode": {}, "btc-mainnet-fullnode": {}, "eth-mainnet-fullnode": {}, "swapi.dev": {}}

func ValidateChain(chain string) (ok bool) {
	_, ok = validChains[chain]
	return
}

func DownloadProviderMetadata(url string, retries int, maxBytes int) (*map[string]any, error) {
	client := resty.New()
	var result map[string]any
	client.SetRetryCount(retries)
	client.SetTimeout(time.Second * 5)
	client.SetHeader("Accept", "application/json")
	resp, err := client.R().ForceContentType("application/json").SetResult(&result).Get(url)

	if err != nil {
		return nil, err
	}

	if len(resp.Body()) > maxBytes {
		return nil, errors.New("DownloadProviderMetadata: max bytes exceeded")
	}

	return &result, nil
}
