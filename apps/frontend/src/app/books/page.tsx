'use client';

import { useState, useEffect } from 'react';
import { useRouter } from 'next/navigation';
import { api, Book } from '../../lib/api';
import { useAuth } from '../../lib/auth';
import { Button } from '../../components/Button';
import { Input } from '../../components/Input';
import { BookCard } from '../../components/BookCard';
import { BookModal } from '../../components/BookModal';
import { DeleteConfirmModal } from '../../components/DeleteConfirmModal';
import { Navbar } from '../../components/Navbar';

export default function BooksPage() {
	const [books, setBooks] = useState<Book[]>([]);
	const [filteredBooks, setFilteredBooks] = useState<Book[]>([]);
	const [searchTerm, setSearchTerm] = useState('');
	const [isLoading, setIsLoading] = useState(true);
	const [error, setError] = useState('');

	// Modal states
	const [isBookModalOpen, setIsBookModalOpen] = useState(false);
	const [editingBook, setEditingBook] = useState<Book | undefined>();
	const [isBookModalLoading, setIsBookModalLoading] = useState(false);

	// Delete modal states
	const [isDeleteModalOpen, setIsDeleteModalOpen] = useState(false);
	const [deletingBook, setDeletingBook] = useState<Book | undefined>();
	const [isDeleteLoading, setIsDeleteLoading] = useState(false);

	const { user } = useAuth();
	const router = useRouter();

	useEffect(() => {
		loadBooks();
	}, []);

	useEffect(() => {
		// Filter books based on search term
		if (searchTerm.trim() === '') {
			setFilteredBooks(books);
		} else {
			const filtered = books.filter(book =>
				book.title.toLowerCase().includes(searchTerm.toLowerCase()) ||
				book.author.toLowerCase().includes(searchTerm.toLowerCase()) ||
				book.genre?.toLowerCase().includes(searchTerm.toLowerCase())
			);
			setFilteredBooks(filtered);
		}
	}, [books, searchTerm]);

	const loadBooks = async () => {
		try {
			setIsLoading(true);
			setError('');
			const data = await api.getBooks();
			setBooks(data);
		} catch (err) {
			setError(err instanceof Error ? err.message : 'Failed to load books');
		} finally {
			setIsLoading(false);
		}
	};

	const handleAddBook = () => {
		setEditingBook(undefined);
		setIsBookModalOpen(true);
	};

	const handleEditBook = (book: Book) => {
		setEditingBook(book);
		setIsBookModalOpen(true);
	};

	const handleSaveBook = async (bookData: Omit<Book, 'id' | 'created_at' | 'updated_at'>) => {
		try {
			setIsBookModalLoading(true);
			setError('');

			if (editingBook) {
				await api.updateBook(editingBook.id, bookData);
			} else {
				await api.createBook(bookData);
			}

			setIsBookModalOpen(false);
			setEditingBook(undefined);
			await loadBooks();
		} catch (err) {
			setError(err instanceof Error ? err.message : 'Failed to save book');
		} finally {
			setIsBookModalLoading(false);
		}
	};

	const handleDeleteBook = (book: Book) => {
		setDeletingBook(book);
		setIsDeleteModalOpen(true);
	};

	const confirmDelete = async () => {
		if (!deletingBook) return;

		try {
			setIsDeleteLoading(true);
			setError('');
			await api.deleteBook(deletingBook.id);
			setIsDeleteModalOpen(false);
			setDeletingBook(undefined);
			await loadBooks();
		} catch (err) {
			setError(err instanceof Error ? err.message : 'Failed to delete book');
		} finally {
			setIsDeleteLoading(false);
		}
	};

	if (!user) {
		router.push('/auth/login');
		return null;
	}

	return (
		<div className="min-h-screen bg-gray-50">
			<Navbar />

			<div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
				{/* Header */}
				<div className="flex flex-col sm:flex-row sm:items-center sm:justify-between mb-8">
					<div>
						<h1 className="text-3xl font-bold text-gray-900">Book Library</h1>
						<p className="text-gray-600 mt-1">
							{books.length} book{books.length !== 1 ? 's' : ''} in collection
						</p>
					</div>

					{user && (
						<Button onClick={handleAddBook} className="mt-4 sm:mt-0">
							Add New Book
						</Button>
					)}
				</div>

				{/* Search */}
				<div className="mb-6">
					<Input
						placeholder="Search books by title, author, or genre..."
						value={searchTerm}
						onChange={(e) => setSearchTerm(e.target.value)}
						className="max-w-md"
					/>
				</div>

				{/* Error message */}
				{error && (
					<div className="mb-6 p-4 bg-red-50 border border-red-200 rounded-lg text-red-700">
						{error}
					</div>
				)}

				{/* Loading state */}
				{isLoading ? (
					<div className="flex justify-center items-center py-12">
						<div className="text-gray-500">Loading books...</div>
					</div>
				) : (
					<>
						{/* Books grid */}
						{filteredBooks.length > 0 ? (
							<div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-6">
								{filteredBooks.map((book) => (
									<BookCard
										key={book.id}
										book={book}
										onEdit={user ? handleEditBook : undefined}
										onDelete={user ? handleDeleteBook : undefined}
										showActions={!!user}
									/>
								))}
							</div>
						) : (
							<div className="text-center py-12">
								<div className="text-gray-500 mb-4">
									{searchTerm ? 'No books found matching your search.' : 'No books in the library yet.'}
								</div>
								{user && !searchTerm && (
									<Button onClick={handleAddBook}>
										Add First Book
									</Button>
								)}
							</div>
						)}
					</>
				)}
			</div>

			{/* Book Modal */}
			<BookModal
				book={editingBook}
				isOpen={isBookModalOpen}
				onClose={() => {
					setIsBookModalOpen(false);
					setEditingBook(undefined);
				}}
				onSave={handleSaveBook}
				isLoading={isBookModalLoading}
			/>

			{/* Delete Confirmation Modal */}
			<DeleteConfirmModal
				isOpen={isDeleteModalOpen}
				onClose={() => {
					setIsDeleteModalOpen(false);
					setDeletingBook(undefined);
				}}
				onConfirm={confirmDelete}
				title="Delete Book"
				message={`Are you sure you want to delete "${deletingBook?.title}"? This action cannot be undone.`}
				isLoading={isDeleteLoading}
			/>
		</div>
	);
}
