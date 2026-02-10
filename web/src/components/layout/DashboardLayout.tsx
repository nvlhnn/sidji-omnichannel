'use client';

import React, { useState, useEffect } from 'react';
import { useRouter, usePathname } from 'next/navigation';
import { Sidebar } from './Sidebar';
import { ComplianceBanner } from './ComplianceBanner';
import { api } from '@/lib/api';
import { User } from '@/lib/types';

export default function DashboardLayout({ children }: { children: React.ReactNode }) {
    const router = useRouter();
    const pathname = usePathname();
    const [user, setUser] = useState<User | null>(null);
    const [organization, setOrganization] = useState<any | null>(null);
    const [activeTab, setActiveTab] = useState<'inbox' | 'contacts' | 'settings'>('inbox');

    useEffect(() => {
        const loadUser = async () => {
            try {
                const response = await api.getMe();
                setUser(response.user);
                setOrganization(response.organization);
            } catch (err) {
                console.error('Failed to load user:', err);
                router.push('/login');
            }
        };
        loadUser();

        if (pathname.includes('/contacts')) setActiveTab('contacts');
        else if (pathname.includes('/settings')) setActiveTab('settings');
        else setActiveTab('inbox');
    }, [pathname, router]);

    const handleLogout = () => {
        api.logout();
        router.push('/login');
    };

    const handleTabChange = (tab: 'inbox' | 'contacts' | 'settings') => {
        if (tab === 'settings') router.push('/settings/channels');
        else if (tab === 'contacts') router.push('/contacts');
        else router.push('/inbox');
        setActiveTab(tab);
    };

    if (!user) {
        return (
            <div className="h-screen flex items-center justify-center bg-[var(--background)]">
                <div className="w-8 h-8 border-4 border-[var(--primary)] border-t-transparent rounded-full animate-spin"></div>
            </div>
        );
    }

    return (
        <div className="h-screen flex flex-col bg-[var(--background)] overflow-hidden">
            <ComplianceBanner organization={organization} />
            <div className="flex-1 flex overflow-hidden">
                <Sidebar
                    user={user}
                    organization={organization}
                    activeTab={activeTab}
                    onTabChange={handleTabChange}
                    onLogout={handleLogout}
                />
                <main className="flex-1 overflow-auto">
                    {children}
                </main>
            </div>
        </div>
    );
}
