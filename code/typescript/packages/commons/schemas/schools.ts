export interface School {
  id: string;
  name: string;
  address: string;
  level: string;
  grade_low: number;
  grade_high: number;
  rating: string;
  latitude: number;
  longitude: number;
}

// Type definition for a single school feature

export interface SchoolFeature {
  type: 'Feature';
  geometry: {
    type: 'Polygon';
    coordinates: number[][] | number[][][] | number[][][][];
  };
  properties: {
    ID?: string;
    name?: string;
    district?: string;
    address?: string;
    city?: string;
    state?: string;
    zipCode?: string;
    level?: 'Elementary' | 'Middle' | 'High' | 'District' | string;
    rating?: number;
    grades?: string;
    students?: number;
    phone?: string;
    website?: string;
    latitude?: number;
    longitude?: number;
    color?: string;
    schools?: School[];
  };
}

// Type definition for the entire GeoJSON structure
export interface SchoolsGeoJSON {
  type: 'FeatureCollection';
  features: SchoolFeature[];
  // Add any metadata that might be in the GeoJSON
  metadata?: {
    generated?: string;
    title?: string;
    description?: string;
    [key: string]: any;
  };
}
