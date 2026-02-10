import React from 'react';
import { cn } from '@/lib/utils';

interface InputProps extends React.InputHTMLAttributes<HTMLInputElement> {
  label?: string;
  error?: string;
  leftIcon?: React.ReactNode;
  rightIcon?: React.ReactNode;
  containerClassName?: string;
}

export const Input = React.forwardRef<HTMLInputElement, InputProps>(
  ({ label, error, leftIcon, rightIcon, className, containerClassName, ...props }, ref) => {
    return (
      <div className={cn('space-y-2 w-full', containerClassName)}>
        {label && (
          <label className="text-[10px] font-black text-[var(--foreground-muted)] uppercase tracking-[0.2em] ml-1">
            {label}
          </label>
        )}
        <div className="relative group">
          {leftIcon && (
            <div className="absolute left-4 top-1/2 -translate-y-1/2 text-[var(--foreground-muted)] group-focus-within:text-[var(--primary)] transition-colors">
              {leftIcon}
            </div>
          )}
          <input
            ref={ref}
            className={cn(
              'w-full h-12 rounded-2xl bg-[var(--background-tertiary)] border border-[var(--border)] focus:border-[var(--primary)] focus:ring-4 focus:ring-[var(--primary)]/10 outline-none transition-all font-bold text-sm px-6',
              leftIcon ? 'pl-12' : '',
              rightIcon ? 'pr-12' : '',
              error ? 'border-red-500 focus:border-red-500 focus:ring-red-500/10' : '',
              className
            )}
            {...props}
          />
          {rightIcon && (
            <div className="absolute right-4 top-1/2 -translate-y-1/2">
              {rightIcon}
            </div>
          )}
        </div>
        {error && <p className="text-[10px] font-bold text-red-500 ml-1">{error}</p>}
      </div>
    );
  }
);

Input.displayName = 'Input';
