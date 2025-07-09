// src/pages/unauthorized.tsx
import { Button } from '@am/commons/components/ui/button';
import { usePageTitle } from '@/hooks/use-page-title';

export function UnauthorizedPage() {
  usePageTitle('Unauthorized Access');

  const redirectToCustomerWebsite = () => {
    window.location.href = 'https://customer-website.com';
  };

  return (
    <div className="min-h-screen flex items-center justify-center bg-gray-50">
      <div className="max-w-md w-full p-8 bg-white shadow-lg rounded-lg">
        <div className="flex justify-center mb-6">{/* <Logo size={60} /> */}</div>

        <h1 className="text-2xl font-bold text-center text-gray-900 mb-4">Unauthorized Access</h1>

        <p className="text-gray-600 text-center mb-6">
          You do not have permission to access this page. Please log in on the customer website.
        </p>

        <div className="flex justify-center">
          <Button onClick={redirectToCustomerWebsite}>Return to customer website</Button>
        </div>
      </div>
    </div>
  );
}
