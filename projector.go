package main

// Go wrapper to `proj` library for WGS84 projections

import (
	"fmt"

	"github.com/pebbe/proj/v5"
)

type (
	utmProjector struct {
		zone        int
		projContext *proj.Context
		projPJ      *proj.PJ
	}
)

// NewUTMProjectorForCoords creates a utmProjector for zone at specified lon and lat.
//
// Reference:
//   UTM Grid Zones of the World compiled by Alan Morton
//   http://www.dmap.co.uk/utmworld.htm
func NewUTMProjectorForCoords(lon, lat float64) (*utmProjector, error) {
	xz, _, err := proj.UTMzone(lon, lat)
	if err != nil {
		return nil, fmt.Errorf("could not determine zone: %v", err)
	}

	return NewUTMProjectorForZone(xz)
}

// NewUTMProjectorFromZone creates and utmProjector for specified longitudinal UTM zone.
//
// Reference:
//   UTM Grid Zones of the World compiled by Alan Morton
//   http://www.dmap.co.uk/utmworld.htm
func NewUTMProjectorForZone(zone int) (*utmProjector, error) {
	ctx := proj.NewContext()

	pj, err := ctx.Create(`
		+proj=pipeline
		+step +proj=unitconvert +xy_in=deg +xy_out=rad
		+step +proj=utm +datum=WGS84 +zone=` + fmt.Sprintf("%d", zone))
	if err != nil {
		return nil, err
	}

	return &utmProjector{
		zone:        zone,
		projContext: ctx,
		projPJ:      pj,
	}, nil
}

// Close releases resources held by projector
func (p *utmProjector) Close() {
	p.projPJ.Close()
	p.projContext.Close()
}

// ToUTMCoord projects WGS84 coordinates to UTM coordinates
func (p *utmProjector) ToUTMCoord(x, y float64) (float64, float64, error) {
	x, y, _, _, err := p.projPJ.Trans(proj.Fwd, x, y, 0, 0)
	if err != nil {
		return 0, 0, err
	}
	return x, y, nil
}

// ToUTMCoords projects WGS84 coordinates to UTM coordinates
func (p *utmProjector) ToUTMCoords(coords [][]float64) ([][]float64, error) {
	if len(coords) == 0 {
		return [][]float64{}, nil
	}

	results := make([][]float64, len(coords))
	for i, coord := range coords {
		x, y, _, _, err := p.projPJ.Trans(proj.Fwd, coord[0], coord[1], 0, 0)
		if err != nil {
			return results, err
		}
		results[i] = []float64{x, y}
	}

	return results, nil
}
