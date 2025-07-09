import { useMemo, useState, useCallback } from 'react';
import { ColumnDef } from '@tanstack/react-table';
import { ArrowUpDown, ChevronDown, Ellipsis, Check } from 'lucide-react';
import {
  DataTable,
  createSortableHeader,
  createTextCell,
  createStatusCell,
} from '@/components/data-table';
import { ComparableProperty } from '@/contexts/property-context';
import { createCurrencyCell, createFloatCell } from '@/components/data-table/columns';

// Define our custom column type
type CustomColumnDef<T> = ColumnDef<T> & {
  isSortable?: boolean;
  headerLabel?: string;
};

interface ComparablePropertiesListProps {
  properties: ComparableProperty[];
  onPropertySelect?: (id: string) => void;
  buttonLabel?: {
    selected: string;
    notSelected: string;
  };
  showExport?: boolean;
  showFilter?: boolean;
  rowClassName?: string;
  selectedRowClassName?: string;
  title?: string;
}

export function ComparablePropertiesList({
  properties,
  onPropertySelect = () => {},
  buttonLabel = { selected: 'Selected', notSelected: 'Select' },
  showExport = true,
  showFilter = false,
  rowClassName = 'border-gray-800 hover:bg-gray-800/50',
  selectedRowClassName = 'bg-primary/20 border-primary/10 hover:bg-primary/50',
  title = 'Comparable Properties',
}: ComparablePropertiesListProps) {
  const [selectAll, setSelectAll] = useState(false);
  const [selectedRows, setSelectedRows] = useState<string[]>([]);
  const [columnsDropdownOpen, setColumnsDropdownOpen] = useState(false);
  const [visibleColumns, setVisibleColumns] = useState<Record<string, boolean>>({
    address: true,
    sqft: true,
    distance: true,
    status: true,
    propertyType: true,
    yearBuilt: true,
    beds: true,
    baths: true,
    closeDate: true,
    amount: true,
  });

  // Toggle select all functionality
  const handleSelectAll = useCallback(() => {
    if (selectedRows.length > 0) {
      setSelectedRows([]);
      setSelectAll(false);
    } else {
      setSelectedRows(properties.map((property) => property.id));
      setSelectAll(true);
    }
  }, [selectedRows, properties]);

  const handleSelectProperty = useCallback(
    (id: string) => {
      onPropertySelect(id);
      properties.forEach((property) => {
        if (property.id === id) {
          if (selectedRows.includes(id)) {
            setSelectedRows(selectedRows.filter((row) => row !== id));
          } else {
            setSelectedRows([...selectedRows, id]);
          }
        }
      });
    },
    [onPropertySelect]
  );

  // Toggle column visibility
  const toggleColumnVisibility = useCallback((columnKey: string) => {
    setVisibleColumns((prev) => ({
      ...prev,
      [columnKey]: !prev[columnKey],
    }));
  }, []);

  // Column definitions for the DataTable
  const allColumns = useMemo<CustomColumnDef<ComparableProperty>[]>(
    () => [
      {
        id: 'select',
        header: () => <input type="checkbox" checked={selectAll} onChange={handleSelectAll} />,
        cell: ({ row }) => (
          <input
            type="checkbox"
            checked={selectedRows.includes(row.original.id)}
            onChange={() => handleSelectProperty(row.original.id)}
          />
        ),
      },
      {
        accessorKey: 'address',
        header: 'Address',
        headerLabel: 'Address',
        isSortable: true,
        cell: createTextCell('address'),
      },
      {
        accessorKey: 'sqft',
        header: 'Sq Ft',
        headerLabel: 'Square Feet',
        isSortable: true,
        cell: createTextCell('sqft'),
      },
      {
        accessorKey: 'distance',
        header: 'Distance',
        headerLabel: 'Distance',
        isSortable: true,
        cell: createFloatCell('distance'),
      },
      {
        accessorKey: 'status',
        header: 'Status',
        headerLabel: 'Status',
        isSortable: true,
        cell: createTextCell('status'),
      },
      {
        accessorKey: 'propertyType',
        header: 'Type',
        headerLabel: 'Property Type',
        isSortable: true,
        cell: createTextCell('propertyType'),
      },
      {
        accessorKey: 'yearBuilt',
        header: 'Built',
        headerLabel: 'Year Built',
        isSortable: true,
        cell: createTextCell('yearBuilt'),
      },
      {
        accessorKey: 'beds',
        header: 'Bedroom',
        headerLabel: 'Bedrooms',
        isSortable: true,
        cell: createTextCell('beds'),
      },
      {
        accessorKey: 'baths',
        header: 'Bathroom',
        headerLabel: 'Bathrooms',
        isSortable: true,
        cell: createTextCell('baths'),
      },
      {
        accessorKey: 'closeDate',
        header: 'Close Date',
        headerLabel: 'Closing Date',
        isSortable: true,
        cell: createTextCell('closeDate'),
      },
      {
        accessorKey: 'price',
        header: 'Amount',
        headerLabel: 'Price Amount',
        isSortable: true,
        cell: createCurrencyCell('price'),
      },
      {
        id: 'actions',
        header: () => <div className="text-gray-400"></div>,
        cell: ({ row }) => {
          const property = row.original;
          return (
            <button
              className={`px-3 py-1 rounded text-xs `}
              onClick={() => onPropertySelect(property.id)}
            >
              <Ellipsis size={14} />
            </button>
          );
        },
      },
    ],
    [onPropertySelect, handleSelectAll, selectedRows]
  );

  // Filter columns based on visibility
  const columns = useMemo(() => {
    return allColumns.filter((column) => {
      // Handle accessorKey properly for TypeScript
      if ('accessorKey' in column) {
        const key = column.accessorKey as string;
        return visibleColumns[key] !== false;
      }
      return true; // Keep columns without accessorKey (like select and actions)
    });
  }, [allColumns, visibleColumns]);

  return (
    <div className="p-6">
      <div className="flex justify-between items-center mb-4">
        <div className="flex items-center space-x-3">
          <button className="flex items-center bg-gray-800 hover:bg-gray-700 text-white px-3 py-2 rounded-md text-sm">
            <span className="mr-2">+</span>
            <span>Save current list</span>
          </button>
          <button className="flex items-center bg-gray-800 hover:bg-gray-700 text-white px-3 py-2 rounded-md text-sm">
            <span>Default list</span>
            <ChevronDown size={14} className="ml-2" />
          </button>
        </div>
        <div className="flex items-center relative">
          <button
            className="flex items-center bg-gray-800 hover:bg-gray-700 text-white px-3 py-2 rounded-md text-sm"
            onClick={() => setColumnsDropdownOpen(!columnsDropdownOpen)}
          >
            <span>Columns</span>
            <ChevronDown size={14} className="ml-2" />
          </button>

          {/* Columns Dropdown Menu */}
          {columnsDropdownOpen && (
            <div className="absolute top-full right-0 mt-1 w-48 bg-gray-800 rounded-md shadow-lg z-50">
              <div className="p-2">
                {allColumns
                  .filter(
                    (col): col is CustomColumnDef<ComparableProperty> =>
                      'accessorKey' in col && !!col.headerLabel
                  )
                  .map((column) => {
                    // Use type assertion to access accessorKey
                    const accessorKey =
                      'accessorKey' in column ? (column.accessorKey as string) : '';
                    return (
                      <button
                        key={accessorKey}
                        className="flex items-center justify-between w-full px-3 py-2 text-sm text-white hover:bg-gray-700 rounded-md"
                        onClick={() => toggleColumnVisibility(accessorKey)}
                      >
                        <span>{column.headerLabel}</span>
                        {visibleColumns[accessorKey] && <Check size={14} />}
                      </button>
                    );
                  })}
              </div>
            </div>
          )}
        </div>
      </div>

      <div className="bg-primary-foreground">
        <DataTable
          columns={columns}
          data={properties}
          rowClassName={rowClassName}
          selectedRowClassName={selectedRowClassName}
          selectRow={onPropertySelect}
        />
      </div>
    </div>
  );
}
