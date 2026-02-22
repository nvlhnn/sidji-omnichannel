'use client';

import { useState, useEffect } from 'react';
import { api } from '@/lib/api';
import { Channel } from '@/lib/types';
import { 
  Instagram, 
  Trash2, 
  Plus, 
  Zap, 
  Bot, 
  MessageCircle, 
  Facebook,
  Shield,
  Activity,
  ChevronRight,
  Music
} from 'lucide-react';
import { AddChannelModal } from '@/components/channels/AddChannelModal';
import Link from 'next/link';

export default function ChannelsPage() {
    const [channels, setChannels] = useState<Channel[]>([]);
    const [loading, setLoading] = useState(true);
    const [showAddModal, setShowAddModal] = useState(false);
    const [isSecure, setIsSecure] = useState(true);
    const [organization, setOrganization] = useState<any>(null);

    useEffect(() => {
        loadData();
        loadFacebookSDK();
        if (typeof window !== 'undefined') {
            setIsSecure(window.location.protocol === 'https:');
        }
    }, []);

    const loadData = async () => {
         try {
            const [channelsRes, authRes] = await Promise.all([
                api.getChannels(),
                api.getMe()
            ]);
            setChannels(channelsRes.data);
            setOrganization(authRes.organization);
        } catch (err) {
            console.error('Failed to load data:', err);
        } finally {
            setLoading(false);
        }
    };
    
    // Limits
    const getLimit = (plan: string) => {
        if (plan === 'growth') return -1; // unlimited
        if (plan === 'scale') return -1;
        return 1; // starter
    };
    const limit = organization ? getLimit(organization.plan) : 1;
    const isLimitReached = limit !== -1 && channels.length >= limit;

    const loadChannels = async () => {
        try {
            const response = await api.getChannels();
            setChannels(response.data);
        } catch (err) {
            console.error('Failed to load channels:', err);
        }
    };

    const loadFacebookSDK = () => {
        if (typeof window === 'undefined') return;
        // @ts-ignore
        if (window.FB) return;
        // @ts-ignore
        window.fbAsyncInit = function() {
            // @ts-ignore
            window.FB.init({
                appId      : process.env.NEXT_PUBLIC_META_APP_ID,
                cookie     : true,
                xfbml      : true,
                version    : 'v18.0'
            });
        };
        (function(d, s, id){
             var js, fjs = d.getElementsByTagName(s)[0];
             if (d.getElementById(id)) {return;}
             js = d.createElement(s); js.id = id;
             // @ts-ignore
             js.src = "https://connect.facebook.net/en_US/sdk.js";
             // @ts-ignore
             fjs.parentNode.insertBefore(js, fjs);
        }(document, 'script', 'facebook-jssdk'));
    };

    const handleDisconnect = async (id: string, name: string) => {
        if (!confirm(`Are you sure you want to disconnect ${name}?`)) return;
        try {
            await api.deleteChannel(id);
            setChannels(prev => prev.filter(c => c.id !== id));
        } catch (err) {
            console.error('Failed to disconnect:', err);
            alert('Failed to disconnect channel');
        }
    };

    const handleActivate = async (id: string, name: string) => {
        try {
            await api.activateChannel(id);
            await loadChannels();
        } catch (err) {
            console.error('Failed to activate:', err);
            alert('Failed to activate channel');
        }
    };

    return (
        <div className="p-4 sm:p-8 max-w-6xl mx-auto space-y-8 animate-fadeIn">
            {/* Security Warning */}
            {!isSecure && (
                <div className="p-4 rounded-2xl bg-amber-500/5 border border-amber-500/20 flex items-center gap-4">
                    <div className="w-10 h-10 rounded-xl bg-amber-500/10 flex items-center justify-center text-amber-500 shrink-0">
                        <Shield size={20} />
                    </div>
                    <div>
                        <h4 className="text-amber-500 font-bold text-sm">HTTP Connection Detected</h4>
                        <p className="text-amber-500/70 text-xs mt-0.5 leading-relaxed">
                            Meta requires HTTPS for secure connections. Fast Connect is disabled. Please use an <strong>HTTPS tunnel</strong> or <strong>Manual Setup</strong>.
                        </p>
                    </div>
                </div>
            )}

            {/* Header */}
            <div className="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-6">
                <div className="space-y-1">
                    <h1 className="text-4xl font-black tracking-tight bg-gradient-to-r from-[var(--foreground)] to-[var(--foreground-muted)] bg-clip-text text-transparent">
                        Channels
                    </h1>
                    <div className="flex items-center gap-2 text-sm text-[var(--foreground-secondary)]">
                        <Activity size={14} className="text-green-500" />
                        <span>
                            {channels.length} {channels.length === 1 ? 'channel' : 'channels'} currently active
                        </span>
                        {limit !== -1 && (
                            <span className="px-2 py-0.5 rounded-full bg-[var(--background-tertiary)] text-[10px] font-bold border border-[var(--border)]">
                                {channels.length}/{limit} USED
                            </span>
                        )}
                    </div>
                </div>
                
                <button 
                    onClick={() => setShowAddModal(true)}
                    disabled={isLimitReached}
                    className={`group relative flex items-center gap-3 px-6 py-3.5 rounded-2xl font-bold transition-all duration-300 shadow-lg shadow-[var(--primary)]/10
                        ${isLimitReached 
                            ? 'bg-[var(--background-tertiary)] text-[var(--foreground-muted)] cursor-not-allowed border border-[var(--border)]' 
                            : 'bg-[var(--primary)] hover:bg-[var(--primary-hover)] text-white hover:scale-[1.02] active:scale-95'
                        }
                    `}
                >
                    <Plus size={20} className="group-hover:rotate-90 transition-transform duration-300" />
                    <span>Connect Channel</span>
                    {isLimitReached && (
                        <div className="absolute -top-3 -right-3 px-3 py-1 bg-amber-500 text-white text-[10px] font-black rounded-full shadow-lg animate-bounce">
                            LIMIT REACHED
                        </div>
                    )}
                </button>
            </div>

            {/* Main Content */}
            <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
                {/* Channels List */}
                <div className="lg:col-span-2 space-y-4">
                    <div className="flex items-center justify-between px-2">
                        <h2 className="text-xs font-black uppercase tracking-[0.2em] text-[var(--foreground-muted)]">Active Connections</h2>
                        <div className="h-[1px] flex-1 bg-[var(--border)] mx-4 opacity-50" />
                    </div>

                    <div className="grid grid-cols-1 gap-4">
                        {loading ? (
                            <div className="py-20 rounded-3xl bg-[var(--background-secondary)]/50 border border-[var(--border)] border-dashed flex flex-col items-center justify-center">
                                <div className="w-12 h-12 border-4 border-[var(--primary)] border-t-transparent rounded-full animate-spin mb-4" />
                                <p className="text-sm font-medium text-[var(--foreground-muted)]">Looking for connections...</p>
                            </div>
                        ) : channels.length === 0 ? (
                            <div className="py-20 rounded-3xl bg-[var(--background-secondary)]/50 border border-[var(--border)] border-dashed text-center space-y-4">
                                <div className="w-20 h-20 bg-[var(--background-tertiary)] rounded-2xl mx-auto flex items-center justify-center text-[var(--foreground-muted)]">
                                    <MessageCircle size={40} strokeWidth={1.5} />
                                </div>
                                <div className="space-y-1">
                                    <h3 className="text-xl font-bold">Inbox is quiet</h3>
                                    <p className="text-sm text-[var(--foreground-secondary)] max-w-xs mx-auto">
                                        You haven't connected any messaging platforms yet. Start by adding WhatsApp or Instagram.
                                    </p>
                                </div>
                                <button 
                                    onClick={() => setShowAddModal(true)}
                                    className="px-6 py-2.5 bg-[var(--background-tertiary)] hover:bg-[var(--background-secondary)] text-[var(--foreground)] font-bold rounded-xl transition-all border border-[var(--border)]"
                                >
                                    Get Connected
                                </button>
                            </div>
                        ) : (
                            channels.map(channel => (
                                <div key={channel.id} className="group relative bg-[var(--background-secondary)]/30 hover:bg-[var(--background-secondary)] p-5 rounded-3xl border border-[var(--border)] hover:border-[var(--primary)]/30 transition-all duration-300">
                                    <div className="flex items-center justify-between gap-4">
                                        <div className="flex items-center gap-5">
                                            <div className={`w-14 h-14 rounded-2xl flex items-center justify-center shadow-inner transition-transform duration-300 group-hover:scale-105
                                                ${channel.type === 'instagram' 
                                                    ? 'bg-[#E1306C11] text-[#E1306C]' 
                                                    : channel.type === 'facebook'
                                                        ? 'bg-[#0866FF11] text-[#0866FF]'
                                                        : channel.type === 'tiktok'
                                                            ? 'bg-black/10 text-black dark:bg-white/10 dark:text-white'
                                                            : 'bg-[#25D36611] text-[#25D366]'
                                                }`}>
                                                {channel.type === 'instagram' && <Instagram size={28} />}
                                                {channel.type === 'facebook' && <Facebook size={28} />}
                                                {channel.type === 'whatsapp' && <MessageCircle size={28} />}
                                                {channel.type === 'tiktok' && <Music size={28} />}
                                            </div>
                                            <div className="space-y-1">
                                                <div className="flex items-center gap-3">
                                                    <h3 className="font-bold text-lg leading-none">{channel.name}</h3>
                                                    <span className={`px-2 py-0.5 rounded-full text-[9px] font-black uppercase tracking-wider
                                                        ${channel.status === 'active' 
                                                            ? 'bg-green-500/10 text-green-500' 
                                                            : 'bg-amber-500/10 text-amber-500'
                                                        }`}>
                                                        {channel.status}
                                                    </span>
                                                </div>
                                                <div className="flex items-center gap-4 text-xs text-[var(--foreground-muted)] font-medium">
                                                    <span className="capitalize">{channel.type}</span>
                                                    <div className="w-1 h-1 rounded-full bg-[var(--border)]" />
                                                    <span>via {channel.provider}</span>
                                                </div>
                                            </div>
                                        </div>

                                        <div className="flex items-center gap-2">
                                            <Link 
                                                href={`/settings/channels/${channel.id}/ai`} 
                                                className="w-10 h-10 rounded-xl flex items-center justify-center hover:bg-[var(--primary)]/10 text-[var(--foreground-muted)] hover:text-[var(--primary)] transition-all"
                                                title="Configure AI"
                                            >
                                                <Bot size={20} />
                                            </Link>
                                            <button 
                                                onClick={() => handleDisconnect(channel.id, channel.name)}
                                                className="w-10 h-10 rounded-xl flex items-center justify-center hover:bg-red-500/10 text-[var(--foreground-muted)] hover:text-red-500 transition-all"
                                                title="Disconnect"
                                            >
                                                <Trash2 size={20} />
                                            </button>
                                        </div>
                                    </div>
                                </div>
                            ))
                        )}
                    </div>
                </div>

                {/* Info Sidebar */}
                <div className="space-y-6">
                    <div className="flex items-center gap-2 px-1">
                        <h2 className="text-xs font-black uppercase tracking-[0.2em] text-[var(--foreground-muted)]">Next Steps</h2>
                        <div className="h-[1px] flex-1 bg-[var(--border)] opacity-30" />
                    </div>

                    <div className="p-6 rounded-3xl bg-gradient-to-br from-[var(--primary)]/10 to-transparent border border-[var(--primary)]/20 space-y-4">
                        <div className="w-10 h-10 rounded-xl bg-[var(--primary)] flex items-center justify-center text-white">
                            <Zap size={20} fill="currentColor" />
                        </div>
                        <div className="space-y-2">
                            <h4 className="font-bold">Connect your business</h4>
                            <p className="text-xs text-[var(--foreground-muted)] leading-relaxed font-medium">
                                Link your Meta accounts to centralize all customer conversations. We support WhatsApp, Instagram, and Facebook Messenger.
                            </p>
                        </div>
                        <Link href="/docs/channels" className="flex items-center justify-between p-3 rounded-xl bg-[var(--background-tertiary)] border border-[var(--border)] group hover:border-[var(--primary)]/50 transition-all">
                           <span className="text-xs font-bold">Read Integration Guide</span>
                           <ChevronRight size={14} className="group-hover:translate-x-1 transition-transform" />
                        </Link>
                    </div>

                    <div className="p-6 rounded-3xl bg-[var(--background-secondary)]/30 border border-[var(--border)] space-y-4">
                        <h4 className="text-sm font-bold">Supported Platforms</h4>
                        <div className="space-y-3">
                            <div className="flex items-center gap-3">
                                <div className="w-8 h-8 rounded-lg bg-green-500/10 text-green-500 flex items-center justify-center"><MessageCircle size={16}/></div>
                                <span className="text-xs font-semibold">WhatsApp Business</span>
                            </div>
                            <div className="flex items-center gap-3">
                                <div className="w-8 h-8 rounded-lg bg-[#E1306C11] text-[#E1306C] flex items-center justify-center"><Instagram size={16}/></div>
                                <span className="text-xs font-semibold">Instagram Messages</span>
                            </div>
                            <div className="flex items-center gap-3">
                                <div className="w-8 h-8 rounded-lg bg-[#0866FF11] text-[#0866FF] flex items-center justify-center"><Facebook size={16}/></div>
                                <span className="text-xs font-semibold">Facebook Messenger</span>
                            </div>
                            <div className="flex items-center gap-3">
                                <div className="w-8 h-8 rounded-lg bg-black/10 text-black dark:bg-white/10 dark:text-white flex items-center justify-center"><Music size={16}/></div>
                                <span className="text-xs font-semibold">TikTok Messages</span>
                            </div>
                        </div>
                    </div>
                </div>
            </div>

            {showAddModal && (
                <AddChannelModal 
                    onClose={() => setShowAddModal(false)} 
                    onSuccess={() => {
                        setShowAddModal(false);
                        loadChannels();
                    }}
                />
            )}
        </div>
    );
}
