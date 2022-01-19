![test](https://github.com/wboayue/geo/workflows/ci/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/wboayue/geo)](https://goreportcard.com/report/github.com/wboayue/geo)
<a href="https://codecov.io/gh/wboayue/geo" target="_blank">
    <img src="https://img.shields.io/codecov/c/github/wboayue/geo?color=%2334D058" alt="Coverage">
</a>

# Overview

Go library for the geometric manipulations in the WGS84 coordinate system. 

It uses the Geos library for geometric calculations and Proj to project between coordinate systems.

# Quick Start

```go
package main

import (
	"github.com/wboayue/geo"
)

func main() {
    // Distance between coordinates

    a := geo.LatLng{Lng: -10.773746, Lat: 6.287188}
    b := geo.LatLng{Lng: -10.774412, Lat: 6.285524}

    distance_M := a.Distance(b)
    // 198.0

    // Buffering
	circle := a.Buffer(200.0)   // buffer by 200 m -> geo.Circle
	circle.AsRegion()           // converts to geo.Region 
}
```

# Dependencies

This library requires the Geos and Proj C libraries. They can be installed as follows:

## OSX
```bash
brew install proj geos pkg-config
```

## Ubuntu
```bash
apt-get install -y libproj-dev libgeos-dev
```
