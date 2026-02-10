'use client';

import React from 'react';
import Link from 'next/link';
import { User, Organization } from '@/lib/types';
import { Avatar } from '../ui/Avatar';
import { LayoutDashboard, Users, Settings, LogOut, MessageSquare, Zap } from 'lucide-react';

interface SidebarProps {
  user: User | null;
  organization: Organization | null;
  activeTab: 'inbox' | 'contacts' | 'settings';
  onTabChange: (tab: 'inbox' | 'contacts' | 'settings') => void;
  onLogout: () => void;
}

export function Sidebar({ user, organization, activeTab, onTabChange, onLogout }: SidebarProps) {
  const navItems = [
    { id: 'inbox' as const, label: 'Transmissions', icon: MessageSquare },
    { id: 'contacts' as const, label: 'Customers', icon: Users },
    { id: 'settings' as const, label: 'System', icon: Settings },
  ];

  const creditLimit = organization?.ai_credits_limit ?? 10;
  const creditsUsed = organization?.ai_credits_used ?? 0;
  const percentage = creditLimit === -1 ? 0 : Math.min(100, (creditsUsed / creditLimit) * 100);
  const isUnlimited = creditLimit === -1;

  let progressColor = 'bg-primary';
  if (percentage > 85) progressColor = 'bg-amber-500';
  if (percentage >= 100) progressColor = 'bg-red-500';

  return (
    <div className="w-[88px] bg-[#06060a] border-r border-white/5 flex flex-col items-center py-8 z-50 relative">
      {/* Dynamic Background Glow */}
      <div className="absolute top-0 left-0 w-full h-[300px] bg-primary/5 blur-[60px] pointer-events-none" />

      {/* Logo Section */}
      <div className="mb-12 relative group">
        <div className="absolute inset-0 bg-primary/20 blur-[20px] rounded-2xl opacity-0 group-hover:opacity-100 transition-opacity duration-700" />
        <div className="w-14 h-14 rounded-2xl bg-primary flex items-center justify-center shadow-2xl relative z-10 animate-float border border-white/10 group-hover:scale-110 transition-transform duration-500">
           <Zap className="w-7 h-7 text-white fill-white/20" />
        </div>
      </div>

      {/* Navigation */}
      <nav className="flex-1 flex flex-col items-center gap-4 w-full px-3">
        {navItems.map((item) => {
          const isActive = activeTab === item.id;
          const Icon = item.icon;
          return (
            <button
              key={item.id}
              onClick={() => onTabChange(item.id)}
              className={`
                relative w-full aspect-square rounded-2xl flex items-center justify-center transition-all duration-500 group
                ${isActive 
                  ? 'glass-elevated bg-white/10 text-white shadow-xl ring-1 ring-white/20' 
                  : 'text-[var(--foreground-muted)] hover:text-white hover:bg-white/5'}
              `}
              title={item.label}
            >
              <Icon className={`w-6 h-6 transition-all duration-500 ${isActive ? 'scale-110 drop-shadow-[0_0_8px_rgba(255,255,255,0.3)]' : 'group-hover:scale-110'}`} />
              {isActive && (
                <div className="absolute left-0 top-1/2 -translate-y-1/2 w-1 h-6 bg-primary rounded-r-full shadow-glow" />
              )}
            </button>
          );
        })}
      </nav>

      {/* Secondary Actions & Info */}
      <div className="mt-auto flex flex-col items-center gap-8 w-full px-4">
        {/* Credits Orb */}
        {organization && (
            <div className="relative group flex flex-col items-center">
                <div className="w-10 h-10 rounded-full glass border border-white/10 flex items-center justify-center overflow-hidden relative shadow-inner">
                    <div 
                        className={`absolute bottom-0 left-0 w-full ${progressColor} transition-all duration-1000 ease-out opacity-80`}
                        style={{ height: `${isUnlimited ? 100 : percentage}%` }}
                    >
                         <div className="absolute top-0 left-0 w-full h-[2px] bg-white/30 shimmer" />
                    </div>
                    <Zap className="w-4 h-4 text-white relative z-10" />
                </div>
                
                {/* Credits Tooltip */}
                <div className="absolute left-full ml-4 top-1/2 -translate-y-1/2 glass-elevated px-4 py-2 rounded-xl text-[10px] font-bold uppercase tracking-widest text-white whitespace-nowrap opacity-0 group-hover:opacity-100 translate-x-4 group-hover:translate-x-0 transition-all duration-300 pointer-events-none z-[100] border border-white/10 shadow-2xl">
                    AI Message Credits: {isUnlimited ? '∞' : `${creditsUsed}/${creditLimit}`}
                </div>
            </div>
        )}



        {/* User Sphere */}
        <div className="flex flex-col items-center gap-4 pb-2">
            {user && (
                <Link href="/settings/profile" className="relative group block">
                    <div className="absolute inset-0 bg-primary/20 blur-[15px] rounded-full opacity-0 group-hover:opacity-100 transition-opacity duration-500" />
                    <Avatar 
                        name={user.name} 
                        src={user.avatar_url} 
                        size="md" 
                        status={user.status} 
                        className="border border-white/10 hover:border-primary/40 transition-all duration-500 cursor-pointer"
                    />
                </Link>
            )}
            
            <button
                onClick={onLogout}
                className="w-12 h-12 rounded-2xl flex items-center justify-center text-[var(--foreground-muted)] hover:text-red-500 hover:bg-red-500/10 transition-all duration-300 group"
                title="Deactivate Session"
            >
                <LogOut className="w-5 h-5 transition-transform group-hover:-translate-x-1" />
            </button>
        </div>
      </div>
    </div>
  );
}


