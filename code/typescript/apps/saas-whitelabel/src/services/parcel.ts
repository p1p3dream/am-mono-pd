import { useQuery } from '@tanstack/react-query';
import { API_CONFIG } from '../config/api';
import type { ParcelGeoJSON } from '@am/commons/schemas';

// Query keys
export const parcelKeys = {
  all: ['parcels'] as const,
  byId: (id: string) => [...parcelKeys.all, id] as const,
};

// API functions
const fetchParcels = async (): Promise<ParcelGeoJSON['features']> => {
  if (API_CONFIG.useMockData) {
    // Dynamically import mock data only when needed
    const parcels = await import('./mockdata/davidson_min_parcels.json');
    return (parcels.default as ParcelGeoJSON).features;
  }

  // Real API call
  const response = await fetch(`${API_CONFIG.apiBaseUrl}/parcels`);
  if (!response.ok) {
    throw new Error(`Failed to fetch parcels: ${response.statusText}`);
  }
  const data = await response.json();
  return data.features;
};

// React Query hooks
export function useParcelsQuery() {
  return useQuery({
    queryKey: parcelKeys.all,
    queryFn: fetchParcels,
  });
}
