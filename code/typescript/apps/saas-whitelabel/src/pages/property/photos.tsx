import { useCallback, useState, useEffect } from 'react';
import { useRouter, useParams } from '@tanstack/react-router';
import { usePageTitle } from '@/hooks/use-page-title';
import { ScrollArea } from '@am/commons/components/ui/scroll-area';
import { useProperty } from '@/contexts/property-context';
import { PropertyImage } from '@/services/properties';
import { Image } from '@am/commons/components/ui/image';
import { PropertyInfoSection } from '@/pages/property/components/info-section';

// Default image to use if no images are available
const defaultImage: PropertyImage = {
  id: 0,
  url: 'https://via.placeholder.com/800x600?text=No+Image+Available',
  alt: 'No Image Available',
};

export function PropertyPhotosPage() {
  usePageTitle('Property Photos');
  const router = useRouter();
  const { propertyId } = useParams({ strict: false });
  const { property, isLoading, error } = useProperty();

  // Use images from the property context
  const propertyImages = property.images;

  // Initialize with the default image
  const [selectedImage, setSelectedImage] = useState<PropertyImage>(defaultImage);

  // Update selected image when property images are available
  useEffect(() => {
    const firstImage = propertyImages?.[0];
    if (firstImage) {
      setSelectedImage(firstImage);
    }
  }, [propertyImages]);

  const handleThumbnailClick = (image: PropertyImage) => {
    setSelectedImage(image);
  };

  const handleCustomCompsChange = (value: boolean) => {
    // Handle custom comps change
    console.log('Custom comps changed:', value);
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
        {/* Main content */}
        <div className="p-6">
          <h2 className="text-2xl font-bold text-white mb-4">Property Photos</h2>

          {/* Main image display with improved centering */}
          <div className="mb-4 rounded-lg overflow-hidden bg-gray-800 h-[500px]">
            <Image
              src={selectedImage.url}
              alt={selectedImage.alt}
              className="max-w-full max-h-full object-contain"
              fallbackClassName="w-full h-full flex items-center justify-center bg-gray-800 text-gray-500"
              useWrapper={true}
              wrapperClassName="w-full h-full flex items-center justify-center"
            />
          </div>

          {/* Thumbnails */}
          <ScrollArea className="w-full">
            <div className="flex space-x-2 pb-2">
              {propertyImages.map((image) => (
                <div
                  key={image.id}
                  className={`w-24 h-24 rounded-md overflow-hidden cursor-pointer flex-shrink-0 ${
                    selectedImage.id === image.id ? 'ring-2 ring-purple-500' : ''
                  }`}
                  onClick={() => handleThumbnailClick(image)}
                >
                  <Image
                    src={image.url}
                    alt={image.alt}
                    className="max-w-full max-h-full object-cover"
                    fallbackClassName="w-full h-full bg-gray-700 flex items-center justify-center"
                    useWrapper={true}
                    wrapperClassName="w-full h-full flex items-center justify-center"
                  />
                </div>
              ))}
            </div>
          </ScrollArea>
        </div>
      </div>
    </>
  );
}
