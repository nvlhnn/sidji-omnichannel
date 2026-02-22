// User Types
export interface User {
    id: string;
    name: string;
    email: string;
    role: 'admin' | 'supervisor' | 'agent';
    avatar_url?: string;
    status: 'online' | 'away' | 'offline';
}

export interface Organization {
    id: string;
    name: string;
    slug: string;
    plan: 'starter' | 'growth' | 'scale' | 'enterprise';
    subscription_status: 'active' | 'past_due' | 'canceled' | 'trial';
    ai_credits_limit: number;
    ai_credits_used: number;
    message_usage_limit: number;
    message_usage_used: number;
    billing_cycle_start: string;
    is_over_limit?: boolean;
    user_count?: number;
    channel_count?: number;
}

export interface AuthResponse {
    user: User;
    organization: Organization;
    access_token: string;
    expires_in: number;
}

// Channel Types
export type ChannelType = 'whatsapp' | 'instagram' | 'facebook' | 'tiktok';

export interface Channel {
    id: string;
    type: ChannelType;
    provider?: string;
    name: string;
    config?: string; // JSON string
    status: 'active' | 'disconnected' | 'pending';
}

// Contact Types
export interface Contact {
    id: string;
    name: string;
    phone?: string;
    email?: string;
    avatar_url?: string;
    whatsapp_id?: string;
    instagram_id?: string;
    facebook_id?: string;
    tiktok_id?: string;
    tags?: string[];
    created_at: string;
    updated_at: string;
}

// Conversation Types
export type ConversationStatus = 'open' | 'pending' | 'resolved' | 'closed';

export interface Conversation {
    id: string;
    status: ConversationStatus;
    channel: Channel;
    contact: Contact;
    assigned_user?: User;
    last_message?: Message;
    last_message_at?: string;
    unread_count: number;
}

export interface ConversationFilter {
    status?: ConversationStatus;
    channel_id?: string;
    assigned_to?: string;
    unassigned?: boolean;
    search?: string;
    page?: number;
    limit?: number;
}

// Message Types
export type SenderType = 'contact' | 'agent' | 'system' | 'ai' | 'note';
export type MessageType = 'text' | 'image' | 'video' | 'audio' | 'document' | 'sticker';
export type MessageStatus = 'pending' | 'sent' | 'delivered' | 'read' | 'failed';

export interface Message {
    id: string;
    conversation_id: string;
    sender_type: SenderType;
    sender_id: string;
    content: string;
    message_type: MessageType;
    media_url?: string;
    media_mime_type?: string;
    media_file_name?: string;
    reply_to_id?: string;
    status: MessageStatus;
    created_at: string;
    sender?: User;
}

export interface MessageList {
    messages: Message[];
    total_count: number;
    has_more: boolean;
}

// WebSocket Events
export type WSEvent =
    | { event: 'message:new'; data: Message }
    | { event: 'message:status'; data: { id: string; status: MessageStatus } }
    | { event: 'conversation:update'; data: Conversation }
    | { event: 'conversation:ai_draft'; data: { conversation_id: string; content: string; created_at: string } }
    | { event: 'typing'; data: { user_id: string; conversation_id: string } };

// API Response Types
export interface PaginatedResponse<T> {
    data: T[];
    total: number;
    page: number;
    limit: number;
}

// Team Types
export interface TeamMember {
    id: string;
    organization_id: string;
    name: string;
    email: string;
    role: 'admin' | 'supervisor' | 'agent';
    avatar_url?: string;
    status: 'active' | 'inactive';
    last_login_at?: string;
    created_at: string;
}

export interface InviteTeamMemberInput {
    email: string;
    name: string;
    role: 'admin' | 'supervisor' | 'agent';
}

// Canned Response Types
export interface CannedResponse {
    id: string;
    shortcut: string;
    title: string;
    content: string;
    category?: string;
    tags?: string[];
    usage_count: number;
    created_at: string;
    updated_at: string;
}

export interface CreateCannedResponseInput {
    shortcut: string;
    title: string;
    content: string;
    category?: string;
    tags?: string[];
}

export type AIMode = 'manual' | 'auto' | 'hybrid';

export interface AIConfig {
    id: string;
    channel_id: string;
    is_enabled: boolean;
    mode: AIMode;
    persona: string;
    handover_timeout_minutes: number;
    created_at: string;
    updated_at: string;
}

export interface UpdateAIConfigInput {
    is_enabled: boolean;
    mode: AIMode;
    persona: string;
    handover_timeout_minutes: number;
}
