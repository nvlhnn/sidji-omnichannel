'use client';

import { useState, useEffect } from 'react';
import { api } from '@/lib/api';
import { Channel, ChannelType, TeamMember, CannedResponse, InviteTeamMemberInput, CreateCannedResponseInput } from '@/lib/types';
import { AddChannelModal } from '../channels/AddChannelModal';

export function SettingsView() {
  const [activeTab, setActiveTab] = useState<'channels' | 'team' | 'canned-responses' | 'profile'>('channels');
  // Channels State
  const [channels, setChannels] = useState<Channel[]>([]);
  const [isLoadingChannels, setIsLoadingChannels] = useState(false);
  const [showAddChannel, setShowAddChannel] = useState(false);
  
  // Team State
  const [teamMembers, setTeamMembers] = useState<TeamMember[]>([]);
  const [isLoadingTeam, setIsLoadingTeam] = useState(false);
  const [showInviteMember, setShowInviteMember] = useState(false);
  
  // Canned Responses State
  const [cannedResponses, setCannedResponses] = useState<CannedResponse[]>([]);
  const [isLoadingCanned, setIsLoadingCanned] = useState(false);
  const [showAddCanned, setShowAddCanned] = useState(false);

  // Load Data based on tab
  useEffect(() => {
    if (activeTab === 'channels') loadChannels();
    if (activeTab === 'team') loadTeamMembers();
    if (activeTab === 'canned-responses') loadCannedResponses();
  }, [activeTab]);

  const loadChannels = async () => {
    try {
      setIsLoadingChannels(true);
      const response = await api.getChannels();
      setChannels(response.data || []);
    } catch (error) {
      console.error('Failed to load channels:', error);
    } finally {
      setIsLoadingChannels(false);
    }
  };

  const loadTeamMembers = async () => {
    try {
      setIsLoadingTeam(true);
      const members = await api.getTeamMembers();
      setTeamMembers(members || []);
    } catch (error) {
      console.error('Failed to load team members:', error);
    } finally {
      setIsLoadingTeam(false);
    }
  };
  
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

  const handleDeleteChannel = async (id: string) => {
    if (!confirm('Are you sure you want to delete this channel?')) return;
    try {
      await api.deleteChannel(id);
      loadChannels();
    } catch (error) {
      console.error('Failed to delete channel:', error);
      alert('Failed to delete channel');
    }
  };

  const handleActivateChannel = async (id: string) => {
    try {
      await api.activateChannel(id);
      loadChannels();
    } catch (error) {
      console.error('Failed to activate channel:', error);
      alert('Failed to activate channel');
    }
  };
  
  const handleRemoveMember = async (id: string) => {
    if (!confirm('Are you sure you want to remove this team member?')) return;
    try {
      await api.removeTeamMember(id);
      loadTeamMembers();
    } catch (error) {
      console.error('Failed to remove member:', error);
      alert('Failed to remove member');
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

  return (
    <div className="flex-1 flex flex-col bg-[var(--background)] overflow-hidden">
      {/* Header */}
      <div className="p-6 border-b border-[var(--border)] bg-[var(--background-secondary)]">
        <h1 className="text-2xl font-bold bg-gradient-to-r from-[var(--foreground)] to-[var(--foreground-secondary)] bg-clip-text text-transparent">
          Settings
        </h1>
        <div className="flex gap-4 mt-6 overflow-x-auto no-scrollbar">
          {(['channels', 'team', 'canned-responses', 'profile'] as const).map((tab) => (
            <button
              key={tab}
              onClick={() => setActiveTab(tab)}
              className={`pb-2 px-1 text-sm font-medium border-b-2 transition-colors whitespace-nowrap
                ${activeTab === tab
                  ? 'border-[var(--primary)] text-[var(--primary)]'
                  : 'border-transparent text-[var(--foreground-muted)] hover:text-[var(--foreground)]'
                }
              `}
            >
              {tab.split('-').map(w => w.charAt(0).toUpperCase() + w.slice(1)).join(' ')}
            </button>
          ))}
        </div>
      </div>

      {/* Content */}
      <div className="flex-1 overflow-y-auto p-6">
        {/* CHANNELS TAB */}
        {activeTab === 'channels' && (
          <div className="max-w-4xl mx-auto">
            <div className="flex justify-between items-center mb-6">
              <h2 className="text-lg font-semibold">Connected Channels</h2>
              <button
                onClick={() => setShowAddChannel(true)}
                className="btn btn-primary flex items-center gap-2"
              >
                <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 4v16m8-8H4" />
                </svg>
                Connect Channel
              </button>
            </div>

            {isLoadingChannels ? (
              <div className="animate-pulse space-y-4">
                {[1, 2, 3].map((i) => (
                  <div key={i} className="h-20 bg-[var(--background-secondary)] rounded-xl" />
                ))}
              </div>
            ) : (
              <div className="grid gap-4">
                {channels.map((channel) => (
                  <div
                    key={channel.id}
                    className="p-4 bg-[var(--background-secondary)] rounded-xl border border-[var(--border)] flex items-center justify-between"
                  >
                    <div className="flex items-center gap-4">
                      <div className={`w-12 h-12 rounded-lg flex items-center justify-center
                        ${channel.type === 'whatsapp' ? 'bg-green-500/10 text-green-500' : 
                          channel.type === 'facebook' ? 'bg-blue-500/10 text-blue-500' : 
                          channel.type === 'tiktok' ? 'bg-black/10 text-black dark:bg-white/10 dark:text-white' :
                          'bg-pink-500/10 text-pink-500'}
                      `}>
                         {/* Channel Icon Logic - kept simple for brevity */}
                         <span className="capitalize text-xs font-bold">{channel.type.substring(0,2)}</span>
                      </div>
                      <div>
                        <h3 className="font-medium text-[var(--foreground)]">{channel.name}</h3>
                        <div className="flex items-center gap-2 mt-1">
                          <span className={`w-2 h-2 rounded-full ${channel.status === 'active' ? 'bg-green-500' : 'bg-yellow-500'}`} />
                          <p className="text-sm text-[var(--foreground-muted)] capitalize">{channel.status}</p>
                        </div>
                      </div>
                    </div>
                    
                    <div className="flex items-center gap-2">
                      {channel.status !== 'active' && (
                        <button
                          onClick={() => handleActivateChannel(channel.id)}
                          className="px-3 py-1.5 text-sm font-medium text-green-500 hover:bg-green-500/10 rounded-lg transition-colors"
                        >
                          Activate
                        </button>
                      )}
                      <button
                        onClick={() => handleDeleteChannel(channel.id)}
                        className="p-2 text-[var(--foreground-muted)] hover:text-red-500 hover:bg-red-500/10 rounded-lg transition-colors"
                        title="Delete Channel"
                      >
                         <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1.5} d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                        </svg>
                      </button>
                    </div>
                  </div>
                ))}

                {channels.length === 0 && (
                  <div className="text-center py-12 bg-[var(--background-secondary)] rounded-xl border border-[var(--border)] border-dashed">
                    <p className="text-[var(--foreground-muted)]">No channels connected yet.</p>
                  </div>
                )}
              </div>
            )}
          </div>
        )}

        {/* TEAM TAB */}
        {activeTab === 'team' && (
          <div className="max-w-4xl mx-auto">
             <div className="flex justify-between items-center mb-6">
              <h2 className="text-lg font-semibold">Team Members</h2>
              <button
                onClick={() => setShowInviteMember(true)}
                className="btn btn-primary flex items-center gap-2"
              >
                <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 4v16m8-8H4" />
                </svg>
                Invite Member
              </button>
            </div>
            
             {isLoadingTeam ? (
              <div className="animate-pulse space-y-4">
                {[1, 2, 3].map((i) => (
                  <div key={i} className="h-16 bg-[var(--background-secondary)] rounded-xl" />
                ))}
              </div>
            ) : (
                <div className="grid gap-4">
                    {teamMembers.map((member) => (
                        <div key={member.id} className="p-4 bg-[var(--background-secondary)] rounded-xl border border-[var(--border)] flex items-center justify-between">
                            <div className="flex items-center gap-4">
                                <div className="w-10 h-10 rounded-full bg-[var(--primary)] text-white flex items-center justify-center font-bold">
                                    {member.name.charAt(0).toUpperCase()}
                                </div>
                                <div>
                                    <h3 className="font-medium text-[var(--foreground)]">{member.name}</h3>
                                    <p className="text-sm text-[var(--foreground-muted)]">{member.email}</p>
                                </div>
                            </div>
                            <div className="flex items-center gap-4">
                                <span className={`px-2 py-1 rounded text-xs font-medium capitalize 
                                    ${member.role === 'admin' ? 'bg-purple-500/10 text-purple-500' : 
                                      member.role === 'supervisor' ? 'bg-blue-500/10 text-blue-500' : 
                                      'bg-gray-500/10 text-gray-500'}`}>
                                    {member.role}
                                </span>
                                <button 
                                    onClick={() => handleRemoveMember(member.id)}
                                    className="text-sm text-red-500 hover:underline"
                                >
                                    Remove
                                </button>
                            </div>
                        </div>
                    ))}
                </div>
            )}
          </div>
        )}

        {/* CANNED RESPONSES TAB */}
        {activeTab === 'canned-responses' && (
             <div className="max-w-4xl mx-auto">
             <div className="flex justify-between items-center mb-6">
              <h2 className="text-lg font-semibold">Canned Responses</h2>
              <button
                onClick={() => setShowAddCanned(true)}
                className="btn btn-primary flex items-center gap-2"
              >
                <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 4v16m8-8H4" />
                </svg>
                Create Response
              </button>
            </div>
            
            {isLoadingCanned ? (
                <div className="animate-pulse space-y-4">
                    <div className="h-16 bg-[var(--background-secondary)] rounded-xl" />
                </div>
            ) : (
                <div className="grid gap-4">
                    {cannedResponses.map((response) => (
                        <div key={response.id} className="p-4 bg-[var(--background-secondary)] rounded-xl border border-[var(--border)]">
                            <div className="flex justify-between items-start mb-2">
                                <div className="flex items-center gap-2">
                                    <span className="font-mono text-sm bg-[var(--background-tertiary)] px-2 py-1 rounded text-[var(--primary)]">
                                        /{response.shortcut}
                                    </span>
                                    <span className="font-medium text-sm text-[var(--foreground)]">{response.title}</span>
                                </div>
                                <button 
                                    onClick={() => handleDeleteCannedResponse(response.id)}
                                    className="text-red-500 hover:bg-red-500/10 p-1 rounded"
                                >
                                    <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" /></svg>
                                </button>
                            </div>
                            <p className="text-[var(--foreground-muted)] text-sm whitespace-pre-wrap">{response.content}</p>
                            {response.category && (
                                <div className="mt-2">
                                    <span className="text-xs text-[var(--foreground-muted)] bg-[var(--background)] px-2 py-0.5 rounded-full border border-[var(--border)]">
                                        {response.category}
                                    </span>
                                </div>
                            )}
                        </div>
                    ))}
                     {cannedResponses.length === 0 && (
                        <div className="text-center py-12 bg-[var(--background-secondary)] rounded-xl border border-[var(--border)] border-dashed">
                            <p className="text-[var(--foreground-muted)]">No canned responses yet.</p>
                        </div>
                    )}
                </div>
            )}
             </div>
        )}

        {activeTab === 'profile' && (
          <div className="text-center py-12">
            <h3 className="text-xl font-semibold mb-2">Profile Settings</h3>
            <p className="text-[var(--foreground-muted)]">Update your personal information and preferences.</p>
          </div>
        )}
      </div>

      {/* MODALS */}
      {showAddChannel && (
        <AddChannelModal onClose={() => setShowAddChannel(false)} onSuccess={() => {
          setShowAddChannel(false);
          loadChannels();
        }} />
      )}
      
      {showInviteMember && (
          <InviteMemberModal onClose={() => setShowInviteMember(false)} onSuccess={() => {
              setShowInviteMember(false);
              loadTeamMembers();
          }} />
      )}
      
      {showAddCanned && (
          <AddCannedResponseModal onClose={() => setShowAddCanned(false)} onSuccess={() => {
              setShowAddCanned(false);
              loadCannedResponses();
          }} />
      )}
    </div>
  );
}



function InviteMemberModal({ onClose, onSuccess }: { onClose: () => void; onSuccess: () => void }) {
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
        <div className="fixed inset-0 bg-black/50 backdrop-blur-sm flex items-center justify-center z-50">
            <div className="bg-[var(--background)] w-full max-w-md rounded-2xl shadow-xl overflow-hidden border border-[var(--border)]">
                <div className="p-6 border-b border-[var(--border)] flex justify-between items-center">
                    <h3 className="text-lg font-semibold">Invite Team Member</h3>
                    <button onClick={onClose} className="text-[var(--foreground-muted)] hover:text-[var(--foreground)]">
                        <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" /></svg>
                    </button>
                </div>
                <form onSubmit={handleSubmit} className="p-6 space-y-4">
                    <div>
                        <label className="block text-sm font-medium mb-1">Name</label>
                        <input type="text" required value={name} onChange={(e) => setName(e.target.value)} className="input w-full" placeholder="John Doe" />
                    </div>
                    <div>
                        <label className="block text-sm font-medium mb-1">Email</label>
                        <input type="email" required value={email} onChange={(e) => setEmail(e.target.value)} className="input w-full" placeholder="john@example.com" />
                    </div>
                    <div>
                        <label className="block text-sm font-medium mb-1">Role</label>
                        <select 
                            value={role} 
                            onChange={(e) => setRole(e.target.value as any)} 
                            className="input w-full bg-[var(--background-tertiary)]"
                        >
                            <option value="agent">Agent</option>
                            <option value="supervisor">Supervisor</option>
                            <option value="admin">Admin</option>
                        </select>
                    </div>
                    <div className="pt-4 flex justify-end gap-2">
                        <button type="button" onClick={onClose} className="px-4 py-2 rounded-lg hover:bg-[var(--background-secondary)] transition-colors">Cancel</button>
                        <button type="submit" disabled={isSubmitting} className="btn btn-primary px-6">{isSubmitting ? 'Inviting...' : 'Send Invitation'}</button>
                    </div>
                </form>
            </div>
        </div>
    );
}

