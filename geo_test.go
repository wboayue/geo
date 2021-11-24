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
		assert.Equal(t, "POINT (-10.773746 6.287188)", coordA.WKT())
	})

	t.Run("geojson", func(t *testing.T) {
		expected := `{ "type": "Point", "coordinates": [-10.774412, 6.285524] }`
		assert.Equal(t, expected, coordB.GeoJSON())
	})
}

func TestCircle(t *testing.T) {
	center := LatLng{Lng: -10.773746, Lat: 6.287188}

	circle := Circle{
		Center: center,
		Radius: 200.0,
	}

	exteriorCoord := LatLng{Lng: -10.776922, Lat: 6.291283}

	t.Run("buffer", func(t *testing.T) {
		buffer := 200.0 // meters

		circleB := circle.Buffer(buffer)

		assert.Equal(t, circle.Center, circleB.Center)
		assert.Equal(t, circle.Radius+buffer, circleB.Radius)
	})

	t.Run("asRegion", func(t *testing.T) {
		region := circle.AsRegion()

		expectedWKT := `POLYGON ((-10.771939 6.287194, -10.771972 6.286841, -10.772074 6.286502, -10.772240 6.286188, -10.772464 6.285914, -10.772737 6.285688, -10.773049 6.285520, -10.773387 6.285416, -10.773740 6.285380, -10.774093 6.285413, -10.774432 6.285515, -10.774745 6.285681, -10.775020 6.285905, -10.775245 6.286178, -10.775413 6.286490, -10.775517 6.286829, -10.775553 6.287182, -10.775520 6.287535, -10.775418 6.287874, -10.775252 6.288188, -10.775028 6.288462, -10.774755 6.288688, -10.774443 6.288856, -10.774105 6.288960, -10.773752 6.288996, -10.773399 6.288963, -10.773060 6.288861, -10.772747 6.288695, -10.772472 6.288471, -10.772247 6.288198, -10.772079 6.287886, -10.771975 6.287547, -10.771939 6.287194))`
		expectedGeoJSON := `{"type":"Polygon","coordinates":[[[-10.771939,6.287194],[-10.771972,6.286841],[-10.772074,6.286502],[-10.772240,6.286188],[-10.772464,6.285914],[-10.772737,6.285688],[-10.773049,6.285520],[-10.773387,6.285416],[-10.773740,6.285380],[-10.774093,6.285413],[-10.774432,6.285515],[-10.774745,6.285681],[-10.775020,6.285905],[-10.775245,6.286178],[-10.775413,6.286490],[-10.775517,6.286829],[-10.775553,6.287182],[-10.775520,6.287535],[-10.775418,6.287874],[-10.775252,6.288188],[-10.775028,6.288462],[-10.774755,6.288688],[-10.774443,6.288856],[-10.774105,6.288960],[-10.773752,6.288996],[-10.773399,6.288963],[-10.773060,6.288861],[-10.772747,6.288695],[-10.772472,6.288471],[-10.772247,6.288198],[-10.772079,6.287886],[-10.771975,6.287547],[-10.771939,6.287194]]]}`

		assert.Equal(t, 33, len(region.Vertices))
		assert.Equal(t, expectedWKT, region.WKT())
		assert.Equal(t, expectedGeoJSON, region.GeoJSON())
	})

	t.Run("contains", func(t *testing.T) {
		assert.True(t, circle.ContainsCoord(center))
		assert.False(t, circle.ContainsCoord(exteriorCoord))
	})
}
