import { Card, CardContent, CardHeader, CardTitle } from '@am/commons/components/ui/card';
import {
  LineChart,
  Line,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  ResponsiveContainer,
} from 'recharts';
import { usePageTitle } from '@/hooks/use-page-title';
import { useComparablePropertiesQuery } from '@/services/property';
import { useProperty } from '@/contexts/property-context';
import { useRouter } from '@tanstack/react-router';
import { PropertyDetails } from '@am/commons';

export function HomePage() {
  const { data: availableProperties = [], isLoading } = useComparablePropertiesQuery('');
  const router = useRouter();

  const chooseProperty = (id: string) => {
    router.navigate({ to: '/property/$propertyId/detail', params: { propertyId: id } });
  };

  usePageTitle('Property Selector');

  if (isLoading) {
    return (
      <div className="flex items-center justify-center min-h-screen">
        <div className="animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-blue-500"></div>
      </div>
    );
  }

  return (
    <div className="container mx-auto px-4 py-8">
      <h1 className="text-3xl font-bold mb-8 text-gray-100">Select a Property (test page)</h1>
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
        {availableProperties.map((property: PropertyDetails) => (
          <Card
            key={property.id}
            className="group hover:shadow-lg transition-all duration-200 cursor-pointer hover:bg-accent/50"
          >
            <CardContent className="p-4">
              <button
                className="w-full text-left space-y-2 hover:bg-accent/50 p-3 rounded-md transition-colors"
                onClick={() => chooseProperty(property.id)}
              >
                <div className="font-medium text-foreground group-hover:text-accent-foreground transition-colors">
                  {property.address}
                </div>
                {property.price && (
                  <div className="text-sm text-muted-foreground group-hover:text-accent-foreground/80">
                    ${property.price.toLocaleString()}
                  </div>
                )}
              </button>
            </CardContent>
          </Card>
        ))}
      </div>
    </div>
  );
}