function AddCannedResponseModal({ onClose, onSuccess }: { onClose: () => void; onSuccess: () => void }) {
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
        <div className="fixed inset-0 bg-black/50 backdrop-blur-sm flex items-center justify-center z-50">
            <div className="bg-[var(--background)] w-full max-w-md rounded-2xl shadow-xl overflow-hidden border border-[var(--border)]">
                <div className="p-6 border-b border-[var(--border)] flex justify-between items-center">
                    <h3 className="text-lg font-semibold">Create Canned Response</h3>
                    <button onClick={onClose} className="text-[var(--foreground-muted)] hover:text-[var(--foreground)]">
                        <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" /></svg>
                    </button>
                </div>
                <form onSubmit={handleSubmit} className="p-6 space-y-4">
                    <div>
                        <label className="block text-sm font-medium mb-1">Title</label>
                        <input type="text" required value={title} onChange={(e) => setTitle(e.target.value)} className="input w-full" placeholder="Customer Greeting" />
                    </div>
                    <div>
                        <label className="block text-sm font-medium mb-1">Shortcut</label>
                        <div className="relative">
                            <span className="absolute left-3 top-1/2 -translate-y-1/2 text-[var(--foreground-muted)]">/</span>
                            <input type="text" required value={shortcut} onChange={(e) => setShortcut(e.target.value)} className="input w-full pl-6" placeholder="hello" />
                        </div>
                        <p className="text-xs text-[var(--foreground-muted)] mt-1">Type taking this shortcut in chat to trigger the response.</p>
                    </div>
                    <div>
                        <label className="block text-sm font-medium mb-1">Response Content</label>
                        <textarea required value={content} onChange={(e) => setContent(e.target.value)} className="input w-full h-32 resize-none" placeholder="Hello! How can I help you today?" />
                    </div>
                    <div>
                        <label className="block text-sm font-medium mb-1">Category (Optional)</label>
                        <input type="text" value={category} onChange={(e) => setCategory(e.target.value)} className="input w-full" placeholder="Greetings" />
                    </div>
                    <div className="pt-4 flex justify-end gap-2">
                        <button type="button" onClick={onClose} className="px-4 py-2 rounded-lg hover:bg-[var(--background-secondary)] transition-colors">Cancel</button>
                        <button type="submit" disabled={isSubmitting} className="btn btn-primary px-6">{isSubmitting ? 'Saving...' : 'Create Response'}</button>
                    </div>
                </form>
            </div>
        </div>
    );
}
