'use client';

import { ChannelType } from '@/lib/types';
import { Instagram, Facebook, MessageCircle } from 'lucide-react';

interface ChannelBadgeProps {
  type: ChannelType;
  size?: 'xs' | 'sm' | 'md';
}

export function ChannelBadge({ type, size = 'sm' }: ChannelBadgeProps) {
  const sizeClasses = {
    xs: 'w-4 h-4 p-0.5 rounded-lg',
    sm: 'w-6 h-6 p-1.5 rounded-xl',
    md: 'w-8 h-8 p-2 rounded-2xl',
  };

  const Icon = type === 'whatsapp' ? MessageCircle : type === 'instagram' ? Instagram : Facebook;

  const bgClasses = {
    whatsapp: 'bg-[#25D366] shadow-[0_0_10px_rgba(37,211,102,0.3)]',
    instagram: 'bg-gradient-to-tr from-[#f09433] via-[#e6683c] to-[#bc1888] shadow-[0_0_10px_rgba(230,104,60,0.3)]',
    facebook: 'bg-[#0866FF] shadow-[0_0_10px_rgba(8,102,255,0.3)]',
  };

  return (
    <div className={`
        inline-flex items-center justify-center border-2 border-[var(--background)] ring-1 ring-white/10
        ${sizeClasses[size]} ${bgClasses[type]}
    `}>
      <Icon className="w-full h-full text-white" />
    </div>
  );
}

