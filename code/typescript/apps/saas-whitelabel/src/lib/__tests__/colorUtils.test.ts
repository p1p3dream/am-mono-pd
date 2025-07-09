import { describe, it, expect, beforeEach } from '@jest/globals';
import { addColorsToFeatures, MAP_COLOR_PALETTE } from '../colorUtils';
import type { GeoFeature } from '../colorUtils';

describe('colorUtils', () => {
  beforeEach(() => {
    jest.clearAllMocks();
  });

  describe('MAP_COLOR_PALETTE', () => {
    it('should have 20 unique colors', () => {
      expect(MAP_COLOR_PALETTE.length).toBe(20);
      const uniqueColors = new Set(MAP_COLOR_PALETTE);
      expect(uniqueColors.size).toBe(20);
    });

    it('should contain valid hex color codes', () => {
      const hexColorRegex = /^#[0-9A-F]{6}$/i;
      MAP_COLOR_PALETTE.forEach((color) => {
        expect(color).toMatch(hexColorRegex);
      });
    });
  });

  describe('addColorsToFeatures', () => {
    const mockFeatures: GeoFeature[] = [
      {
        type: 'Feature',
        geometry: {
          type: 'Point',
          coordinates: [0, 0],
        },
        properties: {},
      },
      {
        type: 'Feature',
        geometry: {
          type: 'Point',
          coordinates: [1, 1],
        },
        properties: {},
      },
      {
        type: 'Feature',
        geometry: {
          type: 'Point',
          coordinates: [2, 2],
        },
        properties: {},
      },
    ];

    it('should add color property to each feature', () => {
      const result = addColorsToFeatures(mockFeatures);

      result.forEach((feature) => {
        expect(feature.properties).toHaveProperty('color');
        expect(MAP_COLOR_PALETTE).toContain(feature.properties.color);
      });
    });

    it('should preserve existing properties', () => {
      const featuresWithProps = mockFeatures.map((feature) => ({
        ...feature,
        properties: {
          ...feature.properties,
          name: 'Test Feature',
        },
      }));

      const result = addColorsToFeatures(featuresWithProps);

      result.forEach((feature) => {
        expect(feature.properties).toHaveProperty('name', 'Test Feature');
        expect(feature.properties).toHaveProperty('color');
      });
    });

    it('should handle empty feature array', () => {
      const result = addColorsToFeatures([]);
      expect(result).toEqual([]);
    });

    it('should distribute colors evenly across features', () => {
      const result = addColorsToFeatures(mockFeatures);
      const colors = result.map((f) => f.properties.color);

      // Check that we have unique colors for each feature
      const uniqueColors = new Set(colors);
      expect(uniqueColors.size).toBe(Math.min(mockFeatures.length, MAP_COLOR_PALETTE.length));
    });
  });
});
