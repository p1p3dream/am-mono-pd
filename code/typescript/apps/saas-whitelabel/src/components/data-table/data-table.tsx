import React, { useState } from 'react';
import {
  ColumnDef,
  SortingState,
  flexRender,
  getCoreRowModel,
  getSortedRowModel,
  useReactTable,
  PaginationState,
  getPaginationRowModel,
  SortDirection,
} from '@tanstack/react-table';
import {
  Table,
  TableBody,
  TableCaption,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@am/commons/components/ui/table';
import { ScrollArea } from '@am/commons/components/ui/scroll-area';
import { cn } from '@/lib/utils';
import { Button } from '@am/commons/components/ui/button';
import {
  ChevronLeft,
  ChevronRight,
  ChevronsLeft,
  ChevronsRight,
  ArrowUp,
  ArrowDown,
} from 'lucide-react';

// Extend the ColumnDef type to include isSortable
declare module '@tanstack/react-table' {
  interface ColumnMeta<TData extends unknown, TValue> {
    isSortable?: boolean;
  }
}

interface DataTableProps<TData, TValue> {
  columns: ColumnDef<TData, TValue>[];
  data: TData[];
  className?: string;
  caption?: string;
  rowClassName?: string;
  selectedRowClassName?: string;
  scrollAreaClassName?: string;
  withScrollArea?: boolean;
  selectRow?: (id: string) => void;
  pageSize?: number;
  selectedRows?: string[];
}

export function DataTable<TData, TValue>({
  columns,
  data,
  className,
  caption,
  rowClassName = 'border-gray-800 hover:bg-gray-800/50',
  selectedRowClassName = 'bg-primary border-primary/80 hover:bg-primary/90',
  scrollAreaClassName,
  withScrollArea = true,
  selectRow,
  selectedRows,
  pageSize = 10,
}: DataTableProps<TData, TValue>) {
  const [sorting, setSorting] = useState<SortingState>([]);
  const [pagination, setPagination] = useState<PaginationState>({
    pageIndex: 0,
    pageSize,
  });

  const table = useReactTable({
    data,
    columns,
    state: {
      sorting,
      pagination,
    },
    onSortingChange: setSorting,
    onPaginationChange: setPagination,
    getCoreRowModel: getCoreRowModel(),
    getSortedRowModel: getSortedRowModel(),
    getPaginationRowModel: getPaginationRowModel(),
  });

  // Function to render sort indicator
  const renderSortIndicator = (sorted: false | SortDirection) => {
    if (!sorted) return null;
    return sorted === 'asc' ? (
      <ArrowUp className="ml-1 h-4 w-4" />
    ) : (
      <ArrowDown className="ml-1 h-4 w-4" />
    );
  };

  const tableContent = (
    <Table>
      {caption && <TableCaption>{caption}</TableCaption>}
      <TableHeader>
        {table.getHeaderGroups().map((headerGroup) => (
          <TableRow key={headerGroup.id} className={rowClassName}>
            {headerGroup.headers.map((header) => {
              // Get the meta information from the column to check if it's sortable
              const isSortable = (header.column.columnDef as any)?.isSortable;

              return (
                <TableHead
                  key={header.id}
                  className={cn(
                    isSortable ? 'cursor-pointer select-none' : '',
                    header.column.getCanSort() ? 'cursor-pointer select-none' : ''
                  )}
                  onClick={isSortable ? () => header.column.toggleSorting() : undefined}
                >
                  <div className="flex items-center">
                    {header.isPlaceholder
                      ? null
                      : flexRender(header.column.columnDef.header, header.getContext())}

                    {/* Show sort indicator if column is sortable */}
                    {isSortable && renderSortIndicator(header.column.getIsSorted())}
                  </div>
                </TableHead>
              );
            })}
          </TableRow>
        ))}
      </TableHeader>
      <TableBody>
        {table.getRowModel().rows.length ? (
          table.getRowModel().rows.map((row) => {
            const isSelected = selectedRows?.includes((row.original as any)?.id);
            const rowId = (row.original as any)?.id;
            return (
              <TableRow
                key={row.id}
                className={isSelected ? selectedRowClassName : rowClassName}
                onClick={() => selectRow && rowId && selectRow(rowId)}
                style={{ cursor: selectRow ? 'pointer' : 'default' }}
              >
                {row.getVisibleCells().map((cell) => (
                  <TableCell key={cell.id}>
                    {flexRender(cell.column.columnDef.cell, cell.getContext())}
                  </TableCell>
                ))}
              </TableRow>
            );
          })
        ) : (
          <TableRow>
            <TableCell colSpan={columns.length} className="h-24 text-center text-gray-400">
              No results.
            </TableCell>
          </TableRow>
        )}
      </TableBody>
    </Table>
  );

  return (
    <div className={className}>
      {withScrollArea ? (
        <ScrollArea className={scrollAreaClassName}>{tableContent}</ScrollArea>
      ) : (
        tableContent
      )}

      {/* Pagination Controls */}
      <div className="flex items-center justify-between space-x-2 py-4">
        <div className="flex-1 text-sm text-gray-400">
          {table.getFilteredRowModel().rows.length > 0 && (
            <>
              Showing{' '}
              <strong>
                {table.getState().pagination.pageIndex * table.getState().pagination.pageSize + 1}-
                {Math.min(
                  (table.getState().pagination.pageIndex + 1) *
                    table.getState().pagination.pageSize,
                  table.getFilteredRowModel().rows.length
                )}
              </strong>{' '}
              of <strong>{table.getFilteredRowModel().rows.length}</strong> items
            </>
          )}
        </div>
        <div className="flex items-center space-x-2">
          <Button
            variant="outline"
            size="sm"
            onClick={() => table.setPageIndex(0)}
            disabled={!table.getCanPreviousPage()}
            className="hidden sm:flex"
          >
            <ChevronsLeft className="h-4 w-4" />
          </Button>
          <Button
            variant="outline"
            size="sm"
            onClick={() => table.previousPage()}
            disabled={!table.getCanPreviousPage()}
          >
            <ChevronLeft className="h-4 w-4" />
          </Button>
          <div className="flex items-center gap-1">
            <span className="text-sm font-medium">Page</span>
            <strong className="text-sm">
              {table.getState().pagination.pageIndex + 1} of {table.getPageCount()}
            </strong>
          </div>
          <Button
            variant="outline"
            size="sm"
            onClick={() => table.nextPage()}
            disabled={!table.getCanNextPage()}
          >
            <ChevronRight className="h-4 w-4" />
          </Button>
          <Button
            variant="outline"
            size="sm"
            onClick={() => table.setPageIndex(table.getPageCount() - 1)}
            disabled={!table.getCanNextPage()}
            className="hidden sm:flex"
          >
            <ChevronsRight className="h-4 w-4" />
          </Button>
        </div>
      </div>
    </div>
  );
}
