import { usePageTitle } from '@/hooks/use-page-title';

export function HomePage() {
  usePageTitle('Home Page');

  return (
    <div className="p-6">
      <h1 className="text-2xl font-bold mb-4">Home Page</h1>
      <p>This is the home page. If you can see this content, routing is working correctly.</p>
    </div>
  );
}
