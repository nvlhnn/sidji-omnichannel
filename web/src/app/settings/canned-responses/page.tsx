'use client';

import { useState, useEffect } from 'react';
import { api } from '@/lib/api';
import { CannedResponse } from '@/lib/types';
import { AddCannedResponseModal } from '@/components/settings/AddCannedResponseModal';
import { MessageSquarePlus, Trash2, Search, Zap } from 'lucide-react';

export default function CannedResponsesPage() {
    const [cannedResponses, setCannedResponses] = useState<CannedResponse[]>([]);
    const [isLoadingCanned, setIsLoadingCanned] = useState(false);
    const [showAddCanned, setShowAddCanned] = useState(false);
    const [searchQuery, setSearchQuery] = useState('');

    useEffect(() => {
        loadCannedResponses();
    }, []);

    const loadCannedResponses = async () => {
        try {
            setIsLoadingCanned(true);
            const responses = await api.getCannedResponses();
            setCannedResponses(responses || []);
        } catch (error) {
            console.error('Failed to load canned responses:', error);
        } finally {
            setIsLoadingCanned(false);
        }
    };

    const handleDeleteCannedResponse = async (id: string) => {
        if (!confirm('Are you sure you want to delete this canned response?')) return;
        try {
            await api.deleteCannedResponse(id);
            loadCannedResponses();
        } catch (error) {
            console.error('Failed to delete canned response:', error);
            alert('Failed to delete canned response');
        }
    };

    const filteredResponses = cannedResponses.filter(response => 
        response.title.toLowerCase().includes(searchQuery.toLowerCase()) || 
        response.shortcut.toLowerCase().includes(searchQuery.toLowerCase()) || 
        response.content.toLowerCase().includes(searchQuery.toLowerCase())
    );

    return (
        <div className="p-6 max-w-5xl mx-auto">
            <div className="flex flex-col md:flex-row justify-between items-start md:items-center gap-4 mb-8">
                <div>
                    <h1 className="text-3xl font-bold font-display bg-gradient-to-r from-[var(--foreground)] to-[var(--foreground-muted)] bg-clip-text text-transparent">
                        Canned Responses
                    </h1>
                    <p className="text-[var(--foreground-secondary)] mt-1">
                        Create shortcuts for frequently used messages.
                    </p>
                </div>
                <button
                    onClick={() => setShowAddCanned(true)}
                    className="btn-primary flex items-center justify-center gap-2 px-6 py-2.5 rounded-xl font-medium transition-all"
                >
                    <MessageSquarePlus className="w-5 h-5" />
                    Create Response
                </button>
            </div>

            {/* Search */}
            <div className="relative mb-6 max-w-md">
                <Search className="absolute left-3 top-1/2 -translate-y-1/2 w-5 h-5 text-[var(--foreground-muted)]" />
                <input 
                    type="text" 
                    placeholder="Search shortcuts..." 
                    value={searchQuery}
                    onChange={(e) => setSearchQuery(e.target.value)}
                    className="input w-full pl-10 bg-[var(--background-secondary)]" 
                />
            </div>

            {isLoadingCanned ? (
                <div className="grid gap-4 md:grid-cols-2">
                    {[1, 2, 3, 4].map((i) => (
                        <div key={i} className="h-32 bg-[var(--background-secondary)] rounded-2xl animate-pulse" />
                    ))}
                </div>
            ) : (
                <div className="grid gap-4 md:grid-cols-1">
                    {filteredResponses.map((response) => (
                        <div 
                            key={response.id} 
                            className="bg-[var(--background-secondary)] p-5 rounded-2xl border border-[var(--border)] group hover:border-[var(--primary)]/30 transition-all flex flex-col justify-between"
                        >
                            <div className="flex justify-between items-start mb-3">
                                <div className="flex items-center gap-3">
                                    <span className="font-mono text-sm font-semibold bg-[var(--primary)]/10 text-[var(--primary)] px-2.5 py-1 rounded-lg border border-[var(--primary)]/20 shadow-sm">
                                        /{response.shortcut}
                                    </span>
                                    <h3 className="font-semibold text-[var(--foreground)]">{response.title}</h3>
                                </div>
                                <button 
                                    onClick={() => handleDeleteCannedResponse(response.id)}
                                    className="p-2 text-[var(--foreground-muted)] hover:text-red-500 hover:bg-red-500/10 rounded-lg transition-all opacity-0 group-hover:opacity-100"
                                    title="Delete Response"
                                >
                                    <Trash2 className="w-4 h-4" />
                                </button>
                            </div>
                            
                            <p className="text-[var(--foreground-secondary)] text-sm whitespace-pre-wrap line-clamp-3 mb-4 leading-relaxed bg-[var(--background-tertiary)]/50 p-3 rounded-xl">
                                {response.content}
                            </p>

                            {response.category && (
                                <div className="flex items-center gap-2 mt-auto">
                                    <span className="text-xs font-medium text-[var(--foreground-muted)] bg-[var(--background-tertiary)] px-2.5 py-1 rounded-full border border-[var(--border)] uppercase tracking-wider">
                                        {response.category}
                                    </span>
                                </div>
                            )}
                        </div>
                    ))}

                    {filteredResponses.length === 0 && (
                        <div className="col-span-full text-center py-20 bg-[var(--background-secondary)] rounded-2xl border border-[var(--border)] border-dashed">
                            <div className="w-16 h-16 mx-auto bg-[var(--background-tertiary)] rounded-2xl flex items-center justify-center mb-4 text-[var(--foreground-muted)]">
                                <Zap className="w-8 h-8" />
                            </div>
                            <h3 className="text-xl font-medium mb-2">
                                {searchQuery ? 'No matching responses' : 'No canned responses yet'}
                            </h3>
                            <p className="text-[var(--foreground-secondary)] max-w-sm mx-auto mb-6">
                                {searchQuery ? 'Try a different search term.' : 'Use shortcuts to reply faster to common questions.'}
                            </p>
                            {!searchQuery && (
                                <button
                                    onClick={() => setShowAddCanned(true)}
                                    className="btn-secondary px-6"
                                >
                                    Create First Response
                                </button>
                            )}
                        </div>
                    )}
                </div>
            )}

            {showAddCanned && (
                <AddCannedResponseModal 
                    onClose={() => setShowAddCanned(false)} 
                    onSuccess={() => {
                        setShowAddCanned(false);
                        loadCannedResponses();
                    }} 
                />
            )}
        </div>
    );
}
