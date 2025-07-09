import { useState } from 'react';
import { useNavigate } from '@tanstack/react-router';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Input } from '@/components/ui/input';
import { usePageTitle } from '@/hooks/use-page-title';

type LoginFormData = {
  email: string;
  password: string;
};

// Mock credentials for demonstration - in a real app, this would be handled by an API
const VALID_CREDENTIALS = {
  email: 'admin@abodemine.com',
  password: '123',
};

export function Login() {
  usePageTitle('Admin Sign In');
  const navigate = useNavigate();
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState('');

  const [formData, setFormData] = useState<LoginFormData>({
    email: '',
    password: '',
  });

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');
    setIsLoading(true);

    try {
      // Check if credentials match (this would be an API call in a real app)
      if (
        formData.email === VALID_CREDENTIALS.email &&
        formData.password === VALID_CREDENTIALS.password
      ) {
        // Store user data in localStorage
        localStorage.setItem(
          '@admin:user',
          JSON.stringify({
            email: formData.email,
            name: 'Admin User',
          })
        );

        // Navigate to admin dashboard using TanStack Router
        navigate({ to: '/admin' });
      } else {
        setError('Invalid email or password');
      }
    } catch (error) {
      console.error('Login error:', error);
      setError('An error occurred. Please try again.');
    } finally {
      setIsLoading(false);
    }
  };

  const handleInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { name, value } = e.target;
    setFormData((prev) => ({
      ...prev,
      [name]: value,
    }));
  };

  return (
    <div className="min-h-screen flex items-center justify-center bg-background">
      <Card className="w-[350px] border border-border">
        <CardHeader className="space-y-3">
          <div className="flex justify-center">{/* <Logo size={60} /> */}</div>
          <CardTitle className="text-2xl text-center text-primary">Admin Access</CardTitle>
          <CardDescription className="text-center">
            Enter your credentials to access the admin panel
          </CardDescription>
        </CardHeader>
        <CardContent>
          <form onSubmit={handleSubmit} className="space-y-4">
            {error && <div className="text-sm text-destructive text-center">{error}</div>}
            <div className="space-y-2">
              <Input
                type="email"
                name="email"
                placeholder="Email address"
                value={formData.email}
                onChange={handleInputChange}
                required
                disabled={isLoading}
              />
            </div>
            <div className="space-y-2">
              <Input
                type="password"
                name="password"
                placeholder="Password"
                value={formData.password}
                onChange={handleInputChange}
                required
                disabled={isLoading}
              />
            </div>
            <Button
              type="submit"
              className="w-full bg-primary hover:bg-primary/90"
              disabled={isLoading}
            >
              {isLoading ? 'Signing in...' : 'Sign in'}
            </Button>
          </form>
        </CardContent>
      </Card>
    </div>
  );
}
