import { useQuery } from '@tanstack/react-query';
import { API_CONFIG } from '@/config/api';
import { SchoolFeature, SchoolsGeoJSON } from '@am/commons';
import schools from './mockdata/filtered_school_districts_with_schools.json';
import { addColorsToFeatures } from '../lib/colorUtils';

// Query keys
export const schoolKeys = {
  all: ['schools'] as const,
  list: () => [...schoolKeys.all, 'list'] as const,
};

// Type-cast the imported JSON to our defined interface
const typedSchools = schools as SchoolsGeoJSON;

// Map school ratings to colors
const getRatingColor = (rating: string | undefined): string => {
  if (!rating) return '#FFD700'; // Yellow for unknown/undefined ratings

  const ratingColors: Record<string, string> = {
    'A+': '#00FF00', // Bright Green
    A: '#32CD32', // Lime Green
    'A-': '#90EE90', // Light Green
    'B+': '#98FB98', // Pale Green
    B: '#FFD700', // Yellow
    'B-': '#FFA500', // Orange
    'C+': '#FF8C00', // Dark Orange
    C: '#FF6B6B', // Light Red
    'C-': '#FF4500', // Orange Red
    'D+': '#FF0000', // Red
    D: '#DC143C', // Crimson
    'D-': '#8B0000', // Dark Red
    F: '#800000', // Maroon
  };

  return ratingColors[rating] || '#FF8C00'; // Default to unknown rating
};

// API functions
export const fetchSchools = async (): Promise<SchoolFeature[]> => {
  if (API_CONFIG.useMockData) {
    // First add base colors to features
    const featuresWithBaseColors = addColorsToFeatures(typedSchools.features);

    // Then update colors based on school ratings
    const featuresWithRatingColors = featuresWithBaseColors.map((feature) => {
      const schools = feature.properties.schools || [];
      if (schools.length === 0) return feature;

      // Calculate average rating for the district
      const ratings = schools
        .map((school) => school.rating)
        .filter((rating): rating is string => !!rating);

      if (ratings.length === 0) return feature;

      // Use the most common rating in the district
      const ratingCounts = ratings.reduce(
        (acc, rating) => {
          acc[rating] = (acc[rating] || 0) + 1;
          return acc;
        },
        {} as Record<string, number>
      );

      const ratingEntries = Object.entries(ratingCounts);
      if (ratingEntries.length === 0) return feature;

      const mostCommonRating = (
        ratingEntries.sort(([, a], [, b]) => b - a)[0] as [string, number]
      )[0];

      return {
        ...feature,
        properties: {
          ...feature.properties,
          color: getRatingColor(mostCommonRating),
        },
      };
    });

    return featuresWithRatingColors;
  }

  // Real API call
  const response = await fetch(`${API_CONFIG.apiBaseUrl}/schools`);
  if (!response.ok) {
    throw new Error('Failed to fetch schools');
  }
  return response.json();
};

// React Query hook
export const useSchoolsQuery = () => {
  return useQuery({
    queryKey: schoolKeys.list(),
    queryFn: fetchSchools,
  });
};
