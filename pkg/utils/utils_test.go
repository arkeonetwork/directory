package utils

import "testing"

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
