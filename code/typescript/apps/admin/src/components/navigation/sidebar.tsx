import { useEffect, useState } from 'react';
import { Link } from '@tanstack/react-router';
import { Button } from '@/components/ui/button';
import {
  NavigationMenu,
  NavigationMenuItem,
  NavigationMenuList,
  navigationMenuTriggerStyle,
} from '@/components/ui/navigation-menu';
import { cn } from '@/lib/utils';

// Type definition for navigation items
export interface NavItem {
  title: string;
  href: string;
  icon: React.ReactNode;
}

interface SidebarProps {
  items: NavItem[];
  className?: string;
  defaultOpen?: boolean;
}

export function Sidebar({ items, className, defaultOpen = true }: SidebarProps) {
  // Use localStorage to persist the sidebar state
  const [isOpen, setIsOpen] = useState(() => {
    // Try to get the stored state from localStorage
    const storedState = localStorage.getItem('admin-sidebar-open');

    // If a stored state exists, use it; otherwise use the defaultOpen prop
    return storedState !== null ? storedState === 'true' : defaultOpen;
  });

  // Update localStorage whenever isOpen changes
  useEffect(() => {
    localStorage.setItem('admin-sidebar-open', isOpen.toString());
  }, [isOpen]);

  // Fixed width values for consistent sizing
  const openWidth = 'w-64';
  const closedWidth = 'w-16';

  return (
    <aside
      className={cn(
        'border-r border-border h-[calc(100vh-4rem)] flex flex-col transition-all duration-300',
        isOpen ? openWidth : closedWidth,
        className
      )}
    >
      <div className="p-4 flex justify-end">
        <Button
          variant="ghost"
          size="icon"
          onClick={() => setIsOpen(!isOpen)}
          aria-label={isOpen ? 'Close sidebar' : 'Open sidebar'}
        >
          <svg
            xmlns="http://www.w3.org/2000/svg"
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            strokeWidth="2"
            strokeLinecap="round"
            strokeLinejoin="round"
            className="h-5 w-5"
          >
            {isOpen ? <path d="m15 18-6-6 6-6" /> : <path d="m9 18 6-6-6-6" />}
          </svg>
        </Button>
      </div>

      <NavigationMenu className="max-w-none w-full justify-start p-4 block">
        <NavigationMenuList className="flex flex-col space-y-2 items-start">
          {items.map((item) => (
            <NavigationMenuItem key={item.href} className="w-full">
              <Link
                to={item.href}
                className={cn(
                  navigationMenuTriggerStyle,
                  'w-full justify-start gap-3',
                  !isOpen && 'px-2'
                )}
                activeProps={{ className: 'text-primary' }}
              >
                {item.icon}
                {isOpen && <span>{item.title}</span>}
              </Link>
            </NavigationMenuItem>
          ))}
        </NavigationMenuList>
      </NavigationMenu>
    </aside>
  );
}
