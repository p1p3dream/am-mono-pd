import * as React from 'react';
import { ChevronRight } from 'lucide-react';
import { cn } from '@/lib/utils';

type TNavElement = HTMLElement;
type TNavProps = React.HTMLAttributes<HTMLElement>;

const Breadcrumb = React.forwardRef<TNavElement, TNavProps>((props, ref) => (
  <nav ref={ref} aria-label="breadcrumb" {...props} />
));

type TOlElement = HTMLOListElement;
type TOlProps = React.OlHTMLAttributes<HTMLOListElement>;

const BreadcrumbList = React.forwardRef<TOlElement, TOlProps>(({ className, ...props }, ref) => (
  <ol
    ref={ref}
    className={cn(
      'flex flex-wrap items-center gap-1.5 break-words text-sm text-muted-foreground',
      className
    )}
    {...props}
  />
));

type TLiElement = HTMLLIElement;
type TLiProps = React.LiHTMLAttributes<HTMLLIElement>;

const BreadcrumbItem = React.forwardRef<TLiElement, TLiProps>(({ className, ...props }, ref) => (
  <li ref={ref} className={cn('inline-flex items-center gap-1.5', className)} {...props} />
));

type TAnchorElement = HTMLAnchorElement;
type TAnchorProps = React.AnchorHTMLAttributes<HTMLAnchorElement>;

const BreadcrumbLink = React.forwardRef<TAnchorElement, TAnchorProps>(
  ({ className, ...props }, ref) => (
    <a ref={ref} className={cn('hover:text-foreground transition-colors', className)} {...props} />
  )
);

type TSpanElement = HTMLSpanElement;
type TSpanProps = React.HTMLAttributes<HTMLSpanElement>;

const BreadcrumbPage = React.forwardRef<TSpanElement, TSpanProps>(
  ({ className, ...props }, ref) => (
    <span
      ref={ref}
      role="link"
      aria-disabled="true"
      aria-current="page"
      className={cn('font-normal text-foreground', className)}
      {...props}
    />
  )
);

type TSeparatorProps = React.HTMLAttributes<HTMLSpanElement>;

const BreadcrumbSeparator = ({ className, children, ...props }: TSeparatorProps) => (
  <span
    role="presentation"
    aria-hidden="true"
    className={cn('[&>svg]:size-3.5', className)}
    {...props}
  >
    {children || <ChevronRight />}
  </span>
);

export {
  Breadcrumb,
  BreadcrumbList,
  BreadcrumbItem,
  BreadcrumbLink,
  BreadcrumbPage,
  BreadcrumbSeparator,
};
