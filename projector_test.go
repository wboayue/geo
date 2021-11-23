package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const float64EqualityThreshold = 1e-2

func TestNewUTMProjectorForCoords(t *testing.T) {
	lat := 7.234933
	lon := -8.977841

	projector, err := NewUTMProjectorForCoords(lon, lat)
	if err != nil {
		t.Fatalf("error creating projector: %v", err)
	}
	defer projector.Close()

	assert.NotNil(t, projector.projContext, "expected proj context")
	assert.NotNil(t, projector.projPJ, "expected proj PJ")
	assert.Equal(t, 29, projector.zone)
}

func TestUTMProjector_ToUTMCoord(t *testing.T) {
	lon := float64(-122)
	lat := float64(37)

	projector, err := NewUTMProjectorForCoords(lon, lat)
	if err != nil {
		t.Fatalf("error creating projector: %v", err)
	}
	defer projector.Close()

	u, v, err := projector.ToUTMCoord(lon, lat)
	if err != nil {
		t.Fatalf("error projecting coordinates: %v", err)
	}

	assert.Equal(t, 10, projector.zone)
	assert.InDelta(t, 588977.32, u, float64EqualityThreshold)
	assert.InDelta(t, 4095339.69, v, float64EqualityThreshold)
}

func BenchmarkUTMProjector_ToUTMCoord(b *testing.B) {
	lon := float64(-122)
	lat := float64(37)

	projector, err := NewUTMProjectorForCoords(lon, lat)
	if err != nil {
		b.Fatalf("error creating projector: %v", err)
	}
	defer projector.Close()

	for i := 0; i < b.N; i++ {
		projector.ToUTMCoord(lon, lat)
	}
}

func TestUTMProjector_ToUTMCoords(t *testing.T) {
	wgs84Points := [][]float64{
		{-43.157150, -22.948968},
		{-43.156936, -22.950410},
		{-43.155841, -22.950331},
		{-43.155219, -22.949383},
		{-43.156270, -22.948652},
		{-43.157150, -22.948968},
	}

	firstPoint := wgs84Points[0]
	projector, err := NewUTMProjectorForCoords(firstPoint[0], firstPoint[1])
	if err != nil {
		t.Fatalf("error creating projector: %v", err)
	}
	defer projector.Close()

	utmPoints, err := projector.ToUTMCoords(wgs84Points)
	if err != nil {
		t.Fatalf("error projecting coordinates: %v", err)
	}

	assert.Equal(t, 23, projector.zone)
	assert.Equal(t, 6, len(utmPoints))

	firstUTMPoint := utmPoints[0]
	assert.InDelta(t, 688951.83, firstUTMPoint[0], float64EqualityThreshold)
	assert.InDelta(t, -2539055.65, firstUTMPoint[1], float64EqualityThreshold)
}

func BenchmarkUTMProjector_ToUTMCoords(b *testing.B) {
	wgs84Points := [][]float64{
		{-43.157150, -22.948968},
		{-43.156936, -22.950410},
		{-43.155841, -22.950331},
		{-43.155219, -22.949383},
		{-43.156270, -22.948652},
		{-43.157150, -22.948968},
	}

	firstPoint := wgs84Points[0]
	projector, err := NewUTMProjectorForCoords(firstPoint[0], firstPoint[1])
	if err != nil {
		b.Fatalf("error creating projector: %v", err)
	}
	defer projector.Close()

	for i := 0; i < b.N; i++ {
		projector.ToUTMCoords(wgs84Points)
	}
}
