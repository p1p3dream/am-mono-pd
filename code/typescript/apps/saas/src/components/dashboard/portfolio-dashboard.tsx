import { useMemo } from 'react';
import {
  LineChart,
  Line,
  BarChart,
  Bar,
  PieChart,
  Pie,
  Cell,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  Legend,
  ResponsiveContainer,
  TooltipProps,
} from 'recharts';
import { usePortfolio } from '@/contexts/portfolio-context';
import { formatCurrency, formatPercentage } from '../../lib/formatters';
import { PortfolioHeader } from '../layout/portfolio-header';

const COLORS = ['#0088FE', '#00C49F', '#FFBB28', '#FF8042', '#8884d8', '#82ca9d'];

export function PortfolioDashboard() {
  const { portfolios } = usePortfolio();

  const portfolioData = useMemo(() => {
    return portfolios.map((portfolio) => ({
      name: portfolio.name.length > 15 ? portfolio.name.substring(0, 15) + '...' : portfolio.name,
      value: portfolio.totalValue,
      properties: portfolio.totalProperties,
      yield: portfolio.uwGrossYield,
      capRate: portfolio.uwNetCapRate,
      totalPurchase: portfolio.totalPurchasePrice,
    }));
  }, [portfolios]);

  const monthlyPerformanceData = useMemo(() => {
    // Mock data for monthly performance chart
    const months = [
      'Jan',
      'Feb',
      'Mar',
      'Apr',
      'May',
      'Jun',
      'Jul',
      'Aug',
      'Sep',
      'Oct',
      'Nov',
      'Dec',
    ];
    return months.map((month, index) => {
      const baseValue = 100 + index * 2;
      return {
        name: month,
        'Princeton, TX': baseValue + Math.random() * 10,
        'DFW 89': baseValue - 5 + Math.random() * 8,
      };
    });
  }, []);

  const propertyOccupancyData = useMemo(() => {
    // Mock data for occupancy stats
    return [
      { name: 'Occupied', value: 195 },
      { name: 'Vacant', value: 26 },
    ];
  }, []);

  const yieldComparisonData = useMemo(() => {
    return portfolios.map((portfolio) => ({
      name: portfolio.name.length > 10 ? portfolio.name.substring(0, 10) + '...' : portfolio.name,
      'Gross Yield': portfolio.uwGrossYield,
      'Net Cap Rate': portfolio.uwNetCapRate,
    }));
  }, [portfolios]);

  // Calculate total portfolio stats
  const totalProperties = useMemo(
    () => portfolios.reduce((sum, p) => sum + p.totalProperties, 0),
    [portfolios]
  );
  const totalValue = useMemo(
    () => portfolios.reduce((sum, p) => sum + p.totalValue, 0),
    [portfolios]
  );
  const avgYield = useMemo(
    () => portfolios.reduce((sum, p) => sum + p.uwGrossYield, 0) / portfolios.length,
    [portfolios]
  );
  const avgCapRate = useMemo(
    () => portfolios.reduce((sum, p) => sum + p.uwNetCapRate, 0) / portfolios.length,
    [portfolios]
  );

  const formatTooltipValue = (value: any, name: string, props: any) => {
    if (typeof value === 'number') {
      return [`${value.toFixed(2)}%`, 'Performance'];
    }
    return [value, 'Performance'];
  };

  const formatPortfolioValue = (value: any) => {
    if (typeof value === 'number') {
      return [formatCurrency(value), 'Portfolio Value'];
    }
    return [value, 'Portfolio Value'];
  };

  const formatPercentTooltip = (value: any) => {
    if (typeof value === 'number') {
      return [`${value.toFixed(3)}%`, ''];
    }
    return [value, ''];
  };

  return (
    <div className="flex flex-col h-full">
      <PortfolioHeader />

      <div className="p-6 overflow-auto">
        <h1 className="text-2xl font-bold mb-6">Portfolio Dashboard</h1>

        {/* Summary Stats */}
        <div className="grid grid-cols-2 md:grid-cols-4 gap-4 mb-8">
          <div className="bg-white p-4 rounded-lg border border-gray-200 shadow-sm">
            <div className="text-sm text-gray-500 mb-1">Total Properties</div>
            <div className="text-2xl font-bold">{totalProperties}</div>
          </div>
          <div className="bg-white p-4 rounded-lg border border-gray-200 shadow-sm">
            <div className="text-sm text-gray-500 mb-1">Total Portfolio Value</div>
            <div className="text-2xl font-bold">{formatCurrency(totalValue)}</div>
          </div>
          <div className="bg-white p-4 rounded-lg border border-gray-200 shadow-sm">
            <div className="text-sm text-gray-500 mb-1">Average Gross Yield</div>
            <div className="text-2xl font-bold text-green-600">{formatPercentage(avgYield)}</div>
          </div>
          <div className="bg-white p-4 rounded-lg border border-gray-200 shadow-sm">
            <div className="text-sm text-gray-500 mb-1">Average Cap Rate</div>
            <div className="text-2xl font-bold text-green-600">{formatPercentage(avgCapRate)}</div>
          </div>
        </div>

        <div className="grid grid-cols-1 lg:grid-cols-2 gap-6 mb-6">
          {/* Monthly Performance Chart */}
          <div className="bg-white p-4 rounded-lg border border-gray-200 shadow-sm">
            <h2 className="text-lg font-semibold mb-4">Monthly Performance</h2>
            <div className="h-80">
              <ResponsiveContainer width="100%" height="100%">
                <LineChart data={monthlyPerformanceData}>
                  <CartesianGrid strokeDasharray="3 3" />
                  <XAxis dataKey="name" />
                  <YAxis />
                  <Tooltip formatter={formatTooltipValue} />
                  <Legend />
                  <Line
                    type="monotone"
                    dataKey="Princeton, TX"
                    stroke="#8884d8"
                    strokeWidth={2}
                    activeDot={{ r: 8 }}
                  />
                  <Line type="monotone" dataKey="DFW 89" stroke="#82ca9d" strokeWidth={2} />
                </LineChart>
              </ResponsiveContainer>
            </div>
          </div>

          {/* Portfolio Distribution */}
          <div className="bg-white p-4 rounded-lg border border-gray-200 shadow-sm">
            <h2 className="text-lg font-semibold mb-4">Portfolio Distribution</h2>
            <div className="h-80">
              <ResponsiveContainer width="100%" height="100%">
                <PieChart>
                  <Pie
                    data={portfolioData}
                    dataKey="value"
                    nameKey="name"
                    cx="50%"
                    cy="50%"
                    outerRadius={80}
                    label={({ name, percent }) => `${name} ${(percent * 100).toFixed(0)}%`}
                  >
                    {portfolioData.map((entry, index) => (
                      <Cell key={`cell-${index}`} fill={COLORS[index % COLORS.length]} />
                    ))}
                  </Pie>
                  <Tooltip formatter={formatPortfolioValue} />
                  <Legend />
                </PieChart>
              </ResponsiveContainer>
            </div>
          </div>
        </div>

        <div className="grid grid-cols-1 lg:grid-cols-2 gap-6 mb-6">
          {/* Property Status */}
          <div className="bg-white p-4 rounded-lg border border-gray-200 shadow-sm">
            <h2 className="text-lg font-semibold mb-4">Property Occupancy</h2>
            <div className="h-80">
              <ResponsiveContainer width="100%" height="100%">
                <PieChart>
                  <Pie
                    data={propertyOccupancyData}
                    dataKey="value"
                    nameKey="name"
                    cx="50%"
                    cy="50%"
                    outerRadius={80}
                    label={({ name, percent }) => `${name} ${(percent * 100).toFixed(0)}%`}
                  >
                    <Cell fill="#4CAF50" />
                    <Cell fill="#FFC107" />
                  </Pie>
                  <Tooltip formatter={(value) => [value, 'Properties']} />
                  <Legend />
                </PieChart>
              </ResponsiveContainer>
            </div>
          </div>

          {/* Yield Comparison */}
          <div className="bg-white p-4 rounded-lg border border-gray-200 shadow-sm">
            <h2 className="text-lg font-semibold mb-4">Yield Comparison</h2>
            <div className="h-80">
              <ResponsiveContainer width="100%" height="100%">
                <BarChart data={yieldComparisonData}>
                  <CartesianGrid strokeDasharray="3 3" />
                  <XAxis dataKey="name" />
                  <YAxis domain={[0, 15]} />
                  <Tooltip formatter={formatPercentTooltip} />
                  <Legend />
                  <Bar dataKey="Gross Yield" fill="#8884d8" />
                  <Bar dataKey="Net Cap Rate" fill="#82ca9d" />
                </BarChart>
              </ResponsiveContainer>
            </div>
          </div>
        </div>

        <div className="bg-white p-4 rounded-lg border border-gray-200 shadow-sm mb-6">
          <h2 className="text-lg font-semibold mb-4">Portfolio Properties</h2>
          <div className="h-80">
            <ResponsiveContainer width="100%" height="100%">
              <BarChart data={portfolioData}>
                <CartesianGrid strokeDasharray="3 3" />
                <XAxis dataKey="name" />
                <YAxis />
                <Tooltip formatter={(value) => [value, 'Properties']} />
                <Legend />
                <Bar dataKey="properties" fill="#8884d8" name="Number of Properties" />
              </BarChart>
            </ResponsiveContainer>
          </div>
        </div>
      </div>
    </div>
  );
}
