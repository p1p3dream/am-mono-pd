import { useEffect, useRef, useCallback, useState } from 'react';
import 'mapbox-gl/dist/mapbox-gl.css';
import { School, Building, Home, Map, BarChart, Loader2 } from 'lucide-react';
import { useMapInstance } from './hooks/useMapInstance';
import { useMapLayers } from './hooks/useMapLayers';
import { SchoolLayer, SCHOOL_LAYER_ID } from './layers/SchoolLayer';
import { BuildingLayer, BUILDING_LAYER_ID } from './layers/BuildingLayer';
import { ParcelLayer, PARCEL_LAYER_ID } from './layers/ParcelLayer';
import { MapControls } from './controls/MapControls';
import { LayerToggles } from './layers/LayerToggle';
import { ComparablePropertiesLayer, COMPARABLE_LAYER_ID } from './layers/ComparablePropertiesLayer';
import { SubdivisionLayer, SUBDIVISION_LAYER_ID } from './layers/SubdivisionLayer';
import { PolygonDrawLayer } from './layers/PolygonDrawLayer';
import { ParcelFeature, PropertyDetails, SchoolFeature, SubdivisionFeature } from '../schemas';

// Define types for our component props and data
export type BuildingProperties = {
  height?: number | string;
  floors?: number | string;
  buildingType?: string;
  yearBuilt?: number | string;
  address?: string;
  name?: string;
  bedrooms?: number | string;
  bathrooms?: number | string;
  sqft?: number | string;
  lotSize?: number | string;
  type?: string;
  hoa?: string;
  built?: number | string;
  stories?: number | string;
  status?: string;
  price?: number | string;
  value?: number | string;
  rent?: number | string;
};

export type BuildingBoundary = {
  type: string;
  geometry: GeoJSON.Geometry;
  properties: Record<string, unknown>;
};

export type BuildingFeature = {
  lng: number;
  lat: number;
  boundary: BuildingBoundary;
  properties: BuildingProperties;
};

interface MapComponentProps {
  isLoading?: boolean;
  initialCenter?: [number, number]; // [longitude, latitude]
  initialZoom?: number;
  height?: string | number;
  width?: string | number;
  showControls?: boolean;
  onBuildingSelect?: (building: BuildingFeature | null) => void;
  selectedBuilding?: BuildingFeature | null;
  className?: string;
  disableInteraction?: boolean;
  showSchools?: boolean;
  showSubdivisions?: boolean;
  showParcels?: boolean;
  showBuildingSelection?: boolean;
  parcelMinZoomLevel?: number; // Minimum zoom level at which parcels will display
  comparableProperties?: PropertyDetails[];
  property?: PropertyDetails;
  getSchools?: () => SchoolFeature[];
  getSubdivisions?: () => SubdivisionFeature[];
  getParcels?: () => ParcelFeature[];
  schoolLayerMode?: 'all' | 'districts' | 'schools';
  onPolygonSelect?: (coordinates: [number, number][]) => void;
}

