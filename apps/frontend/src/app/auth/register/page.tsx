'use client';

import { useState } from 'react';
import { useRouter } from 'next/navigation';
import { Button } from '../../../components/Button';
import { Input } from '../../../components/Input';
import { Card } from '../../../components/Card';
import { useAuth } from '../../../lib/auth';

export default function RegisterPage() {
	const [username, setUsername] = useState('');
	const [email, setEmail] = useState('');
	const [password, setPassword] = useState('');
	const [confirmPassword, setConfirmPassword] = useState('');
	const [error, setError] = useState('');
	const [isLoading, setIsLoading] = useState(false);

	const { register } = useAuth();
	const router = useRouter();

	const handleSubmit = async (e: React.FormEvent) => {
		e.preventDefault();
		setError('');

		if (password !== confirmPassword) {
			setError('Passwords do not match');
			return;
		}

		if (password.length < 6) {
			setError('Password must be at least 6 characters');
			return;
		}

		setIsLoading(true);

		try {
			await register(username, password, email);
			router.push('/auth/login?message=Registration successful');
		} catch (err) {
			setError(err instanceof Error ? err.message : 'Registration failed');
		} finally {
			setIsLoading(false);
		}
	};

	return (
		<div className="min-h-screen flex items-center justify-center bg-gray-50 py-12 px-4 sm:px-6 lg:px-8">
			<div className="max-w-md w-full">
				<Card title="Create Account">
					<form onSubmit={handleSubmit} className="space-y-4">
						<Input
							label="Username"
							type="text"
							value={username}
							onChange={(e) => setUsername(e.target.value)}
							required
							disabled={isLoading}
						/>

						<Input
							label="Email"
							type="email"
							value={email}
							onChange={(e) => setEmail(e.target.value)}
							required
							disabled={isLoading}
						/>

						<Input
							label="Password"
							type="password"
							value={password}
							onChange={(e) => setPassword(e.target.value)}
							required
							disabled={isLoading}
						/>

						<Input
							label="Confirm Password"
							type="password"
							value={confirmPassword}
							onChange={(e) => setConfirmPassword(e.target.value)}
							required
							disabled={isLoading}
						/>

						{error && (
							<div className="text-red-600 text-sm">{error}</div>
						)}

						<Button
							type="submit"
							disabled={isLoading}
							className="w-full"
						>
							{isLoading ? 'Creating Account...' : 'Create Account'}
						</Button>
					</form>

					<div className="mt-4 text-center">
						<p className="text-sm text-gray-600">
							Already have an account?{' '}
							<button
								onClick={() => router.push('/auth/login')}
								className="text-blue-600 hover:text-blue-500 font-medium"
							>
								Sign in
							</button>
						</p>
					</div>
				</Card>
			</div>
		</div>
	);
}
