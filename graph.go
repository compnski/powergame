package main

import (
	"log"
	"math/rand"

	"github.com/Infinities-Within/delaunay"
)

type Point = delaunay.Point2d

type Center struct {
	index int

	point     Point // location
	water     bool
	ocean     bool
	coast     bool
	border    bool
	biome     string
	elevation float64 // 0.0-1.0
	moisture  float64 // 0.0-1.0

	neighbors []*Center
	borders   []*Edge
	corners   []*Corner
}

type Edge struct {
	index    int
	d0, d1   *Center // Delaunay edge
	v0, v1   *Corner // Voronoi edge
	midpoint Point   // halfway between v0,v1
	river    int     // volume of water, or 0
}

type Corner struct {
	index int

	point     Point   // location
	ocean     bool    // ocean
	water     bool    // lake or ocean
	coast     bool    // touches ocean and land polygons
	border    bool    // at the edge of the map
	elevation float64 // 0.0-1.0
	moisture  float64 // 0.0-1.0

	touches   []*Center
	protrudes []*Edge
	adjacent  []*Corner

	river          int     // 0 if no river, or volume of water in river
	downslope      *Corner // pointer to adjacent corner most downhill
	watershed      *Corner // pointer to coastal corner, or null
	watershed_size int
}

const MinPointDistance = 20

func genPoints(n, height, width int) []Point {
	r := rand.New(rand.NewSource(int64(n + height + width)))
	var points []Point
	for len(points) < n {
		var skip = false
		point := delaunay.NewPoint(float64(r.Intn(width-20)+10), float64(r.Intn(height-20)+10))
		for _, p := range points {
			if p.SquaredDistance(point) < MinPointDistance {
				skip = true
				break
			}
		}
		if !skip {
			points = append(points, point)
		}
	}
	return points
}

func nextHalfEdge(e int) int {
	if e%3 == 2 {
		return e - 2
	}
	return e + 1
}

func buildGraph(points []Point) {
	triangles, err := delaunay.Triangulate(points)
	if err != nil {
		log.Fatal("Failed to triangulate points", points, err)
	}
	log.Print(triangles)
}
