import mapboxgl from 'mapbox-gl';
import { MAPBOX_ACCESS_TOKEN } from '@/lib/env';

// Cache for storing previously fetched addresses
interface AddressCache {
  [key: string]: {
    address: string;
    city: string;
    state: string;
    zipCode: string;
    fullAddress: string;
  };
}

// In-memory cache to store previous lookups
const addressCache: AddressCache = {};

// Fetch address from coordinates using Mapbox's Reverse Geocoding API
export async function getAddressFromCoordinates(
  lng: number,
  lat: number
): Promise<{
  address: string;
  city: string;
  state: string;
  zipCode: string;
  fullAddress: string;
}> {
  // Create a cache key from coordinates
  const cacheKey = `${lng.toFixed(6)},${lat.toFixed(6)}`;

  // Check if we already have this address in cache
  if (addressCache[cacheKey]) {
    console.log('Using cached address data for', cacheKey);
    return addressCache[cacheKey];
  }

  try {
    // Call the Mapbox Geocoding API
    const response = await fetch(
      `https://api.mapbox.com/geocoding/v5/mapbox.places/${lng},${lat}.json?access_token=${MAPBOX_ACCESS_TOKEN}&types=address`
    );

    if (!response.ok) {
      throw new Error(`Mapbox API error: ${response.status}`);
    }

    const data = await response.json();

    // Parse the Mapbox response
    const result = parseMapboxResponse(data);

    // Cache the result
    addressCache[cacheKey] = result;

    return result;
  } catch (error) {
    console.error('Error fetching address from coordinates:', error);
    // Fallback to a default response in case of error
    return {
      address: 'Unknown address',
      city: 'Unknown city',
      state: 'Unknown state',
      zipCode: '00000',
      fullAddress: 'Address unavailable',
    };
  }
}

// Helper function to parse Mapbox geocoding response
function parseMapboxResponse(data: any): {
  address: string;
  city: string;
  state: string;
  zipCode: string;
  fullAddress: string;
} {
  // Default values
  let address = 'Unknown address';
  let city = 'Unknown city';
  let state = 'Unknown state';
  let zipCode = '00000';
  let fullAddress = '';

  // Check if we have features
  if (data.features && data.features.length > 0) {
    const feature = data.features[0];
    fullAddress = feature.place_name || '';

    // Extract components from the context array
    const context = feature.context || [];

    // Find address components from the context array
    // Street address is usually in the main feature text
    address = feature.text || 'Unknown address';

    // Extract city, state, and zip from context
    context.forEach((item: any) => {
      const id = item.id || '';
      if (id.startsWith('place')) {
        city = item.text || city;
      } else if (id.startsWith('region')) {
        state = item.text || state;
      } else if (id.startsWith('postcode')) {
        zipCode = item.text || zipCode;
      }
    });
  }

  return {
    address,
    city,
    state,
    zipCode,
    fullAddress,
  };
}
