import { useEffect, useState } from 'react';
import { ChevronDown, ChartColumnBig } from 'lucide-react';
import { useProperty } from '@/contexts/property-context';
import { MapComponent } from '@am/commons';

interface PropertyInfoSectionProps {
  onCustomCompsChange?: (value: boolean) => void;
}

export function PropertyInfoSection({ onCustomCompsChange }: PropertyInfoSectionProps) {
  const { property, comparableProperties, isLoading } = useProperty();
  const [customComps, setCustomComps] = useState(false);

  // Map center and zoom state
  const [mapCenter, setMapCenter] = useState<[number, number]>([
    property.location.lng,
    property.location.lat,
  ]);
  const [mapZoom, setMapZoom] = useState(16);

  useEffect(() => {
    setMapCenter([property.location.lng, property.location.lat]);
    setMapZoom(16);
  }, [property]);

  const handleCustomCompsChange = (value: boolean) => {
    setCustomComps(value);
    if (onCustomCompsChange) {
      onCustomCompsChange(value);
    }
  };

  const Attribute = ({ label, value }: { label: string; value: string | number }) => (
    <div className="flex items-center">
      <ChartColumnBig size={20} className="text-gray-400 mr-3" />

      <div className="flex items-center">
        <div className="text-sm text-gray-400">{label}</div>
        <div className="font-semibold text-white pl-2">{value}</div>
      </div>
    </div>
  );

  return (
    <div className="bg-background text-white">
      <div className="flex h-[350px]">
        {/* Left column - All metrics, toggle, and pricing */}
        <div className="w-1/2 flex flex-col">
          {/* Metrics grid layout */}
          <div className="grid grid-cols-2 gap-x-6 gap-y-4 p-6 flex-grow">
            <Attribute label="Bedrooms" value={property.beds} />
            <Attribute label="Bathrooms" value={property.baths} />
            <Attribute label="Building Sqft" value={property.sqft} />
            <Attribute label="Lot size" value={property.lotSize} />
            <Attribute label="Type" value={property.propertyType} />
            <Attribute label="HOA" value="None" />
            <Attribute label="Year built" value={property.yearBuilt} />
            <Attribute label="Stories" value={2} />
          </div>

          {/* Custom comps toggle and dropdown */}
          <div className="px-6 pb-4 flex justify-between items-center">
            <div className="flex items-center">
              <span className="text-sm mr-4 text-white">Custom comps</span>
              <div
                className={`w-12 h-6 rounded-full p-1 cursor-pointer transition-colors ${customComps ? 'bg-blue-600' : 'bg-gray-700'}`}
                onClick={() => handleCustomCompsChange(!customComps)}
              >
                <div
                  className={`w-4 h-4 rounded-full bg-white transform transition-transform ${customComps ? 'translate-x-6' : ''}`}
                ></div>
              </div>
            </div>

            <div>
              <button className="bg-gray-800 hover:bg-gray-700 text-white px-4 py-2 rounded flex items-center">
                <span>High range rental comps</span>
                <ChevronDown size={16} className="ml-2" />
              </button>
            </div>
          </div>

          {/* Pricing Panel */}
          <div className="mx-6 mb-6 bg-gray-800 rounded-lg overflow-hidden">
            <div className="grid grid-cols-3">
              <div className="p-4">
                <div className="text-sm text-gray-400">List Price</div>
                <div className="text-2xl font-bold text-white">
                  ${property.price.toLocaleString()}
                </div>
              </div>
              <div className="p-4 border-l border-r border-gray-700">
                <div className="text-sm text-gray-400">Value</div>
                <div className="text-2xl font-bold text-white">$327,500</div>
              </div>
              <div className="p-4">
                <div className="text-sm text-gray-400">Rent</div>
                <div className="text-2xl font-bold text-white">
                  $3,100
                  <span className="text-xs text-gray-400 ml-1">/month</span>
                </div>
              </div>
            </div>
          </div>
        </div>

        {/* Right column - Full height map */}
        <div className="w-1/2 h-full">
          <MapComponent
            initialCenter={mapCenter}
            initialZoom={mapZoom}
            height="350px"
            width="100%"
            showControls={false}
            disableInteraction={true}
            comparableProperties={comparableProperties}
            property={property}
            isLoading={isLoading}
          />
        </div>
      </div>
    </div>
  );
}
