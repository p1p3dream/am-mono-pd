import { useEffect, useRef, useState } from 'react';
import { addLayer, addSource, createPopup, formatPopupContent } from '../utils/mapUtils';
import { useMapLayers, MapLayer } from '../hooks/useMapLayers';
import { SubdivisionFeature } from '../../schemas';

interface SubdivisionLayerProps {
  map: mapboxgl.Map | null;
  visible?: boolean;
  onSubdivisionClick?: (subdivision: SubdivisionFeature) => void;
  getSubdivisions: () => SubdivisionFeature[];
}

export const SUBDIVISION_LAYER_ID = 'subdivisions';

// Global variable for tracking initialization
let isSubdivisionLayerInitialized = false;

export function SubdivisionLayer({
  map,
  visible = false,
  onSubdivisionClick,
  getSubdivisions,
}: SubdivisionLayerProps) {
  const [subdivisions, setSubdivisions] = useState<SubdivisionFeature[]>([]);
  const popupRef = useRef<mapboxgl.Popup | null>(null);
  const hoveredSubdivisionId = useRef<string | null>(null);
  const { registerLayer, toggleLayerVisibility } = useMapLayers(map);

  // Fetch subdivisions data
  useEffect(() => {
    const subdivisionData = getSubdivisions();
    setSubdivisions(subdivisionData);
  }, []);

  // Initialize subdivision layers
  useEffect(() => {
    if (!map || !subdivisions.length) return;

    // Add source for subdivisions
    addSource(map, SUBDIVISION_LAYER_ID, {
      type: 'geojson',
      data: {
        type: 'FeatureCollection',
        features: subdivisions.map((subdivision) => {
          return {
            type: 'Feature',
            properties: subdivision.properties,
            geometry: subdivision.geometry,
          } as GeoJSON.Feature;
        }),
      },
    });

    // Add fill layer for subdivision boundaries
    addLayer(map, {
      id: 'subdivisions-fill',
      type: 'fill',
      source: SUBDIVISION_LAYER_ID,
      paint: {
        'fill-color': ['get', 'color'],
        'fill-opacity': 0.4,
        'fill-outline-color': '#ffffff',
      },
    });

    // Add outline layer for subdivision boundaries
    addLayer(map, {
      id: 'subdivisions-outline',
      type: 'line',
      source: SUBDIVISION_LAYER_ID,
      paint: {
        'line-color': '#ffffff',
        'line-width': 1.5,
        'line-opacity': 0.7,
      },
    });

    // Register the subdivision layer with the layer manager
    registerLayer({
      id: SUBDIVISION_LAYER_ID,
      visible,
      sourceId: SUBDIVISION_LAYER_ID,
      layerIds: ['subdivisions-fill', 'subdivisions-outline'],
    });

    // Set initial visibility
    toggleLayerVisibility(SUBDIVISION_LAYER_ID, visible);

    if (isSubdivisionLayerInitialized) return;

    // Add hover effects
    map.on('mouseenter', 'subdivisions-fill', (e) => {
      if (e.features && e.features.length > 0 && e.features[0]) {
        // Change cursor to pointer
        map.getCanvas().style.cursor = 'pointer';

        const feature = e.features[0];
        const props = feature.properties || {};
        const subdivisionId = props.ID as string;

        // Show subdivision info
        showSubdivisionInfo(subdivisionId, props);
      }
    });

    // Add mousemove handler to update popup position
    map.on('mousemove', 'subdivisions-fill', (e) => {
      if (e.features && e.features.length > 0 && e.features[0]) {
        const feature = e.features[0];
        const props = feature.properties || {};
        const subdivisionId = props.ID as string;

        // Show subdivision info
        showSubdivisionInfo(subdivisionId, props);

        // Position popup at mouse position
        if (e.lngLat && popupRef.current) {
          popupRef.current.setLngLat(e.lngLat).addTo(map);
        }
      }
    });

    // Mouse leave event
    map.on('mouseleave', 'subdivisions-fill', () => {
      // Reset cursor
      map.getCanvas().style.cursor = '';

      // Reset hover states
      if (hoveredSubdivisionId.current) {
        map.setFeatureState(
          { source: SUBDIVISION_LAYER_ID, id: hoveredSubdivisionId.current },
          { hover: false }
        );

        // Reset fill opacity
        map.setPaintProperty('subdivisions-fill', 'fill-opacity', 0.4);

        hoveredSubdivisionId.current = null;
      }

      // Remove popup
      if (popupRef.current) {
        popupRef.current.remove();
        popupRef.current = null;
      }
    });

    // Click event for subdivisions (only if handler provided)
    if (onSubdivisionClick) {
      map.on('click', 'subdivisions-fill', (e) => {
        if (e.features && e.features.length > 0 && e.features[0]) {
          const feature = e.features[0];
          // Find the corresponding subdivision feature
          const subdivisionFeature = subdivisions.find(
            (s) => s.properties.ID === feature.properties?.ID
          );
          if (subdivisionFeature) {
            onSubdivisionClick(subdivisionFeature);
          }
        }
      });
    }

    // Use global variable instead of state
    isSubdivisionLayerInitialized = true;
  }, [map, subdivisions, registerLayer, toggleLayerVisibility, visible, onSubdivisionClick]);

  // Update visibility when visibility prop changes
  useEffect(() => {
    if (isSubdivisionLayerInitialized) {
      toggleLayerVisibility(SUBDIVISION_LAYER_ID, visible);
    }
  }, [visible, toggleLayerVisibility]);

  // Helper function to show subdivision info on hover
  const showSubdivisionInfo = (subdivisionId: string, props: any) => {
    if (!map) return;

    // Set hovered state
    if (subdivisionId && hoveredSubdivisionId.current !== subdivisionId) {
      // If we hovered over another subdivision before, reset its style
      if (hoveredSubdivisionId.current) {
        map.setFeatureState(
          { source: SUBDIVISION_LAYER_ID, id: hoveredSubdivisionId.current },
          { hover: false }
        );
      }

      // Set hover state on this subdivision
      hoveredSubdivisionId.current = subdivisionId;
      map.setFeatureState({ source: SUBDIVISION_LAYER_ID, id: subdivisionId }, { hover: true });

      // Increase opacity of hovered subdivision
      map.setPaintProperty('subdivisions-fill', 'fill-opacity', [
        'case',
        ['==', ['get', 'ID'], subdivisionId],
        0.7, // Higher opacity for hovered subdivision
        0.4, // Default opacity for non-hovered subdivisions
      ]);

      // Remove existing popup
      if (popupRef.current) {
        popupRef.current.remove();
        popupRef.current = null;
      }
    }

    // Create popup if it doesn't exist
    if (!popupRef.current) {
      // Format popup content
      const content = formatPopupContent(props, {
        title: 'NAMELSAD',
        fields: [
          {
            key: 'NAME',
            label: 'District Number',
            format: (value) => `District ${value || 'N/A'}`,
          },
          {
            key: 'FPCLASS',
            label: 'Type',
            format: (value) => value?.replace('A nonfunctioning ', '') || 'N/A',
          },
          {
            key: 'STATE',
            label: 'State',
          },
          {
            key: 'LONGITUDE',
            label: 'Coordinates',
            format: (value) =>
              `${value?.toFixed(4) || 'N/A'}, ${props.LATITUDE?.toFixed(4) || 'N/A'}`,
          },
        ],
        noImage: true,
      });

      popupRef.current = createPopup(content, {
        className: 'subdivision-popup',
        offset: 20,
      });
    }
  };

  return null; // This is a logical component, doesn't render anything
}
