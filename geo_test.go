package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLatLng(t *testing.T) {
	coordA := LatLng{Lng: -10.773746, Lat: 6.287188}
	coordB := LatLng{Lng: -10.774412, Lat: 6.285524}

	t.Run("buffer", func(t *testing.T) {
		radius := 200.0 // meters

		circle := coordA.Buffer(radius)

		assert.Equal(t, coordA, circle.Center)
		assert.Equal(t, radius, circle.Radius)
	})

	t.Run("distance", func(t *testing.T) {
		distanceM := coordA.Distance(&coordB)
		assert.InDelta(t, 198.0, distanceM, 1.0)
	})

	t.Run("wkt", func(t *testing.T) {
		assert.Equal(t, "POINT (-10.773746, 6.287188)", coordA.WKT())
	})

	t.Run("geojson", func(t *testing.T) {
		expected := `{ "type": "Point", "coordinates": [-10.774412, 6.285524] }`
		assert.Equal(t, expected, coordB.GeoJSON())
		t.Error(coordB.GeoJSON())
	})
}
