import { useQuery } from '@tanstack/react-query';
import { API_CONFIG } from '../config/api';
import { addColorsToFeatures } from '../lib/colorUtils';
import type { SubdivisionsGeoJSON } from '@am/commons/schemas';

// Query keys
export const subdivisionKeys = {
  all: ['subdivisions'] as const,
  byId: (id: string) => [...subdivisionKeys.all, id] as const,
};

// API functions
const fetchSubdivisions = async (): Promise<SubdivisionsGeoJSON['features']> => {
  if (API_CONFIG.useMockData) {
    // Dynamically import mock data only when needed
    const subdivisions = await import('./mockdata/subdivision_geodata_davidson_47037.json');
    const typedSubdivisions = subdivisions.default as SubdivisionsGeoJSON;
    return addColorsToFeatures(typedSubdivisions.features);
  }

  // Real API call
  const response = await fetch(`${API_CONFIG.apiBaseUrl}/subdivisions`);
  if (!response.ok) {
    throw new Error(`Failed to fetch subdivisions: ${response.statusText}`);
  }
  const data = await response.json();
  return addColorsToFeatures(data.features);
};

// React Query hooks
export function useSubdivisionsQuery() {
  return useQuery({
    queryKey: subdivisionKeys.all,
    queryFn: fetchSubdivisions,
  });
}
