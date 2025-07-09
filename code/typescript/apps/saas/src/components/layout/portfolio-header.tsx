import { Search, Bell, HelpCircle, ChevronDown } from 'lucide-react';
import { usePortfolio } from '@/contexts/portfolio-context';

export function PortfolioHeader() {
  const { portfolio } = usePortfolio();
  const userInitial = 'A'; // This would come from auth context

  return (
    <div className="flex items-center justify-between px-6 py-2 bg-white border-b border-gray-200">
      <div className="relative w-full max-w-lg">
        <div className="absolute inset-y-0 left-0 flex items-center pl-3 pointer-events-none">
          <Search className="w-5 h-5 text-gray-400" />
        </div>
        <input
          type="text"
          className="block w-full py-2 pl-10 pr-3 bg-gray-100 border border-gray-300 rounded-md focus:outline-none focus:ring-blue-500 focus:border-blue-500"
          placeholder="Search"
        />
      </div>

      <div className="flex items-center space-x-4">
        <Bell className="w-5 h-5 text-gray-600" />
        <HelpCircle className="w-5 h-5 text-gray-600" />
        <div className="flex items-center">
          <div className="flex items-center justify-center w-8 h-8 bg-blue-600 rounded-full text-white font-medium mr-1">
            {userInitial}
          </div>
          <ChevronDown className="w-4 h-4 text-gray-600" />
        </div>
      </div>
    </div>
  );
}
