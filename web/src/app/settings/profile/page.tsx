'use client';

import { useState, useEffect } from 'react';
import { api } from '@/lib/api';
import { User } from '@/lib/types';
import { User as UserIcon, Mail, Shield, Save, Loader2, Camera } from 'lucide-react';
import { Button } from '@/components/ui/Button';
import { Input } from '@/components/ui/Input';

import { Avatar } from '@/components/ui/Avatar';

export default function ProfilePage() {
    const [user, setUser] = useState<User | null>(null);
    const [organization, setOrganization] = useState<any>(null);
    const [loading, setLoading] = useState(true);
    const [saving, setSaving] = useState(false);
    const [formData, setFormData] = useState({
        name: '',
        email: '',
        avatar_url: ''
    });

    useEffect(() => {
        loadProfile();
    }, []);

    const loadProfile = async () => {
        try {
            const response = await api.getMe();
            setUser(response.user);
            setOrganization(response.organization);
            setFormData({
                name: response.user.name,
                email: response.user.email,
                avatar_url: response.user.avatar_url || ''
            });
        } catch (err) {
            console.error('Failed to load profile:', err);
        } finally {
            setLoading(false);
        }
    };

    const handleSave = async (e: React.FormEvent) => {
        e.preventDefault();
        if (!user) return;
        setSaving(true);
        try {
            const updated = await api.updateTeamMember(user.id, {
                name: formData.name,
                avatar_url: formData.avatar_url
            });
            setUser(updated);
            alert('Profile updated successfully');
        } catch (err) {
            console.error('Failed to update profile:', err);
            alert('Failed to update profile');
        } finally {
            setSaving(false);
        }
    };

    if (loading) {
        return (
            <div className="flex items-center justify-center min-h-[400px]">
                <Loader2 className="w-8 h-8 animate-spin text-[var(--primary)]" />
            </div>
        );
    }

    return (
        <div className="p-6 max-w-5xl mx-auto animate-in fade-in duration-500">
            <div className="mb-10">
                <h1 className="text-4xl font-black tracking-tight bg-gradient-to-r from-[var(--foreground)] to-[var(--foreground-muted)] bg-clip-text text-transparent">Profile Settings</h1>
                <p className="text-[var(--foreground-muted)] text-sm mt-2 font-medium">Manage your personal identity and workspace association.</p>
            </div>

            <div className="grid grid-cols-1 lg:grid-cols-12 gap-10">
                {/* Left: Identity & Workspace */}
                <div className="lg:col-span-4 flex flex-col gap-6">
                    <div className="bg-[var(--background-secondary)] p-8 rounded-[2rem] border border-[var(--border)] shadow-xl relative overflow-hidden group">
                        <div className="absolute top-0 right-0 w-32 h-32 bg-[var(--primary)]/5 blur-3xl rounded-full"></div>
                        <div className="relative flex flex-col items-center">
                            <div className="relative group cursor-pointer">
                                <Avatar name={formData.name} src={formData.avatar_url} size="lg" />
                                <div className="absolute inset-0 bg-black/40 rounded-full flex items-center justify-center opacity-0 group-hover:opacity-100 transition-all duration-300 scale-90 group-hover:scale-100">
                                    <Camera className="w-8 h-8 text-white" />
                                </div>
                            </div>
                            <div className="mt-6 text-center">
                                <h3 className="text-2xl font-black">{formData.name}</h3>
                                <p className="text-xs font-bold text-[var(--primary)] uppercase tracking-widest mt-1">{user?.role}</p>
                            </div>
                        </div>
                    </div>

                    <div className="bg-gradient-to-br from-[var(--background-secondary)] to-[var(--background-tertiary)] p-6 rounded-[2rem] border border-[var(--border)] shadow-lg">
                        <h4 className="text-[10px] font-black uppercase tracking-[0.2em] text-[var(--foreground-muted)] mb-6">Workspace Details</h4>
                        <div className="space-y-6">
                            <div className="flex items-center gap-4">
                                <div className="w-12 h-12 rounded-2xl bg-[var(--primary)]/10 flex items-center justify-center text-[var(--primary)] border border-[var(--primary)]/20 shadow-inner">
                                    <Shield className="w-6 h-6" />
                                </div>
                                <div className="flex-1 min-w-0">
                                    <div className="text-sm font-black truncate">{organization?.name}</div>
                                    <div className="text-[10px] text-[var(--foreground-muted)] font-bold tracking-tight">Active Workspace</div>
                                </div>
                            </div>
                            <div className="pt-6 border-t border-[var(--border)]">
                                <div className="flex justify-between items-center mb-2">
                                    <span className="text-[10px] font-black text-[var(--foreground-muted)] uppercase">Subscription</span>
                                    <span className="text-[10px] font-black text-[var(--primary)] uppercase px-2 py-0.5 bg-[var(--primary)]/10 rounded-full">{organization?.plan}</span>
                                </div>
                                <div className="w-full h-1.5 bg-[var(--background-tertiary)] rounded-full overflow-hidden">
                                    <div 
                                        className="h-full bg-gradient-to-r from-[var(--primary)] to-purple-500 rounded-full" 
                                        style={{ width: organization?.plan === 'starter' ? '25%' : organization?.plan === 'growth' ? '60%' : '100%' }}
                                    ></div>
                                </div>
                            </div>
                        </div>
                    </div>
                </div>

                {/* Right: Detailed Configuration */}
                <div className="lg:col-span-8 flex flex-col gap-6">
                    <div className="bg-[var(--background-secondary)] rounded-[2rem] border border-[var(--border)] overflow-hidden shadow-2xl">
                        <div className="p-8 border-b border-[var(--border)] bg-[var(--background-tertiary)]/20">
                            <h2 className="text-xl font-black flex items-center gap-3">
                                <UserIcon className="w-6 h-6 text-[var(--primary)]" />
                                Account Information
                            </h2>
                        </div>
                        <form onSubmit={handleSave} className="p-8 space-y-8">
                            <div className="grid grid-cols-1 md:grid-cols-2 gap-8">
                                <Input 
                                    label="Full Name"
                                    leftIcon={<UserIcon className="w-4 h-4" />}
                                    value={formData.name}
                                    onChange={(e) => setFormData({ ...formData, name: e.target.value })}
                                    placeholder="John Doe"
                                />

                                <Input 
                                    label="Email Address"
                                    leftIcon={<Mail className="w-4 h-4" />}
                                    value={formData.email}
                                    disabled
                                    placeholder="email@example.com"
                                    className="bg-[var(--background-tertiary)]/50 cursor-not-allowed"
                                />
                            </div>

                            <Input 
                                label="Avatar Content URL"
                                value={formData.avatar_url}
                                onChange={(e) => setFormData({ ...formData, avatar_url: e.target.value })}
                                placeholder="https://images.unsplash.com/..."
                            />

                            <div className="pt-8 border-t border-[var(--border)] flex justify-between items-center">
                                <div className="text-[10px] font-bold text-[var(--foreground-muted)] max-w-[60%] leading-relaxed">
                                    Changes to your primary identity will be reflected across all conversations and internal audit logs.
                                </div>
                                <Button 
                                    type="submit"
                                    isLoading={saving}
                                    size="xl"
                                    leftIcon={<Save className="w-5 h-5" />}
                                >
                                    Save Profile
                                </Button>
                            </div>
                        </form>
                    </div>

                    <div className="bg-red-500/5 rounded-[2rem] border border-red-500/10 p-8 flex items-center justify-between group hover:bg-red-500/10 transition-colors">
                        <div>
                            <h4 className="text-sm font-black text-red-500">Danger Zone</h4>
                            <p className="text-[10px] font-bold text-[var(--foreground-muted)] mt-1">Once you delete your account, there is no going back. Please be certain.</p>
                        </div>
                        <button className="h-12 px-6 rounded-xl border border-red-500/20 text-red-500 text-xs font-black uppercase tracking-widest hover:bg-red-500 hover:text-white transition-all">
                            Delete Account
                        </button>
                    </div>
                </div>
            </div>
        </div>
    );
}


