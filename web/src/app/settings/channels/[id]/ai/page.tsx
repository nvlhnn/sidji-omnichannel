'use client';

import { useState, useEffect } from 'react';
import { api } from '@/lib/api';
import { AIConfig, AIMode } from '@/lib/types';
import { Bot, Save, Plus, Play, Loader2, Info, Edit2, Trash2, Check, X, Search, ShieldAlert, Instagram, Facebook, MessageCircle } from 'lucide-react';
import { useParams } from 'next/navigation';

export default function AISettingsPage() {
    const params = useParams();
    const channelId = params.id as string;

    const [loading, setLoading] = useState(true);
    const [saving, setSaving] = useState(false);
    const [config, setConfig] = useState<AIConfig | null>(null);
    const [channel, setChannel] = useState<any | null>(null);
    const [organization, setOrganization] = useState<any>(null);
    const [knowledgeItems, setKnowledgeItems] = useState<any[]>([]);
    const [knowledgeText, setKnowledgeText] = useState('');
    const [addingKnowledge, setAddingKnowledge] = useState(false);
    const [editingKnowledgeId, setEditingKnowledgeId] = useState<string | null>(null);
    const [editingText, setEditingText] = useState('');
    const [testQuery, setTestQuery] = useState('');
    const [testResult, setTestResult] = useState<{ reply: string; context: string[] } | null>(null);
    const [testing, setTesting] = useState(false);
    const [searchQuery, setSearchQuery] = useState('');

    const isOverLimit = organization?.is_over_limit;

    useEffect(() => {
        if (channelId) {
            loadInitialData();
            loadKnowledge();
        }
    }, [channelId]);

    const loadInitialData = async () => {
        try {
            const [configData, authData, channelsData] = await Promise.all([
                api.getAIConfig(channelId),
                api.getMe(),
                api.getChannels()
            ]);
            setConfig(configData);
            setOrganization(authData.organization);
            
            // Find the specific channel info
            const currentChannel = channelsData.data.find((c: any) => c.id === channelId);
            setChannel(currentChannel);
        } catch (err) {
            console.error('Failed to load initial data:', err);
        } finally {
            setLoading(false);
        }
    };

    const loadKnowledge = async () => {
        try {
            const data = await api.getAIKnowledge(channelId);
            setKnowledgeItems(data || []);
        } catch (err) {
            console.error('Failed to load knowledge:', err);
        }
    };

    const handleSaveConfig = async () => {
        if (!config || isOverLimit) return;
        setSaving(true);
        try {
            const updated = await api.updateAIConfig(channelId, {
                is_enabled: config.is_enabled,
                mode: config.mode,
                persona: config.persona,
                handover_timeout_minutes: config.handover_timeout_minutes
            });
            setConfig(updated);
            alert('Settings saved successfully');
        } catch (err) {
            console.error('Failed to save config:', err);
            alert('Failed to save settings');
        } finally {
            setSaving(false);
        }
    };

    const handleAddKnowledge = async () => {
        if (!knowledgeText.trim() || isOverLimit) return;
        setAddingKnowledge(true);
        try {
            await api.addAIKnowledge(channelId, knowledgeText);
            setKnowledgeText('');
            loadKnowledge(); // Refresh list
        } catch (err) {
            console.error('Failed to add knowledge:', err);
            alert('Failed to add knowledge');
        } finally {
            setAddingKnowledge(false);
        }
    };

    const handleDeleteKnowledge = async (kid: string) => {
        if (isOverLimit) return;
        if (!confirm('Are you sure you want to delete this knowledge?')) return;
        try {
            await api.deleteAIKnowledge(channelId, kid);
            setKnowledgeItems(prev => prev.filter(item => item.id !== kid));
        } catch (err) {
            console.error('Failed to delete knowledge:', err);
            alert('Failed to delete knowledge');
        }
    };

    const handleUpdateKnowledge = async (kid: string) => {
        if (!editingText.trim() || isOverLimit) return;
        try {
            await api.updateAIKnowledge(channelId, kid, editingText);
            setEditingKnowledgeId(null);
            loadKnowledge(); // Refresh
        } catch (err) {
            console.error('Failed to update knowledge:', err);
            alert('Failed to update knowledge');
        }
    };

    const handleTest = async () => {
        if (!testQuery.trim()) return;
        setTesting(true);
        try {
            const res = await api.testAIReply(channelId, testQuery);
            setTestResult(res);
        } catch (err) {
            console.error('Failed to test AI:', err);
            alert('Failed to test AI');
        } finally {
            setTesting(false);
        }
    };

    const filteredKnowledge = knowledgeItems.filter(item => 
        item.content.toLowerCase().includes(searchQuery.toLowerCase())
    );

    if (loading) {
        return (
            <div className="flex items-center justify-center min-h-[400px]">
                <Loader2 className="w-8 h-8 animate-spin text-[var(--primary)]" />
            </div>
        );
    }

    if (!config) {
        return <div className="p-6">Error loading configuration</div>;
    }

    return (
        <div className="min-h-full bg-[var(--background)] animate-in fade-in duration-500">
            {/* Header section with glass effect */}
            <div className="sticky top-0 z-10 p-6 bg-[var(--background)]/80 backdrop-blur-md border-b border-[var(--border)] mb-8">
                <div className="max-w-7xl mx-auto flex items-center justify-between">
                    <div className="flex items-center gap-5">
                        <div className="w-14 h-14 rounded-2xl bg-gradient-to-br from-indigo-500 to-purple-600 flex items-center justify-center text-white shadow-lg shadow-indigo-500/20 relative group">
                            <Bot className="w-7 h-7" />
                            {channel && (
                                <div className={`absolute -bottom-1 -right-1 p-1.5 rounded-lg border-2 border-[var(--background)] shadow-sm flex items-center justify-center
                                    ${channel.type === 'whatsapp' ? 'bg-[#25D366]' : 
                                      channel.type === 'instagram' ? 'bg-[#E1306C]' : 
                                      'bg-[#0866FF]'}`}
                                >
                                    {channel.type === 'whatsapp' && <MessageCircle className="w-3 h-3 text-white" />}
                                    {channel.type === 'instagram' && <Instagram className="w-3 h-3 text-white" />}
                                    {channel.type === 'facebook' && <Facebook className="w-3 h-3 text-white" />}
                                </div>
                            )}
                        </div>
                        <div>
                            <div className="flex items-center gap-2">
                                <h1 className="text-3xl font-bold tracking-tight">AI Assistant</h1>
                                {channel && (
                                    <span className="px-3 py-1 rounded-full bg-white/5 border border-white/10 text-[10px] font-bold uppercase tracking-widest text-[var(--foreground-muted)] flex items-center gap-2">
                                        Channel: <span className="text-white">{channel.name}</span>
                                    </span>
                                )}
                            </div>
                            <p className="text-[var(--foreground-muted)] text-sm mt-1">Configure how AI interacts with your customers on this specific channel.</p>
                        </div>
                    </div>
                    {isOverLimit && (
                        <div className="flex items-center gap-3 px-4 py-2 bg-red-500/10 border border-red-500/20 rounded-2xl text-red-500 animate-pulse">
                            <ShieldAlert className="w-5 h-5" />
                            <span className="text-sm font-bold">Subscription Limit Reached</span>
                        </div>
                    )}
                </div>
            </div>

            <div className={`max-w-7xl mx-auto px-6 pb-12 ${isOverLimit ? 'relative' : ''}`}>
                {isOverLimit && (
                    <div className="absolute inset-0 z-50 flex items-start justify-center pt-24 bg-[var(--background)]/40 backdrop-blur-[2px] rounded-3xl pointer-events-none">
                        <div className="bg-[var(--background-secondary)] p-8 rounded-3xl border border-[var(--border)] shadow-2xl text-center max-w-md pointer-events-auto">
                            <div className="w-16 h-16 rounded-2xl bg-red-500/10 flex items-center justify-center text-red-500 mx-auto mb-6">
                                <ShieldAlert className="w-8 h-8" />
                            </div>
                            <h3 className="text-xl font-bold mb-2">Settings Locked</h3>
                            <p className="text-[var(--foreground-muted)] text-sm mb-6">
                                Your organization has exceeded its subscription limits. 
                                Please upgrade your plan or remove resources (users/channels) to resume AI configuration.
                            </p>
                            <button 
                                onClick={() => window.location.href = '/settings/billing'}
                                className="w-full h-12 rounded-2xl bg-red-500 hover:bg-red-600 text-white font-bold transition-all"
                            >
                                Upgrade Plan
                            </button>
                        </div>
                    </div>
                )}

                <div className={`grid grid-cols-1 xl:grid-cols-12 gap-8 ${isOverLimit ? 'opacity-40 grayscale-[0.5]' : ''}`}>
                    
                    {/* LEFT COLUMN: Configuration (8 cols on large screens) */}
                    <div className="xl:col-span-8 space-y-8">
                        
                        {/* Main Config Card */}
                        <div className="bg-[var(--background-secondary)] rounded-3xl border border-[var(--border)] overflow-hidden shadow-xl shadow-black/5">
                            <div className="p-8 border-b border-[var(--border)] bg-gradient-to-r from-indigo-500/5 to-purple-500/5">
                                <div className="flex items-center justify-between mb-2">
                                    <h2 className="text-xl font-bold flex items-center gap-2">
                                        <Bot className="w-5 h-5 text-indigo-500" />
                                        Behavior Settings
                                    </h2>
                                    <div className="flex items-center gap-3 px-4 py-2 rounded-2xl bg-[var(--background-tertiary)] border border-[var(--border)]">
                                        <span className={`text-xs font-bold uppercase tracking-wider ${config.is_enabled ? 'text-green-500' : 'text-[var(--foreground-muted)]'}`}>
                                            {config.is_enabled ? 'Online' : 'Offline'}
                                        </span>
                                        <button 
                                            disabled={isOverLimit}
                                            onClick={() => setConfig({ ...config, is_enabled: !config.is_enabled })}
                                            className={`w-10 h-5 rounded-full relative transition-all duration-300 ${
                                                config.is_enabled ? 'bg-green-500' : 'bg-[var(--foreground-muted)]/30'
                                            } ${isOverLimit ? 'cursor-not-allowed opacity-50' : ''}`}
                                        >
                                            <span className={`absolute top-0.5 w-4 h-4 rounded-full bg-white shadow-sm transition-all duration-300 ${
                                                config.is_enabled ? 'left-5.5' : 'left-0.5'
                                            }`} />
                                        </button>
                                    </div>
                                </div>
                                <p className="text-sm text-[var(--foreground-muted)]">Set how your AI assistant should respond to incoming messages.</p>
                            </div>

                            <div className="p-8 space-y-8">
                                <div className="grid grid-cols-1 md:grid-cols-2 gap-8">
                                    <div className="space-y-2">
                                        <label className="text-sm font-bold text-[var(--foreground-muted)] uppercase tracking-wider">Operating Mode</label>
                                        <div className="relative group">
                                            <select 
                                                disabled={isOverLimit}
                                                value={config.mode}
                                                onChange={(e) => setConfig({ ...config, mode: e.target.value as AIMode })}
                                                className="w-full h-12 px-4 rounded-2xl bg-[var(--background-tertiary)] border border-[var(--border)] focus:border-indigo-500 outline-none appearance-none transition-all cursor-pointer font-medium disabled:cursor-not-allowed"
                                            >
                                                <option value="manual">Manual (AI Suggestions on Request)</option>
                                                <option value="auto">Fully Automatic (24/7 Support)</option>
                                                <option value="hybrid">Hybrid (Human-First Transition)</option>
                                            </select>
                                            <div className="absolute right-4 top-1/2 -translate-y-1/2 pointer-events-none text-[var(--foreground-muted)]">
                                                <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 9l-7 7-7-7" /></svg>
                                            </div>
                                        </div>
                                    </div>

                                    {config.mode === 'hybrid' && (
                                        <div className="space-y-2 animate-in slide-in-from-right-4 duration-300">
                                            <label className="text-sm font-bold text-[var(--foreground-muted)] uppercase tracking-wider">Human Handover Policy</label>
                                            <div className="flex items-center gap-3">
                                                <div className="relative flex-1">
                                                    <input 
                                                        disabled={isOverLimit}
                                                        type="number"
                                                        value={config.handover_timeout_minutes}
                                                        onChange={(e) => setConfig({ ...config, handover_timeout_minutes: parseInt(e.target.value) })}
                                                        className="w-full h-12 px-4 rounded-2xl bg-[var(--background-tertiary)] border border-[var(--border)] focus:border-indigo-500 outline-none font-medium disabled:cursor-not-allowed"
                                                    />
                                                    <span className="absolute right-4 top-1/2 -translate-y-1/2 text-[var(--foreground-muted)] text-sm">Minutes</span>
                                                </div>
                                                <div className="p-3 rounded-2xl bg-amber-500/10 text-amber-500 group relative">
                                                    <Info className="w-5 h-5" />
                                                    <div className="absolute bottom-full right-0 mb-2 w-64 p-3 bg-gray-900 text-white text-[10px] rounded-xl opacity-0 group-hover:opacity-100 transition-opacity pointer-events-none shadow-2xl z-20">
                                                        AI will wait this long after your last reply before taking back control of the conversation.
                                                    </div>
                                                </div>
                                            </div>
                                        </div>
                                    )}
                                </div>

                                <div className="space-y-2">
                                    <label className="text-sm font-bold text-[var(--foreground-muted)] uppercase tracking-wider">Persona & Intelligence Context</label>
                                    <textarea 
                                        disabled={isOverLimit}
                                        value={config.persona}
                                        onChange={(e) => setConfig({ ...config, persona: e.target.value })}
                                        className="w-full h-48 px-5 py-4 rounded-3xl bg-[var(--background-tertiary)] border border-[var(--border)] focus:border-indigo-500 outline-none resize-none transition-all leading-relaxed font-medium disabled:cursor-not-allowed"
                                        placeholder="Describe your assistant's personality and goals..."
                                    />
                                    <p className="text-[10px] text-[var(--foreground-muted)] px-1 italic">
                                        Pro-tip: Include things like "Be very polite," "Use emojis selectively," or "Focus on booking appointments."
                                    </p>
                                </div>

                                <div className="pt-4 border-t border-[var(--border)] flex justify-end">
                                    <button 
                                        onClick={handleSaveConfig} 
                                        disabled={saving || isOverLimit} 
                                        className="h-12 px-8 rounded-2xl bg-indigo-600 hover:bg-indigo-500 text-white font-bold flex items-center gap-3 transition-all active:scale-95 disabled:opacity-50 disabled:cursor-not-allowed"
                                    >
                                        {saving ? <Loader2 className="w-5 h-5 animate-spin" /> : <Save className="w-5 h-5" />}
                                        Apply Config Changes
                                    </button>
                                </div>
                            </div>
                        </div>

                        {/* Knowledge Content Card */}
                        <div className="bg-[var(--background-secondary)] rounded-3xl border border-[var(--border)] overflow-hidden shadow-xl shadow-black/5">
                            <div className="p-8 border-b border-[var(--border)] flex items-center justify-between gap-6">
                                <div>
                                    <h2 className="text-xl font-bold flex items-center gap-2">
                                        <Search className="w-5 h-5 text-purple-500" />
                                        Intelligence Base
                                    </h2>
                                    <p className="text-sm text-[var(--foreground-muted)]">Browse and manage data used by the AI to answer queries.</p>
                                </div>
                                <div className="relative w-full max-w-sm hidden md:block">
                                    <Search className="absolute left-4 top-1/2 -translate-y-1/2 w-4 h-4 text-[var(--foreground-muted)]" />
                                    <input 
                                        type="text" 
                                        placeholder="Quick context search..."
                                        value={searchQuery}
                                        onChange={(e) => setSearchQuery(e.target.value)}
                                        className="w-full h-11 pl-11 pr-4 rounded-2xl bg-[var(--background-tertiary)] border border-[var(--border)] text-sm outline-none focus:border-purple-500 transition-all font-medium"
                                    />
                                </div>
                            </div>

                            <div className="divide-y divide-[var(--border)] max-h-[800px] overflow-y-auto custom-scrollbar">
                                {filteredKnowledge.length === 0 ? (
                                    <div className="p-20 text-center flex flex-col items-center">
                                        <div className="w-16 h-16 rounded-full bg-[var(--background-tertiary)] flex items-center justify-center mb-4 text-[var(--foreground-muted)] opacity-50">
                                            <Search className="w-8 h-8" />
                                        </div>
                                        <p className="text-[var(--foreground-muted)] font-medium">No intelligence nodes found for this criteria.</p>
                                    </div>
                                ) : (
                                    filteredKnowledge.map((item) => (
                                        <div key={item.id} className="p-8 hover:bg-indigo-500/[0.02] group transition-all relative">
                                            {editingKnowledgeId === item.id ? (
                                                <div className="space-y-4">
                                                    <textarea 
                                                        disabled={isOverLimit}
                                                        value={editingText}
                                                        onChange={(e) => setEditingText(e.target.value)}
                                                        className="w-full h-32 p-4 rounded-2xl bg-[var(--background)] border border-indigo-500 outline-none text-sm font-medium leading-relaxed disabled:cursor-not-allowed"
                                                    />
                                                    <div className="flex justify-end gap-3">
                                                        <button disabled={isOverLimit} onClick={() => setEditingKnowledgeId(null)} className="h-9 px-4 rounded-xl text-xs font-bold text-[var(--foreground-muted)] hover:bg-[var(--background-tertiary)] disabled:cursor-not-allowed">
                                                            Discard
                                                        </button>
                                                        <button disabled={isOverLimit} onClick={() => handleUpdateKnowledge(item.id)} className="h-9 px-6 rounded-xl text-xs font-bold bg-indigo-600 text-white disabled:cursor-not-allowed">
                                                            Update Node
                                                        </button>
                                                    </div>
                                                </div>
                                            ) : (
                                                <div className="flex items-start justify-between gap-6">
                                                    <div className="flex-1 space-y-3">
                                                        <div className="p-4 rounded-2xl bg-[var(--background-tertiary)] border border-[var(--border)] prose prose-invert max-w-none">
                                                            <p className="text-sm whitespace-pre-wrap leading-relaxed font-medium">
                                                                {item.content}
                                                            </p>
                                                        </div>
                                                        <div className="flex items-center gap-2">
                                                            <div className="w-1 h-1 rounded-full bg-indigo-500" />
                                                            <span className="text-[10px] uppercase font-bold tracking-widest text-[var(--foreground-muted)]">
                                                                Indexed on {new Date(item.created_at).toLocaleDateString()}
                                                            </span>
                                                        </div>
                                                    </div>
                                                    <div className="flex gap-2 opacity-0 group-hover:opacity-100 transition-all translate-x-2 group-hover:translate-x-0">
                                                        <button 
                                                            disabled={isOverLimit}
                                                            onClick={() => { setEditingKnowledgeId(item.id); setEditingText(item.content); }}
                                                            className="w-10 h-10 flex items-center justify-center rounded-xl bg-indigo-500/10 text-indigo-500 hover:bg-indigo-500 hover:text-white transition-all shadow-sm disabled:cursor-not-allowed disabled:opacity-50"
                                                        >
                                                            <Edit2 className="w-4 h-4" />
                                                        </button>
                                                        <button 
                                                            disabled={isOverLimit}
                                                            onClick={() => handleDeleteKnowledge(item.id)}
                                                            className="w-10 h-10 flex items-center justify-center rounded-xl bg-red-500/10 text-red-500 hover:bg-red-500 hover:text-white transition-all shadow-sm disabled:cursor-not-allowed disabled:opacity-50"
                                                        >
                                                            <Trash2 className="w-4 h-4" />
                                                        </button>
                                                    </div>
                                                </div>
                                            )}
                                        </div>
                                    ))
                                )}
                            </div>
                        </div>
                    </div>

                    {/* RIGHT COLUMN: Actions (4 cols) */}
                    <div className="xl:col-span-4 space-y-8">
                        
                        {/* Rapid Knowledge Ingestion */}
                        <div className="bg-gradient-to-br from-indigo-600 to-purple-600 rounded-3xl p-1 shadow-xl shadow-indigo-500/10">
                            <div className="bg-[var(--background-secondary)] rounded-[1.35rem] p-7 space-y-5">
                                <h3 className="text-lg font-bold flex items-center gap-2">
                                    <Plus className="w-5 h-5 text-indigo-500" />
                                    Add Knowledge
                                </h3>
                                <div className="space-y-4">
                                    <textarea 
                                        disabled={isOverLimit}
                                        value={knowledgeText}
                                        onChange={(e) => setKnowledgeText(e.target.value)}
                                        className="w-full h-40 px-4 py-4 rounded-2xl bg-[var(--background-tertiary)] border border-[var(--border)] focus:border-indigo-500 outline-none resize-none text-sm font-medium leading-relaxed disabled:cursor-not-allowed"
                                        placeholder="Input text that the AI should learn. e.g. Business details, pricing, FAQs..."
                                    />
                                    <button 
                                        onClick={handleAddKnowledge}
                                        disabled={addingKnowledge || !knowledgeText.trim() || isOverLimit}
                                        className="w-full h-12 rounded-2xl bg-gradient-to-r from-indigo-500 to-indigo-600 text-white font-bold flex items-center justify-center gap-3 shadow-lg shadow-indigo-500/20 active:scale-95 disabled:opacity-50 transition-all font-medium disabled:cursor-not-allowed"
                                    >
                                        {addingKnowledge ? <Loader2 className="w-5 h-5 animate-spin" /> : <Plus className="w-5 h-5" />}
                                        Commit to Memory
                                    </button>
                                </div>
                            </div>
                        </div>

                        {/* Testing Playground Card */}
                        <div className="bg-[var(--background-secondary)] rounded-3xl border border-[var(--border)] p-7 space-y-6 shadow-xl shadow-black/5 relative overflow-hidden">
                            <div className="absolute top-0 right-0 p-8 opacity-5 pointer-events-none">
                                <Play className="w-24 h-24" />
                            </div>
                            
                            <h3 className="text-lg font-bold flex items-center gap-2">
                                <Play className="w-5 h-5 text-purple-500" />
                                AI Playground
                            </h3>
                            <div className="space-y-5">
                                <div className="space-y-2">
                                    <label className="text-xs font-bold text-[var(--foreground-muted)] uppercase tracking-wider">Test Query</label>
                                    <input 
                                        disabled={isOverLimit}
                                        type="text"
                                        value={testQuery}
                                        onChange={(e) => setTestQuery(e.target.value)}
                                        className="w-full h-12 px-4 rounded-2xl bg-[var(--background-tertiary)] border border-[var(--border)] outline-none text-sm font-medium focus:border-purple-500 transition-all shadow-sm disabled:cursor-not-allowed"
                                        placeholder="Ask a question as a customer..."
                                    />
                                </div>
                                <button 
                                    onClick={handleTest}
                                    disabled={testing || !testQuery.trim() || isOverLimit}
                                    className="w-full h-11 rounded-xl bg-[var(--background-tertiary)] border border-[var(--border)] hover:bg-[var(--border)] font-bold text-sm transition-all active:scale-95 disabled:opacity-50 flex items-center justify-center gap-2 disabled:cursor-not-allowed"
                                >
                                    {testing ? <Loader2 className="w-4 h-4 animate-spin text-purple-500" /> : <Bot className="w-4 h-4 text-purple-500" />}
                                    Simulate AI Response
                                </button>
                            </div>

                            {testResult && (
                                <div className="mt-8 space-y-5 animate-in slide-in-from-bottom-4 duration-500">
                                    <div className="p-5 rounded-2xl bg-indigo-500/5 border border-indigo-500/10 shadow-inner">
                                        <div className="flex items-center gap-2 mb-3">
                                            <div className="w-2 h-2 rounded-full bg-indigo-500 animate-pulse" />
                                            <h4 className="text-[10px] font-bold text-indigo-500 uppercase tracking-widest">Assistant Reply</h4>
                                        </div>
                                        <p className="text-sm font-medium leading-relaxed italic text-[var(--foreground)]">"{testResult.reply}"</p>
                                    </div>
                                    
                                    {testResult.context?.length > 0 && (
                                        <div className="p-5 rounded-2xl bg-[var(--background-tertiary)] border border-[var(--border)]">
                                            <h4 className="text-[10px] font-bold text-[var(--foreground-muted)] uppercase tracking-widest mb-3">Knowledge Sources</h4>
                                            <div className="space-y-3">
                                                {testResult.context.map((ctx: string, i: number) => (
                                                    <div key={i} className="flex gap-3">
                                                        <span className="text-indigo-400 font-mono text-xs pt-1">0{i+1}</span>
                                                        <p className="text-[10px] font-medium text-[var(--foreground-secondary)] leading-relaxed italic line-clamp-2">
                                                            {ctx}
                                                        </p>
                                                    </div>
                                                ))}
                                            </div>
                                        </div>
                                    )}
                                </div>
                            )}
                        </div>
                    </div>
                </div>
            </div>
        </div>
    );
}
