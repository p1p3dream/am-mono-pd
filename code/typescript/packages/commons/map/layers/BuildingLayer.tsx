import { useEffect, useCallback, useState } from 'react';
import mapboxgl from 'mapbox-gl';
import { BuildingFeature, BuildingBoundary } from '../map-component';
import { useMapLayers } from '../hooks/useMapLayers';

interface BuildingLayerProps {
  map: mapboxgl.Map | null;
  onBuildingSelect?: (building: BuildingFeature | null) => void;
  selectedBuilding?: BuildingFeature | null;
  allowSelection?: boolean;
}

export const BUILDING_LAYER_ID = 'selected-building';

// Global variable for tracking initialization
let isBuildingLayerInitialized = false;
let globalAllowSelection = true;
export function BuildingLayer({
  map,
  onBuildingSelect,
  selectedBuilding,
  allowSelection = true,
}: BuildingLayerProps) {
  // Remove local state since we're using a global variable
  // const [isLayerInitialized, setIsLayerInitialized] = useState(false);
  const { registerLayer } = useMapLayers(map);

  // Function to handle building click
  const handleBuildingClick = useCallback(
    (building: any, lngLat: { lng: number; lat: number }) => {
      console.log('building', building, allowSelection);
      if (!building || !map || !allowSelection) return;

      const buildingProps = building.properties || {};

      // Get address from properties or use a placeholder
      const addressFromProps =
        buildingProps.address ||
        buildingProps.addr ||
        buildingProps['addr:full'] ||
        buildingProps['addr:housenumber']
          ? `${buildingProps['addr:housenumber'] || ''} ${buildingProps['addr:street'] || ''}`.trim()
          : '';

      // Create a feature from the clicked building
      const buildingFeature = {
        type: 'Feature' as const,
        geometry: building.geometry,
        properties: {
          ...buildingProps,
          name:
            buildingProps.name || `Building at ${lngLat.lng.toFixed(4)}, ${lngLat.lat.toFixed(4)}`,
          address: addressFromProps,
        },
      };

      // Update the source data
      const source = map.getSource(BUILDING_LAYER_ID) as mapboxgl.GeoJSONSource;
      if (source) {
        source.setData({
          type: 'FeatureCollection',
          features: [buildingFeature],
        });
      }

      // Format the building data
      const buildingData: BuildingFeature = {
        lng: lngLat.lng,
        lat: lngLat.lat,
        boundary: buildingFeature as unknown as BuildingBoundary,
        properties: {
          height: buildingProps.height || 'N/A',
          floors: buildingProps.floors || buildingProps.levels || 'N/A',
          buildingType: buildingProps.type || buildingProps.building || 'N/A',
          yearBuilt: buildingProps.year || buildingProps['start_date'] || 'N/A',
          address: addressFromProps || 'Unknown address',
          name: buildingProps.name || 'Unnamed Building',
          bedrooms: buildingProps.bedrooms || 'N/A',
          bathrooms: buildingProps.bathrooms || 'N/A',
          sqft: buildingProps.sqft || 'N/A',
          lotSize: buildingProps['lot_size'] || 'N/A',
          type: buildingProps.type || 'N/A',
          hoa: buildingProps.hoa || 'None',
          built: buildingProps.built || buildingProps.year || 'N/A',
          stories: buildingProps.stories || buildingProps.floors || 'N/A',
          status: buildingProps.status || 'N/A',
          value: buildingProps.value,
          rent: buildingProps.rent,
        },
      };

      // Call the callback if provided
      if (onBuildingSelect) {
        onBuildingSelect(buildingData);
      }
    },
    [map, onBuildingSelect, allowSelection]
  );

  // Function to clear building selection
  const clearBuildingSelection = useCallback(() => {
    if (!map) return;

    const source = map.getSource(BUILDING_LAYER_ID) as mapboxgl.GeoJSONSource;
    if (source) {
      source.setData({
        type: 'FeatureCollection',
        features: [],
      });
    }

    // Call the callback with null if provided
    if (onBuildingSelect) {
      onBuildingSelect(null);
    }
  }, [map, onBuildingSelect]);

  // Initialize the building layer
  useEffect(() => {
    if (!map || isBuildingLayerInitialized) return;

    // Add source for selected building if it doesn't exist
    if (!map.getSource(BUILDING_LAYER_ID)) {
      map.addSource(BUILDING_LAYER_ID, {
        type: 'geojson',
        data: {
          type: 'FeatureCollection',
          features: [],
        },
      });

      // Add a fill layer for the selected building
      map.addLayer({
        id: 'selected-building-fill',
        type: 'fill',
        source: BUILDING_LAYER_ID,
        paint: {
          'fill-color': '#7CFC00', // Bright green
          'fill-opacity': 0.5,
        },
      });

      // Add an outline layer for the selected building
      map.addLayer({
        id: 'selected-building-outline',
        type: 'line',
        source: BUILDING_LAYER_ID,
        paint: {
          'line-color': '#7CFC00', // Bright green
          'line-width': 2,
        },
      });

      // Register the layer with the layer manager
      registerLayer({
        id: BUILDING_LAYER_ID,
        visible: true,
        sourceId: BUILDING_LAYER_ID,
        layerIds: ['selected-building-fill', 'selected-building-outline'],
      });

      // Handle map clicks
      map.on('click', (e) => {
        if (!globalAllowSelection) return;

        // Prevent event bubbling
        e.preventDefault();
        e.originalEvent.stopPropagation();

        // Get features at the clicked point
        const features = map.queryRenderedFeatures(e.point, {
          layers: ['building'], // This targets Mapbox's built-in building layer
        });

        if (features && features.length > 0 && features[0]) {
          handleBuildingClick(features[0], e.lngLat);
        } else {
          clearBuildingSelection();
        }
      });

      // Use global variable instead of local state
      isBuildingLayerInitialized = true;
    }
  }, [
    map,
    handleBuildingClick,
    clearBuildingSelection,
    registerLayer,
    allowSelection,
    // Remove isLayerInitialized from dependencies
  ]);

  // Update the selected building when it changes
  useEffect(() => {
    if (!map || !isBuildingLayerInitialized || !selectedBuilding || !selectedBuilding.boundary)
      return;

    const source = map.getSource(BUILDING_LAYER_ID) as mapboxgl.GeoJSONSource;
    if (source) {
      source.setData({
        type: 'FeatureCollection',
        features: [selectedBuilding.boundary as any],
      });
    }
  }, [map, selectedBuilding]);

  // Enable/disable selection based on allowSelection prop
  useEffect(() => {
    globalAllowSelection = allowSelection;
  }, [allowSelection]);

  return null; // This is a logical component, doesn't render anything
}
