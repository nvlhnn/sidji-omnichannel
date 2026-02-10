'use client';

import DashboardLayout from '@/components/layout/DashboardLayout';
import { SettingsTabs } from '@/components/settings/SettingsTabs';

export default function SettingsLayout({ children }: { children: React.ReactNode }) {
    return (
        <DashboardLayout>
            <SettingsTabs />
            <div className="flex-1 overflow-y-auto">
                {children}
            </div>
        </DashboardLayout>
    );
}
