import React from 'react';

interface DeleteConfirmModalProps {
	isOpen: boolean;
	onClose: () => void;
	onConfirm: () => void;
	title: string;
	message: string;
	isLoading?: boolean;
}

export const DeleteConfirmModal = ({
	isOpen,
	onClose,
	onConfirm,
	title,
	message,
	isLoading = false
}: DeleteConfirmModalProps) => {
	if (!isOpen) return null;

	return (
		<div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center p-4 z-50">
			<div className="bg-white rounded-lg max-w-md w-full p-6">
				<h3 className="text-lg font-semibold text-gray-900 mb-4">{title}</h3>
				<p className="text-gray-600 mb-6">{message}</p>

				<div className="flex gap-3">
					<button
						onClick={onClose}
						disabled={isLoading}
						className="flex-1 px-4 py-2 border border-gray-300 rounded-lg text-gray-700 hover:bg-gray-50 disabled:opacity-50"
					>
						Cancel
					</button>
					<button
						onClick={onConfirm}
						disabled={isLoading}
						className="flex-1 px-4 py-2 bg-red-600 text-white rounded-lg hover:bg-red-700 disabled:opacity-50"
					>
						{isLoading ? 'Deleting...' : 'Delete'}
					</button>
				</div>
			</div>
		</div>
	);
};
