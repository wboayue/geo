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

func TestRegion(t *testing.T) {
	regionA := Region{
		Vertices: []LatLng{
			{Lng: -10.764369, Lat: 6.284756},
			{Lng: -10.764112, Lat: 6.282111},
			{Lng: -10.761880, Lat: 6.282282},
			{Lng: -10.751838, Lat: 6.295250},
			{Lng: -10.754671, Lat: 6.296871},
			{Lng: -10.764369, Lat: 6.284756},
		},
	}

	regionB := Region{
		Vertices: []LatLng{
			{Lng: -10.764842, Lat: 6.285353},
			{Lng: -10.767159, Lat: 6.284671},
			{Lng: -10.766687, Lat: 6.281557},
			{Lng: -10.764927, Lat: 6.279680},
			{Lng: -10.761451, Lat: 6.280746},
			{Lng: -10.760035, Lat: 6.283263},
			{Lng: -10.764842, Lat: 6.285353},
		},
	}

	coordInA := LatLng{Lng: -10.755701, Lat: 6.292605}
	coordInB := LatLng{Lng: -10.765228, Lat: 6.281173}
	coordInAandB := LatLng{Lng: -10.762996, Lat: 6.283562}
	coordOutsideAandB := LatLng{Lng: -10.763683, Lat: 6.288510}

	t.Run("contains", func(t *testing.T) {
		assert.True(t, regionA.ContainsCoord(coordInA))
		assert.True(t, regionA.ContainsCoord(coordInAandB))
		assert.False(t, regionA.ContainsCoord(coordInB))
		assert.False(t, regionA.ContainsCoord(coordOutsideAandB))
	})

	t.Run("union", func(t *testing.T) {
		union := regionA.Union(&regionB)

		assert.True(t, union.ContainsCoord(coordInA))
		assert.True(t, union.ContainsCoord(coordInAandB))
		assert.True(t, union.ContainsCoord(coordInB))
		assert.False(t, union.ContainsCoord(coordOutsideAandB))
	})

	t.Run("intersection", func(t *testing.T) {
		intersection := regionA.Intersection(&regionB)

		assert.False(t, intersection.ContainsCoord(coordInA))
		assert.True(t, intersection.ContainsCoord(coordInAandB))
		assert.False(t, intersection.ContainsCoord(coordInB))
		assert.False(t, intersection.ContainsCoord(coordOutsideAandB))
	})

	t.Run("convexhull", func(t *testing.T) {
		union := regionA.Union(&regionB)
		convexhull := union.ConvexHull()

		assert.True(t, convexhull.ContainsCoord(coordInA))
		assert.True(t, convexhull.ContainsCoord(coordInAandB))
		assert.True(t, convexhull.ContainsCoord(coordInB))
		assert.False(t, convexhull.ContainsCoord(coordOutsideAandB))
	})
}
