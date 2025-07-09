import { useEffect, useRef, useState, useCallback } from 'react';
import { useMapLayers } from '../hooks/useMapLayers';
import {
  createPopup,
  formatPopupContent,
  loadMapImage,
  addSource,
  addLayer,
} from '../utils/mapUtils';
import { PropertyDetails } from '../../schemas';

interface ComparablePropertiesLayerProps {
  map: mapboxgl.Map | null;
  visible?: boolean;
  allowSelection?: boolean;
  comparableProperties: PropertyDetails[];
  property: PropertyDetails;
}

export const COMPARABLE_LAYER_ID = 'comparable-properties';
export const COMPARABLE_CIRCLES_LAYER_ID = 'comparable-properties-circles';
export const COMPARABLE_BACKGROUND_LAYER_ID = 'comparable-properties-background';
export const COMPARABLE_POLYGONS_LAYER_ID = 'comparable-properties-polygons';
export const CURRENT_PROPERTY = 'current-property';

// More specific name for the global initialization variable
let eventsInitialized = false;
let globalAllowSelection = true;

// Utility function to format price to K format
const formatPriceToK = (price: number): string => {
  if (!price || isNaN(price)) return '$0k';

  // Convert to thousands or millions
  if (price >= 1000000) {
    // For millions, show as $X.XM
    const millions = price / 1000000;
    // If it's a whole number, don't show decimal
    if (millions === Math.floor(millions)) {
      return `$${millions.toFixed(0)}M`;
    }
    return `$${millions.toFixed(1).replace(/\.0$/, '')}M`;
  } else {
    // For thousands, show as $Xk
    const thousands = Math.round(price / 1000);
    return `$${thousands}k`;
  }
};

