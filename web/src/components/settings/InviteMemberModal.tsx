'use client';

import { useState } from 'react';
import { api } from '@/lib/api';

interface InviteMemberModalProps {
    onClose: () => void;
    onSuccess: () => void;
}

export function InviteMemberModal({ onClose, onSuccess }: InviteMemberModalProps) {
    const [name, setName] = useState('');
    const [email, setEmail] = useState('');
    const [role, setRole] = useState<'admin' | 'supervisor' | 'agent'>('agent');
    const [isSubmitting, setIsSubmitting] = useState(false);

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        try {
            setIsSubmitting(true);
            await api.inviteTeamMember({ name, email, role });
            onSuccess();
        } catch (error) {
            console.error('Failed to invite member:', error);
            alert('Failed to invite member.');
        } finally {
            setIsSubmitting(false);
        }
    };

    return (
        <div className="fixed inset-0 bg-black/50 backdrop-blur-sm flex items-center justify-center z-50 animate-fadeIn">
            <div className="bg-[var(--background)] w-full max-w-md rounded-2xl shadow-xl overflow-hidden border border-[var(--border)] animate-slideUp">
                <div className="p-6 border-b border-[var(--border)] flex justify-between items-center">
                    <h3 className="text-lg font-semibold">Invite Team Member</h3>
                    <button onClick={onClose} className="text-[var(--foreground-muted)] hover:text-[var(--foreground)] transition-colors">
                        <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" /></svg>
                    </button>
                </div>
                <form onSubmit={handleSubmit} className="p-6 space-y-4">
                    <div>
                        <label className="block text-sm font-medium mb-1.5">Name</label>
                        <input 
                            type="text" 
                            required 
                            value={name} 
                            onChange={(e) => setName(e.target.value)} 
                            className="input w-full" 
                            placeholder="John Doe" 
                        />
                    </div>
                    <div>
                        <label className="block text-sm font-medium mb-1.5">Email</label>
                        <input 
                            type="email" 
                            required 
                            value={email} 
                            onChange={(e) => setEmail(e.target.value)} 
                            className="input w-full" 
                            placeholder="john@example.com" 
                        />
                    </div>
                    <div>
                        <label className="block text-sm font-medium mb-1.5">Role</label>
                        <select 
                            value={role} 
                            onChange={(e) => setRole(e.target.value as any)} 
                            className="input w-full bg-[var(--background-tertiary)] appearance-none"
                        >
                            <option value="agent">Agent</option>
                            <option value="supervisor">Supervisor</option>
                            <option value="admin">Admin</option>
                        </select>
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
                            {isSubmitting ? 'Sending Invite...' : 'Send Invitation'}
                        </button>
                    </div>
                </form>
            </div>
        </div>
    );
}
