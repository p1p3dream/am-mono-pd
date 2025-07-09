package geom

import (
	"math"
)

// Point represents a geographic coordinate
type Point struct {
	Lat float64
	Lon float64
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
	Lon1 := degreesToRadians(p1.Lon)
	lat2 := degreesToRadians(p2.Lat)
	Lon2 := degreesToRadians(p2.Lon)

	dlat := lat2 - lat1
	dLon := Lon2 - Lon1

	a := math.Sin(dlat/2)*math.Sin(dlat/2) +
		math.Cos(lat1)*math.Cos(lat2)*
			math.Sin(dLon/2)*math.Sin(dLon/2)
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
		pointsCopy[0].Lon != pointsCopy[len(pointsCopy)-1].Lon {
		pointsCopy = append(pointsCopy, pointsCopy[0])
	}

	// Find the centroid of the polygon for reference
	var centroidLat, centroidLon float64
	for _, p := range pointsCopy {
		centroidLat += p.Lat
		centroidLon += p.Lon
	}
	centroidLat /= float64(len(pointsCopy))
	centroidLon /= float64(len(pointsCopy))

	centroid := Point{Lat: centroidLat, Lon: centroidLon}

	// Convert to a local planar coordinate system (in meters)
	localCoords := make([][2]float64, len(pointsCopy))
	for i, p := range pointsCopy {
		// Calculate north-south distance (y)
		yDistance := calculateHaversineDistance(
			Point{Lat: centroid.Lat, Lon: centroid.Lon},
			Point{Lat: p.Lat, Lon: centroid.Lon},
		)
		if p.Lat < centroid.Lat {
			yDistance = -yDistance
		}

		// Calculate east-west distance (x)
		xDistance := calculateHaversineDistance(
			Point{Lat: centroid.Lat, Lon: centroid.Lon},
			Point{Lat: centroid.Lat, Lon: p.Lon},
		)
		if p.Lon < centroid.Lon {
			xDistance = -xDistance
		}

		localCoords[i] = [2]float64{xDistance, yDistance}
	}

	// Calculate area using the Shoelace formula
	area := 0.0
	for i := range len(localCoords) - 1 {
		area += localCoords[i][0] * localCoords[i+1][1]
		area -= localCoords[i+1][0] * localCoords[i][1]
	}
	area = math.Abs(area) / 2.0

	// Convert from square meters to square miles
	// 1 square meter = 3.861e-7 square miles
	squareMiles := area * 3.861e-7

	return squareMiles
}
