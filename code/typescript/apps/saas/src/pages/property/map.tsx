import { useState, useCallback, useEffect, useMemo } from 'react';
import { usePageTitle } from '@/hooks/use-page-title';
import { useProperty } from '@/contexts/property-context';
import { useParams, useRouter } from '@tanstack/react-router';

// Icons
import { Home, School, Users, Droplets, MapPin, Building, Calendar, Layers } from 'lucide-react';
import { MapComponent, BuildingFeature } from '@am/commons';

import { ComparablePropertiesList } from './components/lists';
import { useSchoolsQuery } from '@/services/school';
import { useSubdivisionsQuery } from '@/services/subdivision';

export function MapPage() {
  usePageTitle('Property Map');
  const router = useRouter();
  const { propertyId } = useParams({ strict: false });
  const { property, comparableProperties, isLoading, error } = useProperty();
  const { data: schools = [] } = useSchoolsQuery();
  const { data: subdivisions = [] } = useSubdivisionsQuery();

  // State for the selected building
  const [selectedBuilding, setSelectedBuilding] = useState<BuildingFeature | null>(null);
  const [viewMode, setViewMode] = useState<'buy' | 'rent'>('buy');
  const [activeTab, setActiveTab] = useState<string | null>('Buy');

  // Track when a user has manually selected a building
  const [userSelectedBuilding, setUserSelectedBuilding] = useState<boolean>(false);

  // Keep track of selected building location for centering
  const [selectedLocation, setSelectedLocation] = useState<[number, number] | null>(null);

  // Compute the map center coordinates
  const mapCenter = useMemo<[number, number]>(() => {
    if (!property) return [0, 0];

    // If user selected a building location, prioritize that
    if (userSelectedBuilding && selectedLocation) {
      return selectedLocation;
    }

    // If mapboxData exists and has coordinates, use those
    if (
      property.mapboxData &&
      typeof property.mapboxData.lng === 'number' &&
      typeof property.mapboxData.lat === 'number'
    ) {
      return [property.mapboxData.lng, property.mapboxData.lat];
    }

    // Otherwise fall back to property location
    return [property.location.lng, property.location.lat];
  }, [property, userSelectedBuilding, selectedLocation]);

  // Compute zoom level based on whether we're looking at a building or general area
  const zoomLevel = useMemo<number>(() => {
    if (!property) return 15;

    // For user-selected or property buildings, zoom in closer
    if (
      userSelectedBuilding ||
      (property.mapboxData &&
        typeof property.mapboxData.lng === 'number' &&
        typeof property.mapboxData.lat === 'number')
    ) {
      return 17; // Closer zoom for buildings
    }

    // Otherwise use a wider view
    return 15; // Default zoom for property area
  }, [property, userSelectedBuilding]);

  // Generate a unique key for the map component to force remounts when needed
  const mapKey = useMemo(() => {
    if (!property) return 'map-default';

    // When user selects a building or property changes, create a new key
    if (userSelectedBuilding && selectedLocation) {
      return `map-${property.id}-user-${selectedLocation[0]}-${selectedLocation[1]}`;
    }
    return `map-${property.id}`;
  }, [property?.id, userSelectedBuilding, selectedLocation]);

  // Set the selected building when property changes
  useEffect(() => {
    if (!property) return;

    // Reset user selection when property changes
    setUserSelectedBuilding(false);
    setSelectedLocation(null);

    // When property changes, focus on the new property
    if (property.mapboxData) {
      setSelectedBuilding(property.mapboxData as BuildingFeature);
    } else {
      setSelectedBuilding(null);
    }
  }, [property]);

  // Toggle active tab and overlays
  const toggleTab = (tabName: string) => {
    if (activeTab === tabName) {
      setActiveTab(null);
    } else {
      setActiveTab(tabName);
    }
  };

  // Handling back navigation
  const handleBack = useCallback(() => {
    router.navigate({ to: '/' });
  }, [router]);

  // Show loading state
  if (isLoading) {
    return (
      <div className="flex items-center justify-center h-64">
        <div className="text-xl text-white">Loading property data...</div>
      </div>
    );
  }

  // Show error state
  if (error) {
    return (
      <div className="flex items-center justify-center h-64">
        <div className="text-xl text-red-500">Error: {error}</div>
      </div>
    );
  }

  if (!property) {
    return (
      <div className="flex items-center justify-center h-64">
        <div className="text-xl text-white">No property data available</div>
      </div>
    );
  }

  return (
    <div className="flex flex-col w-full h-full overflow-hidden">
      {/* Main Content */}
      <div className="flex-1 flex flex-col overflow-hidden relative">
        {/* Top Navigation Panel - Not Scrollable */}
        <div className="flex justify-between absolute top-6 left-1/2 -translate-x-1/2 z-10">
          <div className="flex justify-around m-auto rounded-md text-white">
            <div className="flex flex-col items-center">
              <button
                className={`flex flex-col items-center justify-center px-4 py-2 rounded-md cursor-pointer h-[72px] ${
                  activeTab === 'Schools'
                    ? 'bg-purple-400/60 hover:bg-purple-400/80'
                    : 'bg-background/70 hover:bg-background/90'
                }`}
                onClick={() => toggleTab('Schools')}
              >
                <School size={18} />
                <span className="text-xs">Schools</span>
              </button>
            </div>
            <button
              className={`flex flex-col items-center justify-center px-4 py-2 rounded-md cursor-pointer h-[72px] ${
                activeTab === 'Admin'
                  ? 'bg-purple-400/60 hover:bg-purple-400/80'
                  : 'bg-background/70 hover:bg-background/90'
              }`}
              onClick={() => toggleTab('Admin')}
            >
              <Users size={18} />
              <span className={`text-xs`}>Admin/Subdivisions</span>
            </button>
          </div>
        </div>

        {/* Scrollable Content Area */}
        <div className="relative">
          {/* Map Container */}
          <div className="m-6 rounded-xl">
            <MapComponent
              key={mapKey}
              initialCenter={mapCenter}
              initialZoom={zoomLevel}
              height="550px"
              width="100%"
              showControls={true}
              className="w-full"
              showSchools={activeTab === 'Schools'}
              showSubdivisions={activeTab === 'Admin'}
              showParcels={activeTab === 'Parcels'}
              getSchools={() => schools}
              getSubdivisions={() => subdivisions}
              comparableProperties={comparableProperties}
              property={property}
              isLoading={isLoading}
            />
          </div>
        </div>

        {/* Property Information Section - Always visible, matches screenshot design */}
        <div className="p-6 bg-background/80 rounded-md mx-6 mt-4">
          <div className="grid grid-cols-2 md:grid-cols-4 gap-y-4 mb-6">
            <div className="flex items-center">
              <div className="text-gray-400 mr-3">
                <Home size={18} />
              </div>
              <div>
                <span className="text-gray-400 text-xs">Bedrooms</span>
                <div className="text-white font-semibold">{property.beds}</div>
              </div>
            </div>

            <div className="flex items-center">
              <div className="text-gray-400 mr-3">
                <Droplets size={18} />
              </div>
              <div>
                <span className="text-gray-400 text-xs">Bathrooms</span>
                <div className="text-white font-semibold">{property.baths}</div>
              </div>
            </div>

            <div className="flex items-center">
              <div className="text-gray-400 mr-3">
                <Building size={18} />
              </div>
              <div>
                <span className="text-gray-400 text-xs">Building Sqft</span>
                <div className="text-white font-semibold">
                  {property.sqft.toLocaleString()} sq ft
                </div>
              </div>
            </div>

            <div className="flex items-center">
              <div className="text-gray-400 mr-3">
                <MapPin size={18} />
              </div>
              <div>
                <span className="text-gray-400 text-xs">Lot size</span>
                <div className="text-white font-semibold">
                  {(property.lotSize * 43560).toLocaleString()} sq ft
                </div>
              </div>
            </div>

            <div className="flex items-center">
              <div className="text-gray-400 mr-3">
                <Home size={18} />
              </div>
              <div>
                <span className="text-gray-400 text-xs">Type</span>
                <div className="text-white font-semibold">{property.propertyType}</div>
              </div>
            </div>

            <div className="flex items-center">
              <div className="text-gray-400 mr-3">
                <Users size={18} />
              </div>
              <div>
                <span className="text-gray-400 text-xs">HOA</span>
                <div className="text-white font-semibold">None</div>
              </div>
            </div>

            <div className="flex items-center">
              <div className="text-gray-400 mr-3">
                <Calendar size={18} />
              </div>
              <div>
                <span className="text-gray-400 text-xs">Year built</span>
                <div className="text-white font-semibold">{property.yearBuilt}</div>
              </div>
            </div>

            <div className="flex items-center">
              <div className="text-gray-400 mr-3">
                <Layers size={18} />
              </div>
              <div>
                <span className="text-gray-400 text-xs">Stories</span>
                <div className="text-white font-semibold">2</div>
              </div>
            </div>
          </div>

          {/* Price Box */}
          <div className="grid grid-cols-3 gap-4 bg-background p-4 rounded-md border border-gray-800">
            <div>
              <div className="text-gray-400 text-sm">List Price</div>
              <div className="text-white text-2xl font-bold">
                ${property.price.toLocaleString()}
              </div>
            </div>
            <div>
              <div className="text-gray-400 text-sm">Value</div>
              <div className="text-white text-2xl font-bold">$327,500</div>
            </div>
            <div>
              <div className="text-gray-400 text-sm">Rent</div>
              <div className="text-white text-2xl font-bold">
                $2,700<span className="text-sm font-normal text-gray-400">/month</span>
              </div>
            </div>
          </div>
        </div>

        {/* Comparables Summary */}
        {selectedBuilding && (
          <div className="p-6 border-b border-purple-900">
            <div className="bg-background rounded-md p-4">
              {/* Buy/Rent Buttons */}
              <div className="flex mb-4">
                <button
                  className={`py-2 px-6 rounded-l-md text-sm font-medium ${
                    viewMode === 'buy' ? 'bg-purple-900 text-white' : 'bg-gray-800 text-gray-400'
                  }`}
                  onClick={() => setViewMode('buy')}
                >
                  Buy
                </button>
                <button
                  className={`py-2 px-6 rounded-r-md text-sm font-medium ${
                    viewMode === 'rent' ? 'bg-purple-900 text-white' : 'bg-gray-800 text-gray-400'
                  }`}
                  onClick={() => setViewMode('rent')}
                >
                  Rent
                </button>
              </div>

              {/* Price Box */}
              <div className="bg-background rounded-md">
                {viewMode === 'buy' && (
                  <div className="grid grid-cols-5 gap-4">
                    <div>
                      <div className="text-gray-400 text-sm">Current List Price</div>
                      <div className="text-orange-500 text-xl font-bold">$2,125,000</div>
                    </div>
                    <div>
                      <div className="text-gray-400 text-sm">Median MLS Listings</div>
                      <div className="text-white text-xl font-bold">
                        $327,500
                        <span className="text-xs font-normal ml-1">$193/Sqft</span>
                      </div>
                    </div>
                    <div>
                      <div className="text-gray-400 text-sm">
                        Median Nat'l SRF Operators Listings
                      </div>
                      <div className="text-white text-xl font-bold">$327,200</div>
                    </div>
                    <div>
                      <div className="text-gray-400 text-sm">Median 3rd Party Listings</div>
                      <div className="text-white text-xl font-bold">$327,500</div>
                    </div>
                    <div>
                      <div className="text-gray-400 text-sm">Market Rent</div>
                      <div className="text-orange-500 text-xl font-bold">
                        $2900
                        <span className="text-sm font-normal ml-1 text-gray-400">/month</span>
                      </div>
                    </div>
                  </div>
                )}

                {viewMode === 'rent' && (
                  <div className="grid grid-cols-5 gap-4">
                    <div>
                      <div className="text-gray-400 text-sm">Current Rent</div>
                      <div className="text-orange-500 text-xl font-bold">$2,700/mo</div>
                    </div>
                    <div>
                      <div className="text-gray-400 text-sm">Median MLS Listings</div>
                      <div className="text-white text-xl font-bold">$2,900/mo</div>
                    </div>
                    <div>
                      <div className="text-gray-400 text-sm">Median Nat'l SRF Operators</div>
                      <div className="text-white text-xl font-bold">$2,750/mo</div>
                    </div>
                    <div>
                      <div className="text-gray-400 text-sm">Median 3rd Party</div>
                      <div className="text-white text-xl font-bold">$2,800/mo</div>
                    </div>
                    <div>
                      <div className="text-gray-400 text-sm">Market Rent</div>
                      <div className="text-orange-500 text-xl font-bold">
                        $2900
                        <span className="text-sm font-normal ml-1 text-gray-400">/month</span>
                      </div>
                    </div>
                  </div>
                )}
              </div>
            </div>
          </div>
        )}

        {/* Comparables List */}
        <div className="p-6">
          <div>
            <ComparablePropertiesList properties={comparableProperties} />
          </div>
        </div>
      </div>
    </div>
  );
}

export default MapPage;
