package geo

import (
	"fmt"
	"math"
	"strings"

	"github.com/polastre/gogeos/geos"
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
		Vertices Coordinates
	}

	// Region
	LineString struct {
		Vertices Coordinates
	}

	Coordinates []LatLng
	Points      [][]float64
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
	defer projector.Close()

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
	return fmt.Sprintf("POINT (%.6f %.6f)", p.Lng, p.Lat)
}

// GeoJSON generates Geo JSON representation
func (p *LatLng) GeoJSON() string {
	return fmt.Sprintf(`{ "type": "Point", "coordinates": [%.6f, %.6f] }`, p.Lng, p.Lat)
}

// Buffer buffers circle by specified buffer.
func (p *Circle) Buffer(buffer float64) *Circle {
	return &Circle{
		Center: p.Center,
		Radius: p.Radius + buffer,
	}
}

// AsRegion converts the Circle to a Region
func (c *Circle) AsRegion() *Region {
	projector, err := NewUTMProjectorForCoords(c.Center.Lng, c.Center.Lat)
	if err != nil {
		panic(err)
	}
	defer projector.Close()

	x, y, err := projector.ToUTMCoord(c.Center.Lng, c.Center.Lat)
	if err != nil {
		panic(err)
	}

	point, err := geos.FromWKT(fmt.Sprintf("POINT (%.6f %.6f)", x, y))
	if err != nil {
		panic(err)
	}

	buf, err := point.Buffer(c.Radius)
	if err != nil {
		panic(err)
	}

	b, err := buf.Shell()
	if err != nil {
		panic(err)
	}

	coords, err := b.Coords()
	if err != nil {
		panic(err)
	}

	points, err := projector.FromUTMGeosCoords(coords)
	if err != nil {
		panic(err)
	}

	return &Region{
		Vertices: points,
	}
}

// ContainsCoord determines if circle contains coordinate
func (c *Circle) ContainsCoord(coord LatLng) bool {
	return c.AsRegion().ContainsCoord(coord)
}

func (r *Region) Union(other *Region) *Region {
	start := r.Vertices[0]
	projector, err := NewUTMProjectorForCoords(start.Lng, start.Lat)
	if err != nil {
		panic(err)
	}
	defer projector.Close()

	pointsA, err := projector.ToUTMCoordsA(r.Vertices)
	if err != nil {
		panic(err)
	}

	pointsB, err := projector.ToUTMCoordsA(other.Vertices)
	if err != nil {
		panic(err)
	}

	geoA, err := geos.FromWKT(pointsA.WKT())
	if err != nil {
		panic(err)
	}

	geoB, err := geos.FromWKT(pointsB.WKT())
	if err != nil {
		panic(err)
	}

	union, err := geoA.Union(geoB)
	if err != nil {
		panic(err)
	}

	return polygonToRegion(projector, union)
}

func (p Points) WKT() string {
	vertices := make([]string, len(p))
	for i, vertex := range p {
		vertices[i] = fmt.Sprintf("%.6f %.6f", vertex[0], vertex[1])
	}
	exterior := strings.Join(vertices, ", ")
	return fmt.Sprintf("POLYGON ((%s))", exterior)
}

func (r *Region) Intersection(other *Region) *Region {
	start := r.Vertices[0]
	projector, err := NewUTMProjectorForCoords(start.Lng, start.Lat)
	if err != nil {
		panic(err)
	}
	defer projector.Close()

	pointsA, err := projector.ToUTMCoordsA(r.Vertices)
	if err != nil {
		panic(err)
	}

	pointsB, err := projector.ToUTMCoordsA(other.Vertices)
	if err != nil {
		panic(err)
	}

	geoA, err := geos.FromWKT(pointsA.WKT())
	if err != nil {
		panic(err)
	}

	geoB, err := geos.FromWKT(pointsB.WKT())
	if err != nil {
		panic(err)
	}

	intersection, err := geoA.Intersection(geoB)
	if err != nil {
		panic(err)
	}

	return polygonToRegion(projector, intersection)
}

func (r *Region) ConvexHull() *Region {
	start := r.Vertices[0]
	projector, err := NewUTMProjectorForCoords(start.Lng, start.Lat)
	if err != nil {
		panic(err)
	}
	defer projector.Close()

	pointsA, err := projector.ToUTMCoordsA(r.Vertices)
	if err != nil {
		panic(err)
	}

	geoA, err := geos.FromWKT(pointsA.WKT())
	if err != nil {
		panic(err)
	}

	convexHull, err := geoA.ConvexHull()
	if err != nil {
		panic(err)
	}

	return polygonToRegion(projector, convexHull)
}

func polygonToRegion(projector *utmProjector, geom *geos.Geometry) *Region {
	b, err := geom.Shell()
	if err != nil {
		panic(err)
	}

	coords, err := b.Coords()
	if err != nil {
		panic(err)
	}

	points, err := projector.FromUTMGeosCoords(coords)
	if err != nil {
		panic(err)
	}

	return &Region{
		Vertices: points,
	}
}

// WKT generates well known text representation
func (r *Region) WKT() string {
	vertices := make([]string, len(r.Vertices))
	for i, vertex := range r.Vertices {
		vertices[i] = fmt.Sprintf("%.6f %.6f", vertex.Lng, vertex.Lat)
	}
	exterior := strings.Join(vertices, ", ")
	return fmt.Sprintf("POLYGON ((%s))", exterior)
}

// GeoJSON generates Geo JSON representation
func (r *Region) GeoJSON() string {
	vertices := make([]string, len(r.Vertices))
	for i, vertex := range r.Vertices {
		vertices[i] = fmt.Sprintf("[%.6f,%.6f]", vertex.Lng, vertex.Lat)
	}
	exterior := strings.Join(vertices, ",")
	return fmt.Sprintf(`{"type":"Polygon","coordinates":[[%s]]}`, exterior)
}

// ContainsCoord determines if region contains coordinate
func (r *Region) ContainsCoord(coord LatLng) bool {
	start := r.Vertices[0]
	projector, err := NewUTMProjectorForCoords(start.Lng, start.Lat)
	if err != nil {
		panic(err)
	}
	defer projector.Close()

	pointsA, err := projector.ToUTMCoordsA(r.Vertices)
	if err != nil {
		panic(err)
	}

	geoA, err := geos.FromWKT(pointsA.WKT())
	if err != nil {
		panic(err)
	}

	x, y, err := projector.ToUTMCoord(coord.Lng, coord.Lat)
	if err != nil {
		panic(err)
	}

	point, err := geos.FromWKT(fmt.Sprintf("POINT (%.6f %.6f)", x, y))
	if err != nil {
		panic(err)
	}

	in, err := geoA.Contains(point)
	if err != nil {
		panic(err)
	}

	return in
}
