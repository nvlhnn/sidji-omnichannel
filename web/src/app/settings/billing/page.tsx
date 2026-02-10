'use client';

import { useState, useEffect } from 'react';
import { api } from '@/lib/api';
import { Organization, User } from '@/lib/types';
import { Check, Zap, Crown, Shield, Activity } from 'lucide-react';

export default function BillingPage() {
    const [user, setUser] = useState<User | null>(null);
    const [organization, setOrganization] = useState<Organization | null>(null);
    const [loading, setLoading] = useState(true);

    useEffect(() => {
        const loadData = async () => {
            try {
                const response = await api.getMe();
                setUser(response.user);
                setOrganization(response.organization);
            } catch (err) {
                console.error('Failed to load profile:', err);
            } finally {
                setLoading(false);
            }
        };
        loadData();
    }, []);

    const tiers = [
        {
            id: 'starter',
            name: 'Starter',
            price: '$0',
            description: 'Perfect for small teams getting started.',
            features: [
                '1 User Seat',
                '1 Connected Channel',
                '10 AI Replies / Month',
                'Basic Analytics',
                'Community Support'
            ],
            color: 'bg-blue-500',
            icon: Zap
        },
        {
            id: 'growth',
            name: 'Growth',
            price: '$29',
            period: '/mo',
            description: 'Scale your operations with advanced AI.',
            features: [
                '3 User Seats',
                'Unlimited Channels',
                '1,000 AI Replies / Month',
                'Advanced Analytics',
                'Priority Email Support',
                'Custom AI Persona'
            ],
            color: 'bg-purple-500',
            icon: Activity,
            popular: true
        },
        {
            id: 'scale',
            name: 'Scale',
            price: '$99',
            period: '/mo',
            description: 'Maximum power for high-volume teams.',
            features: [
                '10 User Seats',
                'Unlimited Channels',
                'Unlimited AI Replies',
                'Real-time Analytics',
                '24/7 Priority Support',
                'Custom AI Training'
            ],
            color: 'bg-orange-500',
            icon: Crown
        },
        {
            id: 'enterprise',
            name: 'Enterprise',
            price: 'Custom',
            description: 'Bespoke features for large organizations.',
            features: [
                'Unlimited User Seats',
                'Unlimited Channels',
                'Unlimited AI Replies',
                'Dedicated Infrastructure',
                'Custom SLA & Support',
                'On-premise Deployment'
            ],
            color: 'bg-gray-900',
            icon: Shield
        }
    ];

    const handleUpgradePlan = async (tierId: string) => {
        if (!confirm(`Are you sure you want to change to the ${tierId} plan?`)) return;
        setLoading(true);
        try {
            const updated = await api.updateOrganization({ plan: tierId });
            setOrganization(updated);
            alert(`Your plan has been updated to ${tierId}!`);
        } catch (err: any) {
            console.error('Failed to update plan:', err);
            alert('Failed to update plan: ' + err.message);
        } finally {
            setLoading(false);
        }
    };

    if (loading) {
        return (
            <div className="flex items-center justify-center min-h-[400px]">
                <div className="w-8 h-8 border-4 border-[var(--primary)] border-t-transparent rounded-full animate-spin"></div>
            </div>
        );
    }

    return (
        <div className="p-6 max-w-6xl mx-auto pb-20">
            {/* Header Section */}
            <div className="mb-10 text-center max-w-2xl mx-auto">
                <h1 className="text-3xl font-bold font-display mb-3">Simple, Transparent Pricing</h1>
                <p className="text-[var(--foreground-secondary)] text-lg">
                    Manage your subscription and AI credits. Upgrade anytime as your business grows.
                </p>
            </div>

            {/* Subscription Metrics Section */}
            <div className="grid grid-cols-1 md:grid-cols-2 gap-6 mb-12">
                {/* AI Automation Quota */}
                {organization && (
                    <div className="p-8 rounded-3xl bg-[var(--background-secondary)] border border-[var(--border)] shadow-sm hover:shadow-md transition-all">
                        <div className="flex items-center gap-4 mb-6">
                            <div className="w-14 h-14 rounded-2xl bg-blue-500/10 flex items-center justify-center text-blue-600 shadow-sm border border-blue-500/10">
                                <Zap className="w-7 h-7" />
                            </div>
                            <div>
                                <div className="text-[10px] font-bold text-[var(--foreground-muted)] uppercase tracking-[0.2em] mb-1">Automation Quota</div>
                                <h2 className="text-2xl font-bold capitalize">AI Messages</h2>
                            </div>
                        </div>

                        <div className="space-y-4">
                            <div className="flex justify-between items-end mb-1">
                                <span className="text-sm font-medium text-[var(--foreground-secondary)]">Utilization</span>
                                <span className="text-lg font-black">
                                    {organization.ai_credits_limit === -1 
                                        ? '∞' 
                                        : `${organization.ai_credits_used} / ${organization.ai_credits_limit}`}
                                </span>
                            </div>
                            <div className="w-full h-3 bg-[var(--background)] rounded-full overflow-hidden p-0.5 border border-[var(--border)]">
                                <div 
                                    className="h-full bg-blue-600 rounded-full transition-all duration-700 ease-out"
                                    style={{ 
                                        width: organization.ai_credits_limit === -1 
                                            ? '100%' 
                                            : `${Math.min(100, (organization.ai_credits_used / (organization.ai_credits_limit || 1)) * 100)}%` 
                                    }}
                                />
                            </div>
                            <p className="text-[10px] text-[var(--foreground-muted)] font-medium italic">
                                * AI credits refresh on your next billing cycle: {organization.billing_cycle_start ? new Date(organization.billing_cycle_start).toLocaleDateString() : 'N/A'}
                            </p>
                        </div>
                    </div>
                )}

                {/* Messaging Volume Quota */}
                {organization && (
                    <div className="p-8 rounded-3xl bg-[var(--background-secondary)] border border-[var(--border)] shadow-sm hover:shadow-md transition-all">
                        <div className="flex items-center gap-4 mb-6">
                            <div className="w-14 h-14 rounded-2xl bg-green-500/10 flex items-center justify-center text-green-600 shadow-sm border border-green-500/10">
                                <Activity className="w-7 h-7" />
                            </div>
                            <div>
                                <div className="text-[10px] font-bold text-[var(--foreground-muted)] uppercase tracking-[0.2em] mb-1">Channel Capacity</div>
                                <h2 className="text-2xl font-bold capitalize">Message Volume</h2>
                            </div>
                        </div>

                        <div className="space-y-4">
                            <div className="flex justify-between items-end mb-1">
                                <span className="text-sm font-medium text-[var(--foreground-secondary)]">Monthly Volume</span>
                                <span className="text-lg font-black">
                                    {organization.message_usage_limit === -1 
                                        ? 'Unlimited' 
                                        : `${organization.message_usage_used} / ${organization.message_usage_limit}`}
                                </span>
                            </div>
                            <div className="w-full h-3 bg-[var(--background)] rounded-full overflow-hidden p-0.5 border border-[var(--border)]">
                                <div 
                                    className="h-full bg-emerald-600 rounded-full transition-all duration-700 ease-out"
                                    style={{ 
                                        width: organization.message_usage_limit === -1 
                                            ? '100%' 
                                            : `${Math.min(100, (organization.message_usage_used / (organization.message_usage_limit || 1)) * 100)}%` 
                                    }}
                                />
                            </div>
                            <p className="text-[10px] text-[var(--foreground-muted)] font-medium italic">
                                * Standard service chats are included. Volume refers to total conversations initiated.
                            </p>
                        </div>
                    </div>
                )}

            </div>

            {/* Pricing Grid */}
            <div className="grid grid-cols-1 md:grid-cols-4 gap-8">
                {tiers.map((tier) => {
                    const isCurrentPlan = organization?.plan === tier.id;
                    const Icon = tier.icon;

                    return (
                        <div 
                            key={tier.id}
                            className={`relative p-8 rounded-3xl border transition-all duration-300 flex flex-col
                                ${tier.popular 
                                    ? 'border-[var(--primary)] bg-[var(--background-secondary)] shadow-xl scale-105 z-10' 
                                    : 'border-[var(--border)] bg-[var(--background)] hover:border-[var(--primary)]/50'
                                }
                            `}
                        >
                            {tier.popular && (
                                <div className="absolute top-0 left-1/2 -translate-x-1/2 -translate-y-1/2 bg-[var(--primary)] text-white text-xs font-bold px-3 py-1 rounded-full uppercase tracking-wide">
                                    Most Popular
                                </div>
                            )}

                            <div className={`w-12 h-12 rounded-xl ${tier.color} bg-opacity-10 flex items-center justify-center mb-6`}>
                                <Icon className={`w-6 h-6 ${tier.color.replace('bg-', 'text-')}`} />
                            </div>

                            <h3 className="text-xl font-bold mb-2">{tier.name}</h3>
                            <p className="text-sm text-[var(--foreground-secondary)] mb-6 min-h-[40px]">
                                {tier.description}
                            </p>

                            <div className="mb-6">
                                <span className="text-4xl font-bold">{tier.price}</span>
                                {tier.period && <span className="text-[var(--foreground-muted)]">{tier.period}</span>}
                            </div>

                            <button 
                                onClick={() => handleUpgradePlan(tier.id)}
                                className={`w-full py-3 rounded-xl font-medium transition-all mb-8
                                    ${isCurrentPlan 
                                        ? 'bg-[var(--background-secondary)] border border-[var(--primary)] text-[var(--primary)] hover:bg-[var(--primary)]/5'
                                        : tier.popular
                                            ? 'bg-[var(--primary)] text-white hover:bg-[var(--primary)]/90 shadow-lg shadow-[var(--primary)]/25'
                                            : 'bg-[var(--background-tertiary)] text-[var(--foreground)] hover:bg-[var(--background-tertiary)]/80'
                                    }
                                `}
                            >
                                {isCurrentPlan ? 'Current Plan (Refresh)' : 'Upgrade Plan'}
                            </button>

                            <ul className="space-y-4 flex-1">
                                {tier.features.map((feature, i) => (
                                    <li key={i} className="flex items-start gap-3 text-sm">
                                        <div className={`mt-0.5 w-5 h-5 rounded-full flex items-center justify-center shrink-0 ${tier.color.replace('bg-', 'bg-').replace('500', '500/20')}`}>
                                            <Check className={`w-3 h-3 ${tier.color.replace('bg-', 'text-')}`} />
                                        </div>
                                        <span className="text-[var(--foreground-secondary)]">{feature}</span>
                                    </li>
                                ))}
                            </ul>
                        </div>
                    );
                })}
            </div>

            {/* Enterprise / Custom Section */}
            <div className="mt-16 p-8 rounded-3xl bg-gradient-to-br from-gray-900 to-gray-800 text-white flex flex-col md:flex-row items-center justify-between gap-8">
                <div>
                     <h3 className="text-2xl font-bold mb-2">Need a custom solution?</h3>
                     <p className="text-gray-300 max-w-xl">
                        Contact our sales team for custom limits, dedicated infrastructure, and enterprise-grade security features.
                     </p>
                </div>
                <button className="px-8 py-3 bg-white text-gray-900 rounded-xl font-bold hover:bg-gray-100 transition-colors">
                    Contact Sales
                </button>
            </div>
        </div>
    );
}
