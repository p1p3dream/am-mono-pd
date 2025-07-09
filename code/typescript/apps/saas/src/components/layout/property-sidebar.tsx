import { ReactNode, useCallback, useEffect, useState } from 'react';
import { useRouter, useRouterState, useParams } from '@tanstack/react-router';
import { LayoutDashboard, BarChart, Map, Image, User, ArrowLeft } from 'lucide-react';
import { useAuth } from '@/contexts/auth-context';

type SideMenuItem = {
  icon: ReactNode;
  label: string;
  path: string;
  paramsPath?: string; // Path with parameter placeholders for navigation
};

// The sidebar menu items for property pages
export function getPropertyMenuItems(propertyId: string): SideMenuItem[] {
  return [
    {
      icon: <LayoutDashboard size={20} />,
      label: 'Detail',
      path: `/property/${propertyId}/detail`,
      paramsPath: '/property/$propertyId/detail',
    },
    {
      icon: <BarChart size={20} />,
      label: 'CMA',
      path: `/property/${propertyId}/cma`,
      paramsPath: '/property/$propertyId/cma',
    },
    {
      icon: <Map size={20} />,
      label: 'Map',
      path: `/property/${propertyId}/map`,
      paramsPath: '/property/$propertyId/map',
    },
    {
      icon: <Image size={20} />,
      label: 'Photos',
      path: `/property/${propertyId}/photos`,
      paramsPath: '/property/$propertyId/photos',
    },
  ];
}

let initialNavLength = 0;

type PropertySidebarProps = {};

export function PropertySidebar({}: PropertySidebarProps) {
  const router = useRouter();
  const { location } = useRouterState();
  const currentPath = location.pathname;
  const { propertyId = '' } = useParams({ strict: false });
  const { signOut } = useAuth();
  const [menuItems, setMenuItems] = useState<SideMenuItem[]>([]);

  useEffect(() => {
    if (propertyId) {
      setMenuItems(getPropertyMenuItems(propertyId));
    }
  }, [propertyId]);

  useEffect(() => {
    console.log('router history', router.history.length);
    if (initialNavLength === 0) {
      initialNavLength = router.history.length;
    }
  }, []);

  const handleBackNavigation = useCallback(() => {
    router.navigate({ to: '/' });
    return;

    const navCount = router.history.length - initialNavLength;

    // Go back by the number of app navigations (plus 1 to exit the app)
    // This skips all the internal app navigation history at once
    if (navCount > 0) {
      window.history.go(-(navCount + 1));
    }

    // Fallback: If after a delay we're still in the app, try to close the tab
    setTimeout(() => {
      if (window.location.href.includes(window.location.origin)) {
        window.close();
      }
    }, 500);
  }, []);

  const handleLogout = useCallback(() => {
    signOut();
  }, [signOut]);

  // Function to render the side menu item
  const SideMenuItem = useCallback(
    ({
      icon,
      label,
      isActive,
      onClick,
    }: {
      icon: ReactNode;
      label: string;
      isActive: boolean;
      onClick: () => void;
    }) => (
      <div
        className={`flex flex-col items-center justify-center py-3 cursor-pointer text-primary `}
        onClick={onClick}
      >
        <div className={`mb-1 ${isActive ? 'text-primary' : 'text-white'}`}>{icon}</div>
        <div className="text-xs">{label}</div>
      </div>
    ),
    []
  );

  return (
    <div className="w-16  flex flex-col items-center shrink-0 sticky top-0 h-screen">
      {/* Back button at the top */}
      <div className="w-full mt-4 mb-4">
        <SideMenuItem
          icon={<ArrowLeft size={20} />}
          label="Back"
          isActive={false}
          onClick={handleBackNavigation}
        />
      </div>

      <div className="flex-1 flex flex-col w-full">
        {menuItems.map((item) => (
          <SideMenuItem
            key={item.path}
            icon={item.icon}
            label={item.label}
            isActive={currentPath === item.path}
            onClick={() => {
              if (item.paramsPath && propertyId) {
                router.navigate({
                  to: item.paramsPath,
                  params: { propertyId },
                });
              } else {
                router.navigate({ to: item.path });
              }
            }}
          />
        ))}
      </div>

      {/* User button at bottom */}
      <div className="mt-auto mb-4 w-full">
        <SideMenuItem
          icon={<User size={20} />}
          label="Logout"
          isActive={false}
          onClick={handleLogout}
        />
      </div>
    </div>
  );
}
