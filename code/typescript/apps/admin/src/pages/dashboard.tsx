import { Card, CardContent, CardHeader, CardTitle, CardDescription } from '@/components/ui/card';
import { Progress } from '@/components/ui/progress';
import { Home, DollarSign, Percent, CircleDollarSign, ArrowUpRight } from 'lucide-react';
import { ResponsiveContainer, BarChart, CartesianGrid, XAxis, YAxis, Tooltip, Bar } from 'recharts';

const data = [
  { month: 'Jan', revenue: 65000, occupancy: 89 },
  { month: 'Feb', revenue: 72000, occupancy: 91 },
  { month: 'Mar', revenue: 85000, occupancy: 94 },
  { month: 'Apr', revenue: 75000, occupancy: 90 },
  { month: 'May', revenue: 92000, occupancy: 93 },
  { month: 'Jun', revenue: 88000, occupancy: 92 },
];

export function Dashboard() {
  return (
    <div className="space-y-4">
      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Properties Listed</CardTitle>
            <Home className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">245</div>
            <div className="flex items-center text-xs text-green-500">
              <ArrowUpRight className="h-4 w-4" />
              <span>+12 new this month</span>
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Total Revenue</CardTitle>
            <DollarSign className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">$842,250</div>
            <div className="flex items-center text-xs text-green-500">
              <ArrowUpRight className="h-4 w-4" />
              <span>+23% from last month</span>
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Occupancy Rate</CardTitle>
            <Percent className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">92.4%</div>
            <Progress value={92.4} className="h-2" />
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Avg. Property Value</CardTitle>
            <CircleDollarSign className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">$425,000</div>
            <div className="flex items-center text-xs text-green-500">
              <ArrowUpRight className="h-4 w-4" />
              <span>+5.2% market value</span>
            </div>
          </CardContent>
        </Card>
      </div>

      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-7">
        <Card className="col-span-4">
          <CardHeader>
            <CardTitle>Property Performance</CardTitle>
            <CardDescription>Monthly revenue and occupancy trends</CardDescription>
          </CardHeader>
          <CardContent>
            <div className="h-[300px]">
              <ResponsiveContainer width="100%" height="100%">
                <BarChart data={data}>
                  <CartesianGrid strokeDasharray="3 3" />
                  <XAxis dataKey="month" />
                  <YAxis yAxisId="left" />
                  <YAxis yAxisId="right" orientation="right" />
                  <Tooltip />
                  <Bar yAxisId="left" dataKey="revenue" fill="#3b82f6" name="Revenue ($)" />
                  <Bar yAxisId="right" dataKey="occupancy" fill="#10b981" name="Occupancy (%)" />
                </BarChart>
              </ResponsiveContainer>
            </div>
          </CardContent>
        </Card>

        <Card className="col-span-3">
          <CardHeader>
            <CardTitle>Recent Properties</CardTitle>
            <CardDescription>Latest properties added to the system</CardDescription>
          </CardHeader>
          <CardContent>
            <div className="space-y-4">
              {[
                {
                  address: '742 Evergreen Terrace',
                  type: 'Residential',
                  price: '$450,000',
                  status: 'Available',
                },
                {
                  address: '123 Business Ave',
                  type: 'Commercial',
                  price: '$1,200,000',
                  status: 'Under Contract',
                },
                {
                  address: '456 Lake View Rd',
                  type: 'Residential',
                  price: '$380,000',
                  status: 'Available',
                },
                {
                  address: '789 Industrial Park',
                  type: 'Industrial',
                  price: '$2,500,000',
                  status: 'Available',
                },
              ].map((property, index) => (
                <div key={index} className="flex items-center justify-between">
                  <div className="space-y-1">
                    <p className="text-sm font-medium leading-none">{property.address}</p>
                    <p className="text-sm text-muted-foreground">{property.type}</p>
                    <p className="text-xs text-muted-foreground">{property.price}</p>
                  </div>
                  <span
                    className={`text-xs px-2 py-1 rounded-full ${
                      property.status === 'Available'
                        ? 'bg-green-100 text-green-700'
                        : 'bg-orange-100 text-orange-700'
                    }`}
                  >
                    {property.status}
                  </span>
                </div>
              ))}
            </div>
          </CardContent>
        </Card>
      </div>
    </div>
  );
}
