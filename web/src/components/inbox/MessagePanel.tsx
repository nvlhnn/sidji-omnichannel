import { useState, useEffect, useRef, useMemo, useLayoutEffect } from 'react';
import { api } from '@/lib/api';
import { CannedResponse, Message, Conversation } from '@/lib/types';
import { Avatar } from '../ui/Avatar';
import { formatTime, formatDate } from '@/lib/utils';
import { 
  ShieldAlert, 
  Sparkles, 
  Loader2, 
  MoreVertical, 
  Send, 
  Smile, 
  Paperclip, 
  ArrowDown, 
  User as UserIcon,
  Info,
  Trash2,
  X,
  Instagram,
  Facebook,
  MessageCircle,
  Hash
} from 'lucide-react';

interface MessagePanelProps {
  conversation: Conversation;
  messages: Message[];
  onSendMessage: (content: string, isNote?: boolean) => void;
  isLoading?: boolean;
  onDeleteConversation: () => void;
  onDeleteContact: () => void;
  draft?: { content: string; created_at: string } | null;
  onDiscardDraft?: () => void;
  isOverLimit?: boolean;
}

export function MessagePanel({ 
  conversation, 
  messages, 
  onSendMessage, 
  isLoading, 
  onDeleteConversation, 
  onDeleteContact,
  draft,
  onDiscardDraft,
  isOverLimit
}: MessagePanelProps) {
  const [showProfile, setShowProfile] = useState(false);
  const [showOptions, setShowOptions] = useState(false);
  const [showScrollButton, setShowScrollButton] = useState(false);
  const messagesContainerRef = useRef<HTMLDivElement>(null);
  const [inputValue, setInputValue] = useState('');
  const [isSuggesting, setIsSuggesting] = useState(false);

  useLayoutEffect(() => {
    if (messagesContainerRef.current) {
      messagesContainerRef.current.scrollTop = messagesContainerRef.current.scrollHeight;
      setShowScrollButton(false);
      
      // Double-pass ensure instant jump even with dynamic content
      const timeoutId = setTimeout(() => {
        if (messagesContainerRef.current) {
          messagesContainerRef.current.scrollTop = messagesContainerRef.current.scrollHeight;
        }
      }, 0);
      
      return () => clearTimeout(timeoutId);
    }
  }, [messages, conversation.id]);


  const handleScroll = (e: React.UIEvent<HTMLDivElement>) => {
    const { scrollTop, scrollHeight, clientHeight } = e.currentTarget;
    const isBottom = scrollHeight - scrollTop - clientHeight < 100;
    setShowScrollButton(!isBottom);
  };

  const scrollToBottom = () => {
    messagesContainerRef.current?.scrollTo({
      top: messagesContainerRef.current.scrollHeight,
      behavior: 'smooth'
    });
  };

  const groupedMessages = useMemo(() => {
    const groups: { date: string; messages: Message[] }[] = [];
    let currentDate = '';

    messages.forEach((message) => {
      const messageDate = formatDate(message.created_at);
      if (messageDate !== currentDate) {
        currentDate = messageDate;
        groups.push({ date: messageDate, messages: [message] });
      } else {
        groups[groups.length - 1].messages.push(message);
      }
    });

    return groups;
  }, [messages]);

  const handleUseDraft = () => {
    if (draft) {
      setInputValue(draft.content);
      onDiscardDraft?.();
    }
  };

  const handleSuggest = async () => {
    setIsSuggesting(true);
    try {
      await api.suggestAI(conversation.id);
    } catch (err) {
      console.error('Failed to get AI suggestion:', err);
    } finally {
      setIsSuggesting(false);
    }
  };

  const ChannelIcon = () => {
    switch (conversation.channel.type) {
      case 'whatsapp': return <MessageCircle className="w-4 h-4 text-whatsapp" />;
      case 'instagram': return <Instagram className="w-4 h-4 text-instagram" />;
      case 'facebook': return <Facebook className="w-4 h-4 text-facebook" />;
      default: return <Hash className="w-4 h-4 text-[var(--foreground-muted)]" />;
    }
  };

  return (
    <div className="flex-1 flex flex-col h-full bg-[var(--background)] relative overflow-hidden">
       {/* Premium Header */}
       <div className="flex items-center justify-between px-6 py-4 glass relative z-20 border-b border-[var(--border)]">
        <div 
          className="flex items-center gap-4 cursor-pointer group"
          onClick={() => setShowProfile(true)}
        >
          <div className="relative">
            <Avatar 
              name={conversation.contact.name} 
              src={conversation.contact.avatar_url}
              size="md"
              className="border border-white/10 group-hover:border-primary/40 transition-all"
            />
            <div className="absolute -bottom-1 -right-1 bg-[var(--background)] p-1 rounded-full">
              <ChannelIcon />
            </div>
          </div>
          <div>
            <div className="flex items-center gap-2">
              <h3 className="font-bold text-lg leading-tight group-hover:text-primary transition-colors">{conversation.contact.name}</h3>
              <span className="w-2 h-2 rounded-full status-online"></span>
            </div>
            <p className="text-sm text-[var(--foreground-muted)] flex items-center gap-1.5">
              <span className="capitalize">{conversation.channel.type}</span> • 
              <span>{conversation.contact.phone || conversation.contact.email || 'No contact info'}</span>
            </p>
          </div>
        </div>
        
        <div className="flex items-center gap-3 relative">
          <button 
            className="btn-secondary !p-2 rounded-xl group"
            onClick={() => setShowProfile(!showProfile)}
            title="View Profile"
          >
            <UserIcon className="w-5 h-5 text-[var(--foreground-muted)] group-hover:text-primary transition-colors" />
          </button>
          
          <div className="relative">
            <button 
              className="btn-secondary !p-2 rounded-xl group"
              onClick={() => setShowOptions(!showOptions)}
              title="Options"
            >
              <MoreVertical className="w-5 h-5 text-[var(--foreground-muted)] group-hover:text-primary transition-colors" />
            </button>

            {showOptions && (
              <div className="absolute top-full right-0 mt-3 w-56 glass-elevated rounded-2xl py-2 z-50 animate-scale-up border border-[var(--border-light)] overflow-hidden">
                <button className="flex items-center gap-3 w-full text-left px-4 py-2.5 hover:bg-white/5 text-sm transition-colors">
                  <Hash className="w-4 h-4 opacity-70" />
                  Mark as Unread
                </button>
                <button className="flex items-center gap-3 w-full text-left px-4 py-2.5 hover:bg-white/5 text-sm transition-colors">
                  <X className="w-4 h-4 opacity-70" />
                  Close Conversation
                </button>
                <div className="h-[1px] bg-[var(--border)] my-1" />
                <button 
                  onClick={() => { setShowOptions(false); onDeleteConversation(); }}
                  className="flex items-center gap-3 w-full text-left px-4 py-2.5 hover:bg-red-500/10 text-sm text-red-500 transition-colors"
                >
                  <Trash2 className="w-4 h-4" />
                  Delete Conversation
                </button>
              </div>
            )}
          </div>
        </div>
      </div>

      {/* Messages Area */}
      <div 
        ref={messagesContainerRef}
        onScroll={handleScroll}
        className="flex-1 overflow-y-auto p-6 space-y-8" 
        onClick={() => setShowOptions(false)}
      >
        {isLoading ? (
          <div className="flex items-center justify-center h-full">
            <div className="flex flex-col items-center gap-4">
              <div className="relative">
                <div className="w-12 h-12 border-2 border-[var(--primary)] rounded-full border-t-transparent"></div>
              </div>
              <p className="text-[var(--foreground-muted)] font-medium">Establishing secure connection...</p>
            </div>
          </div>
        ) : (
          <>
            {groupedMessages.map((group) => (
              <div key={group.date} className="animate-fadeIn">
                <div className="flex items-center gap-4 mb-8">
                  <div className="flex-1 h-[1px] bg-[var(--border)]"></div>
                  <span className="px-4 py-1.5 rounded-full glass text-[11px] font-bold uppercase tracking-[0.15em] text-[var(--foreground-muted)]">
                    {group.date}
                  </span>
                  <div className="flex-1 h-[1px] bg-[var(--border)]"></div>
                </div>
                
                <div className="space-y-4">
                  {group.messages.map((message) => (
                    <MessageBubble key={message.id} message={message} contact={conversation.contact} />
                  ))}
                </div>
              </div>
            ))}
          </>
        )}
      </div>

      {showScrollButton && (
        <button
          onClick={scrollToBottom}
          className="absolute bottom-28 right-8 p-3 bg-primary text-white rounded-full shadow-premium hover:scale-110 active:scale-95 transition-all z-30 flex items-center justify-center animate-bounce"
          title="Scroll to bottom"
        >
          <ArrowDown className="w-5 h-5" />
        </button>
      )}

      {/* AI Suggestion Bar */}
      {draft && !isOverLimit && (
        <div className="px-6 pb-2">
          <div className="message-bubble-ai px-4 py-3 animate-in slide-in-from-bottom-2 fade-in duration-500 border border-white/10 group">
            <div className="flex items-center justify-between mb-2">
              <div className="flex items-center gap-2">
                <Sparkles className="w-4 h-4 text-white animate-pulse" />
                <span className="text-[10px] font-black uppercase tracking-[0.2em] text-white/90">AI Agent Suggestion</span>
              </div>
              <div className="flex gap-2">
                <button 
                  onClick={onDiscardDraft}
                  className="px-2 py-1 text-[10px] font-bold text-white/60 hover:text-white transition-colors"
                >
                  DISCARD
                </button>
                <button 
                  onClick={handleUseDraft}
                  className="px-4 py-1.5 bg-white text-primary rounded-full text-[10px] font-black uppercase tracking-wider hover:bg-primary-hover hover:text-white transition-all shadow-lg"
                >
                  Use Suggestion
                </button>
              </div>
            </div>
            <p className="text-sm text-white leading-relaxed italic opacity-95">
              "{draft.content}"
            </p>
          </div>
        </div>
      )}

      <MessageInput 
        onSend={onSendMessage} 
        value={inputValue} 
        onChange={setInputValue} 
        disabled={isOverLimit}
        onSuggest={handleSuggest}
        isSuggesting={isSuggesting}
      />

       {/* Premium Sidebar Profile */}
       {showProfile && (
        <div className="absolute inset-y-0 right-0 w-[340px] glass-elevated z-40 animate-in slide-in-from-right duration-300">
             <div className="p-8 h-full flex flex-col">
                <div className="flex items-center justify-between mb-10">
                    <h3 className="font-black text-xs uppercase tracking-[0.2em] text-[var(--foreground-muted)]">Identity Profile</h3>
                    <button onClick={() => setShowProfile(false)} className="p-2 hover:bg-white/5 rounded-full transition-colors text-[var(--foreground-muted)] hover:text-white">
                        <X className="w-5 h-5" />
                    </button>
                </div>
                
                <div className="flex flex-col items-center mb-10">
                    <div className="relative mb-6">
                        <Avatar name={conversation.contact.name} src={conversation.contact.avatar_url} size="xl" className="border-4 border-white/5" />
                        <div className="absolute -bottom-2 -right-2 p-2 bg-primary rounded-full shadow-lg border-2 border-[var(--background)]">
                            <ChannelIcon />
                        </div>
                    </div>
                    <h2 className="text-2xl font-black tracking-tight mb-2">{conversation.contact.name}</h2>
                    <div className="flex items-center gap-2 text-[var(--foreground-muted)] text-sm font-medium">
                        <span className="w-2 h-2 rounded-full bg-success"></span>
                        Active on {conversation.channel.type}
                    </div>
                </div>
                
                <div className="space-y-4 flex-1 overflow-y-auto no-scrollbar pb-6">
                    <div className="glass-card p-5 group hover:bg-white/5 transition-colors">
                        <div className="flex items-center gap-3 mb-3 text-[var(--foreground-muted)]">
                            <Hash className="w-4 h-4" />
                            <span className="text-[10px] font-black uppercase tracking-widest">Digital ID</span>
                        </div>
                        <p className="text-sm font-semibold truncate">{conversation.contact.phone || conversation.contact.email || conversation.contact.id}</p>
                    </div>
                    
                    <div className="glass-card p-5 group hover:bg-white/5 transition-colors">
                        <div className="flex items-center gap-3 mb-3 text-[var(--foreground-muted)]">
                            <Paperclip className="w-4 h-4" />
                            <span className="text-[10px] font-black uppercase tracking-widest">Metadata Tags</span>
                        </div>
                        <div className="flex flex-wrap gap-2">
                            {conversation.contact.tags?.map(tag => (
                                <span key={tag} className="px-3 py-1.5 glass bg-primary/5 text-primary text-[10px] font-black uppercase tracking-wider rounded-lg border-primary/20">{tag}</span>
                            )) || <span className="text-xs text-[var(--foreground-muted)] italic">No descriptive tags</span>}
                        </div>
                    </div>

                    <div className="glass-card p-5 group hover:bg-white/5 transition-colors">
                        <div className="flex items-center gap-3 mb-3 text-[var(--foreground-muted)]">
                            <Info className="w-4 h-4" />
                            <span className="text-[10px] font-black uppercase tracking-widest">System Insights</span>
                        </div>
                        <div className="space-y-3">
                            <div className="flex justify-between items-center text-xs">
                                <span className="text-[var(--foreground-muted)]">Joined</span>
                                <span className="font-medium">{new Date(conversation.contact.created_at).toLocaleDateString()}</span>
                            </div>
                            <div className="flex justify-between items-center text-xs">
                                <span className="text-[var(--foreground-muted)]">Total Reach</span>
                                <span className="font-medium">Direct Channel</span>
                            </div>
                        </div>
                    </div>
                </div>

                <div className="pt-6 border-t border-[var(--border)]">
                    <button 
                        onClick={() => { onDeleteContact(); setShowProfile(false); }}
                        className="w-full flex items-center justify-center gap-2 px-6 py-4 bg-red-500/10 text-red-500 rounded-2xl font-black text-[10px] uppercase tracking-[0.2em] hover:bg-red-500 hover:text-white transition-all border border-red-500/20"
                    >
                        <Trash2 className="w-4 h-4" />
                        Erase Identity Data
                    </button>
                </div>
             </div>
        </div>
      )}
    </div>
  );
}

