import { useCallback, useState, useEffect } from 'react';
import { useRouter, useParams } from '@tanstack/react-router';
import { usePageTitle } from '@/hooks/use-page-title';
import { useProperty } from '@/contexts/property-context';
import { PropertyInfoSection } from '@/pages/property/components/info-section';

export function PropertyDetailPage() {
  usePageTitle('Property Detail');
  const router = useRouter();
  const { propertyId } = useParams({ strict: false });
  const { property, isLoading, error } = useProperty();
  const [showFullDescription, setShowFullDescription] = useState(false);
  const [customComps, setCustomComps] = useState(false);
  const [activeTab, setActiveTab] = useState('overview');

  const handleBack = useCallback(() => {
    if (propertyId) {
      router.navigate({ to: '/property/$propertyId/map', params: { propertyId } });
    } else {
      router.navigate({ to: '/' });
    }
  }, [router, propertyId]);

  // Function to toggle the button style based on active state
  const getTabStyle = (tab: string) => {
    return activeTab === tab
      ? 'bg-gray-800 text-white'
      : 'text-gray-400 hover:text-white hover:bg-gray-800';
  };

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
      {/* Use the shared PropertyInfoSection component */}
      <PropertyInfoSection onCustomCompsChange={handleCustomCompsChange} />

      <div className="flex flex-col w-full bg-background text-white">
        {/* Main Content */}
        <div className="p-6 border-b border-gray-800">
          <h2 className="text-2xl font-bold text-white mb-2">Property Overview</h2>
        </div>

        <div className="p-6">
          {/* Tabs */}
          <div className="flex border-b border-gray-700 mb-6">
            <button
              className={`px-4 py-2 ${getTabStyle('overview')}`}
              onClick={() => setActiveTab('overview')}
            >
              Overview
            </button>
            <button
              className={`px-4 py-2 ${getTabStyle('analytics')}`}
              onClick={() => setActiveTab('analytics')}
            >
              Analytics
            </button>
            <button
              className={`px-4 py-2 ${getTabStyle('reports')}`}
              onClick={() => setActiveTab('reports')}
            >
              Reports
            </button>
            <button
              className={`px-4 py-2 ${getTabStyle('notifications')}`}
              onClick={() => setActiveTab('notifications')}
            >
              Notifications
            </button>
          </div>

          {/* Description */}
          <div className="mb-8">
            <h3 className="text-xl font-semibold text-white mb-3">Description</h3>
            <div className="bg-gray-800 p-4 rounded-md">
              <p className="text-gray-300">
                {showFullDescription
                  ? property.description
                  : `${property.description.substring(0, 200)}...`}
              </p>
              <button
                className="text-blue-400 mt-2"
                onClick={() => setShowFullDescription(!showFullDescription)}
              >
                {showFullDescription ? 'Show Less' : 'Show More'}
              </button>
            </div>
          </div>

          {/* Features */}
          <div className="mb-8">
            <h3 className="text-xl font-semibold text-white mb-3">Features</h3>
            <div className="bg-gray-800 p-4 rounded-md">
              <div className="grid grid-cols-2 gap-y-2">
                {property.features.map((feature, index) => (
                  <div key={index} className="flex items-center">
                    <div className="w-2 h-2 bg-blue-500 rounded-full mr-2"></div>
                    <span className="text-gray-300">{feature}</span>
                  </div>
                ))}
              </div>
            </div>
          </div>
        </div>
      </div>
    </>
  );
}
