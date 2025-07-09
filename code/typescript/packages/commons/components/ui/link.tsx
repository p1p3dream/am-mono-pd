import { Link as TanStackLink } from '@tanstack/react-router';
import { ComponentPropsWithoutRef, forwardRef } from 'react';
import { cn } from '@/lib/utils';

export interface LinkProps extends ComponentPropsWithoutRef<typeof TanStackLink> {
  className?: string;
}

export const Link = forwardRef<HTMLAnchorElement, LinkProps>(
  ({ className, children, ...props }, ref) => {
    return (
      <TanStackLink ref={ref} className={cn(className)} {...props}>
        {children}
      </TanStackLink>
    );
  }
);

Link.displayName = 'Link';
