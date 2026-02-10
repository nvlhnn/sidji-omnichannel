import { AuthResponse, Channel, ChannelType, Contact, Conversation, ConversationFilter, Message, MessageList, PaginatedResponse, User, TeamMember, InviteTeamMemberInput, CannedResponse, CreateCannedResponseInput, AIConfig, UpdateAIConfigInput } from './types';

const API_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080/api';

class ApiClient {
    private token: string | null = null;

    setToken(token: string | null) {
        this.token = token;
        if (token) {
            localStorage.setItem('access_token', token);
        } else {
            localStorage.removeItem('access_token');
        }
    }

    getToken(): string | null {
        if (this.token) return this.token;
        if (typeof window !== 'undefined') {
            this.token = localStorage.getItem('access_token');
        }
        return this.token;
    }

    private async request<T>(
        endpoint: string,
        options: RequestInit = {}
    ): Promise<T> {
        const headers: HeadersInit = {
            'Content-Type': 'application/json',
            ...options.headers,
        };

        const token = this.getToken();
        if (token) {
            (headers as Record<string, string>)['Authorization'] = `Bearer ${token}`;
        }

        const response = await fetch(`${API_URL}${endpoint}`, {
            ...options,
            headers,
        });

        if (!response.ok) {
            const error = await response.json().catch(() => ({ error: 'Request failed' }));
            throw new Error(error.error || 'Request failed');
        }

        return response.json();
    }

    // Auth
    async register(data: {
        email: string;
        password: string;
        name: string;
        organization_name: string;
    }): Promise<AuthResponse> {
        const response = await this.request<AuthResponse>('/auth/register', {
            method: 'POST',
            body: JSON.stringify(data),
        });
        this.setToken(response.access_token);
        return response;
    }

    async login(email: string, password: string): Promise<AuthResponse> {
        const response = await this.request<AuthResponse>('/auth/login', {
            method: 'POST',
            body: JSON.stringify({ email, password }),
        });
        this.setToken(response.access_token);
        return response;
    }

    async getMe(): Promise<AuthResponse> {
        return this.request<AuthResponse>('/auth/me');
    }

    logout() {
        this.setToken(null);
    }

    // Conversations
    async getConversations(filter?: ConversationFilter): Promise<PaginatedResponse<Conversation>> {
        const params = new URLSearchParams();
        if (filter) {
            Object.entries(filter).forEach(([key, value]) => {
                if (value !== undefined && value !== '') {
                    params.append(key, String(value));
                }
            });
        }
        const query = params.toString();
        return this.request<PaginatedResponse<Conversation>>(
            `/conversations${query ? `?${query}` : ''}`
        );
    }

    async getConversation(id: string): Promise<Conversation> {
        return this.request<Conversation>(`/conversations/${id}`);
    }

    async assignConversation(id: string, userId: string): Promise<void> {
        await this.request(`/conversations/${id}/assign`, {
            method: 'POST',
            body: JSON.stringify({ user_id: userId }),
        });
    }


    async updateConversationStatus(id: string, status: string): Promise<void> {
        await this.request(`/conversations/${id}/status`, {
            method: 'PATCH',
            body: JSON.stringify({ status }),
        });
    }

    async deleteConversation(id: string): Promise<void> {
        await this.request(`/conversations/${id}`, {
            method: 'DELETE',
        });
    }

    async suggestAI(id: string): Promise<{ content: string }> {
        return this.request<{ content: string }>(`/conversations/${id}/suggest`, {
            method: 'POST',
        });
    }

    // Channels
    async getChannels(): Promise<PaginatedResponse<Channel>> {
        const response = await this.request<{ data: Channel[] }>('/channels');
        // adapt response to match PaginatedResponse structure if needed, or update return type
        // backend returns { data: [...] } for channels based on channel_handler.go
        return {
            data: response.data,
            total: response.data.length,
            page: 1,
            limit: 100
        };
    }

    async createChannel(data: {
        type: ChannelType;
        provider?: string;
        name: string;
        access_token: string;
        phone_number_id?: string;
        ig_user_id?: string;
        facebook_page_id?: string;
        business_account_id?: string;
    }): Promise<Channel> {
        return this.request<Channel>('/channels', {
            method: 'POST',
            body: JSON.stringify(data),
        });
    }




    async deleteChannel(id: string): Promise<void> {
        await this.request(`/channels/${id}`, {
            method: 'DELETE',
        });
    }

    async activateChannel(id: string): Promise<void> {
        await this.request(`/channels/${id}/activate`, {
            method: 'POST',
        });
    }

    // Contacts
    async getContacts(page = 1, limit = 20, search = ''): Promise<PaginatedResponse<Contact>> {
        const params = new URLSearchParams({
            page: String(page),
            limit: String(limit),
        });
        if (search) params.append('search', search);

        return this.request<PaginatedResponse<Contact>>(
            `/contacts?${params.toString()}`
        );
    }

    async deleteContact(id: string): Promise<void> {
        await this.request(`/contacts/${id}`, {
            method: 'DELETE',
        });
    }

