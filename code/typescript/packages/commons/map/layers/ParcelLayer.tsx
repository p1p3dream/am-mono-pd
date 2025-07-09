import { useEffect, useRef, useState } from 'react';
import { createPopup, formatPopupContent, addSource, addLayer } from '../utils/mapUtils';
import { useMapLayers, MapLayer } from '../hooks/useMapLayers';
import { ParcelFeature, PropertyDetails } from '@am/commons/schemas';
import mapboxgl from 'mapbox-gl';

interface ParcelLayerProps {
  map: mapboxgl.Map | null;
  visible?: boolean;
  onParcelClick?: (parcel: ParcelFeature) => void;
  minZoomLevel?: number; // Optional prop to override default minimum zoom level
  getParcels: () => ParcelFeature[];
  property?: PropertyDetails;
  comparableProperties?: PropertyDetails[];
}

export const PARCEL_LAYER_ID = 'parcels';

// Global variable for tracking initialization
let isParcelLayerInitialized = false;

// Default minimum zoom level at which parcels will be displayed
const DEFAULT_MIN_ZOOM_LEVEL = 8;

export function ParcelLayer({
  map,
  visible = false,
  onParcelClick,
  minZoomLevel = DEFAULT_MIN_ZOOM_LEVEL,
  getParcels,
  property,
  comparableProperties,
}: ParcelLayerProps) {
  const [parcels, setParcels] = useState<ParcelFeature[]>([]);
  const hoveredParcelId = useRef<string | null>(null);
  const { registerLayer, toggleLayerVisibility, showPopup, removePopup } = useMapLayers(map);

  // Calculate effective visibility based on zoom level and visible prop
  const isEffectivelyVisible = visible;

  // Fetch parcels data
  useEffect(() => {
    const parcelData = getParcels();
    setParcels(parcelData);
  }, []);

  // Initialize parcel layers
  useEffect(() => {
    if (!map || !parcels.length) return;

    const parcelGeoJSON: GeoJSON.FeatureCollection = {
      type: 'FeatureCollection',
      features: parcels as GeoJSON.Feature[],
    };

    // Add source for parcels
    addSource(map, PARCEL_LAYER_ID, {
      type: 'geojson',
      data: parcelGeoJSON,
    });

    // Check if source was just added or already existed
    const sourceExists = map.getSource(PARCEL_LAYER_ID);

    if (sourceExists) {
      // If the source exists but we need to update the data
      if (isParcelLayerInitialized) {
        const source = map.getSource(PARCEL_LAYER_ID) as mapboxgl.GeoJSONSource;
        if (source) {
          source.setData(parcelGeoJSON);
        }
      }
    }

    // Add fill layer for parcel boundaries
    addLayer(map, {
      id: 'parcels-fill',
      type: 'fill',
      source: PARCEL_LAYER_ID,
      paint: {
        'fill-color': '#4CAF50', // Green color for all parcels
        'fill-opacity': 0.5,
        'fill-outline-color': '#ffffff',
      },
      filter: ['all', ['>=', ['zoom'], minZoomLevel]],
    });

    // Add outline layer for parcel boundaries
    addLayer(map, {
      id: 'parcels-outline',
      type: 'line',
      source: PARCEL_LAYER_ID,
      paint: {
        'line-color': '#ffffff',
        'line-width': 1,
      },
      filter: ['all', ['>=', ['zoom'], minZoomLevel]],
    });

    // Register the parcel layer with the layer manager
    registerLayer({
      id: PARCEL_LAYER_ID,
      visible: isEffectivelyVisible,
      sourceId: PARCEL_LAYER_ID,
      layerIds: ['parcels-fill', 'parcels-outline'],
    });

    // Set initial visibility
    toggleLayerVisibility(PARCEL_LAYER_ID, isEffectivelyVisible);

    if (isParcelLayerInitialized) return;

    // Add hover effects
    map.on('mouseenter', 'parcels-fill', (e) => {
      if (e.features && e.features.length > 0 && e.features[0]) {
        // Change cursor to pointer
        map.getCanvas().style.cursor = 'pointer';

        const feature = e.features[0];
        const props = feature.properties || {};
        const parcelId = props.ID as string;

        // Show parcel info
        showParcelInfo(parcelId, props);
      }
    });

    // Add mousemove handler to update popup position
    map.on('mousemove', 'parcels-fill', (e) => {
      if (e.features && e.features.length > 0 && e.features[0]) {
        const feature = e.features[0];
        const props = feature.properties || {};
        const parcelId = props.ID as string;

        // Show parcel info
        showParcelInfo(parcelId, props);

        // Position popup at mouse position
        if (e.lngLat && hoveredParcelId.current) {
          map.setFeatureState(
            { source: PARCEL_LAYER_ID, id: hoveredParcelId.current },
            { hover: false }
          );
          map.setPaintProperty('parcels-fill', 'fill-opacity', 0.5);
          hoveredParcelId.current = null;
          showParcelInfo(parcelId, props);
        }
      }
    });

    // Mouse leave event
    map.on('mouseleave', 'parcels-fill', () => {
      // Reset cursor
      map.getCanvas().style.cursor = '';

      // Reset hover states
      if (hoveredParcelId.current) {
        map.setFeatureState(
          { source: PARCEL_LAYER_ID, id: hoveredParcelId.current },
          { hover: false }
        );

        // Reset fill opacity
        map.setPaintProperty('parcels-fill', 'fill-opacity', 0.5);

        hoveredParcelId.current = null;
      }

      // Remove popup using the hook
      removePopup();
    });

    // Click event for parcels (only if handler provided)
    if (onParcelClick) {
      map.on('click', 'parcels-fill', (e) => {
        if (e.features && e.features.length > 0 && e.features[0]) {
          const feature = e.features[0];
          // Find the corresponding parcel feature
          const parcelFeature = parcels.find((p, index) => index === parcels.indexOf(p));
          if (parcelFeature) {
            onParcelClick(parcelFeature);
          }
        }
      });
    }

    // Use global variable instead of state
    isParcelLayerInitialized = true;
  }, [
    map,
    parcels,
    registerLayer,
    toggleLayerVisibility,
    isEffectivelyVisible,
    onParcelClick,
    showPopup,
    removePopup,
  ]);

  // Update visibility when visibility prop changes or zoom level changes
  useEffect(() => {
    if (isParcelLayerInitialized) {
      toggleLayerVisibility(PARCEL_LAYER_ID, isEffectivelyVisible);
    }
  }, [isEffectivelyVisible, toggleLayerVisibility]);

  // Helper function to show parcel info on hover
  const showParcelInfo = (parcelId: string, props: any) => {
    if (!map) return;

    // Set hovered state
    if (parcelId && hoveredParcelId.current !== parcelId) {
      // If we hovered over another parcel before, reset its style
      if (hoveredParcelId.current) {
        map.setFeatureState(
          { source: PARCEL_LAYER_ID, id: hoveredParcelId.current },
          { hover: false }
        );
      }

      // Set hover state on this parcel
      hoveredParcelId.current = parcelId;
      map.setFeatureState({ source: PARCEL_LAYER_ID, id: parcelId }, { hover: true });

      // Increase opacity of hovered parcel
      map.setPaintProperty('parcels-fill', 'fill-opacity', [
        'case',
        ['==', ['get', 'ID'], parcelId],
        0.8, // Higher opacity for hovered parcel
        0.5, // Default opacity for non-hovered parcels
      ]);

      // Check if there's a matching comparable property
      const matchingProperty = comparableProperties?.find((prop) => prop.GeoID === parcelId);

      if (matchingProperty) {
        // Find the feature in the comparable properties layer
        const bounds = map.getBounds();
        if (!bounds) return;

        const features = map.queryRenderedFeatures(bounds.toArray(), {
          layers: ['comparable-properties', 'comparable-properties-circles'],
        });

        const matchingFeature = features.find(
          (feature) => feature.properties?.id === matchingProperty.id
        );

        if (matchingFeature && matchingFeature.geometry.type === 'Point') {
          const coordinates = matchingFeature.geometry.coordinates as [number, number];
          const props = matchingFeature.properties || {};

          // Create a detailed popup
          const content = formatPopupContent(props, {
            title: 'address',
            imageUrl: 'mainImage',
            fields: [
              {
                key: 'city',
                format: (value) => `${value || ''}, ${props.state || ''} ${props.zipCode || ''}`,
              },
              {
                key: 'price',
                format: (value) => `$${value ? Number(value).toLocaleString() : 'N/A'}`,
              },
              {
                key: 'beds',
                format: (value) =>
                  `${value || 'N/A'} Beds | ${props.baths || 'N/A'} Baths | ${props.sqft || 'N/A'} sqft`,
              },
              {
                key: 'yearBuilt',
                label: 'Built',
              },
              {
                key: 'propertyType',
              },
              {
                key: 'status',
                format: (value) => value || 'Status unavailable',
              },
            ],
          });

          // Show popup using the hook
          showPopup({
            content,
            coordinates,
            options: {
              closeButton: true,
              className: 'property-popup-detailed',
            },
          });
        }
      } else {
        // Show default parcel popup if no matching property found
        const content = formatPopupContent(props, {
          title: 'Parcel',
          fields: [
            {
              key: 'ID',
              label: 'Parcel ID',
              format: (value) => value || 'ID Unavailable',
            },
          ],
        });

        // Show popup using the hook
        showPopup({
          content,
          coordinates: [props.LONGITUDE || 0, props.LATITUDE || 0],
          options: {
            className: 'parcel-popup',
          },
        });
      }
    }
  };

  return null; // This is a logical component, doesn't render anything
}
