'use client';

import { createContext, useContext, useEffect, useState, ReactNode } from 'react';
import { api, User } from './api';

interface AuthContextType {
	user: User | null;
	token: string | null;
	login: (username: string, password: string) => Promise<void>;
	register: (username: string, password: string, email: string) => Promise<void>;
	logout: () => void;
	isLoading: boolean;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export const useAuth = () => {
	const context = useContext(AuthContext);
	if (context === undefined) {
		throw new Error('useAuth must be used within an AuthProvider');
	}
	return context;
};

interface AuthProviderProps {
	children: ReactNode;
}

export const AuthProvider = ({ children }: AuthProviderProps) => {
	const [user, setUser] = useState<User | null>(null);
	const [token, setToken] = useState<string | null>(null);
	const [isLoading, setIsLoading] = useState(true);

	useEffect(() => {
		// Check for stored token on initialization
		const storedToken = localStorage.getItem('token');
		const storedUser = localStorage.getItem('user');

		if (storedToken && storedUser) {
			setToken(storedToken);
			setUser(JSON.parse(storedUser));
		}

		setIsLoading(false);
	}, []);

	const login = async (username: string, password: string) => {
		try {
			const response = await api.login({ username, password });

			setToken(response.token);
			setUser(response.user);

			localStorage.setItem('token', response.token);
			localStorage.setItem('user', JSON.stringify(response.user));
		} catch (error) {
			throw error;
		}
	};

	const register = async (username: string, password: string, email: string) => {
		try {
			await api.register({ username, password, email });
		} catch (error) {
			throw error;
		}
	};

	const logout = () => {
		setUser(null);
		setToken(null);
		localStorage.removeItem('token');
		localStorage.removeItem('user');
	};

	const value = {
		user,
		token,
		login,
		register,
		logout,
		isLoading,
	};

	return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
};
