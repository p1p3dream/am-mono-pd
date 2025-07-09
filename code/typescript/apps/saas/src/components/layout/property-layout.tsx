import { ReactNode } from 'react';
import { PropertySidebar } from './property-sidebar';
import { ScrollArea } from '@am/commons/components/ui/scroll-area';
import { PropertyHeader } from './property-header';

type PropertyLayoutProps = {
  children: ReactNode;
  onBack?: () => void;
  showBackButton?: boolean;
  hideScrollArea?: boolean;
  showExploreComps?: boolean;
  onExploreCompsClick?: () => void;
};

export function PropertyLayout({
  children,
  showExploreComps = true,
  onExploreCompsClick,
}: PropertyLayoutProps) {
  return (
    <div className="flex h-screen bg-background text-white overflow-hidden">
      {/* Property Sidebar */}
      <PropertySidebar />

      {/* Main Content */}
      <div className="flex-1 flex flex-col overflow-hidden">
        {/* Content Area */}
        <ScrollArea className="flex-1 p-4">
          {/* Property Header */}
          <PropertyHeader
            showExploreComps={showExploreComps}
            onExploreCompsClick={onExploreCompsClick}
          />
          <div className="border border-gray-500 rounded-md p-4">{children}</div>
        </ScrollArea>
      </div>
    </div>
  );
}
