import { Link, Outlet, useLocation, useNavigate, useRouteContext } from '@tanstack/react-router';
import { Button } from '@/components/ui/button';
import {
  Breadcrumb,
  BreadcrumbItem,
  BreadcrumbList,
  BreadcrumbPage,
} from '@/components/ui/breadcrumb';
import { Sidebar, mainNavItems } from '@/components/navigation';
import { ThemeToggle } from '@/components/ui/theme-toggle';

const getBreadcrumbItems = (pathname: string) => {
  const parts = pathname.split('/').filter(Boolean);

  if (parts.length === 1 && parts[0] === 'admin') {
    return [{ label: 'Dashboard', path: '/admin' }];
  }

  return parts.map((part, index) => {
    const path = `/${parts.slice(0, index + 1).join('/')}`;
    return {
      label: part.charAt(0).toUpperCase() + part.slice(1),
      path,
    };
  });
};

export function AdminLayout() {
  const location = useLocation();
  const navigate = useNavigate();
  const { user } = useRouteContext({ from: '/admin' });
  const breadcrumbItems = getBreadcrumbItems(location.pathname);

  const handleSignOut = () => {
    localStorage.removeItem('@admin:user');
    navigate({ to: '/login', replace: true });
  };

  return (
    <div className="min-h-screen flex flex-col">
      <header className="h-16 border-b bg-card flex items-center justify-between px-6 sticky top-0 z-10">
        <Breadcrumb>
          <BreadcrumbList>
            {breadcrumbItems.map((item, index) => (
              <BreadcrumbItem key={item.path}>
                {index === breadcrumbItems.length - 1 ? (
                  <BreadcrumbPage>{item.label}</BreadcrumbPage>
                ) : (
                  <Link to={item.path as '/admin' | '/admin/users' | '/admin/settings'}>
                    {item.label}
                  </Link>
                )}
              </BreadcrumbItem>
            ))}
          </BreadcrumbList>
        </Breadcrumb>
        <div className="flex items-center gap-2">
          {user && <span className="text-sm mr-2">{user.name || user.email}</span>}
          <ThemeToggle />
          <Button variant="outline" size="sm" onClick={handleSignOut}>
            <svg
              xmlns="http://www.w3.org/2000/svg"
              viewBox="0 0 24 24"
              fill="none"
              stroke="currentColor"
              strokeWidth="2"
              strokeLinecap="round"
              strokeLinejoin="round"
              className="h-4 w-4 mr-2"
            >
              <path d="M9 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h4" />
              <polyline points="16 17 21 12 16 7" />
              <line x1="21" x2="9" y1="12" y2="12" />
            </svg>
            Sign out
          </Button>
        </div>
      </header>

      <div className="flex flex-1 overflow-hidden">
        <Sidebar items={mainNavItems} />

        <main className="flex-1 p-6 overflow-auto">
          <div className="mb-6">
            <Outlet />
          </div>
        </main>
      </div>
    </div>
  );
}
