import { usePageTitle } from '@/hooks/use-page-title';

export function DebugPage() {
  usePageTitle('Debug Page');

  return (
    <div className="p-6 bg-white m-6 rounded-lg shadow">
      <h1 className="text-2xl font-bold mb-4">Debug Page</h1>
      <p className="mb-4">If you can see this, routing is working properly!</p>

      <div className="p-4 bg-gray-100 rounded-lg mb-4">
        <h2 className="text-lg font-semibold mb-2">Current Path</h2>
        <code className="block bg-gray-800 text-white p-2 rounded">{window.location.pathname}</code>
      </div>

      <div className="p-4 bg-gray-100 rounded-lg">
        <h2 className="text-lg font-semibold mb-2">Browser Info</h2>
        <ul className="list-disc pl-6">
          <li>User Agent: {navigator.userAgent}</li>
          <li>
            Window Size: {window.innerWidth}x{window.innerHeight}
          </li>
        </ul>
      </div>
    </div>
  );
}
