'use client';

import { useState } from 'react';
import { useRouter } from 'next/navigation';
import { useAuth } from '../../../lib/auth';

export default function LoginPage() {
	const [username, setUsername] = useState('');
	const [password, setPassword] = useState('');
	const [error, setError] = useState('');
	const [isLoading, setIsLoading] = useState(false);

	const { login } = useAuth();
	const router = useRouter();

	const handleSubmit = async (e: React.FormEvent) => {
		e.preventDefault();
		setError('');
		setIsLoading(true);

		try {
			await login(username, password);
			router.push('/books');
		} catch (err) {
			setError(err instanceof Error ? err.message : 'Login failed');
		} finally {
			setIsLoading(false);
		}
	};

	return (
		<div className="min-h-screen flex items-center justify-center bg-gray-50 py-12 px-4 sm:px-6 lg:px-8">
			<div className="max-w-md w-full">
				<div className="bg-white rounded-lg shadow-md border border-gray-200 p-6">
					<h2 className="text-lg font-semibold text-gray-900 mb-6">Sign In to Book Library</h2>

					<form onSubmit={handleSubmit} className="space-y-4">
						<div className="space-y-1">
							<label className="block text-sm font-medium text-gray-700">Username</label>
							<input
								type="text"
								value={username}
								onChange={(e) => setUsername(e.target.value)}
								required
								disabled={isLoading}
								className="w-full px-3 py-2 border border-gray-300 rounded-lg shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent disabled:bg-gray-50 disabled:text-gray-500"
							/>
						</div>

						<div className="space-y-1">
							<label className="block text-sm font-medium text-gray-700">Password</label>
							<input
								type="password"
								value={password}
								onChange={(e) => setPassword(e.target.value)}
								required
								disabled={isLoading}
								className="w-full px-3 py-2 border border-gray-300 rounded-lg shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent disabled:bg-gray-50 disabled:text-gray-500"
							/>
						</div>

						{error && (
							<div className="text-red-600 text-sm">{error}</div>
						)}

						<button
							type="submit"
							disabled={isLoading}
							className="w-full px-4 py-2 rounded font-semibold transition-colors duration-200 focus:outline-none focus:ring-2 focus:ring-offset-2 bg-blue-600 hover:bg-blue-700 text-white focus:ring-blue-500 disabled:opacity-50 disabled:cursor-not-allowed"
						>
							{isLoading ? 'Signing in...' : 'Sign In'}
						</button>
					</form>

					<div className="mt-4 text-center">
						<p className="text-sm text-gray-600">
							Don&apos;t have an account?{' '}
							<button
								onClick={() => router.push('/auth/register')}
								className="text-blue-600 hover:text-blue-500 font-medium"
							>
								Sign up
							</button>
						</p>
					</div>

					<div className="mt-6 p-4 bg-gray-50 rounded-lg">
						<p className="text-sm text-gray-600 mb-2">Demo accounts:</p>
						<div className="text-xs text-gray-500 space-y-1">
							<div>Admin: admin / admin123</div>
							<div>User: user / user123</div>
						</div>
					</div>
				</div>
			</div>
		</div>
	);
}
