import { ColumnDef } from '@tanstack/react-table';
import { ChevronDown, ChevronUp } from 'lucide-react';

/**
 * Creates a sortable column header with up/down icons
 */
export function createSortableHeader<T>(label: string) {
  return ({ column }: { column: any }) => (
    <button
      className="flex items-center text-left text-white"
      onClick={() => column.toggleSorting(column.getIsSorted() === 'asc')}
    >
      {label}
      {column.getIsSorted() === 'asc' ? (
        <ChevronUp className="ml-2 h-4 w-4" />
      ) : column.getIsSorted() === 'desc' ? (
        <ChevronDown className="ml-2 h-4 w-4" />
      ) : null}
    </button>
  );
}

/**
 * Creates a standard text cell with white text
 */
export function createTextCell<T>(key: keyof T) {
  return ({ row }: { row: any }) => <div className="text-white">{row.getValue(key)}</div>;
}

/**
 * Creates a status badge cell
 */
export function createStatusCell<T>(key: keyof T) {
  return ({ row }: { row: any }) => (
    <span className="bg-primary text-white text-xs px-2 py-1 rounded-full">
      {row.getValue(key)}
    </span>
  );
}

/**
 * Creates a checkbox cell
 */
export function createCheckboxCell<T>(key: keyof T, selectHandler?: (id: string) => void) {
  return ({ row }: { row: any }) => (
    <input
      type="checkbox"
      checked={row.original[key] as boolean}
      onChange={() => selectHandler && selectHandler(row.original.id)}
      onClick={(e) => e.stopPropagation()}
    />
  );
}

export function createFloatCell<T>(key: keyof T) {
  // Format the value to 2 decimal places
  return ({ row }: { row: any }) => (
    <div className="text-white">{row.getValue(key).toFixed(2)}</div>
  );
}

export function createCurrencyCell<T>(key: keyof T) {
  return ({ row }: { row: any }) => (
    <div className="text-white">
      {row.getValue(key).toLocaleString('en-US', { style: 'currency', currency: 'USD' })}
    </div>
  );
}
