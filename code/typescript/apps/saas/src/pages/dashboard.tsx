import { usePageTitle } from '@/hooks/use-page-title';
import { PortfolioDashboard } from '@/components/dashboard/portfolio-dashboard';

export function DashboardPage() {
  usePageTitle('Portfolio Dashboard');
  return <PortfolioDashboard />;
}
