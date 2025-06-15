'use client';

import { useEffect } from 'react';
import { useRouter } from 'next/navigation';
import { Button } from '../components/Button';
import { useAuth } from '../lib/auth';

export default function Home() {
  const { user, isLoading } = useAuth();
  const router = useRouter();

  useEffect(() => {
    if (!isLoading) {
      if (user) {
        router.push('/books');
      } else {
        router.push('/auth/login');
      }
    }
  }, [user, isLoading, router]);

  if (isLoading) {
    return (
      <div className="min-h-screen flex items-center justify-center bg-gray-50">
        <div className="text-gray-500">Loading...</div>
      </div>
    );
  }

  return (
    <div className="min-h-screen flex items-center justify-center bg-gray-50">
      <div className="text-center max-w-md">
        <h1 className="text-4xl font-bold text-gray-900 mb-4">📚</h1>
        <h1 className="text-3xl font-bold text-gray-900 mb-4">Book Library</h1>
        <p className="text-gray-600 mb-8">
          A modern library management system built with Go and Next.js
        </p>

        <div className="space-y-4">
          <Button
            onClick={() => router.push('/auth/login')}
            className="w-full"
          >
            Sign In
          </Button>
          <Button
            variant="secondary"
            onClick={() => router.push('/auth/register')}
            className="w-full"
          >
            Create Account
          </Button>
        </div>

        <div className="mt-8 p-4 bg-blue-50 rounded-lg">
          <p className="text-sm text-blue-800 mb-2">Demo Accounts:</p>
          <div className="text-xs text-blue-600 space-y-1">
            <div>Admin: admin / admin123</div>
            <div>User: user / user123</div>
          </div>
        </div>
      </div>
    </div>
  );
}
