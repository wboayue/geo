package main

import (
	"fmt"
	"math"
)

type (
	// LatLng represents a point in WGS84 coordinates
	LatLng struct {
		Lat float64
		Lng float64
	}

	// Circle
	Circle struct {
		Center LatLng
		Radius float64 // meters
	}

	// Region
	Region struct {
		Vertices []LatLng
	}

	// Region
	LineString struct {
		Vertices []LatLng
	}
)

// Buffer buffers point by specified buffer creating a Circle
func (p *LatLng) Buffer(buffer float64) *Circle {
	return &Circle{
		Center: *p,
		Radius: buffer,
	}
}

// Distance calculates the distance between two LatLng points
func (p *LatLng) Distance(other *LatLng) float64 {
	projector, err := NewUTMProjectorForCoords(p.Lng, p.Lat)
	if err != nil {
		panic(err)
	}

	x, y, err := projector.ToUTMCoord(p.Lng, p.Lat)
	if err != nil {
		panic(err)
	}

	u, v, err := projector.ToUTMCoord(other.Lng, other.Lat)
	if err != nil {
		panic(err)
	}

	return math.Sqrt(math.Pow(x-u, 2) + math.Pow(y-v, 2))
}

// WKT generates well known text representation
func (p *LatLng) WKT() string {
	return fmt.Sprintf("POINT (%.6f, %.6f)", p.Lng, p.Lat)
}

// GeoJSON generates Geo JSON representation
func (p *LatLng) GeoJSON() string {
	return fmt.Sprintf(`{ "type": "Point", "coordinates": [%.6f, %.6f] }`, p.Lng, p.Lat)
}

func (p *Circle) Buffer(buffer float64) *Circle {
	return &Circle{
		Center: p.Center,
		Radius: p.Radius + buffer,
	}
}

func (p *Region) Union(other *Region) *Region {
	return nil
}

func (p *Region) Intersection(other *Region) *Region {
	return nil
}

func (p *Region) ConvexHull(other *Region) *Region {
	return nil
}
