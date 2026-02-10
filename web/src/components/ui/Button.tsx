import React from 'react';
import { cn } from '@/lib/utils';
import { Loader2 } from 'lucide-react';

interface ButtonProps extends React.ButtonHTMLAttributes<HTMLButtonElement> {
  variant?: 'primary' | 'secondary' | 'outline' | 'ghost' | 'danger';
  size?: 'sm' | 'md' | 'lg' | 'xl';
  isLoading?: boolean;
  leftIcon?: React.ReactNode;
  rightIcon?: React.ReactNode;
}

export const Button = React.forwardRef<HTMLButtonElement, ButtonProps>(
  ({ className, variant = 'primary', size = 'md', isLoading, leftIcon, rightIcon, children, disabled, ...props }, ref) => {
    const variants = {
      primary: 'bg-gradient-to-r from-[var(--primary)] to-purple-600 text-white shadow-lg shadow-[var(--primary)]/20 hover:shadow-[var(--primary)]/40 active:scale-[0.98]',
      secondary: 'bg-[var(--background-tertiary)] text-[var(--foreground)] hover:bg-[var(--background-secondary)] active:scale-[0.98]',
      outline: 'bg-transparent border border-[var(--border)] text-[var(--foreground)] hover:border-[var(--primary)]/50 active:scale-[0.98]',
      ghost: 'bg-transparent text-[var(--foreground-muted)] hover:text-[var(--foreground)] hover:bg-[var(--background-tertiary)] active:scale-[0.98]',
      danger: 'bg-red-500/10 text-red-500 hover:bg-red-500 hover:text-white active:scale-[0.98]',
    };

    const sizes = {
      sm: 'h-9 px-4 text-xs rounded-xl',
      md: 'h-11 px-6 text-sm rounded-xl',
      lg: 'h-13 px-8 text-base rounded-2xl font-bold',
      xl: 'h-14 px-10 text-lg rounded-2xl font-black tracking-tight',
    };

    return (
      <button
        ref={ref}
        disabled={isLoading || disabled}
        className={cn(
          'inline-flex items-center justify-center gap-2 transition-all duration-200 focus:outline-none focus:ring-2 focus:ring-[var(--primary)]/50 disabled:opacity-50 disabled:cursor-not-allowed disabled:active:scale-100',
          variants[variant],
          sizes[size],
          className
        )}
        {...props}
      >
        {isLoading && <Loader2 className="w-4 h-4 animate-spin shrink-0" />}
        {!isLoading && leftIcon && <span className="shrink-0">{leftIcon}</span>}
        {children}
        {!isLoading && rightIcon && <span className="shrink-0">{rightIcon}</span>}
      </button>
    );
  }
);

Button.displayName = 'Button';
