import { usePortfolio } from '@/contexts/portfolio-context';
import { useRouter } from '@tanstack/react-router';
import { Building } from 'lucide-react';
import { PortfolioHeader } from '../layout/portfolio-header';

export function PortfolioSelector() {
  const { portfolios } = usePortfolio();
  const router = useRouter();

  return (
    <div className="flex flex-col h-full">
      <PortfolioHeader />

      <div className="p-8 max-w-6xl mx-auto w-full">
        <h1 className="text-2xl font-bold mb-6">Select a Portfolio</h1>

        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          {portfolios.map((portfolio) => (
            <div
              key={portfolio.id}
              className="bg-white p-6 rounded-lg border border-gray-200 shadow-sm hover:shadow-md transition-shadow cursor-pointer"
              onClick={() => {
                router.navigate({
                  to: '/portfolio/$portfolioId',
                  params: { portfolioId: portfolio.id },
                });
              }}
            >
              <div className="flex items-center gap-3 mb-4">
                <div className="p-2 bg-blue-100 text-blue-700 rounded-full">
                  <Building size={20} />
                </div>
                <h3 className="font-medium text-lg">{portfolio.name}</h3>
              </div>

              <div className="space-y-3">
                <div className="grid grid-cols-2 gap-2">
                  <div className="bg-gray-50 p-2 rounded">
                    <div className="text-xs text-gray-500">Properties</div>
                    <div className="font-semibold">{portfolio.totalProperties}</div>
                  </div>
                  <div className="bg-gray-50 p-2 rounded">
                    <div className="text-xs text-gray-500">MSA</div>
                    <div className="font-semibold">{portfolio.msa}</div>
                  </div>
                </div>

                <div className="grid grid-cols-2 gap-2">
                  <div className="bg-gray-50 p-2 rounded">
                    <div className="text-xs text-gray-500">Gross Yield</div>
                    <div className="font-semibold text-green-600">
                      {portfolio.uwGrossYield.toFixed(3)}%
                    </div>
                  </div>
                  <div className="bg-gray-50 p-2 rounded">
                    <div className="text-xs text-gray-500">Net Cap Rate</div>
                    <div className="font-semibold text-green-600">
                      {portfolio.uwNetCapRate.toFixed(3)}%
                    </div>
                  </div>
                </div>

                <div className="bg-gray-50 p-2 rounded">
                  <div className="text-xs text-gray-500">Total Value</div>
                  <div className="font-semibold">
                    ${portfolio.totalValue.toLocaleString(undefined, { maximumFractionDigits: 0 })}
                  </div>
                </div>
              </div>
            </div>
          ))}
        </div>
      </div>
    </div>
  );
}
