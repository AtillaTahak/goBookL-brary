const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';

export interface Book {
	id: number;
	title: string;
	author: string;
	year: number;
	genre?: string;
	isbn?: string;
	created_at: string;
	updated_at: string;
}

export interface User {
	id: number;
	username: string;
	email: string;
	role: string;
}

export interface LoginRequest {
	username: string;
	password: string;
}

export interface RegisterRequest {
	username: string;
	password: string;
	email: string;
}

export interface AuthResponse {
	token: string;
	user: User;
}

class ApiClient {
	private baseURL: string;

	constructor(baseURL: string) {
		this.baseURL = baseURL;
	}

	private getHeaders(includeAuth = false): HeadersInit {
		const headers: HeadersInit = {
			'Content-Type': 'application/json',
		};

		if (includeAuth) {
			const token = localStorage.getItem('token');
			if (token) {
				headers.Authorization = `Bearer ${token}`;
			}
		}

		return headers;
	}

	private async handleResponse<T>(response: Response): Promise<T> {
		if (!response.ok) {
			const error = await response.json().catch(() => ({ error: 'An error occurred' }));
			throw new Error(error.error || `HTTP ${response.status}`);
		}

		const contentType = response.headers.get('content-type');
		if (contentType && contentType.includes('application/json')) {
			return response.json();
		}

		return null as T;
	}

	// Authentication
	async login(credentials: LoginRequest): Promise<AuthResponse> {
		const response = await fetch(`${this.baseURL}/auth/login`, {
			method: 'POST',
			headers: this.getHeaders(),
			body: JSON.stringify(credentials),
		});

		return this.handleResponse<AuthResponse>(response);
	}

	async register(userData: RegisterRequest): Promise<{ message: string }> {
		const response = await fetch(`${this.baseURL}/auth/register`, {
			method: 'POST',
			headers: this.getHeaders(),
			body: JSON.stringify(userData),
		});

		return this.handleResponse<{ message: string }>(response);
	}

	// Books
	async getBooks(search?: string): Promise<Book[]> {
		const url = new URL(`${this.baseURL}/books`);
		if (search) {
			url.searchParams.append('search', search);
		}

		const response = await fetch(url.toString(), {
			headers: this.getHeaders(),
		});

		return this.handleResponse<Book[]>(response);
	}

	async getBook(id: number): Promise<Book> {
		const response = await fetch(`${this.baseURL}/books/${id}`, {
			headers: this.getHeaders(),
		});

		return this.handleResponse<Book>(response);
	}

	async createBook(book: Omit<Book, 'id' | 'created_at' | 'updated_at'>): Promise<Book> {
		const response = await fetch(`${this.baseURL}/books`, {
			method: 'POST',
			headers: this.getHeaders(true),
			body: JSON.stringify(book),
		});

		return this.handleResponse<Book>(response);
	}

	async updateBook(id: number, book: Partial<Omit<Book, 'id' | 'created_at' | 'updated_at'>>): Promise<Book> {
		const response = await fetch(`${this.baseURL}/books/${id}`, {
			method: 'PUT',
			headers: this.getHeaders(true),
			body: JSON.stringify(book),
		});

		return this.handleResponse<Book>(response);
	}

	async deleteBook(id: number): Promise<void> {
		const response = await fetch(`${this.baseURL}/books/${id}`, {
			method: 'DELETE',
			headers: this.getHeaders(true),
		});

		return this.handleResponse<void>(response);
	}
}

export const api = new ApiClient(API_BASE_URL);
