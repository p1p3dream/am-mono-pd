import { ReactNode } from 'react';
import { PortfolioSidebar } from './portfolio-sidebar';
import { Outlet } from '@tanstack/react-router';

interface PortfolioLayoutProps {
  children?: ReactNode;
}

export function PortfolioLayout({ children }: PortfolioLayoutProps) {
  return (
    <div className="flex h-screen bg-gray-50 overflow-hidden">
      <PortfolioSidebar />
      <div className="flex-1 overflow-auto">{children || <Outlet />}</div>
    </div>
  );
}

export function PortfolioLayoutWithContext() {
  return <PortfolioLayout></PortfolioLayout>;
}
