import { useEffect, useRef, useState } from 'react';
import { useMapLayers } from '../hooks/useMapLayers';
import { addSource, addLayer } from '../utils/mapUtils';

export const POLYGON_DRAW_LAYER_ID = 'polygon-draw';

interface PolygonDrawLayerProps {
  map: mapboxgl.Map | null;
  visible?: boolean;
  onPolygonComplete?: (coordinates: [number, number][]) => void;
}

export function PolygonDrawLayer({
  map,
  visible = false,
  onPolygonComplete,
}: PolygonDrawLayerProps) {
  const [isDrawing, setIsDrawing] = useState(false);
  const [startPoint, setStartPoint] = useState<[number, number] | null>(null);
  const [currentPoint, setCurrentPoint] = useState<[number, number] | null>(null);
  const { registerLayer, toggleLayerVisibility } = useMapLayers(map);
  const popupRef = useRef<mapboxgl.Popup | null>(null);

  // Handle map drag control
  useEffect(() => {
    if (!map) return;

    // Store original drag state
    const originalDragState = map.dragPan.isEnabled();

    // Disable drag when selection mode is active
    if (isDrawing) {
      map.dragPan.disable();
    } else {
      map.dragPan.enable();
    }

    // Restore original drag state when component unmounts
    return () => {
      if (originalDragState) {
        map.dragPan.enable();
      } else {
        map.dragPan.disable();
      }
    };
  }, [map, isDrawing]);

  // Initialize layers only once when map is ready
  useEffect(() => {
    if (!map) return;

    // Add source for the rectangle
    addSource(map, POLYGON_DRAW_LAYER_ID, {
      type: 'geojson',
      data: {
        type: 'Feature',
        properties: {},
        geometry: {
          type: 'Polygon',
          coordinates: [[]],
        },
      },
    });

    // Add fill layer for the rectangle
    addLayer(map, {
      id: 'polygon-fill',
      type: 'fill',
      source: POLYGON_DRAW_LAYER_ID,
      paint: {
        'fill-color': '#0080ff',
        'fill-opacity': 0.3,
        'fill-outline-color': '#ffffff',
      },
    });

    // Add outline layer for the rectangle
    addLayer(map, {
      id: 'polygon-outline',
      type: 'line',
      source: POLYGON_DRAW_LAYER_ID,
      paint: {
        'line-color': '#ffffff',
        'line-width': 2,
        'line-dasharray': [2, 2],
      },
    });

    // Register the layers
    registerLayer({
      id: POLYGON_DRAW_LAYER_ID,
      visible,
      sourceId: POLYGON_DRAW_LAYER_ID,
      layerIds: ['polygon-fill', 'polygon-outline'],
    });

    // Set initial visibility
    toggleLayerVisibility(POLYGON_DRAW_LAYER_ID, visible);
  }, [map, registerLayer, toggleLayerVisibility, visible]);

  // Handle drawing events
  useEffect(() => {
    if (!map || !isDrawing) return;

    // Add mousedown handler to start drawing
    const handleMouseDown = (e: mapboxgl.MapMouseEvent) => {
      const { lng, lat } = e.lngLat;
      setStartPoint([lng, lat]);
      setCurrentPoint([lng, lat]);
    };

    // Add mousemove handler to update rectangle
    const handleMouseMove = (e: mapboxgl.MapMouseEvent) => {
      if (!startPoint) return;

      const { lng, lat } = e.lngLat;
      setCurrentPoint([lng, lat]);

      // Calculate rectangle coordinates
      const coordinates: [number, number][] = [
        [startPoint[0], startPoint[1]],
        [lng, startPoint[1]],
        [lng, lat],
        [startPoint[0], lat],
        [startPoint[0], startPoint[1]],
      ];

      // Update the rectangle source
      const source = map.getSource(POLYGON_DRAW_LAYER_ID) as mapboxgl.GeoJSONSource;
      if (source) {
        source.setData({
          type: 'Feature',
          properties: {},
          geometry: {
            type: 'Polygon',
            coordinates: [coordinates],
          },
        });
      }
    };

    // Add mouseup handler to complete rectangle
    const handleMouseUp = () => {
      if (!startPoint || !currentPoint) return;

      // Calculate final rectangle coordinates
      const coordinates: [number, number][] = [
        [startPoint[0], startPoint[1]],
        [currentPoint[0], startPoint[1]],
        [currentPoint[0], currentPoint[1]],
        [startPoint[0], currentPoint[1]],
        [startPoint[0], startPoint[1]],
      ];

      // Update the rectangle source with final coordinates
      const source = map.getSource(POLYGON_DRAW_LAYER_ID) as mapboxgl.GeoJSONSource;
      if (source) {
        source.setData({
          type: 'Feature',
          properties: {},
          geometry: {
            type: 'Polygon',
            coordinates: [coordinates],
          },
        });
      }

      // Call the completion callback
      if (onPolygonComplete) {
        onPolygonComplete(coordinates);
      }

      // Reset points for next selection
      setStartPoint(null);
      setCurrentPoint(null);

      // Clear the rectangle after a short delay
      setTimeout(() => {
        if (source) {
          source.setData({
            type: 'Feature',
            properties: {},
            geometry: {
              type: 'Polygon',
              coordinates: [[]],
            },
          });
        }
      }, 100);
    };

    // Add event listeners
    map.on('mousedown', handleMouseDown);
    map.on('mousemove', handleMouseMove);
    map.on('mouseup', handleMouseUp);

    // Cleanup
    return () => {
      map.off('mousedown', handleMouseDown);
      map.off('mousemove', handleMouseMove);
      map.off('mouseup', handleMouseUp);
    };
  }, [map, isDrawing, startPoint, currentPoint, onPolygonComplete]);

  // Update visibility when prop changes
  useEffect(() => {
    if (map) {
      toggleLayerVisibility(POLYGON_DRAW_LAYER_ID, visible);
    }
  }, [visible, toggleLayerVisibility, map]);

  // Expose methods to parent component
  useEffect(() => {
    if (map) {
      // @ts-ignore - Adding custom methods to map object
      map.startPolygonDrawing = () => {
        setIsDrawing(true);
        setStartPoint(null);
        setCurrentPoint(null);
      };

      // @ts-ignore
      map.stopPolygonDrawing = () => {
        setIsDrawing(false);
        setStartPoint(null);
        setCurrentPoint(null);
      };
    }
  }, [map]);

  return null;
}
