import * as React from 'react';
import * as NavigationMenuPrimitive from '@radix-ui/react-navigation-menu';
import { cn } from '@/lib/utils';

type NavigationMenuProps = React.ComponentPropsWithoutRef<typeof NavigationMenuPrimitive.Root>;
type NavigationMenuRef = React.ElementRef<typeof NavigationMenuPrimitive.Root>;

const NavigationMenu = React.forwardRef<NavigationMenuRef, NavigationMenuProps>(
  ({ className, children, ...props }, ref) => {
    return (
      <NavigationMenuPrimitive.Root
        ref={ref}
        className={cn('relative z-10 flex max-w-max flex-1 items-center justify-center', className)}
        {...props}
      >
        {children}
      </NavigationMenuPrimitive.Root>
    );
  }
);
NavigationMenu.displayName = NavigationMenuPrimitive.Root.displayName;

type NavigationMenuListProps = React.ComponentPropsWithoutRef<typeof NavigationMenuPrimitive.List>;
type NavigationMenuListRef = React.ElementRef<typeof NavigationMenuPrimitive.List>;

const NavigationMenuList = React.forwardRef<NavigationMenuListRef, NavigationMenuListProps>(
  ({ className, ...props }, ref) => {
    return (
      <NavigationMenuPrimitive.List
        ref={ref}
        className={cn(
          'group flex flex-1 list-none items-center justify-center space-x-1',
          className
        )}
        {...props}
      />
    );
  }
);
NavigationMenuList.displayName = NavigationMenuPrimitive.List.displayName;

const NavigationMenuItem = NavigationMenuPrimitive.Item;

const navigationMenuTriggerStyle =
  'group inline-flex h-10 w-max items-center justify-center rounded-md bg-background px-4 py-2 text-sm font-medium transition-colors hover:bg-accent hover:text-accent-foreground focus:bg-accent focus:text-accent-foreground focus:outline-hidden disabled:pointer-events-none disabled:opacity-50 data-active:bg-accent/50 data-[state=open]:bg-accent/50';

const NavigationMenuLink = NavigationMenuPrimitive.Link;

export {
  NavigationMenu,
  NavigationMenuList,
  NavigationMenuItem,
  navigationMenuTriggerStyle,
  NavigationMenuLink,
};
