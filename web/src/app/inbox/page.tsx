'use client';

import { useState, useEffect, useCallback, useRef } from 'react';
import { useRouter } from 'next/navigation';
import { ConversationList } from '@/components/inbox/ConversationList';
import { MessagePanel } from '@/components/inbox/MessagePanel';
import { api } from '@/lib/api';
import { useWebSocket } from '@/lib/websocket';
import { Conversation, Message, User, WSEvent, AuthResponse } from '@/lib/types';
import DashboardLayout from '@/components/layout/DashboardLayout';

export default function InboxPage() {
  const router = useRouter();
  const [authData, setAuthData] = useState<AuthResponse | null>(null);
  const [conversations, setConversations] = useState<Conversation[]>([]);
  const [selectedConversation, setSelectedConversation] = useState<Conversation | null>(null);
  const selectedConversationRef = useRef<Conversation | null>(null);

  // Keep ref in sync
  useEffect(() => {
    selectedConversationRef.current = selectedConversation;
  }, [selectedConversation]);

  const [messages, setMessages] = useState<Message[]>([]);
  const [isLoadingConversations, setIsLoadingConversations] = useState(true);
  const [isLoadingMessages, setIsLoadingMessages] = useState(false);
  const [filter, setFilter] = useState<string>('all');
  const [searchQuery, setSearchQuery] = useState('');

  const [drafts, setDrafts] = useState<Record<string, { content: string; created_at: string }>>({});

  // Handle WebSocket messages
  const handleWSMessage = useCallback((event: WSEvent) => {
    const currentSelected = selectedConversationRef.current;

    switch (event.event) {
      case 'conversation:ai_draft':
        setDrafts(prev => ({
          ...prev,
          [event.data.conversation_id]: event.data
        }));
        break;

      case 'message:new':
        if (currentSelected && event.data.conversation_id === currentSelected.id) {
          setMessages((prev) => {
            if (prev.find(m => m.id === event.data.id)) return prev;
            return [...prev, event.data];
          });
        }
        setConversations((prev) => {
          const updatedConv = prev.find(c => c.id === event.data.conversation_id);
          if (updatedConv) {
            const isSelected = currentSelected?.id === updatedConv.id;
            const isContact = event.data.sender_type === 'contact';
            
            const newConv = { 
              ...updatedConv, 
              last_message: event.data, 
              unread_count: updatedConv.unread_count + (isSelected || !isContact ? 0 : 1), 
              last_message_at: event.data.created_at 
            };
            return [newConv, ...prev.filter(c => c.id !== event.data.conversation_id)];
          }
          return prev;
        });
        break;

      case 'message:status':
        if (currentSelected) {
            setMessages((prev) => 
                prev.map(msg => 
                    msg.id === event.data.id ? { ...msg, status: event.data.status } : msg
                )
            );
        }
        break;

      case 'conversation:update':
        setConversations((prev) => {
          const exists = prev.find(c => c.id === event.data.id);
          if (exists) {
            return [event.data, ...prev.filter(c => c.id !== event.data.id)];
          }
          return [event.data, ...prev];
        });
        break;
    }
  }, []);

  useWebSocket(handleWSMessage);

  // Load user data
  useEffect(() => {
    const loadUser = async () => {
      try {
        const data = await api.getMe();
        setAuthData(data);
      } catch {
        router.push('/login');
      }
    };
    loadUser();
  }, [router]);

  // Load conversations
  useEffect(() => {
    const loadConversations = async () => {
      try {
        setIsLoadingConversations(true);
        const response = await api.getConversations({
          status: filter !== 'all' ? filter as never : undefined,
          search: searchQuery || undefined,
        });
        setConversations(response.data);
      } catch (error) {
        console.error('Failed to load conversations:', error);
      } finally {
        setIsLoadingConversations(false);
      }
    };

    if (authData) {
      loadConversations();
    }
  }, [authData, filter, searchQuery]);

  // Load messages when conversation is selected
  useEffect(() => {
    const loadMessages = async () => {
      if (!selectedConversation) {
        setMessages([]);
        return;
      }

      try {
        setIsLoadingMessages(true);
        const response = await api.getMessages(selectedConversation.id);
        setMessages(response.messages);
      } catch (error) {
        console.error('Failed to load messages:', error);
      } finally {
        setIsLoadingMessages(false);
      }
    };

    loadMessages();
  }, [selectedConversation]);

  const handleSendMessage = async (content: string, isNote: boolean = false) => {
    if (!selectedConversation) return;

    try {
      const message = await api.sendMessage(selectedConversation.id, content, 'text', isNote);
      setMessages((prev) => {
        if (prev.find(m => m.id === message.id)) return prev;
        return [...prev, message];
      });
      
      setConversations((prev) => {
        const updatedConv = prev.find(c => c.id === selectedConversation.id);
        if (updatedConv) {
          const newConv = { 
            ...updatedConv, 
            last_message: message, 
            unread_count: 0,
            last_message_at: message.created_at 
          };
          return [newConv, ...prev.filter(c => c.id !== selectedConversation.id)];
        }
        return prev;
      });

      setDrafts(prev => {
        const next = { ...prev };
        delete next[selectedConversation.id];
        return next;
      });
    } catch (error) {
      console.error('Failed to send message:', error);
      alert(error instanceof Error ? error.message : 'Failed to send message');
    }
  };

  const handleConversationSelect = (conversation: Conversation) => {
    setSelectedConversation(conversation);
    setConversations((prev) =>
      prev.map((conv) =>
        conv.id === conversation.id ? { ...conv, unread_count: 0 } : conv
      )
    );
  };

  const handleDeleteConversation = async () => {
    if (!selectedConversation) return;
    if (confirm('Are you sure you want to delete this conversation? This cannot be undone.')) {
      try {
        await api.deleteConversation(selectedConversation.id);
        setConversations(prev => prev.filter(c => c.id !== selectedConversation.id));
        setSelectedConversation(null);
      } catch (err) {
        console.error('Failed to delete conversation:', err);
        alert('Failed to delete conversation');
      }
    }
  };

  const handleDeleteContact = async () => {
    if (!selectedConversation) return;
    if (confirm(`Are you sure you want to delete contact ${selectedConversation.contact.name}? All their conversations will be deleted.`)) {
      try {
        await api.deleteContact(selectedConversation.contact.id);
        setConversations(prev => prev.filter(c => c.contact.id !== selectedConversation.contact.id));
        setSelectedConversation(null);
      } catch (err) {
        console.error('Failed to delete contact:', err);
        alert('Failed to delete contact');
      }
    }
  };

  if (!authData) {
    return (
      <div className="h-screen flex items-center justify-center bg-[var(--background)]">
        <div className="animate-pulse flex flex-col items-center">
          <div className="w-12 h-12 border-3 border-[var(--primary)] border-t-transparent rounded-full animate-spin"></div>
          <p className="mt-4 text-[var(--foreground-muted)]">Loading...</p>
        </div>
      </div>
    );
  }

  return (
    <DashboardLayout>
      <div className="h-full flex bg-[var(--background)] overflow-hidden">
        {/* Left Sidebar: Conversation List */}
        <div className="w-[380px] border-r border-[var(--border)] flex flex-col bg-[var(--background-secondary)] relative z-30">
          <div className="p-6 border-b border-[var(--border)] flex flex-col gap-6">
            <div className="flex items-center justify-between">
              <h1 className="text-3xl font-bold tracking-tighter text-gradient">
                Inbox
              </h1>
              <div className="flex items-center gap-2">
                <div className="w-2 h-2 rounded-full bg-emerald-500 animate-pulse" />
                <span className="text-[10px] font-bold uppercase tracking-widest text-[var(--foreground-muted)]">Live Flow</span>
              </div>
            </div>
            
            <div className="relative w-full group">
              <div className="absolute left-4 top-1/2 -translate-y-1/2 w-5 h-5 text-[var(--foreground-muted)] group-focus-within:text-primary transition-colors duration-300">
                <svg fill="none" stroke="currentColor" viewBox="0 0 24 24" className="w-full h-full">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2.5} d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
                </svg>
              </div>
              <input
                type="text"
                placeholder="Scan transmissions..."
                value={searchQuery}
                onChange={(e) => setSearchQuery(e.target.value)}
                className="input w-full !pl-12 !py-3 !bg-white/5 !rounded-2xl border-white/5 focus:!border-primary/50 focus:!bg-white/10 transition-all font-medium"
              />
            </div>

            <div className="flex gap-2 overflow-x-auto pb-1 no-scrollbar scroll-smooth">
              {['all', 'open', 'pending', 'resolved'].map((f) => {
                const isActive = filter === f;
                return (
                  <button
                    key={f}
                    onClick={() => setFilter(f)}
                    className={`
                      px-5 py-2 rounded-xl text-[10px] font-black uppercase tracking-[0.2em] whitespace-nowrap transition-all flex-shrink-0 border
                      ${isActive
                        ? 'bg-primary text-white border-primary shadow-lg shadow-primary/20 scale-105' 
                        : 'bg-white/5 text-[var(--foreground-muted)] border-white/5 hover:border-white/10 hover:bg-white/10'
                      }
                    `}
                  >
                    {f}
                  </button>
                );
              })}
            </div>
          </div>

          <div className="flex-1 overflow-hidden flex flex-col">
            {isLoadingConversations ? (
              <div className="flex-1 flex flex-col items-center justify-center gap-4">
                <div className="w-12 h-12 border-2 border-primary border-t-transparent rounded-full" />
                <p className="text-[10px] font-black uppercase tracking-widest text-[var(--foreground-muted)]">Synchronizing...</p>
              </div>
            ) : (
              <ConversationList
                conversations={conversations}
                selectedId={selectedConversation?.id}
                onSelect={handleConversationSelect}
              />
            )}
          </div>
        </div>

        {/* Main Content: Message Panel */}
        <div className="flex-1 flex overflow-hidden relative">
          {selectedConversation ? (
            <MessagePanel
              conversation={selectedConversation}
              messages={messages}
              onSendMessage={handleSendMessage}
              isLoading={isLoadingMessages}
              onDeleteConversation={handleDeleteConversation}
              onDeleteContact={handleDeleteContact}
              isOverLimit={authData.organization.is_over_limit}
              draft={selectedConversation ? drafts[selectedConversation.id] : null}
              onDiscardDraft={() => {
                if (selectedConversation) {
                  setDrafts(prev => {
                    const next = { ...prev };
                    delete next[selectedConversation.id];
                    return next;
                  });
                }
              }}
            />
          ) : (
            <EmptyState />
          )}

          {/* Subtle Background Glow */}
          <div className="absolute top-0 right-0 w-[500px] h-[500px] bg-primary/5 blur-[120px] rounded-full -mr-64 -mt-64 pointer-events-none" />
          <div className="absolute bottom-0 left-0 w-[300px] h-[300px] bg-purple-500/5 blur-[100px] rounded-full -ml-32 -mb-32 pointer-events-none" />
        </div>
      </div>
    </DashboardLayout>
  );
}

