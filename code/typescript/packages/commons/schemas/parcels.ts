export interface ParcelFeature {
  type: 'Feature';
  geometry: {
    type: 'Polygon';
    coordinates: number[][] | number[][][] | number[][][][];
  };
}

export interface ParcelGeoJSON {
  type: 'FeatureCollection';
  features: ParcelFeature[];
}
