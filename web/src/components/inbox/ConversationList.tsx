'use client';

import { Conversation } from '@/lib/types';
import { Avatar } from '../ui/Avatar';
import { formatDistanceToNow } from '@/lib/utils';
import { Instagram, Facebook, MessageCircle, Hash, UserCircle, Music } from 'lucide-react';

interface ConversationListProps {
  conversations: Conversation[];
  selectedId?: string;
  onSelect: (conversation: Conversation) => void;
}

export function ConversationList({ conversations, selectedId, onSelect }: ConversationListProps) {
  if (conversations.length === 0) {
    return (
      <div className="flex-1 flex flex-col items-center justify-center p-12 text-center animate-fadeIn">
        <div className="w-24 h-24 mb-6 rounded-3xl glass-elevated flex items-center justify-center animate-float">
          <MessageCircle className="w-10 h-10 text-[var(--foreground-muted)] opacity-50" />
        </div>
        <h3 className="text-lg font-bold tracking-tight mb-2">No active transmissions</h3>
        <p className="text-sm text-[var(--foreground-muted)] max-w-[200px]">
          New customer signals will materialize here automatically.
        </p>
      </div>
    );
  }

  return (
    <div className="flex-1 overflow-y-auto no-scrollbar py-2">
      {conversations.map((conversation, index) => {
        const isSelected = selectedId === conversation.id;
        const hasUnread = conversation.unread_count > 0;
        
        return (
          <div
            key={conversation.id}
            onClick={() => onSelect(conversation)}
            className={`
              relative group mx-3 mb-2 p-4 cursor-pointer transition-all duration-300 rounded-[20px] overflow-hidden
              ${isSelected 
                ? 'glass-elevated bg-white/5 ring-1 ring-white/10 shadow-premium' 
                : 'hover:bg-white/5 border border-transparent hover:border-[var(--border)] active:scale-[0.98]'
              }
              ${hasUnread && !isSelected ? 'bg-primary/5' : ''}
              animate-slideIn
            `}
            style={{ 
              animationDelay: `${index * 40}ms`,
            }}
          >
            {/* Shimmer Effect on Selected */}
            {isSelected && <div className="absolute inset-0 shimmer pointer-events-none opacity-30" />}

            <div className="flex items-start gap-4 relative z-10">
              {/* Avatar with Channel Overlay */}
              <div className="relative flex-shrink-0">
                <Avatar 
                  name={conversation.contact.name} 
                  src={conversation.contact.avatar_url}
                  size="md" 
                  className={`transition-transform duration-500 group-hover:scale-110 ${isSelected ? 'border-2 border-primary' : ''}`}
                />
                <div className={`
                    absolute -bottom-1 -right-1 p-1 rounded-full shadow-lg border-2 border-[var(--background)]
                    ${conversation.channel.type === 'whatsapp' ? 'bg-[#25D366]' : 
                      conversation.channel.type === 'instagram' ? 'bg-[#E1306C]' : 
                      conversation.channel.type === 'facebook' ? 'bg-[#0866FF]' :
                      conversation.channel.type === 'tiktok' ? 'bg-black dark:bg-white' :
                      'bg-[#0866FF]'}
                `}>
                  <ChannelIcon type={conversation.channel.type} />
                </div>
              </div>

              {/* Content Container */}
              <div className="flex-1 min-w-0">
                <div className="flex items-center justify-between mb-1.5">
                  <h4 className={`text-sm font-bold truncate tracking-tight transition-colors ${isSelected ? 'text-white' : 'text-[var(--foreground-secondary)] group-hover:text-white'}`}>
                    {conversation.contact.name}
                  </h4>
                  {conversation.last_message_at && (
                    <span className="text-[10px] font-bold uppercase tracking-widest text-[var(--foreground-muted)] opacity-70">
                      {formatDistanceToNow(conversation.last_message_at).toUpperCase()}
                    </span>
                  )}
                </div>
                
                <div className="flex items-center justify-between gap-2">
                  <p className={`text-sm truncate leading-relaxed ${hasUnread ? 'font-bold text-white' : 'text-[var(--foreground-muted)] line-clamp-1'}`}>
                    {conversation.last_message ? (
                      <>
                        {conversation.last_message.sender_type === 'agent' && <span className="opacity-50 font-black text-[10px] uppercase tracking-wider">YOU: </span>}
                        {conversation.last_message.content}
                      </>
                    ) : (
                      <span className="italic opacity-50">Transmission pending...</span>
                    )}
                  </p>
                  
                  {hasUnread && (
                    <div className="flex-shrink-0 w-5 h-5 rounded-full bg-primary flex items-center justify-center shadow-lg shadow-primary/40">
                      <span className="text-[10px] font-bold text-white">
                        {conversation.unread_count > 9 ? '!' : conversation.unread_count}
                      </span>
                    </div>
                  )}
                </div>

                {/* Meta Info Area */}
                <div className="flex items-center gap-3 mt-3">
                  <StatusBadge status={conversation.status} />
                  
                  {conversation.assigned_user && (
                    <div className="flex items-center gap-1.5 px-2 py-1 rounded-lg glass bg-white/5 border border-white/5 transition-all group-hover:border-white/10">
                        <UserCircle className="w-3 h-3 text-[var(--foreground-muted)]" />
                        <span className="text-[9px] font-bold uppercase tracking-widest text-[var(--foreground-muted)]">
                            {conversation.assigned_user.name.split(' ')[0]}
                        </span>
                    </div>
                  )}
                </div>
              </div>
            </div>

            {/* Unread Indicator Dot */}
            {hasUnread && !isSelected && (
                <div className="absolute top-4 right-4 w-2 h-2 rounded-full bg-primary shadow-glow" />
            )}
          </div>
        );
      })}
    </div>
  );
}

function ChannelIcon({ type }: { type: string }) {
    const props = { className: "w-2.5 h-2.5 text-white" };
    switch (type) {
        case 'whatsapp': return <MessageCircle {...props} />;
        case 'instagram': return <Instagram {...props} />;
        case 'facebook': return <Facebook {...props} />;
        case 'tiktok': return <Music {...props} />;
        default: return <Hash {...props} />;
    }
}

function StatusBadge({ status }: { status: Conversation['status'] }) {
  const styles = {
    open: 'bg-emerald-500/10 text-emerald-500 border-emerald-500/20',
    pending: 'bg-amber-500/10 text-amber-500 border-amber-500/20',
    resolved: 'bg-blue-500/10 text-blue-500 border-blue-500/20',
    closed: 'bg-white/5 text-[var(--foreground-muted)] border-white/10',
  };

  return (
    <span className={`text-[9px] font-black uppercase tracking-[0.15em] px-2.5 py-1 rounded-lg border leading-none ${styles[status]}`}>
      {status}
    </span>
  );
}

