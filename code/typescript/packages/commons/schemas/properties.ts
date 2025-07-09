// Define the property image type
export interface PropertyImage {
  id: number;
  url: string;
  alt: string;
}

// Define the property details type
export interface PropertyDetails {
  id: string;
  address: string;
  city: string;
  state: string;
  zipCode: string;
  fullAddress: string;
  status: string;
  price: number;
  beds: number;
  baths: number;
  sqft: number;
  lotSize: number;
  yearBuilt: number;
  propertyType: string;
  description: string;
  features: string[];
  mainImage: string;
  images: PropertyImage[];
  location: {
    lat: number;
    lng: number;
  };
  mapboxData?: any; // Make mapboxData optional since it might not exist initially

  // Additional MLS specific fields
  mlsNumber?: string;
  mlsSource?: string;
  listingDate?: string;
  originalListingDate?: string;
  originalListingPrice?: number;
  daysOnMarket?: number;
  schoolElementary?: string;
  schoolMiddle?: string;
  schoolHigh?: string;
  schoolDistrict?: string;
  taxAmount?: number;
  specialListingConditions?: string;
  constructionMaterials?: string;
  basementFeatures?: string;
  garageSpaces?: number;
  lotSizeAcres?: number;
  lotDimensions?: string;
  listingAgentName?: string;
  listingAgentPhone?: string;
  listingAgentEmail?: string;
  listingOfficeName?: string;
  listingOfficePhone?: string;
  newConstruction?: boolean;
  photoURLPrefix?: string;
  photoKey?: string;
  photosCount?: number;
  hasDotMarker?: boolean;
  GeoID?: string;
  parcelGeometry?: {
    type: 'Polygon';
    coordinates: number[][] | number[][][] | number[][][][];
  };
  parcelProperties?: {
    ID: string;
    Address: string;
    City: string;
    State: string;
  };
}