    // Messages
    async getMessages(conversationId: string, limit = 50): Promise<MessageList> {
        return this.request<MessageList>(
            `/conversations/${conversationId}/messages?limit=${limit}`
        );
    }

    async sendMessage(
        conversationId: string,
        content: string,
        messageType: string = 'text',
        isNote: boolean = false
    ): Promise<Message> {
        return this.request<Message>(`/conversations/${conversationId}/messages`, {
            method: 'POST',
            body: JSON.stringify({ content, message_type: messageType, is_note: isNote }),
        });
    }

    // Organization
    async updateOrganization(data: { name?: string; plan?: string }): Promise<any> {
        return this.request<any>('/team/organization', {
            method: 'PATCH',
            body: JSON.stringify(data),
        });
    }

    // Team Management
    async getTeamMembers(): Promise<TeamMember[]> {
        const response = await this.request<{ data: TeamMember[] }>('/team/members');
        return response.data;
    }

    async inviteTeamMember(data: InviteTeamMemberInput): Promise<TeamMember> {
        return this.request<TeamMember>('/team/members', {
            method: 'POST',
            body: JSON.stringify(data),
        });
    }

    async updateTeamMember(id: string, data: { name?: string; avatar_url?: string }): Promise<User> {
        return this.request<User>(`/team/members/${id}`, {
            method: 'PATCH',
            body: JSON.stringify(data),
        });
    }

    async updateTeamMemberRole(id: string, role: string): Promise<void> {
        await this.request(`/team/members/${id}/role`, {
            method: 'PATCH',
            body: JSON.stringify({ role }),
        });
    }

    async removeTeamMember(id: string): Promise<void> {
        await this.request(`/team/members/${id}`, {
            method: 'DELETE',
        });
    }

    // Canned Responses
    async getCannedResponses(): Promise<CannedResponse[]> {
        const response = await this.request<{ data: CannedResponse[] }>('/canned-responses');
        return response.data;
    }



    async createCannedResponse(data: CreateCannedResponseInput): Promise<CannedResponse> {
        return this.request<CannedResponse>('/canned-responses', {
            method: 'POST',
            body: JSON.stringify(data),
        });
    }

    async deleteCannedResponse(id: string): Promise<void> {
        await this.request(`/canned-responses/${id}`, {
            method: 'DELETE',
        });
    }
    async connectInstagram(accessToken: string, selectedId?: string): Promise<any> {
        return this.request<any>('/channels/instagram/connect', {
            method: 'POST',
            body: JSON.stringify({
                access_token: accessToken,
                selected_id: selectedId
            }),
        });
    }

    async connectWhatsApp(accessToken: string, selectedId?: string): Promise<any> {
        return this.request<any>('/channels/whatsapp/connect', {
            method: 'POST',
            body: JSON.stringify({
                access_token: accessToken,
                selected_id: selectedId
            }),
        });
    }

    async connectFacebook(accessToken: string, selectedId?: string): Promise<any> {
        return this.request<any>('/channels/facebook/connect', {
            method: 'POST',
            body: JSON.stringify({
                access_token: accessToken,
                selected_id: selectedId
            }),
        });
    }

    async discoverMeta(accessToken: string): Promise<{
        facebook_pages: { id: string, name: string }[],
        instagram_users: { id: string, name: string, parent_id: string }[],
        whatsapp_numbers: { id: string, display_name: string, waba_id: string, waba_name: string }[]
    }> {
        return this.request<any>('/channels/discover/meta', {
            method: 'POST',
            body: JSON.stringify({ access_token: accessToken }),
        });
    }

    // AI
    async getAIConfig(channelId: string): Promise<AIConfig> {
        return this.request<AIConfig>(`/channels/${channelId}/ai`);
    }

    async updateAIConfig(channelId: string, data: UpdateAIConfigInput): Promise<AIConfig> {
        return this.request<AIConfig>(`/channels/${channelId}/ai`, {
            method: 'PUT',
            body: JSON.stringify(data),
        });
    }

    async getAIKnowledge(channelId: string): Promise<any[]> {
        return this.request<any[]>(`/channels/${channelId}/ai/knowledge`);
    }

    async addAIKnowledge(channelId: string, content: string): Promise<void> {
        await this.request(`/channels/${channelId}/ai/knowledge`, {
            method: 'POST',
            body: JSON.stringify({ content }),
        });
    }

    async updateAIKnowledge(channelId: string, knowledgeId: string, content: string): Promise<void> {
        await this.request(`/channels/${channelId}/ai/knowledge/${knowledgeId}`, {
            method: 'PUT',
            body: JSON.stringify({ content }),
        });
    }

    async deleteAIKnowledge(channelId: string, knowledgeId: string): Promise<void> {
        await this.request(`/channels/${channelId}/ai/knowledge/${knowledgeId}`, {
            method: 'DELETE',
        });
    }

    async testAIReply(channelId: string, query: string): Promise<{ reply: string; context: string[] }> {
        return this.request<{ reply: string; context: string[] }>(`/channels/${channelId}/ai/test`, {
            method: 'POST',
            body: JSON.stringify({ query }),
        });
    }
}

export const api = new ApiClient();
