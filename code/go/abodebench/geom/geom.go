package main

import (
	"fmt"
	"math"
)

// Point represents a geographic coordinate
type Point struct {
	Lat float64
	Lng float64
}

// degreesToRadians converts degrees to radians
func degreesToRadians(degrees float64) float64 {
	return degrees * math.Pi / 180.0
}

// calculateHaversineDistance calculates the distance between two points on Earth
// using the Haversine formula, returning distance in meters
func calculateHaversineDistance(p1, p2 Point) float64 {
	// Earth's radius in meters
	earthRadius := 6371000.0

	lat1 := degreesToRadians(p1.Lat)
	lng1 := degreesToRadians(p1.Lng)
	lat2 := degreesToRadians(p2.Lat)
	lng2 := degreesToRadians(p2.Lng)

	dlat := lat2 - lat1
	dlng := lng2 - lng1

	a := math.Sin(dlat/2)*math.Sin(dlat/2) +
		math.Cos(lat1)*math.Cos(lat2)*
			math.Sin(dlng/2)*math.Sin(dlng/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return earthRadius * c
}

// CalculatePolygonArea calculates the area of a polygon in square miles
// using the Shoelace formula (Gauss's area formula) with geographic coordinates
func CalculatePolygonArea(points []Point) float64 {
	if len(points) < 3 {
		return 0 // A polygon must have at least 3 points
	}

	// Make sure the polygon is closed (first point = last point)
	pointsCopy := make([]Point, len(points))
	copy(pointsCopy, points)
	if pointsCopy[0].Lat != pointsCopy[len(pointsCopy)-1].Lat ||
		pointsCopy[0].Lng != pointsCopy[len(pointsCopy)-1].Lng {
		pointsCopy = append(pointsCopy, pointsCopy[0])
	}

	// For more accurate calculations on the surface of the Earth,
	// we'll project the points onto a plane using a local projection
	// This is simplified and works best for small to medium polygons

	// Find the centroid of the polygon for reference
	var centroidLat, centroidLng float64
	for _, p := range pointsCopy {
		centroidLat += p.Lat
		centroidLng += p.Lng
	}
	centroidLat /= float64(len(pointsCopy))
	centroidLng /= float64(len(pointsCopy))

	centroid := Point{Lat: centroidLat, Lng: centroidLng}

	// Convert to a local planar coordinate system (in meters)
	localCoords := make([][2]float64, len(pointsCopy))
	for i, p := range pointsCopy {
		// Calculate north-south distance (y)
		yDistance := calculateHaversineDistance(
			Point{Lat: centroid.Lat, Lng: centroid.Lng},
			Point{Lat: p.Lat, Lng: centroid.Lng},
		)
		if p.Lat < centroid.Lat {
			yDistance = -yDistance
		}

		// Calculate east-west distance (x)
		xDistance := calculateHaversineDistance(
			Point{Lat: centroid.Lat, Lng: centroid.Lng},
			Point{Lat: centroid.Lat, Lng: p.Lng},
		)
		if p.Lng < centroid.Lng {
			xDistance = -xDistance
		}

		localCoords[i] = [2]float64{xDistance, yDistance}
	}

	// Calculate area using the Shoelace formula
	area := 0.0
	for i := 0; i < len(localCoords)-1; i++ {
		area += localCoords[i][0] * localCoords[i+1][1]
		area -= localCoords[i+1][0] * localCoords[i][1]
	}
	area = math.Abs(area) / 2.0

	// Convert from square meters to square miles
	// 1 square meter = 3.861e-7 square miles
	squareMiles := area * 3.861e-7

	return squareMiles
}

func main() {
	// Example: Calculate area of a polygon
	polygon := []Point{
		{Lat: 34.052234, Lng: -118.243683},
		{Lat: 34.152234, Lng: -118.343683},
		{Lat: 34.102234, Lng: -118.193683},
		{Lat: 34.052234, Lng: -118.243683}, // Closing point
	}

	area := CalculatePolygonArea(polygon)
	fmt.Printf("Polygon area: %.2f square miles\n", area)

	// Example 2: Calculate area of a different polygon
	polygon2 := []Point{
		{Lat: 37.7749, Lng: -122.4194}, // San Francisco
		{Lat: 37.3382, Lng: -121.8863}, // San Jose
		{Lat: 37.8715, Lng: -122.2730}, // Berkeley
	}

	area2 := CalculatePolygonArea(polygon2)
	fmt.Printf("Polygon 2 area: %.2f square miles\n", area2)

	// Example 3: Calculate area of a different polygon
	polygon3 := []Point{
		{Lat: 35.1, Lng: -118.243683},
		{Lat: 35.2, Lng: -118.343683},
		{Lat: 35.3, Lng: -118.193683},
		{Lat: 35.1, Lng: -118.243683},
	}

	area3 := CalculatePolygonArea(polygon3)
	fmt.Printf("Polygon 3 area: %.2f square miles\n", area3)

	// Example 4: Calculate area of polygon from coordinates
	polygon4 := []Point{
		{Lat: 34.052234, Lng: -118.243683},
		{Lat: 34.152234, Lng: -118.343683},
		{Lat: 34.102234, Lng: -118.193683},
		{Lat: 34.052234, Lng: -118.243683},
	}

	area4 := CalculatePolygonArea(polygon4)
	fmt.Printf("Polygon 4 area: %.2f square miles\n", area4)

	// Example 5: Calculate area of polygon from provided coordinates
	polygon5 := []Point{
		{Lat: 34.052235, Lng: -118.243683},
		{Lat: 34.152235, Lng: -118.343683},
		{Lat: 34.102235, Lng: -118.193683},
		{Lat: 34.052235, Lng: -118.243683},
	}

	area5 := CalculatePolygonArea(polygon5)
	fmt.Printf("Polygon 5 area: %.2f square miles\n", area5)

	// Example 6: Calculate area of polygon from WKT format coordinates
	polygon6 := []Point{
		{Lat: 34.05237622604136, Lng: -118.25947584667969},
		{Lat: 34.059239834420666, Lng: -118.27540106335451},
		{Lat: 34.070937647321635, Lng: -118.266778703125},
		{Lat: 34.072227033405305, Lng: -118.2622886137085},
		{Lat: 34.071181715995074, Lng: -118.25943944656687},
		{Lat: 34.06981936757878, Lng: -118.25668272534179},
		{Lat: 34.06511561924042, Lng: -118.25172781506348},
		{Lat: 34.06349246254429, Lng: -118.2501301244812},
		{Lat: 34.06186927475736, Lng: -118.24956240216065},
		{Lat: 34.05237622604136, Lng: -118.25947584667969}, // Closing point
	}

	area6 := CalculatePolygonArea(polygon6)
	fmt.Printf("Polygon 6 area: %.2f square miles\n", area6)

	// Example 7: Calculate area of rectangular polygon from WKT format
	// Note: In WKT format, coordinates are (longitude latitude) pairs
	polygon7 := []Point{
		{Lat: 34.05237622604136, Lng: -118.27540106335451},
		{Lat: 34.05237622604136, Lng: -118.24956240216065},
		{Lat: 34.072227033405305, Lng: -118.24956240216065},
		{Lat: 34.072227033405305, Lng: -118.27540106335451},
		{Lat: 34.05237622604136, Lng: -118.27540106335451}, // Closing point
	}

	area7 := CalculatePolygonArea(polygon7)
	fmt.Printf("Polygon 7 area: %.2f square miles\n", area7)

	// Example 8: Calculate area of polygon from LINESTRING WKT format
	// Note: In WKT format, coordinates are (longitude latitude) pairs
	polygon8 := []Point{
		{Lat: 34.04564731453473, Lng: -118.29468886903943},
		{Lat: 34.096552943203235, Lng: -118.29434554628553},
		{Lat: 34.09669509485089, Lng: -118.2270542865199},
		{Lat: 34.04493612550349, Lng: -118.22688262514295},
		{Lat: 34.04564731453473, Lng: -118.29468886903943}, // Closing point
	}

	area8 := CalculatePolygonArea(polygon8)
	fmt.Printf("Polygon 8 area: %.2f square miles\n", area8)
}
