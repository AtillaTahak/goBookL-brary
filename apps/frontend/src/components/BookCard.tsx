import React from 'react';
import { Book } from '../lib/api';
import { Card } from './Card';
import { Button } from './Button';

interface BookCardProps {
	book: Book;
	onEdit?: (book: Book) => void;
	onDelete?: (book: Book) => void;
	showActions?: boolean;
}

export const BookCard = ({ book, onEdit, onDelete, showActions = false }: BookCardProps) => {
	return (
		<Card className="h-full">
			<div className="flex flex-col h-full">
				<div className="flex-1">
					<h3 className="text-lg font-semibold text-gray-900 mb-2">{book.title}</h3>
					<p className="text-gray-600 mb-1">by {book.author}</p>
					<p className="text-gray-500 mb-2">Published: {book.year}</p>
					{book.genre && (
						<p className="text-sm text-gray-500 mb-2">Genre: {book.genre}</p>
					)}
					{book.isbn && (
						<p className="text-xs text-gray-400">ISBN: {book.isbn}</p>
					)}
				</div>

				{showActions && (onEdit || onDelete) && (
					<div className="flex gap-2 mt-4 pt-4 border-t border-gray-200">
						{onEdit && (
							<Button
								variant="secondary"
								onClick={() => onEdit(book)}
								className="flex-1"
							>
								Edit
							</Button>
						)}
						{onDelete && (
							<Button
								variant="secondary"
								onClick={() => onDelete(book)}
								className="flex-1 bg-red-50 text-red-600 hover:bg-red-100"
							>
								Delete
							</Button>
						)}
					</div>
				)}
			</div>
		</Card>
	);
};