function MessageBubble({ message, contact }: { message: Message; contact: Conversation['contact'] }) {
  const isAgent = message.sender_type === 'agent';
  const isAI = message.sender_type === 'ai';
  const isNote = message.sender_type === 'note';
  const isSelf = isAgent || isAI || isNote;
  
  const getMediaUrl = (path?: string) => {
    if (!path) return '';
    if (path.startsWith('http')) return path;
    const apiUrl = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';
    const baseUrl = apiUrl.replace(/\/api$/, ''); 
    return `${baseUrl}${path}`;
  };

  const mediaUrl = getMediaUrl(message.media_url);

  return (
    <div className={`flex items-end gap-3 ${isSelf ? 'justify-end' : 'justify-start'} group/bubble`}>
      {!isSelf && (
        <Avatar name={contact.name} src={contact.avatar_url} size="xs" className="mb-1" />
      )}
      
      <div className={`max-w-[75%] flex flex-col ${isSelf ? 'items-end' : 'items-start'}`}>
        <div
          className={`px-5 py-3 !max-w-fit shadow-md relative group/bubble-content ${
            isAI 
              ? 'message-bubble-ai text-white pt-6' 
              : isNote
                ? 'message-bubble-note pt-6'
                : isAgent 
                  ? 'message-bubble-agent text-white' 
                  : 'message-bubble-contact text-[var(--foreground)]'
          } transition-transform duration-200 hover:scale-[1.01]`}
        >
          {isNote && (
            <div className="absolute top-2 left-3 group/note-info">
                <Hash className="w-3 h-3 opacity-40" />
                <div className="absolute bottom-full left-0 mb-2 px-2 py-1 glass bg-black/90 rounded-lg text-[9px] font-bold text-white opacity-0 group-hover/note-info:opacity-100 transition-opacity pointer-events-none whitespace-nowrap">
                  Internal Note
                </div>
            </div>
          )}
          
          {isAI && (
            <div className="absolute top-2 left-3 group/ai-info">
                <Sparkles className="w-3 h-3 text-white/40" />
                <div className="absolute bottom-full left-0 mb-2 px-2 py-1 glass bg-black/90 rounded-lg text-[9px] font-bold text-white opacity-0 group-hover/ai-info:opacity-100 transition-opacity pointer-events-none whitespace-nowrap">
                  AI Assistant
                </div>
            </div>
          )}

          {message.message_type === 'image' && mediaUrl && (
            <div className="mb-3 rounded-xl overflow-hidden glass border border-white/10 group/media">
                <img 
                src={mediaUrl} 
                alt="Uploaded media" 
                className="max-h-[300px] w-auto object-cover cursor-zoom-in group-hover/media:scale-105 transition-transform duration-500"
                onClick={() => window.open(mediaUrl, '_blank')}
                />
            </div>
          )}

          {message.message_type === 'video' && mediaUrl && (
            <video 
              src={mediaUrl} 
              controls 
              className="rounded-xl mb-3 max-w-full glass border border-white/10"
            />
          )}

          {message.message_type === 'audio' && mediaUrl && (
            <div className="mb-3 glass rounded-2xl p-2 border border-white/10">
                <audio 
                src={mediaUrl} 
                controls 
                className="h-10 w-[240px]"
                />
            </div>
          )}

          {message.message_type === 'document' && mediaUrl && (
            <div className="flex items-center gap-4 p-4 glass-card mb-3 border border-white/10 group/doc hover:bg-white/5 transition-colors cursor-pointer" onClick={() => window.open(mediaUrl, '_blank')}>
                <div className="p-3 bg-primary/20 rounded-2xl text-primary group-hover/doc:scale-110 transition-transform">
                    <Paperclip className="w-6 h-6" />
                </div>
                <div className="overflow-hidden">
                    <p className="text-sm font-bold truncate pr-4">{message.media_file_name || 'Attached Asset'}</p>
                    <p className="text-[10px] font-black uppercase tracking-widest opacity-60">Digital Document</p>
                </div>
            </div>
          )}

          <p className="text-sm leading-relaxed whitespace-pre-wrap break-words">{message.content}</p>
        </div>
        
        <div className={`flex items-center gap-2 mt-2 transition-opacity duration-300 ${isSelf ? 'flex-row-reverse' : ''}`}>
          <span className="text-[10px] font-bold text-[var(--foreground-muted)] uppercase tracking-widest opacity-60">
            {formatTime(message.created_at)}
          </span>
          {isSelf && !isNote && (
            <div className="flex items-center">
                <MessageStatusIcon status={message.status} />
            </div>
          )}
          {isAI && <span className="px-2 py-0.5 glass bg-white/10 rounded-full text-[8px] font-black uppercase tracking-widest text-white/80">SYNTHETIC</span>}
        </div>
      </div>
      
      {isSelf && message.sender && (
        <Avatar name={message.sender.name} src={message.sender.avatar_url} size="xs" className="mb-1" />
      )}
    </div>
  );
}

