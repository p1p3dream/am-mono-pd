import { useEffect, useRef, useState, useMemo } from 'react';
import { createPopup, formatPopupContent, addSource, addLayer } from '../utils/mapUtils';
import { useMapLayers, MapLayer } from '../hooks/useMapLayers';
import { SchoolFeature } from '../../schemas';
import { Switch } from '../../components/ui/switch';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '../../components/ui/select';

interface SchoolLayerProps {
  map: mapboxgl.Map | null;
  visible?: boolean;
  onSchoolClick?: (school: SchoolFeature) => void;
  getSchools: () => SchoolFeature[];
  layerMode?: 'all' | 'districts' | 'schools';
}

export const SCHOOL_LAYER_ID = 'schools';
export const DISTRICT_LAYER_ID = 'districts';

// Global variable for tracking initialization
let isSchoolLayerInitialized = false;

export function SchoolLayer({
  map,
  visible = false,
  onSchoolClick,
  getSchools,
  layerMode = 'all',
}: SchoolLayerProps) {
  const [schools, setSchools] = useState<SchoolFeature[]>([]);
  const popupRef = useRef<mapboxgl.Popup | null>(null);
  const hoveredDistrictId = useRef<string | null>(null);
  const hoveredSchoolId = useRef<string | null>(null);
  const clickedRef = useRef<string | null>(null);
  const { registerLayer, toggleLayerVisibility } = useMapLayers(map);
  const [showDistricts, setShowDistricts] = useState(true);
  const [showSchools, setShowSchools] = useState(true);
  const [selectedLevel, setSelectedLevel] = useState<string>('all');
  const [modalPosition, setModalPosition] = useState({ x: 20, y: 0 });
  const [isDragging, setIsDragging] = useState(false);
  const [dragStart, setDragStart] = useState({ x: 0, y: 0 });

  // Fetch schools data
  useEffect(() => {
    const schoolData = getSchools();
    setSchools(schoolData);
  }, []);

  // Initialize school layers
  useEffect(() => {
    if (!map || !schools.length) return;

    const schoolGeoJSON: GeoJSON.FeatureCollection = {
      type: 'FeatureCollection',
      features: schools.map((school) => {
        return {
          type: 'Feature',
          properties: school.properties,
          geometry: school.geometry,
        } as GeoJSON.Feature;
      }),
    };

    // Create a separate GeoJSON for school points
    const schoolPointsGeoJSON: GeoJSON.FeatureCollection = {
      type: 'FeatureCollection',
      features: schools.flatMap((school) => {
        const schoolsList =
          typeof school.properties.schools === 'string'
            ? JSON.parse(school.properties.schools)
            : school.properties.schools || [];

        // Filter schools by level if selectedLevel is not 'all'
        const filteredSchoolsList =
          selectedLevel === 'all'
            ? schoolsList
            : schoolsList.filter((s: any) => s.level === selectedLevel);

        return filteredSchoolsList.map((s: any) => ({
          type: 'Feature',
          properties: {
            id: s.id,
            name: s.name,
            rating: s.rating,
            address: s.address,
            website: s.website,
            level: s.level,
            grade_low: s.grade_low,
            grade_high: s.grade_high,
            grades: s.grades,
            districtId: school.properties.ID,
            districtName: school.properties.name,
          },
          geometry: {
            type: 'Point',
            coordinates: [s.longitude, s.latitude],
          },
        }));
      }),
    };

    // Add source for school districts
    addSource(map, DISTRICT_LAYER_ID, {
      type: 'geojson',
      data: schoolGeoJSON,
    });

    // Add source for school points
    addSource(map, SCHOOL_LAYER_ID, {
      type: 'geojson',
      data: schoolPointsGeoJSON,
    });

    // Check if sources exist and update data if needed
    if (isSchoolLayerInitialized) {
      const districtSource = map.getSource(DISTRICT_LAYER_ID) as mapboxgl.GeoJSONSource;
      if (districtSource) {
        districtSource.setData(schoolGeoJSON);
      }
      const schoolSource = map.getSource(SCHOOL_LAYER_ID) as mapboxgl.GeoJSONSource;
      if (schoolSource) {
        schoolSource.setData(schoolPointsGeoJSON);
      }
    }

    // Add fill layer for school boundaries
    addLayer(map, {
      id: 'districts-fill',
      type: 'fill',
      source: DISTRICT_LAYER_ID,
      paint: {
        'fill-color': ['get', 'color'],
        'fill-opacity': 0.5,
        'fill-outline-color': '#ffffff',
      },
    });

    // Add outline layer for school boundaries
    addLayer(map, {
      id: 'districts-outline',
      type: 'line',
      source: DISTRICT_LAYER_ID,
      paint: {
        'line-color': '#ffffff',
        'line-width': 1,
      },
    });

    // Add school points layer
    addLayer(map, {
      id: 'schools-points',
      type: 'circle',
      source: SCHOOL_LAYER_ID,
      paint: {
        'circle-radius': 6,
        'circle-color': '#FFD700',
        'circle-stroke-width': 2,
        'circle-stroke-color': '#ffffff',
        'circle-opacity': 0.8,
      },
    });

    // Register the layers with the layer manager
    registerLayer({
      id: DISTRICT_LAYER_ID,
      visible,
      sourceId: DISTRICT_LAYER_ID,
      layerIds: ['districts-fill', 'districts-outline'],
    });

    registerLayer({
      id: SCHOOL_LAYER_ID,
      visible,
      sourceId: SCHOOL_LAYER_ID,
      layerIds: ['schools-points'],
    });

    // Set initial visibility
    toggleLayerVisibility(DISTRICT_LAYER_ID, visible);
    toggleLayerVisibility(SCHOOL_LAYER_ID, visible);

    if (isSchoolLayerInitialized) return;

    // Add hover effects for both boundaries and points
    const hoverLayers = ['districts-fill', 'schools-points'];
    hoverLayers.forEach((layerId) => {
      // Add hover effects
      map.on('mouseenter', layerId, (e) => {
        if (e.features && e.features.length > 0 && e.features[0]) {
          // Change cursor to pointer
          map.getCanvas().style.cursor = 'pointer';

          if (clickedRef.current) {
            return;
          }

          const feature = e.features[0];
          const props = feature.properties || {};

          if (layerId === 'schools-points') {
            showInfo(props.id, props, 'school');
          } else if (!hoveredSchoolId.current) {
            showInfo(props.ID, props, 'district');
          }
        }
      });

      // Add mousemove handler to update popup position
      map.on('mousemove', layerId, (e) => {
        if (clickedRef.current) {
          return;
        }

        if (e.features && e.features.length > 0 && e.features[0]) {
          const feature = e.features[0];
          const props = feature.properties || {};

          if (layerId === 'schools-points') {
            showInfo(props.id, props, 'school');
          } else if (!hoveredSchoolId.current) {
            showInfo(props.ID, props, 'district');
          }

          // Position popup at mouse position
          if (e.lngLat && popupRef.current) {
            popupRef.current.setLngLat(e.lngLat).addTo(map);
          }
        }
      });

      // Mouse leave event
      map.on('mouseleave', layerId, () => {
        // Reset cursor
        map.getCanvas().style.cursor = '';

        // Reset hover states
        if (layerId === 'schools-points') {
          hoveredSchoolId.current = null;
        } else {
          hoveredDistrictId.current = null;
          // Reset fill opacity
          map.setPaintProperty('districts-fill', 'fill-opacity', 0.5);
        }

        if (clickedRef.current) {
          return;
        }

        // Remove popup
        if (popupRef.current) {
          popupRef.current.remove();
          popupRef.current = null;
        }
      });

      // Click event for schools
      map.on('click', layerId, (e) => {
        if (e.features && e.features.length > 0 && e.features[0]) {
          const feature = e.features[0];
          const props = feature.properties || {};

          if (layerId === 'schools-points') {
            // For school points, freeze the popup

            showInfo(props.id, props, 'school');
            popupRef.current?.setLngLat(e.lngLat).addTo(map);
            clickedRef.current = props.id;
          }

          // Find and call onSchoolClick for school features
          if (onSchoolClick) {
            const schoolFeature = schools.find((s) => s.properties.ID === feature.properties?.ID);
            if (schoolFeature) {
              onSchoolClick(schoolFeature);
            }
          }
        }
      });
    });

    // Add click handler to map to clear clicked school when clicking outside
    map.on('click', (e) => {
      // Check if click was on a school layer
      const features = map.queryRenderedFeatures(e.point, {
        layers: ['schools-points'],
      });

      if (features.length === 0) {
        // Click was outside school layers, clear clicked school
        if (popupRef.current && clickedRef.current) {
          popupRef.current.remove();
          popupRef.current = null;
        }
        clickedRef.current = null;
      }
    });

    // Use global variable instead of state
    isSchoolLayerInitialized = true;
  }, [map, schools, registerLayer, toggleLayerVisibility, visible, onSchoolClick, selectedLevel]);

  // Update visibility when visibility prop or layerMode changes
  useEffect(() => {
    if (isSchoolLayerInitialized) {
      if (visible) {
        // Show layers based on layerMode
        toggleLayerVisibility(DISTRICT_LAYER_ID, layerMode === 'all' || layerMode === 'districts');
        toggleLayerVisibility(SCHOOL_LAYER_ID, layerMode === 'all' || layerMode === 'schools');
      } else {
        // Hide all layers when not visible
        toggleLayerVisibility(DISTRICT_LAYER_ID, false);
        toggleLayerVisibility(SCHOOL_LAYER_ID, false);
      }
    }
  }, [visible, toggleLayerVisibility, layerMode]);

  // Helper function to show info on hover/click
  const showInfo = (id: string, props: any, type: 'school' | 'district') => {
    if (!map) return;

    if (type === 'school' && id !== hoveredSchoolId.current) {
      hoveredSchoolId.current = id;
      hoveredDistrictId.current = null;
      if (popupRef.current) {
        popupRef.current.remove();
        popupRef.current = null;
      }
    } else if (!hoveredSchoolId.current && id && hoveredDistrictId.current !== id) {
      // If we hovered over another district before, reset its style
      if (hoveredDistrictId.current) {
        map.setFeatureState(
          { source: DISTRICT_LAYER_ID, id: hoveredDistrictId.current },
          { hover: false }
        );
      }

      // Set hover state on this district
      hoveredDistrictId.current = id;
      map.setFeatureState({ source: DISTRICT_LAYER_ID, id }, { hover: true });

      // Increase opacity of hovered district
      map.setPaintProperty('districts-fill', 'fill-opacity', [
        'case',
        ['==', ['get', 'ID'], id],
        0.8, // Higher opacity for hovered district
        0.5, // Default opacity for non-hovered districts
      ]);

      // Remove existing popup
      if (popupRef.current) {
        popupRef.current.remove();
        popupRef.current = null;
      }
    }

    // Create popup if it doesn't exist
    if (!popupRef.current) {
      let content = '';

      if (type === 'school') {
        content = `
          <div style="font-family: sans-serif; padding: 8px; max-width: 300px; color: #fff;">
            <div style="font-weight: bold; margin-bottom: 4px; font-size: 1.1em;">
              <a href="${props.website}" target="_blank" style="color: #fff; text-decoration: none;">${props.name || 'Unnamed School'}</a>
            </div>
            <div style="display: flex; justify-content: space-between; align-items: center; margin: 4px 0;">
              <span>Rating:</span>
              <span style="color: ${props.rating ? getRatingColor(props.rating) : '#FFD700'}; text-shadow: -1px -1px 0 #000, 1px -1px 0 #000, -1px 1px 0 #000, 1px 1px 0 #000; font-weight: bold; font-size: 1.4em;">${props.rating || 'N/A'}</span>
            </div>
            <div style="margin: 4px 0;">
              <span style="color: #ccc;">Address:</span>
              <span>${props.address || 'Address not available'}</span>
            </div>
            <div style="margin: 4px 0;">
              <span style="color: #ccc;">Level:</span>
              <span>${props.level || 'Type not specified'}</span>
            </div>
            <div style="margin: 4px 0;">
              <span style="color: #ccc;">Grades:</span>
              <span>${props.grade_low || 'Grades not specified'} - ${props.grade_high || 'Grades not specified'}</span>
            </div>
            <div style="margin-top: 8px; font-size: 0.9em; color: #ccc; border-top: 1px solid rgba(255,255,255,0.2); padding-top: 8px;">
              ${props.districtName || 'School District'}
            </div>
          </div>
        `;
      } else {
        const schools =
          typeof props.schools === 'string' ? JSON.parse(props.schools) : props.schools || [];
        const ratings = schools.map((s: any) => s.rating).filter((r: string) => r);
        const averageRating =
          ratings.length > 0
            ? ratings.reduce((acc: string, curr: string) => {
                const ratingValues: Record<string, number> = {
                  'A+': 12,
                  A: 11,
                  'A-': 10,
                  'B+': 9,
                  B: 8,
                  'B-': 7,
                  'C+': 6,
                  C: 5,
                  'C-': 4,
                  'D+': 3,
                  D: 2,
                  'D-': 1,
                  F: 0,
                };
                return acc + ratingValues[curr] || 0;
              }, 0) / ratings.length
            : null;

        const getLetterGrade = (avg: number): string => {
          if (avg >= 11.5) return 'A+';
          if (avg >= 11) return 'A';
          if (avg >= 10.5) return 'A-';
          if (avg >= 9.5) return 'B+';
          if (avg >= 9) return 'B';
          if (avg >= 8.5) return 'B-';
          if (avg >= 7.5) return 'C+';
          if (avg >= 7) return 'C';
          if (avg >= 6.5) return 'C-';
          if (avg >= 5.5) return 'D+';
          if (avg >= 5) return 'D';
          if (avg >= 4.5) return 'D-';
          return 'F';
        };

        const districtRating = averageRating !== null ? getLetterGrade(averageRating) : null;

        content = `
          <div style="font-family: sans-serif; padding: 8px; max-width: 250px; color: #fff;">
            <div style="font-weight: bold; margin-bottom: 4px;">${props.name || 'School District'}</div>
            <div style="display: flex; justify-content: space-between; align-items: center;">
              <span>Total Schools:</span>
              <span>${schools.length}</span>
            </div>
            <div style="display: flex; justify-content: space-between; align-items: center; margin-top: 4px;">
              <span>Average Rating:</span>
              <span style="color: ${districtRating ? getRatingColor(districtRating) : '#FFD700'}; text-shadow: -1px -1px 0 #000, 1px -1px 0 #000, -1px 1px 0 #000, 1px 1px 0 #000; font-weight: bold; font-size: 1.4em;">${districtRating || 'N/A'}</span>
            </div>
          </div>
        `;
      }

      popupRef.current = createPopup(content, {
        className: 'school-popup',
      });
    }
  };

  // Helper function to get color based on rating
  const getRatingColor = (rating: string): string => {
    const ratingColors: Record<string, string> = {
      'A+': '#00FF00', // Bright Green
      A: '#32CD32', // Lime Green
      'A-': '#90EE90', // Light Green
      'B+': '#98FB98', // Pale Green
      B: '#FFD700', // Yellow
      'B-': '#FFA500', // Orange
      'C+': '#FF8C00', // Dark Orange
      C: '#FF6B6B', // Light Red
      'C-': '#FF4500', // Orange Red
      'D+': '#FF0000', // Red
      D: '#DC143C', // Crimson
      'D-': '#8B0000', // Dark Red
      F: '#800000', // Maroon
    };
    return ratingColors[rating] || '#FFD700'; // Default to yellow for unknown ratings
  };

  // Update modal position when map container is available
  useEffect(() => {
    if (!map) return;

    const mapContainer = map.getContainer();
    if (!mapContainer) return;

    const updatePosition = () => {
      const rect = mapContainer.getBoundingClientRect();
      setModalPosition({
        x: 20,
        y: rect.height - 200, // 200px from bottom
      });
    };

    // Initial position
    updatePosition();

    // Update on map resize
    map.on('resize', updatePosition);
    return () => {
      map.off('resize', updatePosition);
    };
  }, [map]);

  // Handle drag events
  const handleMouseDown = (e: React.MouseEvent) => {
    if (!map) return;
    const mapContainer = map.getContainer();
    if (!mapContainer) return;

    setIsDragging(true);
    const rect = mapContainer.getBoundingClientRect();
    setDragStart({
      x: e.clientX - modalPosition.x,
      y: e.clientY - modalPosition.y,
    });
  };

  const handleMouseMove = (e: MouseEvent) => {
    if (!isDragging || !map) return;

    const mapContainer = map.getContainer();
    if (!mapContainer) return;

    const rect = mapContainer.getBoundingClientRect();
    const newX = e.clientX - dragStart.x;
    const newY = e.clientY - dragStart.y;

    // Constrain to map bounds
    setModalPosition({
      x: Math.max(0, Math.min(newX, rect.width - 200)), // 200px is modal width
      y: Math.max(0, Math.min(newY, rect.height - 100)), // 200px is modal height
    });
  };

  const handleMouseUp = () => {
    setIsDragging(false);
  };

  useEffect(() => {
    if (isDragging) {
      window.addEventListener('mousemove', handleMouseMove);
      window.addEventListener('mouseup', handleMouseUp);
    }
    return () => {
      window.removeEventListener('mousemove', handleMouseMove);
      window.removeEventListener('mouseup', handleMouseUp);
    };
  }, [isDragging]);

  // Get unique school levels
  const schoolLevels = useMemo(() => {
    const levels = new Set<string>();
    schools.forEach((school) => {
      const schoolsList =
        typeof school.properties.schools === 'string'
          ? JSON.parse(school.properties.schools)
          : school.properties.schools || [];
      schoolsList.forEach((s: any) => {
        if (s.level) levels.add(s.level);
      });
    });
    return Array.from(levels);
  }, [schools]);

  // Add effect to reinitialize layers when selectedLevel changes
  useEffect(() => {
    if (!map) return;
    isSchoolLayerInitialized = false; // Reset initialization flag
    // Reinitialize layers with new filter
    const schoolData = getSchools();
    setSchools(schoolData);
  }, [map, selectedLevel, getSchools]);

  return visible ? (
    <div
      className="absolute z-50 bg-background/90 rounded-lg shadow-lg p-4 min-w-[200px] border border-gray-500 rounded-md p-4"
      style={{
        left: modalPosition.x,
        top: modalPosition.y,
        cursor: isDragging ? 'grabbing' : 'grab',
        maxHeight: '200px',
        overflowY: 'auto',
      }}
      onMouseDown={handleMouseDown}
    >
      <div className="space-y-4">
        <div className="flex items-center justify-between">
          <span className="text-sm font-medium">Districts</span>
          <Switch
            checked={showDistricts}
            onCheckedChange={(checked) => {
              setShowDistricts(checked);
              toggleLayerVisibility(DISTRICT_LAYER_ID, checked);
            }}
          />
        </div>
        <div className="flex items-center justify-between">
          <span className="text-sm font-medium">Schools</span>
          <Switch
            checked={showSchools}
            onCheckedChange={(checked) => {
              setShowSchools(checked);
              toggleLayerVisibility(SCHOOL_LAYER_ID, checked);
            }}
          />
        </div>
        <div className="space-y-2">
          <span className="text-sm font-medium">School Level</span>
          <Select value={selectedLevel} onValueChange={setSelectedLevel}>
            <SelectTrigger>
              <SelectValue placeholder="Select level" />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="all">All Levels</SelectItem>
              {schoolLevels.map((level) => (
                <SelectItem key={level} value={level}>
                  {level}
                </SelectItem>
              ))}
            </SelectContent>
          </Select>
        </div>
      </div>
    </div>
  ) : null;
}
