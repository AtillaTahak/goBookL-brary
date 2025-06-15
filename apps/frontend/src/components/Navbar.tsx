'use client';

import { useAuth } from '../lib/auth';
import { useRouter } from 'next/navigation';
import { Button } from './Button';

export const Navbar = () => {
	const { user, logout } = useAuth();
	const router = useRouter();

	const handleLogout = () => {
		logout();
		router.push('/auth/login');
	};

	return (
		<nav className="bg-white shadow-sm border-b border-gray-200">
			<div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
				<div className="flex justify-between items-center h-16">
					<div className="flex items-center">
						<h1 className="text-xl font-bold text-gray-900">📚 Book Library</h1>
					</div>

					<div className="flex items-center space-x-4">
						{user ? (
							<>
								<span className="text-sm text-gray-600">
									Welcome, <span className="font-medium">{user.username}</span>
									{user.role === 'admin' && (
										<span className="ml-1 px-2 py-1 text-xs bg-blue-100 text-blue-800 rounded">
											Admin
										</span>
									)}
								</span>
								<Button
									variant="secondary"
									onClick={handleLogout}
								>
									Logout
								</Button>
							</>
						) : (
							<div className="space-x-2">
								<Button
									variant="secondary"
									onClick={() => router.push('/auth/login')}
								>
									Login
								</Button>
								<Button
									onClick={() => router.push('/auth/register')}
								>
									Register
								</Button>
							</div>
						)}
					</div>
				</div>
			</div>
		</nav>
	);
};
