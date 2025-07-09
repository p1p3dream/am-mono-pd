// Common color palette for map features
export const MAP_COLOR_PALETTE = [
  '#FF5733', // Bright Red-Orange
  '#3498DB', // Bright Blue
  '#2ECC71', // Bright Green
  '#9B59B6', // Purple
  '#F1C40F', // Yellow
  '#E74C3C', // Red
  '#1ABC9C', // Turquoise
  '#34495E', // Navy Blue
  '#D35400', // Orange
  '#8E44AD', // Violet
  '#16A085', // Dark Turquoise
  '#2980B9', // Ocean Blue
  '#27AE60', // Dark Green
  '#E67E22', // Carrot Orange
  '#C0392B', // Dark Red
  '#7D3C98', // Dark Purple
  '#2C3E50', // Dark Gray-Blue
  '#F39C12', // Orange-Yellow
  '#00BFFF', // Deep Sky Blue
  '#FF1493', // Deep Pink
];

// Generic interface for any feature that can have a color property
export interface GeoFeature {
  type: 'Feature';
  geometry: {
    type: string;
    coordinates: any;
  };
  properties: Record<string, any>;
}

/**
 * Adds colors to GeoJSON features
 *
 * @param features - Array of GeoJSON features
 * @returns The same features with color properties added
 */
export function addColorsToFeatures<T extends GeoFeature>(features: T[]): T[] {
  // Helper function to shuffle array (Fisher-Yates algorithm)
  const shuffleArray = <U>(array: U[]): U[] => {
    const shuffled = [...array];
    for (let i = shuffled.length - 1; i > 0; i--) {
      const j = Math.floor(Math.random() * (i + 1));
      // Swap elements using type assertion to handle undefined
      const temp = shuffled[i] as U;
      shuffled[i] = shuffled[j] as U;
      shuffled[j] = temp;
    }
    return shuffled;
  };

  // Shuffle features to ensure random color distribution
  // This makes adjacent features unlikely to have the same color
  const shuffledFeatures = shuffleArray(features);

  // Assign colors based on index in the shuffled array
  const featuresWithColors = shuffledFeatures.map((feature, index) => {
    // Assign color based on index in the array (ensures variation)
    const colorIndex = index % MAP_COLOR_PALETTE.length;
    const color = MAP_COLOR_PALETTE[colorIndex];

    // Return a new feature object with modified properties
    return {
      ...feature,
      properties: {
        ...feature.properties,
        color: color,
      },
    };
  });

  return featuresWithColors;
}
