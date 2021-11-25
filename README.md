![test](https://github.com/wboayue/geo/workflows/ci/badge.svg)

# Overview

Go library for the geometric manipulations in the WGS84 coordinate system. 

It uses the Geos library for geometric calculations and Proj to project between coordinate systems.

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
