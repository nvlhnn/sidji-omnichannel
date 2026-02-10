'use client';

import { usePathname, useRouter } from 'next/navigation';

export function SettingsTabs() {
    const router = useRouter();
    const pathname = usePathname();

    const tabs = [
        { id: 'channels', label: 'Channels', path: '/settings/channels' },
        { id: 'team', label: 'Team', path: '/settings/team' },
        { id: 'canned-responses', label: 'Canned Responses', path: '/settings/canned-responses' },
        { id: 'billing', label: 'Billing & Plans', path: '/settings/billing' },
    ];

    return (
        <div className="border-b border-[var(--border)] bg-[var(--background-secondary)] px-6 pt-6">
            <h1 className="text-2xl font-bold text-[var(--foreground)] mb-6">
                Settings
            </h1>
            <div className="flex gap-4 overflow-x-auto no-scrollbar">
                {tabs.map((tab) => (
                    <button
                        key={tab.id}
                        onClick={() => router.push(tab.path)}
                        className={`pb-2 px-1 text-sm font-medium border-b-2 transition-colors whitespace-nowrap
                            ${pathname.includes(tab.path)
                                ? 'border-[var(--primary)] text-[var(--primary)]'
                                : 'border-transparent text-[var(--foreground-muted)] hover:text-[var(--foreground)]'
                            }
                        `}
                    >
                        {tab.label}
                    </button>
                ))}
            </div>
        </div>
    );
}
