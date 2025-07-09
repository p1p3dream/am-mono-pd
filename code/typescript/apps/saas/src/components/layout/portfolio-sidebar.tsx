import { ReactNode, useCallback } from 'react';
import { useRouter, useRouterState } from '@tanstack/react-router';
import {
  ChevronDown,
  Home,
  File,
  Building,
  Package,
  ListChecks,
  Upload,
  Save,
  Search,
  Construction,
  Layers,
} from 'lucide-react';
import { useAuth } from '@/contexts/auth-context';
import { usePortfolio } from '@/contexts/portfolio-context';

interface SideNavSectionProps {
  title: string;
  children: ReactNode;
  isExpanded?: boolean;
}

const SideNavSection = ({ title, children, isExpanded = true }: SideNavSectionProps) => {
  return (
    <div className="mb-2">
      <div className="flex items-center justify-between px-4 py-2 text-sm font-medium text-gray-500">
        <span>{title}</span>
        <ChevronDown size={16} />
      </div>
      {isExpanded && <div>{children}</div>}
    </div>
  );
};

interface SideNavItemProps {
  icon: ReactNode;
  label: string;
  path: string;
  isActive?: boolean;
  onClick?: () => void;
  disabled?: boolean;
}

const SideNavItem = ({
  icon,
  label,
  path,
  isActive = false,
  onClick,
  disabled = false,
}: SideNavItemProps) => {
  const router = useRouter();

  const handleClick = () => {
    if (disabled) return;

    if (onClick) {
      onClick();
    } else {
      router.navigate({ to: path });
    }
  };

  return (
    <div
      className={`flex items-center px-4 py-2 text-sm cursor-pointer hover:bg-gray-100 ${
        isActive ? 'bg-gray-100 font-medium' : ''
      } ${disabled ? 'opacity-50 cursor-not-allowed' : ''}`}
      onClick={handleClick}
    >
      <div className="mr-3 text-gray-600">{icon}</div>
      <div>{label}</div>
    </div>
  );
};

export function PortfolioSidebar() {
  const router = useRouter();
  const { location } = useRouterState();
  const currentPath = location.pathname;
  const { signOut } = useAuth();
  const { portfolios } = usePortfolio();

  const handleLogout = useCallback(() => {
    signOut();
  }, [signOut]);

  return (
    <div className="w-64 border-r border-gray-200 h-screen bg-white overflow-y-auto flex flex-col">
      <div className="p-4 border-b border-gray-200">
        <div className="flex items-center">
          <div className="font-bold text-xl">AbodeMine</div>
          <div className="ml-auto">
            <ChevronDown size={16} />
          </div>
        </div>
      </div>

      <div className="flex-1 overflow-y-auto">
        <SideNavItem
          icon={<Home size={18} />}
          label="Dashboard"
          path="/dashboard"
          isActive={currentPath === '/dashboard'}
        />

        <SideNavItem
          icon={<File size={18} />}
          label="Underwriting assumptions"
          path="/assumptions"
          isActive={currentPath === '/assumptions'}
          disabled={true}
        />

        <SideNavSection title="Portfolios">
          {portfolios.map((portfolio) => (
            <SideNavItem
              key={portfolio.id}
              icon={<Building size={18} />}
              label={portfolio.name}
              path={`/portfolio/${portfolio.id}`}
              isActive={currentPath.includes(`/portfolio/${portfolio.id}`)}
            />
          ))}
        </SideNavSection>

        <SideNavSection title="Acquisitions">
          <SideNavItem
            icon={<Package size={18} />}
            label="MLS buy boxes"
            path="/"
            isActive={false}
            disabled={true}
          />
          <SideNavItem
            icon={<ListChecks size={18} />}
            label="Marketing lists"
            path="/"
            isActive={false}
            disabled={true}
          />
          <SideNavItem
            icon={<Upload size={18} />}
            label="Uploaded lists"
            path="/"
            isActive={false}
            disabled={true}
          />
          <SideNavItem
            icon={<Save size={18} />}
            label="Saved homes"
            path="/"
            isActive={false}
            disabled={true}
          />
        </SideNavSection>

        <SideNavSection title="Market search">
          <SideNavItem
            icon={<Search size={18} />}
            label="Princeton, TX"
            path="/"
            isActive={false}
            disabled={true}
          />
          <SideNavItem
            icon={<Search size={18} />}
            label="DFW 2,2"
            path="/"
            isActive={false}
            disabled={true}
          />
          <SideNavItem
            icon={<Search size={18} />}
            label="Austin, TX"
            path="/"
            isActive={false}
            disabled={true}
          />
          <SideNavItem
            icon={<Search size={18} />}
            label="Buy box 4"
            path="/"
            isActive={false}
            disabled={true}
          />
          <SideNavItem
            icon={<Search size={18} />}
            label="Nashville"
            path="/"
            isActive={false}
            disabled={true}
          />
        </SideNavSection>

        <SideNavSection title="Built for rent">
          <SideNavItem
            icon={<Construction size={18} />}
            label="Projects"
            path="/"
            isActive={false}
            disabled={true}
          />
        </SideNavSection>
      </div>

      <div className="p-4 border-t border-gray-200">
        <SideNavItem
          icon={<Layers size={18} />}
          label="Invite teammates"
          path="/"
          onClick={() => alert('Invite teammates')}
        />
      </div>
    </div>
  );
}
