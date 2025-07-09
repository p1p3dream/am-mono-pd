import { useQuery } from '@tanstack/react-query';
import { API_CONFIG } from '../config/api';
import type { PropertyDetails } from '@am/commons/schemas';
import type { ComparableProperty } from '../contexts/property-context';

// Query keys
export const propertyKeys = {
  all: ['property'] as const,
  detail: (id: string) => [...propertyKeys.all, id] as const,
  comparables: (id: string) => [...propertyKeys.all, id, 'comparables'] as const,
};

// Helper function to calculate distance between two points
function calculateDistance(lat1: number, lon1: number, lat2: number, lon2: number): number {
  const R = 3958.8; // Radius of Earth in miles
  const toRadians = (deg: number): number => (deg * Math.PI) / 180;

  const phi1 = toRadians(lat1);
  const phi2 = toRadians(lat2);
  const deltaPhi = toRadians(lat2 - lat1);
  const deltaLambda = toRadians(lon2 - lon1);

  const a =
    Math.sin(deltaPhi / 2) ** 2 + Math.cos(phi1) * Math.cos(phi2) * Math.sin(deltaLambda / 2) ** 2;
  const c = 2 * Math.atan2(Math.sqrt(a), Math.sqrt(1 - a));

  return R * c;
}

// API functions
const fetchPropertyById = async (id: string): Promise<PropertyDetails> => {
  if (API_CONFIG.useMockData) {
    // Dynamically import mock data only when needed
    const { mockPropertyData, mockComparableProperties } = await import(
      './mockdata/property-context-mock'
    );

    if (mockPropertyData.id === id) {
      return mockPropertyData;
    }

    const foundProperty = mockComparableProperties.find((p) => p.id === id);
    if (!foundProperty) {
      throw new Error(`Property with ID ${id} not found`);
    }

    return foundProperty;
  }

  // Real API call
  const response = await fetch(`${API_CONFIG.apiBaseUrl}/properties/${id}`);
  if (!response.ok) {
    throw new Error(`Failed to fetch property: ${response.statusText}`);
  }
  return response.json();
};

const fetchComparableProperties = async (propertyId: string): Promise<ComparableProperty[]> => {
  if (API_CONFIG.useMockData) {
    // Dynamically import mock data only when needed
    const { mockComparableProperties } = await import('./mockdata/property-context-mock');
    const property = await fetchPropertyById(propertyId);

    return mockComparableProperties.map((p) => ({
      ...p,
      distance: calculateDistance(
        property.location.lat,
        property.location.lng,
        p.location.lat,
        p.location.lng
      ),
      closeDate: p.listingDate,
    }));
  }

  // Real API call
  const response = await fetch(`${API_CONFIG.apiBaseUrl}/properties/${propertyId}/comparables`);
  if (!response.ok) {
    throw new Error(`Failed to fetch comparable properties: ${response.statusText}`);
  }
  return response.json();
};

// React Query hooks
export function usePropertyQuery(id: string) {
  return useQuery({
    queryKey: propertyKeys.detail(id),
    queryFn: () => fetchPropertyById(id),
  });
}

export function useComparablePropertiesQuery(propertyId: string) {
  return useQuery({
    queryKey: propertyKeys.comparables(propertyId),
    queryFn: () => fetchComparableProperties(propertyId),
  });
}
