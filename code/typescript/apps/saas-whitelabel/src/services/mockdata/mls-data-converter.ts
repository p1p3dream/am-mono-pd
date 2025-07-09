import { PropertyDetails, PropertyImage } from '@am/commons/schemas';

/**
 * Convert MLS JSON data to PropertyDetails format
 * @param mlsData The raw MLS listing data from the JSON file
 * @returns PropertyDetails object with data mapped from MLS format
 */
export function convertMlsDataToPropertyDetails(mlsData: any): PropertyDetails {
  // Create basic image array from PhotosCount
  const images: PropertyImage[] = [];

  // If we have image data, create placeholder images
  if (
    mlsData.PhotosCount &&
    mlsData.PhotosCount > 0 &&
    mlsData.PhotoKey &&
    mlsData.PhotoURLPrefix
  ) {
    // In a real implementation, we would use the actual photo URLs from the MLS data
    // For Nashville MLS, the actual URL pattern might look like:
    // https://d2sa73fkuaa019.cloudfront.net/listings/mtrmls/[PhotoKey]/[PhotoNumber].jpg

    // Try to use real photo URLs if available in the MLS data
    for (let i = 0; i < Math.min(mlsData.PhotosCount, 10); i++) {
      // In a real implementation, you'd build the actual URLs based on the PhotoKey and PhotoURLPrefix
      // This is a placeholder implementation
      images.push({
        id: i,
        url: `${mlsData.PhotoURLPrefix}${mlsData.PhotoKey}/photo_${i + 1}.jpg`,
        alt: `Property image ${i + 1}`,
      });
    }
  }

  // Extract features from various MLS fields
  const features: string[] = [];
  if (mlsData.ConstructionMaterials)
    features.push(`Construction: ${mlsData.ConstructionMaterials}`);
  if (mlsData.BasementFeatures) features.push(`Basement: ${mlsData.BasementFeatures}`);
  if (mlsData.GarageSpaces) features.push(`Garage Spaces: ${mlsData.GarageSpaces}`);

  // Add property type as a feature
  if (mlsData.MLSPropertyType) features.push(mlsData.MLSPropertyType);
  if (mlsData.MLSPropertySubType) features.push(mlsData.MLSPropertySubType);

  // Add year built as a feature if available
  if (mlsData.YearBuilt) features.push(`Built in ${mlsData.YearBuilt}`);

  // Determine main image safely
  const mainImage = images && images.length > 0 && images[0] ? images[0].url : '';

  // Map MLS data to PropertyDetails format
  const propertyDetails: PropertyDetails = {
    id: mlsData.ATTOM_ID?.toString() || mlsData.MLSListingID?.toString() || '',
    address: mlsData.PropertyAddressFull || mlsData.MLSListingAddress || '',
    city: mlsData.PropertyAddressCity || mlsData.MLSListingCity || '',
    state: mlsData.PropertyAddressState || mlsData.MLSListingState || '',
    zipCode: mlsData.PropertyAddressZIP?.toString() || mlsData.MLSListingZip || '',
    fullAddress: mlsData.PropertyAddressFull || mlsData.MLSListingAddress || '',
    status: mlsData.ListingStatus || 'Unknown',
    price: mlsData.LatestListingPrice || 0,
    beds: mlsData.BedroomsTotal || 0,
    baths: (mlsData.BathroomsFull || 0) + (mlsData.BathroomsHalf || 0) * 0.5,
    sqft: mlsData.LivingAreaSquareFeet || 0,
    lotSize: mlsData.LotSizeSquareFeet || 0,
    yearBuilt: mlsData.YearBuilt || 0,
    propertyType: mlsData.MLSPropertyType || mlsData.ATTOMPropertySubType || 'Unknown',
    description: mlsData.PublicListingRemarks || '',
    features: features,
    mainImage: mainImage,
    images: images,
    location: {
      lat: mlsData.Latitude || 0,
      lng: mlsData.Longitude || 0,
    },
    parcelGeometry: mlsData.parcelGeometry || mlsData.parcel_geometry || null,
    parcelProperties: mlsData.parcelProperties || mlsData.parcel_properties || null,

    // Additional MLS specific fields
    mlsNumber: mlsData.MLSNumber || '',
    mlsSource: mlsData.MLSSource || '',
    listingDate: mlsData.ListingDate || '',
    originalListingDate: mlsData.OriginalListingDate || '',
    originalListingPrice: mlsData.OriginalListingPrice || 0,
    daysOnMarket: mlsData.DaysOnMarket || 0,
    schoolElementary: mlsData.SchoolElementary || '',
    schoolMiddle: mlsData.SchoolMiddle || '',
    schoolHigh: mlsData.SchoolHigh || '',
    schoolDistrict:
      mlsData.SchoolElementaryDistrict ||
      mlsData.SchoolMiddleDistrict ||
      mlsData.SchoolHighDistrict ||
      '',
    taxAmount: mlsData.TaxAmount || 0,
    specialListingConditions: mlsData.SpecialListingConditions || '',
    constructionMaterials: mlsData.ConstructionMaterials || '',
    basementFeatures: mlsData.BasementFeatures || '',
    garageSpaces: mlsData.GarageSpaces || 0,
    lotSizeAcres: mlsData.LotSizeAcres || 0,
    lotDimensions: mlsData.LotDimensions || '',
    listingAgentName: mlsData.ListingAgentFullName || '',
    listingAgentPhone: mlsData.ListingAgentPreferredPhone || '',
    listingAgentEmail: mlsData.ListingAgentEmail || '',
    listingOfficeName: mlsData.ListingOfficeName || '',
    listingOfficePhone: mlsData.ListingOfficePhone || '',
    newConstruction: mlsData.NewConstructionYN === 'Y',
    photoURLPrefix: mlsData.PhotoURLPrefix || '',
    photoKey: mlsData.PhotoKey || '',
    photosCount: mlsData.PhotosCount || 0,
    hasDotMarker: false, // Default to false, will be overridden in property-context-mock for selected properties
  };

  return propertyDetails;
}

/**
 * Convert an array of MLS data to an array of PropertyDetails
 * @param mlsDataArray Array of MLS listing data from the JSON file
 * @returns Array of PropertyDetails objects
 */
export function convertMlsDataArrayToPropertyDetails(mlsDataArray: any[]): PropertyDetails[] {
  return mlsDataArray.map((mlsData) => convertMlsDataToPropertyDetails(mlsData));
}
