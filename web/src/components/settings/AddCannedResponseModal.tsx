'use client';

import { useState } from 'react';
import { api } from '@/lib/api';

interface AddCannedResponseModalProps {
    onClose: () => void;
    onSuccess: () => void;
}

export function AddCannedResponseModal({ onClose, onSuccess }: AddCannedResponseModalProps) {
    const [title, setTitle] = useState('');
    const [shortcut, setShortcut] = useState('');
    const [content, setContent] = useState('');
    const [category, setCategory] = useState('');
    const [isSubmitting, setIsSubmitting] = useState(false);

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        try {
            setIsSubmitting(true);
            await api.createCannedResponse({ title, shortcut, content, category });
            onSuccess();
        } catch (error) {
            console.error('Failed to create response:', error);
            alert('Failed to create canned response.');
        } finally {
            setIsSubmitting(false);
        }
    };

    return (
        <div className="fixed inset-0 bg-black/50 backdrop-blur-sm flex items-center justify-center z-50 animate-fadeIn">
            <div className="bg-[var(--background)] w-full max-w-md rounded-2xl shadow-xl overflow-hidden border border-[var(--border)] animate-slideUp">
                <div className="p-6 border-b border-[var(--border)] flex justify-between items-center">
                    <h3 className="text-lg font-semibold">Create Canned Response</h3>
                    <button onClick={onClose} className="text-[var(--foreground-muted)] hover:text-[var(--foreground)] transition-colors">
                        <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" /></svg>
                    </button>
                </div>
                <form onSubmit={handleSubmit} className="p-6 space-y-4">
                    <div>
                        <label className="block text-sm font-medium mb-1.5">Title</label>
                        <input 
                            type="text" 
                            required 
                            value={title} 
                            onChange={(e) => setTitle(e.target.value)} 
                            className="input w-full" 
                            placeholder="Customer Greeting" 
                        />
                    </div>
                    <div>
                        <label className="block text-sm font-medium mb-1.5">Shortcut</label>
                        <div className="relative">
                            <span className="absolute left-3 top-1/2 -translate-y-1/2 text-[var(--foreground-muted)] font-mono text-sm">/</span>
                            <input 
                                type="text" 
                                required 
                                value={shortcut} 
                                onChange={(e) => setShortcut(e.target.value)} 
                                className="input w-full pl-7 font-mono text-sm" 
                                placeholder="hello" 
                            />
                        </div>
                        <p className="text-xs text-[var(--foreground-muted)] mt-1.5 flex items-center gap-1">
                            <svg className="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" /></svg>
                            Type this shortcut in chat to trigger the response
                        </p>
                    </div>
                    <div>
                        <label className="block text-sm font-medium mb-1.5">Response Content</label>
                        <textarea 
                            required 
                            value={content} 
                            onChange={(e) => setContent(e.target.value)} 
                            className="input w-full h-32 resize-none py-3" 
                            placeholder="Hello! How can I help you today?" 
                        />
                    </div>
                    <div>
                        <label className="block text-sm font-medium mb-1.5">Category (Optional)</label>
                        <input 
                            type="text" 
                            value={category} 
                            onChange={(e) => setCategory(e.target.value)} 
                            className="input w-full" 
                            placeholder="Greetings" 
                        />
                    </div>
                    <div className="pt-4 flex justify-end gap-3">
                        <button 
                            type="button" 
                            onClick={onClose} 
                            className="px-4 py-2 rounded-xl text-sm font-medium hover:bg-[var(--background-secondary)] transition-colors"
                        >
                            Cancel
                        </button>
                        <button 
                            type="submit" 
                            disabled={isSubmitting} 
                            className="btn-primary px-6 py-2 rounded-xl text-sm"
                        >
                            {isSubmitting ? 'Saving...' : 'Create Response'}
                        </button>
                    </div>
                </form>
            </div>
        </div>
    );
}
