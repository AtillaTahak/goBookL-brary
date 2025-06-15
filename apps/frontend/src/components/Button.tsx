import React from "react";

type ButtonProps = {
	children: React.ReactNode;
	onClick?: () => void;
	variant?: 'primary' | 'secondary';
	disabled?: boolean;
	type?: 'button' | 'submit' | 'reset';
	className?: string;
};

export const Button = ({
	children,
	onClick,
	variant = 'primary',
	disabled = false,
	type = 'button',
	className = ''
}: ButtonProps) => {
	const base = 'px-4 py-2 rounded font-semibold transition-colors duration-200 focus:outline-none focus:ring-2 focus:ring-offset-2';
	const styles = {
		primary: 'bg-blue-600 hover:bg-blue-700 text-white focus:ring-blue-500',
		secondary: 'bg-gray-200 hover:bg-gray-300 text-gray-900 focus:ring-gray-500'
	};

	return (
		<button
			type={type}
			onClick={onClick}
			disabled={disabled}
			className={`${base} ${styles[variant]} ${disabled ? 'opacity-50 cursor-not-allowed' : ''} ${className}`}
		>
			{children}
		</button>
	);
};
