'use client';

import { Shield, Zap, AlertTriangle } from 'lucide-react';
import Link from 'next/link';

interface ComplianceBannerProps {
    organization: any;
}

export function ComplianceBanner({ organization }: ComplianceBannerProps) {
    if (!organization?.is_over_limit) return null;

    return (
        <div className="bg-amber-600 text-white py-2 px-4 flex items-center justify-between gap-4 animate-fadeIn sticky top-0 z-50">
            <div className="flex items-center gap-3 overflow-hidden">
                <AlertTriangle className="w-5 h-5 shrink-0" />
                <p className="text-sm font-medium truncate">
                    Your account is over its limit ({organization.user_count} members, {organization.channel_count} channels). 
                    Some features are restricted until you upgrade or remove resources.
                </p>
            </div>
            <Link 
                href="/settings/billing"
                className="whitespace-nowrap bg-white text-amber-600 px-4 py-1 rounded-full text-xs font-bold hover:bg-amber-50 transition-colors flex items-center gap-1.5"
            >
                <Zap className="w-3.5 h-3.5 fill-current" />
                Upgrade Now
            </Link>
        </div>
    );
}
