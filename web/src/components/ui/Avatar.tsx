'use client';

interface AvatarProps {
  src?: string;
  name: string;
  size?: 'xs' | 'sm' | 'md' | 'lg' | 'xl' | '2xl';
  status?: 'online' | 'away' | 'offline';
  className?: string;
}

const sizeClasses = {
  xs: 'w-6 h-6 text-[10px]',
  sm: 'w-8 h-8 text-xs',
  md: 'w-10 h-10 text-sm',
  lg: 'w-12 h-12 text-base',
  xl: 'w-16 h-16 text-xl',
  '2xl': 'w-24 h-24 text-2xl',
};

const statusSizeClasses = {
  xs: 'w-2 h-2',
  sm: 'w-2.5 h-2.5',
  md: 'w-3 h-3',
  lg: 'w-3.5 h-3.5',
  xl: 'w-4 h-4',
  '2xl': 'w-5 h-5',
};

export function Avatar({ src, name, size = 'md', status, className = '' }: AvatarProps) {
  const initials = name
    .split(' ')
    .map((n) => n[0])
    .filter(Boolean)
    .join('')
    .toUpperCase()
    .slice(0, 2);

  const colors = [
    'bg-violet-600',
    'bg-blue-600',
    'bg-emerald-600',
    'bg-orange-600',
    'bg-pink-600',
    'bg-indigo-600',
  ];

  // Generate consistent color based on name
  const colorIndex = name.split('').reduce((acc, char) => acc + char.charCodeAt(0), 0) % colors.length;

  return (
    <div className={`relative inline-flex shrink-0 rounded-full ${className}`}>
      {src ? (
        <img
          src={src}
          alt={name}
          className={`${sizeClasses[size]} rounded-full object-cover border border-white/10`}
        />
      ) : (
        <div
          className={`${sizeClasses[size]} rounded-full ${colors[colorIndex]} 
            flex items-center justify-center font-bold text-white border border-white/10 shadow-inner`}
        >
          {initials}
        </div>
      )}
      {status && (
        <span
          className={`absolute bottom-0 right-0 ${statusSizeClasses[size]} rounded-full 
            status-${status} ring-2 ring-[var(--background)]`}
        />
      )}
    </div>
  );
}