function EmptyState() {
  return (
    <div className="flex-1 flex flex-col items-center justify-center bg-[var(--background)] p-12 transition-all duration-700 animate-fadeIn">
      <div className="relative group grayscale hover:grayscale-0 transition-all duration-700">
        <div className="absolute inset-0 bg-primary/20 blur-[60px] rounded-full group-hover:bg-primary/40 transition-all duration-700" />
        <div className="relative w-40 h-40 mb-12 rounded-[40px] glass-elevated flex items-center justify-center border-white/10 animate-float shadow-2xl">
          <svg
            className="w-16 h-16 text-primary group-hover:scale-110 transition-transform duration-700"
            fill="none"
            stroke="currentColor"
            viewBox="0 0 24 24"
          >
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              strokeWidth={1}
              d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z"
            />
          </svg>
        </div>
      </div>
      <h2 className="text-3xl font-black tracking-tighter mb-4 text-gradient">Select a Transmission</h2>
      <p className="text-[var(--foreground-muted)] text-center max-w-sm leading-relaxed font-medium">
        Choose a signal from the sidebar to establish a secure link. 
        Your omnichannel communications are ready for processing.
      </p>
      
      <div className="mt-12 flex gap-4">
         <div className="px-4 py-2 glass rounded-xl border-white/5 flex items-center gap-3">
            <div className="w-2 h-2 rounded-full bg-emerald-500 shadow-[0_0_8px_rgba(16,185,129,0.5)]" />
            <span className="text-[10px] font-black uppercase tracking-widest text-[var(--foreground-muted)]">WhatsApp Hub</span>
         </div>
         <div className="px-4 py-2 glass rounded-xl border-white/5 flex items-center gap-3">
            <div className="w-2 h-2 rounded-full bg-pink-500 shadow-[0_0_8px_rgba(236,72,153,0.5)]" />
            <span className="text-[10px] font-black uppercase tracking-widest text-[var(--foreground-muted)]">Instagram Flow</span>
         </div>
      </div>
    </div>
  );
}

