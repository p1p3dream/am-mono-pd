import React, { createContext, useContext, useState, ReactNode, useCallback } from 'react';
import { PropertyDetails } from '@am/commons/schemas';

interface Portfolio {
  id: string;
  name: string;
  totalProperties: number;
  totalPurchasePrice: number;
  totalValue: number;
  avgPPU: number;
  avgRent: number;
  uwGrossYield: number;
  uwNetCapRate: number;
  msa: string;
  properties: PropertyDetails[];
}

// Define the context type
interface PortfolioContextType {
  portfolio: Portfolio | undefined;
  portfolios: Portfolio[];
  selectPortfolio: (id: string) => void;
  isLoading: boolean;
  error: string | null;
}

// Create the context with default values
const PortfolioContext = createContext<PortfolioContextType | undefined>(undefined);

// Provider props type
interface PortfolioProviderProps {
  children: ReactNode;
}

// Create the provider component
export function PortfolioProvider({ children }: PortfolioProviderProps) {
  const [selectedPortfolioId, setSelectedPortfolioId] = useState<string>('');
  const [isLoading, setIsLoading] = useState<boolean>(false);
  const [error, setError] = useState<string | null>(null);

  // Mock data for now - would come from an API
  const mockPortfolios: Portfolio[] = [
    {
      id: 'portfolio-1',
      name: 'Princeton, TX 2-4 bds',
      totalProperties: 221,
      totalPurchasePrice: 1242124,
      totalValue: 1500722,
      avgPPU: 241215,
      avgRent: 2300,
      uwGrossYield: 11.376,
      uwNetCapRate: 6.532,
      msa: 'Princeton, TX',
      properties: [],
    },
    {
      id: 'portfolio-2',
      name: 'DFW 89',
      totalProperties: 89,
      totalPurchasePrice: 987654,
      totalValue: 1100000,
      avgPPU: 210000,
      avgRent: 1950,
      uwGrossYield: 10.5,
      uwNetCapRate: 5.8,
      msa: 'Dallas-Fort Worth, TX',
      properties: [],
    },
  ];

  const selectPortfolio = useCallback((id: string) => {
    setSelectedPortfolioId(id);
  }, []);

  // Find the selected portfolio
  const portfolio = mockPortfolios.find((p) => p.id === selectedPortfolioId);

  // Value to be provided to consumers
  const value = {
    portfolio,
    portfolios: mockPortfolios,
    selectPortfolio,
    isLoading,
    error,
  };

  return <PortfolioContext.Provider value={value}>{children}</PortfolioContext.Provider>;
}

// Custom hook to use the portfolio context
export function usePortfolio() {
  const context = useContext(PortfolioContext);
  if (context === undefined) {
    throw new Error('usePortfolio must be used within a PortfolioProvider');
  }

  return context;
}

// Helper component to handle portfolio loading
export function PortfolioLoader({ children }: { children: ReactNode }) {
  const { isLoading, error } = usePortfolio();

  if (isLoading) {
    return <div>Loading portfolio data...</div>;
  }

  if (error) {
    return <div>Error loading portfolio: {error}</div>;
  }

  return <>{children}</>;
}
