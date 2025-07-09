import { ReactNode } from 'react';
import { Outlet } from '@tanstack/react-router';
import { PropertyProvider, useProperty } from '../../contexts/property-context';
import { PropertyLayout } from './property-layout';

function PropertyContent() {
  const { isLoading, error } = useProperty();

  if (isLoading) {
    return (
      <div className="flex items-center justify-center h-full">
        <div className="text-xl text-white">Loading property data...</div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="flex items-center justify-center h-full">
        <div className="text-xl text-red-500">Error: {error}</div>
      </div>
    );
  }

  return <Outlet />;
}

export function PropertyLayoutWithContext() {
  return (
    <PropertyProvider>
      <PropertyLayout>
        <PropertyContent />
      </PropertyLayout>
    </PropertyProvider>
  );
}
