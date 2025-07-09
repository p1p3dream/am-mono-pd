import { useCallback, useState, useEffect } from 'react';
import { useRouter, useParams } from '@tanstack/react-router';
import { usePageTitle } from '@/hooks/use-page-title';
import { BarChart2 } from 'lucide-react';
import { ScrollArea } from '@am/commons/components/ui/scroll-area';
import { useProperty } from '@/contexts/property-context';
import { PropertyInfoSection } from '@/pages/property/components/info-section';
import { ComparablePropertiesList } from './components/lists';

export function PropertyCMAPage() {
  usePageTitle('Property CMA');
  const router = useRouter();
  const { propertyId } = useParams({ strict: false });
  const { property, comparableProperties, selectProperty, isLoading, error } = useProperty();
  const [customComps, setCustomComps] = useState(false);

  const handleBack = useCallback(() => {
    if (propertyId) {
      router.navigate({ to: '/property/$propertyId/map', params: { propertyId } });
    } else {
      router.navigate({ to: '/' });
    }
  }, [router, propertyId]);

  const handleCustomCompsChange = (value: boolean) => {
    setCustomComps(value);
  };

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

  return (
    <>
      <PropertyInfoSection onCustomCompsChange={handleCustomCompsChange} />
      <div className="bg-background text-white">
        <div className="p-6 border-b border-gray-800">
          <h1 className="text-2xl font-bold text-white mb-2">Comparative Market Analysis</h1>
          <p className="text-gray-300">
            View and compare market data for similar properties in the area.
          </p>
        </div>

        {/* Map and Comparables Section */}
        <ScrollArea className="flex-1">
          <div className="flex flex-col ">
            {/* Map Section */}
            <div className="w-full p-6 border-r border-gray-800">
              {/* Summary Metrics */}
              <div className="bg-gray-800 p-4 rounded-md">
                <h3 className="text-lg font-semibold mb-3">Market Analysis Summary</h3>
                <div className="grid grid-cols-2 gap-4">
                  <div>
                    <div className="text-gray-400 text-sm">Average Listing Price</div>
                    <div className="text-xl font-bold">$767,000</div>
                  </div>
                  <div>
                    <div className="text-gray-400 text-sm">Average Price per SqFt</div>
                    <div className="text-xl font-bold">$325</div>
                  </div>
                  <div>
                    <div className="text-gray-400 text-sm">Days on Market (Avg)</div>
                    <div className="text-xl font-bold">34</div>
                  </div>
                  <div>
                    <div className="text-gray-400 text-sm">Year Built (Avg)</div>
                    <div className="text-xl font-bold">1992</div>
                  </div>
                </div>
                <button className="mt-4 flex items-center text-blue-400 hover:text-blue-300">
                  <BarChart2 size={16} className="mr-2" />
                  <span>View Detailed Report</span>
                </button>
              </div>
            </div>
          </div>
          {/* Comparable Properties List */}
          <div className="w-full ">
            <ComparablePropertiesList
              properties={comparableProperties}
              onPropertySelect={selectProperty}
              rowClassName="border-gray-800 hover:bg-gray-800/50"
              showExport={false}
              showFilter={true}
              title="Comparable Properties"
              buttonLabel={{ selected: 'Selected', notSelected: 'Select' }}
            />
          </div>
        </ScrollArea>
      </div>
    </>
  );
}
