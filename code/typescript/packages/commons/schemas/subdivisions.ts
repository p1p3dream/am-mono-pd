// Type definition for a single subdivision feature
export interface SubdivisionFeature {
  type: 'Feature';
  geometry: {
    type: 'Polygon';
    coordinates: number[][] | number[][][] | number[][][][];
  };
  properties: {
    ID?: string;
    STATEFPCD?: string;
    COUNTYFPCD?: string;
    COUSUBFPCD?: string;
    STATE?: string;
    STATEID?: string;
    COUNTYID?: string;
    GNISCD?: string;
    NAME?: string;
    NAMELSAD?: string;
    LSADCD?: string;
    LSAD?: string;
    FPCLASSCD?: string;
    FPCLASS?: string;
    MTFCCD?: string;
    MTFC?: string;
    FUNCSTATCD?: string;
    FUNCSTAT?: string;
    LONGITUDE?: number;
    LATITUDE?: number;
    color?: string;
  };
}

// Type definition for the entire GeoJSON structure
export interface SubdivisionsGeoJSON {
  type: 'FeatureCollection';
  name: string;
  crs: {
    type: string;
    properties: {
      name: string;
    };
  };
  features: SubdivisionFeature[];
}