export function MapComponent({
  isLoading = false,
  initialCenter = [-75.1652, 39.9526], // Default to Philadelphia
  initialZoom = 15,
  height = '550px',
  width = '100%',
  showControls = true,
  onBuildingSelect,
  selectedBuilding,
  className = '',
  disableInteraction = false,
  showSchools = false,
  showSubdivisions = false,
  showParcels = false,
  showBuildingSelection = false,
  parcelMinZoomLevel = 15,
  comparableProperties,
  property,
  getSchools,
  getSubdivisions,
  getParcels,
  schoolLayerMode = 'all',
  onPolygonSelect,
}: MapComponentProps) {
  // Get map instance from our custom hook
  const {
    mapInstanceRef,
    isMapInitialized,
    flyToLocation,
    initializeMap,
    cleanup,
    setMapContainer,
  } = useMapInstance({
    initialCenter,
    initialZoom,
    disableInteraction,
  });

  // Get layer management functions
  const { toggleLayerVisibility, layers } = useMapLayers(mapInstanceRef.current);

  // Track previous prop values to detect actual changes
  const prevPropsRef = useRef({
    initialCenter,
    initialZoom,
    disableInteraction,
    showSchools,
    showSubdivisions,
    showParcels,
  });

  // Zoom in and out handlers
  const handleZoomIn = useCallback(() => {
    if (mapInstanceRef.current) {
      mapInstanceRef.current.zoomIn();
    }
  }, []);

  const handleZoomOut = useCallback(() => {
    if (mapInstanceRef.current) {
      mapInstanceRef.current.zoomOut();
    }
  }, []);

  const handlePinClick = useCallback(() => {
    if (mapInstanceRef.current) {
      mapInstanceRef.current.flyTo({
        center: initialCenter,
      });
    }
  }, [initialCenter]);

  // Initialize the map on component mount
  useEffect(() => {
    const containerElement = document.getElementById('map-container');
    if (!containerElement) return;
    if (isLoading) return;

    // Initialize the map
    initializeMap(containerElement)
      .then(() => {
        // Store the initial props for comparison
        prevPropsRef.current = {
          initialCenter,
          initialZoom,
          disableInteraction,
          showSchools,
          showSubdivisions,
          showParcels,
        };
      })
      .catch((error) => {
        console.error('Error initializing map:', error);
      });

    // Clean up on unmount
    return () => {
      if (mapInstanceRef.current) {
        cleanup();
      }
    };
  }, [
    initializeMap,
    cleanup,
    initialCenter,
    initialZoom,
    disableInteraction,
    showSchools,
    showSubdivisions,
    showParcels,
  ]);

  const allowBuildingSelection = !(showSubdivisions || showSchools) && !disableInteraction;

  // Check for changes to initialCenter and initialZoom
  useEffect(() => {
    if (!mapInstanceRef.current) return;

    // Only fly to location if it actually changed
    if (
      prevPropsRef.current.initialCenter[0] !== initialCenter[0] ||
      prevPropsRef.current.initialCenter[1] !== initialCenter[1] ||
      prevPropsRef.current.initialZoom !== initialZoom
    ) {
      // Fly to the new location with animation
      flyToLocation(initialCenter, initialZoom);

      // Update previous values
      prevPropsRef.current.initialCenter = initialCenter;
      prevPropsRef.current.initialZoom = initialZoom;
    }
  }, [initialCenter, initialZoom, flyToLocation]);

  // Generate layer toggles from available layers
  const layerToggleData = layers.map((layer) => {
    // Customize layer toggle display based on layer ID
    if (layer.id === SCHOOL_LAYER_ID) {
      return {
        id: layer.id,
        label: 'Schools',
        visible: layer.visible,
        icon: <School size={18} />,
      };
    } else if (layer.id === BUILDING_LAYER_ID) {
      return {
        id: layer.id,
        label: 'Selected Building',
        visible: layer.visible,
        icon: <Building size={18} />,
      };
    } else if (layer.id === COMPARABLE_LAYER_ID) {
      return {
        id: layer.id,
        label: 'Comparable Properties',
        visible: layer.visible,
        icon: <Home size={18} />,
      };
    } else if (layer.id === SUBDIVISION_LAYER_ID) {
      return {
        id: layer.id,
        label: 'Subdivisions',
        visible: layer.visible,
        icon: <Map size={18} />,
      };
    } else if (layer.id === PARCEL_LAYER_ID) {
      return {
        id: layer.id,
        label: 'Parcels',
        visible: layer.visible,
        icon: <BarChart size={18} />,
      };
    }

    // Default for any other layers
    return {
      id: layer.id,
      label: layer.id,
      visible: layer.visible,
    };
  });

  const [isDrawingPolygon, setIsDrawingPolygon] = useState(false);

  const handlePolygonDraw = useCallback(() => {
    if (!mapInstanceRef.current) return;

    setIsDrawingPolygon(!isDrawingPolygon);
    if (!isDrawingPolygon) {
      // @ts-ignore - Custom method added to map
      mapInstanceRef.current.startPolygonDrawing();
    } else {
      // @ts-ignore - Custom method added to map
      mapInstanceRef.current.stopPolygonDrawing();
    }
  }, [isDrawingPolygon]);

  if (isLoading) {
    return (
      <div className="w-full h-full flex items-center justify-center">
        <Loader2 className="w-8 h-8 animate-spin" />
      </div>
    );
  }

  return (
    <div className={`map-wrapper relative ${className}`} style={{ height: height }}>
      <div
        className={`w-full ${disableInteraction ? 'cursor-default' : ''}`}
        style={{
          height: height, //`calc(${height} + 150px)`,
          minHeight: height,
          width: width,
          position: 'relative',
        }}
      >
        <div
          id="map-container"
          className="w-full h-full"
          style={{ position: 'absolute', overflow: 'visible' }}
          ref={(el) => setMapContainer(el)}
        />
      </div>

      {/* Layer toggles */}
      {isMapInitialized && layers.length > 0 && (
        <div className="absolute top-4 left-4 z-10">
          <LayerToggles layers={layerToggleData} onToggle={toggleLayerVisibility} />
        </div>
      )}

      {/* Map Controls */}
      {showControls && (
        <MapControls
          onPinClick={handlePinClick}
          onZoomIn={handleZoomIn}
          onZoomOut={handleZoomOut}
          onPolygonDraw={handlePolygonDraw}
          isDrawingPolygon={isDrawingPolygon}
        />
      )}

      {/* School Layer */}
      {getSchools && (
        <SchoolLayer
          map={mapInstanceRef.current}
          visible={showSchools}
          getSchools={getSchools}
          layerMode={schoolLayerMode}
        />
      )}

      {/* Subdivision Layer */}
      {getSubdivisions && (
        <SubdivisionLayer
          map={mapInstanceRef.current}
          visible={showSubdivisions}
          getSubdivisions={getSubdivisions}
        />
      )}

      {/* Parcel Layer */}
      {getParcels && (
        <ParcelLayer
          map={mapInstanceRef.current}
          visible={showParcels}
          minZoomLevel={parcelMinZoomLevel}
          getParcels={getParcels}
        />
      )}

      {/* Building Selection Layer */}
      {showBuildingSelection && (
        <BuildingLayer
          map={mapInstanceRef.current}
          onBuildingSelect={onBuildingSelect}
          selectedBuilding={selectedBuilding}
          allowSelection={allowBuildingSelection}
        />
      )}

      {/* Comparable Properties Layer */}
      {comparableProperties && property && (
        <ComparablePropertiesLayer
          map={mapInstanceRef.current}
          visible={true}
          allowSelection={allowBuildingSelection}
          comparableProperties={comparableProperties}
          property={property}
        />
      )}

      {/* Polygon Draw Layer */}
      <PolygonDrawLayer
        map={mapInstanceRef.current}
        visible={isDrawingPolygon}
        onPolygonComplete={onPolygonSelect}
      />
    </div>
  );
}
