import { useState, useCallback, useEffect, useRef } from 'react';
import mapboxgl from 'mapbox-gl';

export interface MapLayer {
  id: string;
  visible: boolean;
  sourceId: string;
  layerIds: string[];
}

export interface PopupContent {
  title?: string;
  content: string;
  coordinates: [number, number];
  options?: mapboxgl.PopupOptions;
}

export function useMapLayers(map: mapboxgl.Map | null, initialLayers: MapLayer[] = []) {
  const [layers, setLayers] = useState<MapLayer[]>(initialLayers);
  const popupRef = useRef<mapboxgl.Popup | null>(null);

  // Create or update a popup
  const showPopup = useCallback(
    ({ title, content, coordinates, options = {} }: PopupContent) => {
      if (!map) return;

      // Remove existing popup if any
      if (popupRef.current) {
        popupRef.current.remove();
      }

      // Create new popup
      popupRef.current = new mapboxgl.Popup({
        closeButton: true,
        ...options,
      });

      // Set content and position
      popupRef.current.setLngLat(coordinates).setHTML(content).addTo(map);
    },
    [map]
  );

  // Remove current popup
  const removePopup = useCallback(() => {
    if (popupRef.current) {
      popupRef.current.remove();
      popupRef.current = null;
    }
  }, []);

  // Get current popup reference
  const getPopupRef = useCallback(() => popupRef.current, []);

  // Register a new layer
  const registerLayer = useCallback((layer: MapLayer) => {
    setLayers((prevLayers) => {
      // Check if the layer already exists
      const existingLayerIndex = prevLayers.findIndex((l) => l.id === layer.id);
      if (existingLayerIndex >= 0) {
        // Update existing layer
        const updatedLayers = [...prevLayers];
        updatedLayers[existingLayerIndex] = {
          ...updatedLayers[existingLayerIndex],
          ...layer,
        };
        return updatedLayers;
      }
      // Add new layer
      return [...prevLayers, layer];
    });
  }, []);

  // Toggle layer visibility
  const toggleLayerVisibility = useCallback(
    (layerId: string, visible?: boolean) => {
      if (!map) return;

      setLayers((prevLayers) => {
        return prevLayers.map((layer) => {
          if (layer.id === layerId) {
            const newVisibility = visible !== undefined ? visible : !layer.visible;

            // Toggle the visibility of all layers in this group
            layer.layerIds.forEach((id) => {
              if (map.getLayer(id)) {
                map.setLayoutProperty(id, 'visibility', newVisibility ? 'visible' : 'none');
              }
            });

            return {
              ...layer,
              visible: newVisibility,
            };
          }
          return layer;
        });
      });
    },
    [map]
  );

  // Get layer by ID
  const getLayer = useCallback(
    (layerId: string) => {
      return layers.find((layer) => layer.id === layerId);
    },
    [layers]
  );

  // Update layer properties (like source data)
  const updateLayerSource = useCallback(
    (layerId: string, data: GeoJSON.FeatureCollection) => {
      if (!map) return;

      const layer = getLayer(layerId);
      if (!layer) return;

      const source = map.getSource(layer.sourceId) as mapboxgl.GeoJSONSource;
      if (source) {
        source.setData(data);
      }
    },
    [map, getLayer]
  );

  // Initialize all layers on mount or when map changes
  useEffect(() => {
    if (!map) return;

    // Set all layers to their initial visibility
    layers.forEach((layer) => {
      layer.layerIds.forEach((id) => {
        if (map.getLayer(id)) {
          map.setLayoutProperty(id, 'visibility', layer.visible ? 'visible' : 'none');
        }
      });
    });
  }, [map, layers]);

  return {
    layers,
    registerLayer,
    toggleLayerVisibility,
    getLayer,
    updateLayerSource,
    showPopup,
    removePopup,
    getPopupRef,
  };
}
