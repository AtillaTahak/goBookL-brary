import React, { useState, useEffect } from 'react';
import { Book } from '../lib/api';
import { Button } from './Button';
import { Input } from './Input';
import { Card } from './Card';

interface BookModalProps {
	book?: Book;
	isOpen: boolean;
	onClose: () => void;
	onSave: (book: Omit<Book, 'id' | 'created_at' | 'updated_at'>) => void;
	isLoading?: boolean;
}

export const BookModal = ({ book, isOpen, onClose, onSave, isLoading = false }: BookModalProps) => {
	const [title, setTitle] = useState('');
	const [author, setAuthor] = useState('');
	const [year, setYear] = useState<number>(new Date().getFullYear());
	const [genre, setGenre] = useState('');
	const [isbn, setIsbn] = useState('');

	useEffect(() => {
		if (book) {
			setTitle(book.title);
			setAuthor(book.author);
			setYear(book.year);
			setGenre(book.genre || '');
			setIsbn(book.isbn || '');
		} else {
			setTitle('');
			setAuthor('');
			setYear(new Date().getFullYear());
			setGenre('');
			setIsbn('');
		}
	}, [book, isOpen]);

	const handleSubmit = (e: React.FormEvent) => {
		e.preventDefault();
		onSave({
			title: title.trim(),
			author: author.trim(),
			year,
			genre: genre.trim() || undefined,
			isbn: isbn.trim() || undefined,
		});
	};

	if (!isOpen) return null;

	return (
		<div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center p-4 z-50">
			<div className="bg-white rounded-lg max-w-md w-full max-h-[90vh] overflow-y-auto">
				<Card title={book ? 'Edit Book' : 'Add New Book'}>
					<form onSubmit={handleSubmit} className="space-y-4">
						<Input
							label="Title *"
							type="text"
							value={title}
							onChange={(e) => setTitle(e.target.value)}
							required
							disabled={isLoading}
						/>

						<Input
							label="Author *"
							type="text"
							value={author}
							onChange={(e) => setAuthor(e.target.value)}
							required
							disabled={isLoading}
						/>

						<Input
							label="Year *"
							type="number"
							value={year}
							onChange={(e) => setYear(parseInt(e.target.value))}
							required
							min="1000"
							max={new Date().getFullYear() + 10}
							disabled={isLoading}
						/>

						<Input
							label="Genre"
							type="text"
							value={genre}
							onChange={(e) => setGenre(e.target.value)}
							disabled={isLoading}
						/>

						<Input
							label="ISBN"
							type="text"
							value={isbn}
							onChange={(e) => setIsbn(e.target.value)}
							disabled={isLoading}
						/>

						<div className="flex gap-3 pt-4">
							<Button
								type="button"
								variant="secondary"
								onClick={onClose}
								disabled={isLoading}
								className="flex-1"
							>
								Cancel
							</Button>
							<Button
								type="submit"
								disabled={isLoading || !title.trim() || !author.trim()}
								className="flex-1"
							>
								{isLoading ? 'Saving...' : book ? 'Update' : 'Add Book'}
							</Button>
						</div>
					</form>
				</Card>
			</div>
		</div>
	);
};
