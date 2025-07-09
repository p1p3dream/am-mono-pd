import mapboxgl from 'mapbox-gl';
import { MAPBOX_ACCESS_TOKEN } from '@/lib/env';

// Configure Mapbox access token
mapboxgl.accessToken = MAPBOX_ACCESS_TOKEN;

type MapOptions = {
  style?: string;
  center?: [number, number];
  zoom?: number;
  pitch?: number;
  bearing?: number;
  preserveDrawingBuffer?: boolean;
  antialias?: boolean;
  attributionControl?: boolean;
  renderWorldCopies?: boolean;
  boxZoom?: boolean;
  logoPosition?: 'top-left' | 'top-right' | 'bottom-left' | 'bottom-right';
  dragPan?: boolean;
  scrollZoom?: boolean;
  doubleClickZoom?: boolean;
  touchZoomRotate?: boolean;
};

type MapControls = {
  dragPan?: boolean;
  scrollZoom?: boolean;
  doubleClickZoom?: boolean;
  touchZoomRotate?: boolean;
};

/**
 * A simplified singleton for managing a Mapbox GL map instance.
 * The map instance persists across component mounts/unmounts.
 */
class MapSingleton {
  private static instance: mapboxgl.Map | null = null;
  private static container: HTMLElement | null = null;
  private static previousOptions: MapOptions | null = null;

  /**
   * Initialize the map instance or reuse existing instance
   */
  public static init(container: HTMLElement, options: MapOptions = {}): Promise<mapboxgl.Map> {
    // If we already have an instance
    if (this.instance) {
      try {
        // Move the map to the new container if needed
        if (this.container !== container) {
          const mapDiv = this.instance.getContainer();
          // Convert NodeList to array and iterate
          container.appendChild(mapDiv);
          this.container = container;
          this.instance.resize();
        }

        // Update map options if they have changed
        if (
          options.zoom !== this.previousOptions?.zoom ||
          options.center?.[0] !== this.previousOptions?.center?.[0] ||
          options.center?.[1] !== this.previousOptions?.center?.[1]
        ) {
          this.previousOptions = options;
          this.instance.flyTo({
            center: options.center,
            zoom: options.zoom,
            essential: true,
          });
        }

        // Update interaction controls
        if (
          options.dragPan !== undefined ||
          options.scrollZoom !== undefined ||
          options.doubleClickZoom !== undefined ||
          options.touchZoomRotate !== undefined
        ) {
          this.updateControls({
            dragPan: options.dragPan,
            scrollZoom: options.scrollZoom,
            doubleClickZoom: options.doubleClickZoom,
            touchZoomRotate: options.touchZoomRotate,
          });
        }

        return Promise.resolve(this.instance);
      } catch (error) {
        console.error('Error reusing map instance:', error);
        // If something goes wrong, remove the instance and create a new one
        if (this.instance) {
          this.instance.remove();
          this.instance = null;
        }
      }
    }

    // Create a new map instance
    this.container = container;

    return new Promise((resolve, reject) => {
      try {
        // Create map instance with options
        const mapOptions: mapboxgl.MapboxOptions = {
          container,
          style: options.style || 'mapbox://styles/mapbox/dark-v11',
          center: options.center || [-86.808815, 36.182115],
          zoom: options.zoom || 15,
          pitch: options.pitch || 0,
          bearing: options.bearing || 0,
          preserveDrawingBuffer:
            options.preserveDrawingBuffer !== undefined ? options.preserveDrawingBuffer : true,
          antialias: options.antialias !== undefined ? options.antialias : true,
          attributionControl:
            options.attributionControl !== undefined ? options.attributionControl : false,
          renderWorldCopies:
            options.renderWorldCopies !== undefined ? options.renderWorldCopies : true,
          boxZoom: options.boxZoom !== undefined ? options.boxZoom : false,
          logoPosition: options.logoPosition || 'top-left',
          dragPan: options.dragPan,
          scrollZoom: options.scrollZoom,
          doubleClickZoom: options.doubleClickZoom,
          touchZoomRotate: options.touchZoomRotate,
        };

        // Create new map
        const mapInstance = new mapboxgl.Map(mapOptions);
        console.log('Creating new map instance');

        // Set up map load handler
        mapInstance.on('load', () => {
          this.instance = mapInstance;
          resolve(mapInstance);
        });

        // Handle errors
        mapInstance.on('error', (error) => {
          console.error('Map error:', error);
          reject(error);
        });
      } catch (error) {
        console.error('Error creating map:', error);
        reject(error);
      }
    });
  }

  /**
   * Updates map interaction controls
   */
  public static updateControls(controls: MapControls): void {
    if (!this.instance) return;

    const map = this.instance;

    if (controls.dragPan !== undefined) {
      controls.dragPan ? map.dragPan.enable() : map.dragPan.disable();
    }

    if (controls.scrollZoom !== undefined) {
      controls.scrollZoom ? map.scrollZoom.enable() : map.scrollZoom.disable();
    }

    if (controls.doubleClickZoom !== undefined) {
      controls.doubleClickZoom ? map.doubleClickZoom.enable() : map.doubleClickZoom.disable();
    }

    if (controls.touchZoomRotate !== undefined) {
      controls.touchZoomRotate ? map.touchZoomRotate.enable() : map.touchZoomRotate.disable();
    }
  }

  /**
   * Get the current map instance
   */
  public static getInstance(): mapboxgl.Map | null {
    return this.instance;
  }

  /**
   * Check if map is available
   */
  public static hasInstance(): boolean {
    return this.instance !== null;
  }

  /**
   * Clean up the map instance (only use when application is shutting down)
   */
  public static cleanup(): void {
    if (this.instance) {
      const map = this.instance;
      this.instance = null;
      this.container = null;
      this.previousOptions = null;
      map.remove();
    }
  }
}

export default MapSingleton;
