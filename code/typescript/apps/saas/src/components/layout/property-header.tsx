import { ReactNode } from 'react';
import { useRouter, useParams } from '@tanstack/react-router';
import { Download, Heart, BarChart2 } from 'lucide-react';
import { useProperty } from '@/contexts/property-context';
import { Image } from '@am/commons/components/ui/image';
import { Skeleton } from '@am/commons/components/ui/skeleton';

interface PropertyHeaderProps {
  showExploreComps?: boolean;
  onExploreCompsClick?: () => void;
  children?: ReactNode;
}

export function PropertyHeader({
  showExploreComps = true,
  onExploreCompsClick,
  children,
}: PropertyHeaderProps) {
  const router = useRouter();
  const { propertyId = '' } = useParams({ strict: false });
  const { property, isLoading } = useProperty();

  const handleExploreCompsClick = () => {
    if (onExploreCompsClick) {
      onExploreCompsClick();
    } else if (propertyId) {
      router.navigate({
        to: '/property/$propertyId/map',
        params: { propertyId },
      });
    } else {
      router.navigate({ to: '/' });
    }
  };

  if (isLoading) {
    return (
      <div className="px-6 py-4 flex justify-between items-start">
        <div className="flex items-start">
          <Skeleton className="w-20 h-16 rounded mr-4" />
          <div>
            <Skeleton className="h-8 w-48 mb-2" />
            <Skeleton className="h-6 w-64 mb-2" />
            <div className="flex items-center mt-1">
              <Skeleton className="h-4 w-16" />
            </div>
          </div>
        </div>
        <div className="flex space-x-2">
          {showExploreComps && <Skeleton className="h-9 w-32" />}
          {children}
          <Skeleton className="h-9 w-24" />
          <Skeleton className="h-9 w-24" />
        </div>
      </div>
    );
  }

  if (!property) {
    return null;
  }

  return (
    <div className="px-6 py-4 flex justify-between items-start">
      <div className="flex items-start">
        <div className="w-20 h-16 rounded overflow-hidden mr-4">
          <Image
            src={property.mainImage}
            alt="Property"
            className="max-w-full max-h-full object-cover"
            fallbackClassName="w-full h-full bg-gray-800"
            useWrapper={true}
            wrapperClassName="w-full h-full flex items-center justify-center"
          />
        </div>
        <div>
          <h1 className="text-2xl font-bold text-white">Property Detail</h1>
          <p className="text-white text-lg">{property.fullAddress}</p>
          <div className="flex items-center mt-1">
            <span className="text-gray-400 mr-2">Status</span>
            <span className="bg-blue-500 text-white text-xs px-2 py-1 rounded-full">
              {property.status}
            </span>
          </div>
        </div>
      </div>
      <div className="flex space-x-2">
        {showExploreComps && (
          <button
            className="flex items-center bg-indigo-600 hover:bg-indigo-700 text-white px-4 py-2 rounded-md transition-colors"
            onClick={handleExploreCompsClick}
          >
            <BarChart2 size={16} className="mr-2" />
            Explore Comps
          </button>
        )}
        {children}
        <button className="flex items-center px-4 py-2 border border-gray-700 rounded-md bg-gray-800 hover:bg-gray-700">
          <span className="mr-2">Download</span>
        </button>
        <button className="flex items-center px-4 py-2 border border-gray-700 rounded-md bg-gray-800 hover:bg-gray-700">
          <Heart size={16} className="mr-2" />
          <span>Favorite</span>
        </button>
      </div>
    </div>
  );
}
