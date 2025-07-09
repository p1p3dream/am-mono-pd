import { useMemo, useState, useEffect } from 'react';
import { PortfolioHeader } from '../layout/portfolio-header';
import { usePortfolio } from '@/contexts/portfolio-context';
import { CheckCircle, Filter, Mail, Download, Plus } from 'lucide-react';
import { formatCurrency, formatPercentage } from '../../lib/formatters';
import { useParams } from '@tanstack/react-router';

interface PropertyTableItem {
  id: string;
  address: string;
  beds: number;
  baths: number;
  sqft: number;
  rent: number;
  yield: number;
  capRate: number;
  leaseEnd: string;
}

export function PortfolioDetail() {
  const { portfolioId = '' } = useParams({ strict: false });
  const { portfolio, selectPortfolio } = usePortfolio();
  const [selectedRows, setSelectedRows] = useState<string[]>([]);

  useEffect(() => {
    if (portfolioId) {
      selectPortfolio(portfolioId);
    }
  }, [portfolioId, selectPortfolio]);

  const mockProperties: PropertyTableItem[] = useMemo(
    () => [
      {
        id: '1',
        address: '8080 Railroad St.',
        beds: 2,
        baths: 2,
        sqft: 1942,
        rent: 2000,
        yield: 11.376,
        capRate: 6.532,
        leaseEnd: '4/29/2025',
      },
      {
        id: '2',
        address: '7529 E. Pecan St.',
        beds: 2,
        baths: 2,
        sqft: 1942,
        rent: 2000,
        yield: 11.376,
        capRate: 6.532,
        leaseEnd: '4/29/2025',
      },
      {
        id: '3',
        address: '3890 Poplar Dr.',
        beds: 2,
        baths: 2,
        sqft: 1942,
        rent: 2000,
        yield: 11.376,
        capRate: 6.532,
        leaseEnd: '4/29/2025',
      },
      {
        id: '4',
        address: '7529 E. Pecan St.',
        beds: 2,
        baths: 2,
        sqft: 1942,
        rent: 2000,
        yield: 11.376,
        capRate: 6.532,
        leaseEnd: '4/29/2025',
      },
      {
        id: '5',
        address: '3605 Parker Rd.',
        beds: 2,
        baths: 2,
        sqft: 1942,
        rent: 2000,
        yield: 11.376,
        capRate: 6.532,
        leaseEnd: '4/29/2025',
      },
      {
        id: '6',
        address: '775 Rolling Green Rd.',
        beds: 2,
        baths: 2,
        sqft: 1942,
        rent: 2000,
        yield: 11.376,
        capRate: 6.532,
        leaseEnd: '4/29/2025',
      },
      {
        id: '7',
        address: '8558 Green Rd.',
        beds: 2,
        baths: 2,
        sqft: 1942,
        rent: 2000,
        yield: 11.376,
        capRate: 6.532,
        leaseEnd: '4/29/2025',
      },
      {
        id: '8',
        address: '775 Rolling Green Rd.',
        beds: 2,
        baths: 2,
        sqft: 1942,
        rent: 2000,
        yield: 11.376,
        capRate: 3.75,
        leaseEnd: '4/29/2025',
      },
      {
        id: '9',
        address: '7529 E. Pecan St.',
        beds: 2,
        baths: 2,
        sqft: 1942,
        rent: 2000,
        yield: 11.376,
        capRate: 3.75,
        leaseEnd: '4/29/2025',
      },
      {
        id: '10',
        address: '8558 Green Rd.',
        beds: 2,
        baths: 2,
        sqft: 1942,
        rent: 2000,
        yield: 11.376,
        capRate: 3.75,
        leaseEnd: '4/29/2025',
      },
      {
        id: '11',
        address: '3890 Poplar Dr.',
        beds: 2,
        baths: 2,
        sqft: 1942,
        rent: 2000,
        yield: 11.376,
        capRate: 3.75,
        leaseEnd: '4/29/2025',
      },
      {
        id: '12',
        address: '8080 Railroad St.',
        beds: 2,
        baths: 2,
        sqft: 1942,
        rent: 2000,
        yield: 11.376,
        capRate: 3.75,
        leaseEnd: '4/29/2025',
      },
      {
        id: '13',
        address: '3605 Parker Rd.',
        beds: 2,
        baths: 2,
        sqft: 1942,
        rent: 2000,
        yield: 11.376,
        capRate: 3.75,
        leaseEnd: '4/29/2025',
      },
    ],
    []
  );

  const toggleRowSelection = (id: string) => {
    setSelectedRows((prev) =>
      prev.includes(id) ? prev.filter((rowId) => rowId !== id) : [...prev, id]
    );
  };

  const toggleAllRows = () => {
    if (selectedRows.length === mockProperties.length) {
      setSelectedRows([]);
    } else {
      setSelectedRows(mockProperties.map((p) => p.id));
    }
  };

  if (!portfolio) {
    return <div>No portfolio selected</div>;
  }

  return (
    <div className="flex flex-col h-full">
      <PortfolioHeader />

      <div className="p-6 flex flex-col h-full">
        <h1 className="text-2xl font-bold mb-6">{portfolio.name}</h1>

        {/* Stats grid */}
        <div className="grid grid-cols-2 lg:grid-cols-4 gap-4 mb-6">
          <StatCard title="Total properties" value={portfolio.totalProperties.toString()} />
          <StatCard title="MSA" value={portfolio.msa} />
          <StatCard
            title="Total purchase price"
            value={formatCurrency(portfolio.totalPurchasePrice)}
          />
          <StatCard title="UW gross yield" value={formatPercentage(portfolio.uwGrossYield)} />
          <StatCard title="Avg PPU" value={formatCurrency(portfolio.avgPPU)} />
          <StatCard title="Avg rent" value={formatCurrency(portfolio.avgRent)} />
          <StatCard title="Today's portfolio value" value={formatCurrency(portfolio.totalValue)} />
          <StatCard title="UW net cap rate" value={formatPercentage(portfolio.uwNetCapRate)} />
        </div>

        {/* Table actions */}
        <div className="flex justify-between mb-4">
          <div className="flex space-x-2">
            <button className="p-2 text-gray-600 rounded hover:bg-gray-100">
              <CheckCircle size={20} />
            </button>
            <button className="flex items-center px-3 py-1 text-gray-600 border border-gray-300 rounded hover:bg-gray-100">
              <Filter size={16} className="mr-1" />
              <span>Filters</span>
            </button>
            <button className="flex items-center px-3 py-1 text-gray-600 border border-gray-300 rounded hover:bg-gray-100">
              Columns
            </button>
          </div>
          <div className="flex space-x-2">
            <button className="flex items-center px-3 py-1 text-gray-600 border border-gray-300 rounded hover:bg-gray-100">
              <Mail size={16} className="mr-1" />
              <span>Share</span>
            </button>
            <button className="flex items-center px-3 py-1 text-gray-600 border border-gray-300 rounded hover:bg-gray-100">
              <Download size={16} className="mr-1" />
              <span>Export</span>
            </button>
            <button className="flex items-center px-3 py-1 text-gray-600 border border-gray-300 rounded hover:bg-gray-100">
              <Plus size={16} className="mr-1" />
              <span>Add</span>
            </button>
          </div>
        </div>

        {/* Properties table */}
        <div className="border border-gray-200 rounded-md overflow-hidden flex-1 bg-white">
          <table className="min-w-full divide-y divide-gray-200">
            <thead className="bg-gray-50">
              <tr>
                <th className="px-4 py-3 text-left">
                  <input
                    type="checkbox"
                    className="h-4 w-4"
                    checked={selectedRows.length === mockProperties.length}
                    onChange={toggleAllRows}
                  />
                </th>
                <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  Address
                </th>
                <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  Beds
                </th>
                <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  Baths
                </th>
                <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  Sq ft
                </th>
                <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  Rent
                </th>
                <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  Yield
                </th>
                <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  Cap rate
                </th>
                <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  Lease end
                </th>
              </tr>
            </thead>
            <tbody className="bg-white divide-y divide-gray-200">
              {mockProperties.map((property) => (
                <tr key={property.id} className="hover:bg-gray-50">
                  <td className="px-4 py-3">
                    <input
                      type="checkbox"
                      className="h-4 w-4"
                      checked={selectedRows.includes(property.id)}
                      onChange={() => toggleRowSelection(property.id)}
                    />
                  </td>
                  <td className="px-4 py-3 text-sm text-gray-900">{property.address}</td>
                  <td className="px-4 py-3 text-sm text-gray-900">{property.beds}</td>
                  <td className="px-4 py-3 text-sm text-gray-900">{property.baths}</td>
                  <td className="px-4 py-3 text-sm text-gray-900">
                    {property.sqft.toLocaleString()}
                  </td>
                  <td className="px-4 py-3 text-sm text-gray-900">
                    {formatCurrency(property.rent)}
                  </td>
                  <td className="px-4 py-3 text-sm text-gray-900 text-green-600">
                    {formatPercentage(property.yield)}
                  </td>
                  <td className="px-4 py-3 text-sm text-gray-900 text-green-600">
                    {formatPercentage(property.capRate)}
                  </td>
                  <td className="px-4 py-3 text-sm text-gray-900">{property.leaseEnd}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </div>
    </div>
  );
}

interface StatCardProps {
  title: string;
  value: string;
}

function StatCard({ title, value }: StatCardProps) {
  return (
    <div className="bg-white p-4 rounded-md border border-gray-200">
      <div className="text-sm text-gray-500 mb-1">{title}</div>
      <div className="text-lg font-semibold">{value}</div>
    </div>
  );
}