export function ComparablePropertiesLayer({
  map,
  visible = true,
  allowSelection = true,
  comparableProperties,
  property,
}: ComparablePropertiesLayerProps) {
  const [isInitialized, setIsInitialized] = useState(false);
  const { registerLayer, updateLayerSource, showPopup, removePopup } = useMapLayers(map);

  // Initialize the layer
  useEffect(() => {
    // Wait for map to load
    const initializeLayer = async () => {
      if (!map) return;

      try {
        // Load the star icon
        await loadMapImage(map, '/star_icon.png', 'star-icon');

        // Add source for comparable properties
        addSource(map, COMPARABLE_LAYER_ID, {
          type: 'geojson',
          data: {
            type: 'FeatureCollection',
            features: [],
          },
        });

        // Add source for the property
        addSource(map, CURRENT_PROPERTY, {
          type: 'geojson',
          data: {
            type: 'FeatureCollection',
            features: [],
          },
        });

        // Add source for parcel polygons
        addSource(map, `${COMPARABLE_LAYER_ID}-polygons`, {
          type: 'geojson',
          data: {
            type: 'FeatureCollection',
            features: [],
          },
        });

        // Add a layer for the current property (star icon)
        addLayer(map, {
          id: CURRENT_PROPERTY,
          type: 'symbol',
          source: CURRENT_PROPERTY,
          layout: {
            'icon-image': 'star-icon',
            'icon-size': 0.8,
            'icon-allow-overlap': true,
            'icon-ignore-placement': true,
          },
          paint: {
            'icon-opacity': 1.0,
          },
        });

        // Add parcel polygons layer
        addLayer(map, {
          id: COMPARABLE_POLYGONS_LAYER_ID,
          type: 'fill',
          source: `${COMPARABLE_LAYER_ID}-polygons`,
          paint: {
            'fill-color': '#4ade80', // Green color
            'fill-opacity': 0.3,
            'fill-outline-color': '#22c55e',
            'fill-antialias': true,
          },
        });

        // Add background circles for text - add this FIRST so it's below the text
        addLayer(map, {
          id: COMPARABLE_BACKGROUND_LAYER_ID,
          type: 'circle',
          source: COMPARABLE_LAYER_ID,
          paint: {
            'circle-radius': 20,
            'circle-color': '#ffffff',
            'circle-opacity': 0.95,
            'circle-stroke-width': 1.5,
            'circle-stroke-color': '#e5e7eb',
          },
          filter: [
            'any',
            ['all', ['>=', ['zoom'], 14], ['==', ['get', 'showAsCircle'], true]],
            ['all', ['>=', ['zoom'], 10], ['==', ['get', 'showAsCircle'], false]],
          ],
        });

        // Add circle layer for lower zoom levels
        addLayer(map, {
          id: COMPARABLE_CIRCLES_LAYER_ID,
          type: 'circle',
          source: COMPARABLE_LAYER_ID,
          paint: {
            'circle-radius': {
              base: 10,
              stops: [
                [8, 1],
                [14, 20],
              ],
            },
            'circle-color': '#3b82f6', // Blue color
            'circle-opacity': 0.8,
            'circle-stroke-width': 2,
            'circle-stroke-color': '#ffffff',
          },
          filter: [
            'any',
            [
              'all',
              ['>=', ['zoom'], 6],
              ['<', ['zoom'], 14],
              ['==', ['get', 'showAsCircle'], true],
            ],
            [
              'all',
              ['>=', ['zoom'], 6],
              ['<', ['zoom'], 10],
              ['==', ['get', 'showAsCircle'], false],
            ],
          ],
        });

        // Add text layer for property prices LAST so it's on top
        addLayer(map, {
          id: COMPARABLE_LAYER_ID,
          type: 'symbol',
          source: COMPARABLE_LAYER_ID,
          layout: {
            'text-field': ['get', 'formattedPrice'],
            'text-size': 13,
            'text-font': ['Open Sans Bold', 'Arial Unicode MS Bold'],
            'text-offset': [0, 0],
            'text-anchor': 'center',
            'text-allow-overlap': true,
            'text-ignore-placement': true,
            'text-letter-spacing': 0.02,
            'symbol-z-order': 'source',
            'symbol-sort-key': ['-', 0, ['get', 'price']], // Higher priced properties appear on top
          },
          paint: {
            'text-color': '#3b82f6', // Blue text
            'text-opacity': 1,
          },
          filter: [
            'any',
            ['all', ['>=', ['zoom'], 14], ['==', ['get', 'showAsCircle'], true]],
            ['all', ['>=', ['zoom'], 10], ['==', ['get', 'showAsCircle'], false]],
          ],
        });

        setupEventHandlers();
        setIsInitialized(true);
      } catch (error) {
        console.error('Error initializing comparable properties layer:', error);
      }
    };

    // Register the layers
    registerLayer({
      id: COMPARABLE_LAYER_ID,
      visible,
      sourceId: COMPARABLE_LAYER_ID,
      layerIds: [
        COMPARABLE_POLYGONS_LAYER_ID, // Add polygons layer first (bottom)
        COMPARABLE_LAYER_ID,
        COMPARABLE_CIRCLES_LAYER_ID,
        COMPARABLE_BACKGROUND_LAYER_ID,
      ],
    });

    registerLayer({
      id: `${COMPARABLE_LAYER_ID}-polygons`,
      visible,
      sourceId: `${COMPARABLE_LAYER_ID}-polygons`,
      layerIds: [COMPARABLE_POLYGONS_LAYER_ID],
    });

    registerLayer({
      id: CURRENT_PROPERTY,
      visible,
      sourceId: CURRENT_PROPERTY,
      layerIds: [CURRENT_PROPERTY],
    });

    // Set up event handlers
    if (!map) return;

    initializeLayer();
  }, [map, registerLayer, visible]);

  // Set up event handlers for property interactions
  const setupEventHandlers = useCallback(() => {
    if (!map || eventsInitialized) return;

    const layers = [
      CURRENT_PROPERTY,
      COMPARABLE_LAYER_ID,
      COMPARABLE_CIRCLES_LAYER_ID,
      COMPARABLE_POLYGONS_LAYER_ID,
    ];

    // Mouse enter handler
    layers.forEach((layerId) => {
      let popupTimeout: NodeJS.Timeout | null = null;

      map.on('mouseenter', layerId, () => {
        if (!globalAllowSelection) return;

        map.getCanvas().style.cursor = 'pointer';

        // Highlight the polygon on hover
        if (layerId === COMPARABLE_POLYGONS_LAYER_ID) {
          map.setPaintProperty(COMPARABLE_POLYGONS_LAYER_ID, 'fill-opacity', 0.5);
        }

        // Clear any pending timeout when mouse enters
        if (popupTimeout) {
          clearTimeout(popupTimeout);
          popupTimeout = null;
        }
      });

      map.on('mouseleave', layerId, () => {
        map.getCanvas().style.cursor = '';

        // Reset polygon opacity on mouse leave
        if (layerId === COMPARABLE_POLYGONS_LAYER_ID) {
          map.setPaintProperty(COMPARABLE_POLYGONS_LAYER_ID, 'fill-opacity', 0.3);
        }
      });

      // Click event for detailed popup
      map.on('click', layerId, (e) => {
        if (!globalAllowSelection) return;

        e.preventDefault();
        e.originalEvent.stopPropagation();

        if (e.features && e.features.length > 0 && e.features[0]) {
          const feature = e.features[0];
          const props = feature.properties || {};

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

          // Get the coordinates from the feature
          let coordinates: [number, number] = [0, 0];

          if (feature.geometry.type === 'Point') {
            // For point features, use the point coordinates
            coordinates = (feature.geometry as GeoJSON.Point).coordinates as [number, number];
          } else if (feature.geometry.type === 'Polygon') {
            // For polygon features, use the click coordinates
            coordinates = e.lngLat.toArray() as [number, number];
          }

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
      });
    });

    map.on('click', (e) => {
      if (!globalAllowSelection) return;

      const features = map.queryRenderedFeatures(e.point, {
        layers: layers,
      });

      // if nothing returned, remove the popup
      if (!features || features.length === 0) {
        removePopup();
      }
    });
    eventsInitialized = true;
  }, [map, showPopup, removePopup]);

  // Update properties when the list changes
  useEffect(() => {
    if (!map || !isInitialized || !comparableProperties.length) return;

    // Convert to GeoJSON features
    const pointFeatures = comparableProperties
      .filter((cp) => cp.id !== property.id)
      .map((property, index) => {
        if (!property.location) return null;

        // Use location coordinates
        let coordinates: [number, number] = [property.location.lng, property.location.lat];

        // Format price as "600k" format
        const formattedPrice = formatPriceToK(property.price);

        // Determine if this property should be shown as a circle at lower zoom levels
        const showAsCircle = index % 3 !== 0;

        return {
          type: 'Feature',
          id: property.id,
          geometry: {
            type: 'Point',
            coordinates: coordinates,
          },
          properties: {
            id: property.id,
            address: property.address,
            city: property.city,
            state: property.state,
            zipCode: property.zipCode,
            status: property.status,
            price: property.price,
            formattedPrice: formattedPrice,
            beds: property.beds,
            baths: property.baths,
            sqft: property.sqft,
            yearBuilt: property.yearBuilt,
            propertyType: property.propertyType,
            mainImage: property.mainImage,
            showAsCircle: showAsCircle,
          },
        };
      })
      .filter(Boolean) as GeoJSON.Feature[];

    // Create polygon features
    const polygonFeatures = comparableProperties
      .map((property) => {
        if (!property.parcelGeometry) return null;

        // Format price as "600k" format
        const formattedPrice = formatPriceToK(property.price);

        return {
          type: 'Feature',
          id: property.id,
          geometry: property.parcelGeometry as GeoJSON.Polygon,
          properties: {
            id: property.id,
            address: property.address,
            city: property.city,
            state: property.state,
            zipCode: property.zipCode,
            status: property.status,
            price: property.price,
            formattedPrice: formattedPrice,
            beds: property.beds,
            baths: property.baths,
            sqft: property.sqft,
            yearBuilt: property.yearBuilt,
            propertyType: property.propertyType,
            mainImage: property.mainImage,
            showAsCircle: false,
          },
        };
      })
      .filter(Boolean) as GeoJSON.Feature[];

    // Update both sources
    updateLayerSource(COMPARABLE_LAYER_ID, {
      type: 'FeatureCollection',
      features: pointFeatures,
    });

    updateLayerSource(`${COMPARABLE_LAYER_ID}-polygons`, {
      type: 'FeatureCollection',
      features: polygonFeatures,
    });
  }, [comparableProperties, map, property, updateLayerSource, isInitialized]);

  useEffect(() => {
    if (!map || !isInitialized) return;

    // Update the source data
    updateLayerSource(CURRENT_PROPERTY, {
      type: 'FeatureCollection',
      features: [
        {
          type: 'Feature',
          id: property.id,
          geometry: {
            type: 'Point',
            coordinates: [property.location.lng, property.location.lat],
          },
          properties: {
            id: property.id,
            address: property.address,
            city: property.city,
            state: property.state,
            zipCode: property.zipCode,
            status: property.status,
            price: property.price,
            beds: property.beds,
            baths: property.baths,
            sqft: property.sqft,
            yearBuilt: property.yearBuilt,
            propertyType: property.propertyType,
            mainImage: property.mainImage,
          },
        },
      ],
    });
  }, [property, map, updateLayerSource, isInitialized]);

  useEffect(() => {
    globalAllowSelection = allowSelection;
    // if disabled, remove the popup
    if (!allowSelection) {
      removePopup();
    }
  }, [allowSelection, removePopup]);

  return null; // This is a logical component, doesn't render anything
}
