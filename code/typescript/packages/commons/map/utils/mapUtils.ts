import mapboxgl from 'mapbox-gl';

// Helper to create a popup for features
export function createPopup(content: string, options?: mapboxgl.PopupOptions): mapboxgl.Popup {
  return new mapboxgl.Popup({
    closeButton: false,
    closeOnClick: false,
    offset: 15,
    ...options,
  }).setHTML(content);
}
// Helper to format popup content
export function formatPopupContent(
  props: Record<string, any>,
  config: {
    title?: string;
    fields?: Array<{
      key: string;
      label?: string;
      format?: (value: any) => string;
    }>;
    imageUrl?: string;
    noImage?: boolean;
  }
): string {
  const { title, fields = [], imageUrl, noImage = false } = config;

  const titleContent = title
    ? `<div style="font-weight: bold; margin-bottom: 4px;">${
        props[title] || 'Name Unavailable'
      }</div>`
    : '';

  // Generate a house placeholder SVG (from Lucide)
  const housePlaceholderSvg = `
    <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="lucide lucide-home" style="color: #888;">
      <path d="m3 9 9-7 9 7v11a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2z"></path>
      <polyline points="9 22 9 12 15 12 15 22"></polyline>
    </svg>
  `;

  // Create the image content with skeleton loader and house placeholder
  const imageContent = `
    <div style="position: relative; width: 100%; padding-top: 75%; /* 4:3 aspect ratio */ margin-bottom: 8px; background-color: #2a2a2a; border-radius: 4px; overflow: hidden;">
      <div style="position: absolute; top: 0; left: 0; width: 100%; height: 100%; display: flex; align-items: center; justify-content: center;">
        ${housePlaceholderSvg}
      </div>
      ${
        imageUrl && props[imageUrl]
          ? `<div class="skeleton" style="position: absolute; top: 0; left: 0; width: 100%; height: 100%; background-color: #3a3a3a; animation: pulse 1.5s infinite ease-in-out;"></div>
         <img src="${props[imageUrl]}" style="position: absolute; top: 0; left: 0; width: 100%; height: 100%; object-fit: cover; border-radius: 4px; opacity: 0; transition: opacity 0.3s ease;" onload="this.style.opacity='1'; this.previousElementSibling.style.display='none';">`
          : ''
      }
    </div>
    <style>
      @keyframes pulse {
        0% { opacity: 0.6; }
        50% { opacity: 0.8; }
        100% { opacity: 0.6; }
      }
    </style>
  `;

  const fieldsContent = fields
    .map((field) => {
      const value = props[field.key];
      const formattedValue = field.format ? field.format(value) : value;
      const label = field.label ? `${field.label}: ` : '';

      return value ? `<div>${label}${formattedValue}</div>` : '';
    })
    .join('');

  return `
    <div style="font-family: sans-serif; padding: 8px; max-width: 250px; color: #fff; ">
      ${noImage ? '' : imageContent}
      ${titleContent}
      ${fieldsContent}
    </div>
  `;
}

// Helper to add a feature as selected
export function highlightFeature(
  map: mapboxgl.Map,
  source: string,
  featureId: string | number,
  highlightProperty: string = 'hover',
  value: boolean = true
): void {
  map.setFeatureState({ source, id: featureId }, { [highlightProperty]: value });
}

// Helper to load an image into the map
export function loadMapImage(map: mapboxgl.Map, url: string, id: string): Promise<void> {
  return new Promise((resolve, reject) => {
    // Check if image is already loaded
    if (map.hasImage(id)) {
      resolve();
      return;
    }

    map.loadImage(url, (error, image) => {
      if (error) {
        reject(error);
        return;
      }

      map.addImage(id, image as ImageBitmap);
      resolve();
    });
  });
}

// Helper to create a point feature from coordinates
export function createPointFeature(
  lng: number,
  lat: number,
  properties: Record<string, any> = {}
): GeoJSON.Feature<GeoJSON.Point> {
  return {
    type: 'Feature',
    geometry: {
      type: 'Point',
      coordinates: [lng, lat],
    },
    properties,
  };
}

// Helper to calculate centroid of a polygon
export function calculatePolygonCentroid(
  coordinates: number[][][] // Polygon coordinates
): [number, number] | null {
  // For simple polygon, use the first ring of coordinates
  if (!coordinates || !Array.isArray(coordinates[0]) || coordinates[0].length === 0) {
    return null;
  }

  const ring = coordinates[0];

  let sumX = 0;
  let sumY = 0;
  let validPoints = 0;

  for (let i = 0; i < ring.length; i++) {
    const point = ring[i];
    if (
      Array.isArray(point) &&
      point.length >= 2 &&
      typeof point[0] === 'number' &&
      typeof point[1] === 'number'
    ) {
      sumX += point[0];
      sumY += point[1];
      validPoints++;
    }
  }

  return validPoints > 0 ? [sumX / validPoints, sumY / validPoints] : null;
}

export function addLayer(map: mapboxgl.Map, layer: mapboxgl.Layer) {
  if (!map) return;

  if (!map.getLayer(layer.id)) {
    map.addLayer(layer);
  }
}

export function addSource(
  map: mapboxgl.Map,
  sourceId: string,
  source: mapboxgl.GeoJSONSourceSpecification
) {
  if (!map) return;

  const currentSource = map.getSource(sourceId) as mapboxgl.GeoJSONSource;

  if (!currentSource) {
    map.addSource(sourceId, source);
  } else {
    if (source) {
      currentSource.setData(source.data as GeoJSON.FeatureCollection);
    }
  }
}

// Helper to check if a point is inside a polygon
export function pointInPolygon(
  point: [number, number],
  polygon: number[][] | number[][][] | number[][][][]
): boolean {
  if (!polygon || !Array.isArray(polygon) || polygon.length === 0) {
    return false;
  }

  // Handle multi-polygon case by checking the first polygon
  const coordinates =
    Array.isArray(polygon[0]) && Array.isArray(polygon[0][0])
      ? polygon[0]
      : (polygon as number[][]);

  if (!coordinates || coordinates.length < 3) {
    return false;
  }

  let inside = false;
  for (let i = 0, j = coordinates.length - 1; i < coordinates.length; j = i++) {
    const coordI = coordinates[i] as number[];
    const coordJ = coordinates[j] as number[];

    if (!coordI || !coordJ || coordI.length < 2 || coordJ.length < 2) {
      continue;
    }

    const xi = coordI[0];
    const yi = coordI[1];
    const xj = coordJ[0];
    const yj = coordJ[1];

    if (
      typeof xi !== 'number' ||
      typeof yi !== 'number' ||
      typeof xj !== 'number' ||
      typeof yj !== 'number'
    ) {
      continue;
    }

    const intersect =
      yi > point[1] !== yj > point[1] && point[0] < ((xj - xi) * (point[1] - yi)) / (yj - yi) + xi;

    if (intersect) inside = !inside;
  }

  return inside;
}
