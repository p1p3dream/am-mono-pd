import { useRef, useState, useEffect, useCallback } from 'react';
import mapboxgl, { EasingOptions } from 'mapbox-gl';
import MapSingleton from '../map-singleton';

interface UseMapInstanceOptions {
  initialCenter?: [number, number];
  initialZoom?: number;
  disableInteraction?: boolean;
}

export function useMapInstance({
  initialCenter = [-75.1652, 39.9526], // Default to Philadelphia
  initialZoom = 15,
  disableInteraction = false,
}: UseMapInstanceOptions = {}) {
  // Reference to the map container element
  const mapContainer = useRef<HTMLElement | null>(null);

  // Store map instance in a ref to maintain access across renders
  const mapInstanceRef = useRef<mapboxgl.Map | null>(null);

  // Track if the map is initialized
  const [isMapInitialized, setIsMapInitialized] = useState(false);

  // Track map coordinates and zoom
  const [lng, setLng] = useState(initialCenter[0]);
  const [lat, setLat] = useState(initialCenter[1]);
  const [zoom, setZoom] = useState(initialZoom);

  // Track if changes are from user interaction vs prop changes
  const userInteractionRef = useRef(false);

  // Method to fly to a location
  const flyToLocation = useCallback(
    (
      center: [number, number],
      zoom: number,
      options?: { duration?: number; animate?: boolean }
    ) => {
      const map = mapInstanceRef.current;
      if (!map) return;

      // Default options
      const flyOptions: EasingOptions = {
        center,
        zoom,
        essential: true,
        duration: options?.duration || 1500, // Default animation duration
        animate: options?.animate !== false, // Default to animate
      };

      // Use flyTo for smooth animation
      map.flyTo(flyOptions);
    },
    []
  );

  // Setup map styles to hide logos and set cursor styles
  const setupMapStyles = useCallback((disableInteractions: boolean) => {
    // Remove any existing style element
    const existingStyle = document.getElementById('mapbox-custom-styles');
    if (existingStyle) {
      existingStyle.remove();
    }

    // Add a style to hide the logo and set cursor styles
    const style = document.createElement('style');
    style.id = 'mapbox-custom-styles';
    style.innerHTML = `
      .mapboxgl-ctrl-logo {
        display: none !important;
      }
      .mapboxgl-ctrl-attrib {
        display: none !important;
      }
      ${
        disableInteractions
          ? `
      .mapboxgl-canvas-container,
      .mapboxgl-canvas,
      .mapboxgl-map,
      .mapboxgl-canvas-container.mapboxgl-interactive,
      .mapboxgl-canvas-container.mapboxgl-interactive canvas.mapboxgl-canvas {
        cursor: default !important;
      }
      .mapboxgl-canvas-container.mapboxgl-interactive:active,
      .mapboxgl-canvas-container.mapboxgl-interactive:hover {
        cursor: default !important;
      }`
          : ''
      }
    `;
    document.head.appendChild(style);
  }, []);

  // Initialize the map
  const initializeMap = useCallback(
    (container: HTMLElement) => {
      if (!container) return Promise.reject('No container provided');

      // Configure map options with proper typing
      const mapOptions = {
        style: 'mapbox://styles/mapbox/dark-v11',
        center: initialCenter,
        zoom: initialZoom,
        pitch: 0, // Top-down view
        bearing: 0,
        preserveDrawingBuffer: true,
        antialias: true,
        attributionControl: false,
        renderWorldCopies: true,
        boxZoom: false,
        logoPosition: 'top-left' as 'top-left',
        dragPan: !disableInteraction,
        scrollZoom: !disableInteraction,
        doubleClickZoom: !disableInteraction,
        touchZoomRotate: !disableInteraction,
      };

      return MapSingleton.init(container, mapOptions).then((map) => {
        mapInstanceRef.current = map;
        setupMapStyles(disableInteraction);

        // Update local state to match initialCenter and initialZoom
        setLng(initialCenter[0]);
        setLat(initialCenter[1]);
        setZoom(initialZoom);

        // Handle map move events
        map.on('moveend', () => {
          const center = map.getCenter();
          setLng(center.lng);
          setLat(center.lat);
          setZoom(map.getZoom());

          // Reset user interaction flag after a short delay
          setTimeout(() => {
            userInteractionRef.current = false;
          }, 100);
        });

        map.on('dragstart', () => {
          userInteractionRef.current = true;
        });

        map.on('zoomstart', () => {
          userInteractionRef.current = true;
        });

        setIsMapInitialized(true);
        return map;
      });
    },
    [initialCenter, initialZoom, disableInteraction, setupMapStyles]
  );

  // Cleanup function for map
  const cleanup = useCallback(() => {
    mapContainer.current = null;
  }, []);

  return {
    mapContainer,
    mapInstanceRef,
    isMapInitialized,
    lng,
    lat,
    zoom,
    userInteractionRef,
    flyToLocation,
    setupMapStyles,
    initializeMap,
    cleanup,
    setMapContainer: (element: HTMLElement | null) => {
      mapContainer.current = element;
    },
  };
}
