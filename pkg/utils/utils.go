package utils

import (
	"errors"
	"math"
	"strconv"
	"strings"

	"github.com/ArkeoNetwork/directory/pkg/types"
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

func IsNearEqual(a float64, b float64, epsilon float64) bool {
	return math.Abs(a-b) <= epsilon
}

// see arkeo-protocol/common/chain.go
var validChains = map[string]struct{}{"arkeo-mainnet-fullnode": {}, "btc-mainnet-fullnode": {}, "eth-mainnet-fullnode": {}, "swapi.dev": {}}

func ValidateChain(chain string) (ok bool) {
	_, ok = validChains[chain]
	return
}
