// src/pages/auth/login.tsx
import { useState } from 'react';
import { useRouter } from '@tanstack/react-router';
import { Button } from '@am/commons/components/ui/button';
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from '@am/commons/components/ui/card';
import { Input } from '@am/commons/components/ui/input';
import { usePageTitle } from '@/hooks/use-page-title';
import { useAuth } from '@/contexts/auth-context';

type LoginFormData = {
  email: string;
  password: string;
};

export function LoginPage() {
  usePageTitle('Sign In');
  const router = useRouter();
  const { signIn, authState } = useAuth();
  const [formData, setFormData] = useState<LoginFormData>({
    email: '',
    password: '',
  });
  const [error, setError] = useState('');

  const isLoading = authState === 'loading';

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');

    try {
      const success = await signIn(formData.email, formData.password);

      if (success) {
        router.navigate({ to: '/property/cma', replace: true });
      } else {
        setError('Invalid email or password');
      }
    } catch (error) {
      console.error('Error during sign in:', error);
      setError('Authentication failed');
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
    <div className="min-h-screen flex items-center justify-center bg-linear-to-br from-primary/5 via-secondary/5 to-primary/5">
      <Card className="w-[350px]">
        <CardHeader className="space-y-3">
          <div className="flex justify-center">{/* <Logo size={60} /> */}</div>
          <CardTitle className="text-2xl text-center text-primary">Welcome Back</CardTitle>
          <CardDescription className="text-center">
            Enter your credentials to access
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
