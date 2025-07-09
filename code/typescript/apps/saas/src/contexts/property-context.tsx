import React, { createContext, useContext, useState, ReactNode, useCallback } from 'react';
import type { ComparableProperty as ComparablePropertyType } from '../services/mockdata/property-context-mock';
import { useParams } from '@tanstack/react-router';
import { PropertyDetails } from '@am/commons/schemas';
import { useComparablePropertiesQuery, usePropertyQuery } from '@/services/property';

// Re-export the ComparableProperty type
export type ComparableProperty = ComparablePropertyType;

// Define the context type
interface PropertyContextType {
  property: PropertyDetails | undefined;
  comparableProperties: ComparableProperty[];
  selectProperty: (id: string) => void;
  updateMapboxData: (mapboxData: any) => void;
  updatePropertyDetails: (details: Partial<PropertyDetails>) => void;
  isLoading: boolean;
  error: string | null;
}

// Create the context with default values
const PropertyContext = createContext<PropertyContextType | undefined>(undefined);

// Provider props type
interface PropertyProviderProps {
  children: ReactNode;
}

// Create the provider component
export function PropertyProvider({ children }: PropertyProviderProps) {
  const { propertyId = '' } = useParams({ strict: false });
  const [selectedPropertyId, setSelectedPropertyId] = useState<string>(propertyId);

  // Use React Query hooks
  const {
    data: property,
    isLoading: isPropertyLoading,
    error: propertyError,
  } = usePropertyQuery(selectedPropertyId);

  const {
    data: comparableProperties = [],
    isLoading: isComparablesLoading,
    error: comparablesError,
  } = useComparablePropertiesQuery(selectedPropertyId);

  const selectProperty = useCallback((id: string) => {
    setSelectedPropertyId(id);
  }, []);

  const updateMapboxData = useCallback(
    (mapboxData: any) => {
      if (!property) return;
      // In a real app, you would make an API call here
      console.log('Updating mapbox data:', mapboxData);
    },
    [property]
  );

  const updatePropertyDetails = useCallback(
    (details: Partial<PropertyDetails>) => {
      if (!property) return;
      // In a real app, you would make an API call here
      console.log('Updating property details:', details);
    },
    [property]
  );

  // Value to be provided to consumers
  const value = {
    property,
    comparableProperties,
    selectProperty,
    updateMapboxData,
    updatePropertyDetails,
    isLoading: isPropertyLoading || isComparablesLoading,
    error: propertyError?.message || comparablesError?.message || null,
  };

  return <PropertyContext.Provider value={value}>{children}</PropertyContext.Provider>;
}

// Custom hook to use the property context
export function useProperty() {
  const context = useContext(PropertyContext);
  if (context === undefined) {
    throw new Error('useProperty must be used within a PropertyProvider');
  }

  return context;
}

// Helper component to handle property loading from URL
export function PropertyLoader({ children }: { children: ReactNode }) {
  const { propertyId = '' } = useParams({ strict: false });
  const { isLoading, error } = useProperty();

  if (isLoading) {
    return <div>Loading property data...</div>;
  }

  if (error) {
    return <div>Error loading property: {error}</div>;
  }

  return <>{children}</>;
}
