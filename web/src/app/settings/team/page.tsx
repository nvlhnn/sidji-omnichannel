'use client';

import { useState, useEffect } from 'react';
import { api } from '@/lib/api';
import { TeamMember } from '@/lib/types';
import { InviteMemberModal } from '@/components/settings/InviteMemberModal';
import { Users, UserPlus, Shield, User, Trash2 } from 'lucide-react';

export default function TeamPage() {
    const [teamMembers, setTeamMembers] = useState<TeamMember[]>([]);
    const [isLoadingTeam, setIsLoadingTeam] = useState(false);
    const [showInviteMember, setShowInviteMember] = useState(false);
    const [organization, setOrganization] = useState<any>(null);

    useEffect(() => {
        loadData();
    }, []);

    const loadData = async () => {
        try {
            setIsLoadingTeam(true);
            const [members, authRes] = await Promise.all([
                api.getTeamMembers(),
                api.getMe()
            ]);
            setTeamMembers(members || []);
            setOrganization(authRes.organization);
        } catch (error) {
            console.error('Failed to load data:', error);
        } finally {
            setIsLoadingTeam(false);
        }
    };

    const handleRemoveMember = async (id: string) => {
        if (!confirm('Are you sure you want to remove this team member? This cannot be undone.')) return;
        try {
            await api.removeTeamMember(id);
            loadData();
        } catch (error) {
            console.error('Failed to remove member:', error);
            alert('Failed to remove member');
        }
    };

    // Limits
    const getLimit = (plan: string) => {
        if (plan === 'growth') return 3;
        if (plan === 'scale') return 10;
        return 1; // starter
    };
    const limit = organization ? getLimit(organization.plan) : 1;
    const isLimitReached = teamMembers.length >= limit;

    return (
        <div className="p-6 max-w-5xl mx-auto">
            <div className="flex flex-col md:flex-row justify-between items-start md:items-center gap-4 mb-8">
                <div>
                    <h1 className="text-3xl font-bold font-display bg-gradient-to-r from-[var(--foreground)] to-[var(--foreground-muted)] bg-clip-text text-transparent">
                        Team Members
                    </h1>
                    <p className="text-[var(--foreground-secondary)] mt-1">
                        Manage your team access and roles. ({teamMembers.length}/{limit} seats used)
                    </p>
                </div>
                {isLimitReached ? (
                     <div className="flex items-center gap-2 px-4 py-2 bg-yellow-500/10 text-yellow-600 rounded-xl text-sm font-medium border border-yellow-500/20">
                        <Shield className="w-4 h-4" />
                        Upgrade to Invite More
                     </div>
                ) : (
                    <button
                        onClick={() => setShowInviteMember(true)}
                        className="btn-primary flex items-center justify-center gap-2 px-6 py-2.5 rounded-xl font-medium transition-all"
                    >
                        <UserPlus className="w-5 h-5" />
                        Invite Member
                    </button>
                )}
            </div>

            {isLoadingTeam ? (
                <div className="space-y-4">
                    {[1, 2, 3].map((i) => (
                        <div key={i} className="h-20 bg-[var(--background-secondary)] rounded-2xl animate-pulse" />
                    ))}
                </div>
            ) : (
                <div className="grid gap-4">
                    {teamMembers.map((member) => (
                        <div 
                            key={member.id} 
                            className="bg-[var(--background-secondary)] p-5 rounded-2xl border border-[var(--border)] group hover:border-[var(--primary)]/30 transition-all flex items-center justify-between"
                        >
                            <div className="flex items-center gap-5">
                                <div className={`w-12 h-12 rounded-xl flex items-center justify-center text-lg font-bold shadow-inner ${
                                    member.role === 'admin' 
                                        ? 'bg-purple-500/10 text-purple-500 ring-1 ring-purple-500/20' 
                                        : 'bg-[var(--primary)]/10 text-[var(--primary)] ring-1 ring-[var(--primary)]/20'
                                }`}>
                                    {member.name.charAt(0).toUpperCase()}
                                </div>
                                <div>
                                    <div className="flex items-center gap-2">
                                        <h3 className="font-semibold text-lg text-[var(--foreground)]">{member.name}</h3>
                                        {member.role === 'admin' && (
                                            <span className="flex items-center gap-1 px-2 py-0.5 rounded-full bg-purple-500/10 text-purple-500 text-[10px] font-bold uppercase tracking-wider border border-purple-500/20">
                                                <Shield className="w-3 h-3" />
                                                Admin
                                            </span>
                                        )}
                                    </div>
                                    <p className="text-sm text-[var(--foreground-muted)] mt-0.5">{member.email}</p>
                                </div>
                            </div>

                            <div className="flex items-center gap-4">
                                <div className="hidden md:block text-right mr-4">
                                    <span className={`text-xs font-medium uppercase tracking-wider ${
                                        member.role === 'admin' ? 'text-purple-500' : 'text-[var(--foreground-muted)]'
                                    }`}>
                                        {member.role} Role
                                    </span>
                                </div>
                                <button 
                                    onClick={() => handleRemoveMember(member.id)}
                                    className="p-3 text-[var(--foreground-muted)] hover:text-red-500 hover:bg-red-500/10 rounded-xl transition-all opacity-0 group-hover:opacity-100"
                                    title="Remove Member"
                                >
                                    <Trash2 className="w-5 h-5" />
                                </button>
                            </div>
                        </div>
                    ))}

                    {teamMembers.length === 0 && (
                        <div className="text-center py-20 bg-[var(--background-secondary)] rounded-2xl border border-[var(--border)] border-dashed">
                            <div className="w-16 h-16 mx-auto bg-[var(--background-tertiary)] rounded-2xl flex items-center justify-center mb-4 text-[var(--foreground-muted)]">
                                <Users className="w-8 h-8" />
                            </div>
                            <h3 className="text-xl font-medium mb-2">No team members yet</h3>
                            <p className="text-[var(--foreground-secondary)] max-w-sm mx-auto mb-6">
                                Invite your colleagues to collaborate on customer conversations.
                            </p>
                            <button
                                onClick={() => setShowInviteMember(true)}
                                className="btn-secondary px-6"
                            >
                                Invite First Member
                            </button>
                        </div>
                    )}
                </div>
            )}

            {showInviteMember && (
                <InviteMemberModal 
                    onClose={() => setShowInviteMember(false)} 
                    onSuccess={() => {
                        setShowInviteMember(false);
                        loadData();
                    }} 
                />
            )}
        </div>
    );
}