function MessageStatusIcon({ status }: { status: Message['status'] }) {
  const icons = {
    pending: <Loader2 className="w-3 h-3 text-[var(--foreground-muted)]" />,
    sent: <div className="flex"><Send className="w-3 h-3 text-[var(--foreground-muted)]" /></div>,
    delivered: (
      <div className="flex">
        <Send className="w-3 h-3 text-[var(--foreground-muted)]" />
        <Send className="w-3 h-3 text-[var(--foreground-muted)] -ml-1.5" />
      </div>
    ),
    read: (
      <div className="flex">
        <Send className="w-3 h-3 text-primary" />
        <Send className="w-3 h-3 text-primary -ml-1.5" />
      </div>
    ),
    failed: <ShieldAlert className="w-3 h-3 text-[var(--error)]" />,
  };

  return icons[status] || null;
}

function MessageInput({ 
  onSend, 
  value, 
  onChange, 
  disabled, 
  onSuggest, 
  isSuggesting 
}: { 
  onSend: (content: string, isNote: boolean) => void; 
  value: string; 
  onChange: (val: string) => void; 
  disabled?: boolean;
  onSuggest?: () => void;
  isSuggesting?: boolean;
}) {
  const [cannedResponses, setCannedResponses] = useState<CannedResponse[]>([]);
  const [showCanned, setShowCanned] = useState(false);
  const [filteredCanned, setFilteredCanned] = useState<CannedResponse[]>([]);
  const [isNoteMode, setIsNoteMode] = useState(false);
  const inputRef = useRef<HTMLInputElement>(null);

  useEffect(() => {
    api.getCannedResponses().then(setCannedResponses).catch(console.error);
  }, []);

  useEffect(() => {
    if (value.startsWith('/') && !disabled) {
        const query = value.slice(1).toLowerCase();
        const filtered = cannedResponses.filter(c => 
            c.shortcut.toLowerCase().includes(query) || 
            c.content.toLowerCase().includes(query) ||
            c.title.toLowerCase().includes(query)
        );
        setFilteredCanned(filtered);
        setShowCanned(true);
    } else {
        setShowCanned(false);
    }
  }, [value, cannedResponses, disabled]);

  const handleSubmit = (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    if (disabled) return;
    if (value.trim()) {
      onSend(value.trim(), isNoteMode);
      onChange('');
      setShowCanned(false);
    }
  };

  const insertCanned = (content: string) => {
    if (disabled) return;
    onChange(content);
    setShowCanned(false);
    inputRef.current?.focus();
  };

  return (
    <div className="px-6 py-6 relative">
       {showCanned && filteredCanned.length > 0 && (
          <div className="absolute bottom-full left-6 mb-4 w-[340px] glass-elevated rounded-2xl shadow-premium overflow-hidden z-50 animate-scale-up max-h-72 overflow-y-auto no-scrollbar border border-[var(--border-light)]">
              <div className="p-3 bg-white/5 border-b border-[var(--border)]">
                 <p className="text-[10px] font-black uppercase tracking-[0.2em] text-[var(--foreground-muted)] px-1">Quick Responses</p>
              </div>
              {filteredCanned.map(cr => (
                  <button 
                    key={cr.id}
                    onClick={() => insertCanned(cr.content)}
                    className="w-full text-left p-4 hover:bg-white/5 border-b border-[var(--border)] last:border-0 transition-all flex flex-col gap-1 active:scale-[0.99]"
                  >
                      <div className="flex justify-between items-center">
                          <span className="font-black text-xs text-primary uppercase tracking-widest">/{cr.shortcut}</span>
                          <span className="text-[10px] font-bold text-[var(--foreground-muted)]">{cr.title}</span>
                      </div>
                      <p className="text-xs text-[var(--foreground-muted)] truncate opacity-80">{cr.content}</p>
                  </button>
              ))}
          </div>
      )}

      {disabled && (
          <div className="absolute inset-x-6 inset-y-6 glass bg-amber-500/5 z-10 flex items-center justify-center pointer-events-none rounded-[20px] border border-amber-500/20">
             <div className="flex items-center gap-3 px-6 py-2.5 rounded-full bg-[var(--background)] shadow-premium border border-amber-500/30">
                <ShieldAlert className="w-4 h-4 text-amber-500" />
                <span className="text-[10px] font-black text-amber-500 uppercase tracking-[0.2em]">Usage Limit Reached • Upgrade Plan</span>
             </div>
          </div>
      )}

      <form onSubmit={handleSubmit} className={`flex items-center gap-4 ${disabled ? 'opacity-40 grayscale pointer-events-none' : ''}`}>
        <div className="flex items-center gap-2">
            <button
            type="button"
            className="btn-secondary !p-3 rounded-2xl group active:scale-95 transition-all"
            title="Attach Assets"
            >
                <Paperclip className="w-5 h-5 text-[var(--foreground-muted)] group-hover:text-primary transition-colors" />
            </button>
            <button
            type="button"
            onClick={() => setIsNoteMode(!isNoteMode)}
            disabled={disabled}
            title={isNoteMode ? "Synthetic Message" : "Secure Internal Note"}
            className={`!p-3 rounded-2xl transition-all active:scale-95 ${isNoteMode ? 'bg-amber-500/10 text-amber-500 border border-amber-500/30' : 'btn-secondary text-[var(--foreground-muted)] hover:text-primary'}`}
            >
                {isNoteMode ? <Hash className="w-5 h-5" /> : <UserIcon className="w-5 h-5" />}
            </button>
        </div>

        <div className="flex-1 relative group/input">
          <input
            ref={inputRef}
            type="text"
            name="message"
            disabled={disabled}
            value={value}
            onChange={(e) => onChange(e.target.value)}
            placeholder={disabled ? "Messaging Restricted" : isNoteMode ? "Synthesizing internal note..." : "Enter transmission (or /shortcuts)..."}
            className={`input w-full !py-4 !px-6 !pr-32 !rounded-[20px] font-medium !text-base transition-all ${isNoteMode ? '!border-amber-500/50 !focus:shadow-amber-500/10' : ''}`}
            autoComplete="off"
            maxLength={500}
          />
          
          <div className="absolute right-4 top-1/2 -translate-y-1/2 flex items-center gap-3">
            {!disabled && (
                <>
                    <span className={`text-[9px] font-black tracking-widest ${value.length >= 450 ? 'text-red-500 scale-110' : 'text-[var(--foreground-muted)]'} transition-all`}>
                      {value.length}/500
                    </span>
                    <button
                        type="button"
                        onClick={onSuggest}
                        disabled={isSuggesting}
                        title="Quantum AI Suggestion"
                        className={`p-2 rounded-xl transition-all duration-500 ${isSuggesting ? 'text-primary animate-spin' : 'text-primary/60 hover:text-primary hover:scale-125 hover:rotate-12 bg-primary/5'}`}
                    >
                        {isSuggesting ? (
                            <Loader2 className="w-5 h-5" />
                        ) : (
                            <Sparkles className="w-5 h-5" />
                        )}
                    </button>
                    <button
                        type="button"
                        className="p-1 rounded-lg text-[var(--foreground-muted)] hover:text-white transition-colors"
                    >
                        <Smile className="w-5 h-5" />
                    </button>
                </>
            )}
          </div>
        </div>

        <button
          type="submit"
          disabled={disabled || !value.trim()}
          className={`${isNoteMode ? 'bg-amber-600 hover:bg-amber-500 shadow-amber-500/20' : 'btn-primary'} !p-4 !rounded-[20px] disabled:opacity-30 disabled:cursor-not-allowed transition-all active:scale-90 shadow-premium`}
        >
          <Send className={`w-6 h-6 text-white ${value.trim() ? 'animate-pulse' : ''}`} />
        </button>
      </form>
    </div>
  );
}

