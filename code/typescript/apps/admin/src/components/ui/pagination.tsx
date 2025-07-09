import * as React from 'react';
import { ChevronLeft, ChevronRight, MoreHorizontal } from 'lucide-react';
import { cn } from '@/lib/utils';
import { Button } from '@/components/ui/button';

type TPaginationProps = {
  total: number;
  page: number;
  limit: number;
  onChange: (page: number) => void;
};

type TPaginationButtonElement = HTMLButtonElement;
type TPaginationButtonProps = React.ButtonHTMLAttributes<HTMLButtonElement>;

const PaginationButton = React.forwardRef<TPaginationButtonElement, TPaginationButtonProps>(
  ({ className, ...props }, ref) => (
    <Button
      ref={ref}
      variant="outline"
      size="sm"
      className={cn('w-8 h-8 p-0', className)}
      {...props}
    />
  )
);

export function Pagination({ total, page, limit, onChange }: TPaginationProps) {
  const totalPages = Math.ceil(total / limit);

  const generatePages = () => {
    const pages = [];
    const showEllipsis = totalPages > 7;

    if (!showEllipsis) {
      for (let i = 1; i <= totalPages; i++) {
        pages.push(i);
      }
      return pages;
    }

    // Always show first page
    pages.push(1);

    if (page > 3) {
      pages.push('ellipsis-start');
    }

    // Show pages around current page
    for (let i = Math.max(2, page - 1); i <= Math.min(totalPages - 1, page + 1); i++) {
      pages.push(i);
    }

    if (page < totalPages - 2) {
      pages.push('ellipsis-end');
    }

    // Always show last page
    if (totalPages > 1) {
      pages.push(totalPages);
    }

    return pages;
  };

  return (
    <div className="flex items-center justify-center space-x-2">
      <PaginationButton onClick={() => onChange(page - 1)} disabled={page === 1}>
        <ChevronLeft className="h-4 w-4" />
      </PaginationButton>

      {generatePages().map((pageNumber) => {
        if (pageNumber === 'ellipsis-start' || pageNumber === 'ellipsis-end') {
          return <MoreHorizontal key={pageNumber} className="h-4 w-4" />;
        }

        return (
          <PaginationButton
            key={pageNumber}
            onClick={() => onChange(pageNumber as number)}
            className={cn(
              pageNumber === page && 'bg-primary text-primary-foreground hover:bg-primary/90'
            )}
          >
            {pageNumber}
          </PaginationButton>
        );
      })}

      <PaginationButton onClick={() => onChange(page + 1)} disabled={page === totalPages}>
        <ChevronRight className="h-4 w-4" />
      </PaginationButton>
    </div>
  );
}
