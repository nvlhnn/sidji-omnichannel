'use client';

import React, { useState } from 'react';
import { api } from '@/lib/api';
import { ChannelType } from '@/lib/types';
import { 
  X, 
  MessageCircle, 
  Instagram, 
  Facebook, 
  Music,
  Zap, 
  Settings, 
  ChevronRight,
  Info,
  CheckCircle2,
  Lock
} from 'lucide-react';

interface AddChannelModalProps {
  onClose: () => void;
  onSuccess: () => void;
  initialType?: ChannelType;
}

type ConnectMode = 'fast' | 'manual';

export function AddChannelModal({ onClose, onSuccess, initialType = 'whatsapp' }: AddChannelModalProps) {
  const [type, setType] = useState<ChannelType>(initialType);
  const [mode, setMode] = useState<ConnectMode>('fast');
  const [name, setName] = useState('');
  const [accessToken, setAccessToken] = useState('');
  const [phoneNumberId, setPhoneNumberId] = useState('');
  const [businessAccountId, setBusinessAccountId] = useState('');
  const [igUserId, setIgUserId] = useState('');
  const [facebookPageId, setFacebookPageId] = useState('');
  const [tiktokId, setTiktokId] = useState('');
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [discoveryResult, setDiscoveryResult] = useState<{
    facebook_pages: { id: string, name: string }[],
    instagram_users: { id: string, name: string, parent_id: string }[],
    whatsapp_numbers: { id: string, display_name: string, waba_id: string, waba_name: string }[]
  } | null>(null);
  const [step, setStep] = useState<'init' | 'discovery' | 'manual'>('init');

  const handleMetaLogin = () => {
    // @ts-ignore
    if (!window.FB) {
      alert('Facebook SDK not loaded. Please check your internet connection and blockers.');
      return;
    }

    const scope = 'instagram_basic,instagram_manage_messages,pages_manage_metadata,pages_show_list,pages_messaging';

    const loginOptions = type === 'whatsapp' 
      ? {
          config_id: '1540306283537446', // REPLACE WITH YOUR ACTUAL CONFIG ID
          response_type: 'code',
          override_default_response_type: true,
          extras: {
            setup: {}
          }
        }
      : type === 'facebook'
        ? { scope: 'pages_manage_metadata,pages_show_list,pages_messaging' }
        : { scope };

    // @ts-ignore
    window.FB.login((response) => {
      if (response.authResponse) {
        const token = response.authResponse.accessToken || response.authResponse.code;
        setAccessToken(token);
        handleDiscovery(token);
      }
    }, loginOptions);
  };

  const handleDiscovery = async (token: string) => {
    try {
      setIsSubmitting(true);
      const result = await api.discoverMeta(token);
      setDiscoveryResult(result);
      
      // Auto-connect if only one account total across all types
      const totalAccounts = result.facebook_pages.length + result.instagram_users.length + result.whatsapp_numbers.length;
      
      if (totalAccounts === 0) {
        alert("No accounts found with this Meta login. Make sure you've authorized the correct pages/accounts.");
        setIsSubmitting(false);
        return;
      }

      if (totalAccounts === 1) {
        if (result.facebook_pages.length === 1) handleFacebookConnect(token, result.facebook_pages[0].id);
        else if (result.instagram_users.length === 1) handleInstagramConnect(token, result.instagram_users[0].id);
        else if (result.whatsapp_numbers.length === 1) handleWhatsAppConnect(token, result.whatsapp_numbers[0].id);
      } else {
        setStep('discovery');
        setIsSubmitting(false);
      }
    } catch (error) {
      console.error('Discovery failed:', error);
      alert('Failed to discover accounts. Please try again.');
      setIsSubmitting(false);
    }
  };

  const handleFacebookConnect = async (token: string, selectedId?: string) => {
    try {
      setIsSubmitting(true);
      await api.connectFacebook(token, selectedId);
      onSuccess();
    } catch (error) {
      console.error('Failed to connect Facebook:', error);
      alert('Failed to connect Facebook.');
    } finally {
      setIsSubmitting(false);
    }
  };

  const handleInstagramConnect = async (token: string, selectedId?: string) => {
    try {
      setIsSubmitting(true);
      await api.connectInstagram(token, selectedId);
      onSuccess();
    } catch (error) {
      console.error('Failed to connect Instagram:', error);
      alert('Failed to connect Instagram.');
    } finally {
      setIsSubmitting(false);
    }
  };

  const handleWhatsAppConnect = async (token: string, selectedId?: string) => {
    try {
      setIsSubmitting(true);
      await api.connectWhatsApp(token, selectedId);
      onSuccess();
    } catch (error) {
      console.error('Failed to connect WhatsApp:', error);
      alert('Failed to connect WhatsApp.');
    } finally {
      setIsSubmitting(false);
    }
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      setIsSubmitting(true);
      const channel = await api.createChannel({
        type,
        provider: 'meta',
        name,
        access_token: accessToken,
        phone_number_id: type === 'whatsapp' ? phoneNumberId : undefined,
        business_account_id: type === 'whatsapp' ? businessAccountId : undefined,
        ig_user_id: type === 'instagram' ? igUserId : undefined,
        facebook_page_id: type === 'facebook' ? facebookPageId : undefined,
        tiktok_open_id: type === 'tiktok' ? tiktokId : undefined,
      });
      await api.activateChannel(channel.id);
      onSuccess();
    } catch (error) {
      console.error('Failed to create channel:', error);
      alert('Failed to create channel. Please check your inputs.');
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <div className="fixed inset-0 bg-[#0a0a0f]/80 backdrop-blur-md flex items-center justify-center z-[100] p-4">
      <div className="bg-[var(--background)] w-full max-w-lg rounded-3xl shadow-2xl overflow-hidden border border-[var(--border)] animate-fadeIn flex flex-col max-h-[90vh]">
        {/* Header */}
        <div className="px-6 py-5 border-b border-[var(--border)] flex justify-between items-center bg-[var(--background-secondary)]/50">
          <div>
            <h3 className="text-xl font-bold tracking-tight">Connect Channel</h3>
            <p className="text-xs text-[var(--foreground-muted)] mt-0.5">Add a new communication line to Sidji</p>
          </div>
          <button 
            onClick={onClose} 
            className="w-10 h-10 rounded-full flex items-center justify-center hover:bg-[var(--background-tertiary)] text-[var(--foreground-muted)] hover:text-[var(--foreground)] transition-all duration-200"
          >
            <X size={20} />
          </button>
        </div>
        
        <div className="flex-1 overflow-y-auto custom-scrollbar">
          <div className="p-6 space-y-8">
            {/* Step 1: Select Platform */}
            <div className="space-y-4">
              <div className="flex items-center gap-2 px-1">
                <span className="w-6 h-6 rounded-full bg-[var(--primary)]/10 text-[var(--primary)] flex items-center justify-center text-[10px] font-bold">1</span>
                <label className="text-sm font-semibold uppercase tracking-wider text-[var(--foreground-secondary)]">Select Platform</label>
              </div>
              <div className="grid grid-cols-2 gap-3">
                {[
                  { id: 'whatsapp', label: 'WhatsApp', icon: MessageCircle, color: 'whatsapp', brand: '#25D366' },
                  { id: 'instagram', label: 'Instagram', icon: Instagram, color: 'instagram', brand: '#E1306C' },
                  { id: 'facebook', label: 'Facebook', icon: Facebook, color: 'facebook', brand: '#0866FF' },
                  { id: 'tiktok', label: 'TikTok', icon: Music, color: 'tiktok', brand: '#000000' }
                ].map((p) => {
                  const Icon = p.icon;
                  const isActive = type === p.id;
                  return (
                    <button
                      key={p.id}
                      type="button"
                      onClick={() => setType(p.id as ChannelType)}
                      className={`relative overflow-hidden group p-4 rounded-2xl border transition-all duration-300 flex flex-col items-center gap-3
                        ${isActive 
                          ? 'border-transparent shadow-lg scale-[1.02]' 
                          : 'border-[var(--border)] bg-[var(--background-secondary)]/30 hover:bg-[var(--background-secondary)]'
                        }
                      `}
                      style={isActive ? { background: `linear-gradient(135deg, ${p.brand}22, ${p.brand}11)`, border: `1px solid ${p.brand}55` } : {}}
                    >
                      <div className={`w-12 h-12 rounded-xl flex items-center justify-center transition-transform duration-300 group-hover:scale-110
                        ${isActive ? 'text-white' : 'text-[var(--foreground-muted)] bg-[var(--background-tertiary)]'}
                      `}
                      style={isActive ? { background: p.brand } : {}}
                      >
                        <Icon size={24} />
                      </div>
                      <span className={`text-sm font-bold ${isActive ? 'text-[var(--foreground)]' : 'text-[var(--foreground-muted)]'}`}>
                        {p.label}
                      </span>
                      {isActive && (
                        <div className="absolute top-2 right-2">
                          <CheckCircle2 size={16} style={{ color: p.brand }} />
                        </div>
                      )}
                    </button>
                  );
                })}
              </div>
            </div>

            {/* Step 2: Connection Method */}
            <div className="space-y-4">
              <div className="flex items-center gap-2 px-1">
                <span className="w-6 h-6 rounded-full bg-[var(--primary)]/10 text-[var(--primary)] flex items-center justify-center text-[10px] font-bold">2</span>
                <label className="text-sm font-semibold uppercase tracking-wider text-[var(--foreground-secondary)]">Setup Method</label>
              </div>

              <div className="flex p-1 bg-[var(--background-secondary)] rounded-2xl border border-[var(--border)]">
                <button
                  type="button"
                  onClick={() => setMode('fast')}
                  className={`flex-1 py-3 px-4 rounded-xl text-sm font-semibold flex items-center justify-center gap-2 transition-all duration-200
                    ${mode === 'fast' ? 'bg-[var(--primary)] text-white shadow-md' : 'text-[var(--foreground-muted)] hover:text-[var(--foreground)]'}
                  `}
                >
                  <Zap size={16} />
                  Fast Connect
                </button>
                <button
                  type="button"
                  onClick={() => setMode('manual')}
                  className={`flex-1 py-3 px-4 rounded-xl text-sm font-semibold flex items-center justify-center gap-2 transition-all duration-200
                    ${mode === 'manual' ? 'bg-[var(--primary)] text-white shadow-md' : 'text-[var(--foreground-muted)] hover:text-[var(--foreground)]'}
                  `}
                >
                  <Settings size={16} />
                  Advanced
                </button>
              </div>

              {mode === 'fast' ? (
                <div className="animate-fadeIn space-y-4">
                  {step === 'init' ? (
                    type === 'tiktok' ? (
                      <div className="bg-gradient-to-br from-gray-900/5 to-black/5 dark:from-gray-100/5 dark:to-white/5 border border-black/10 dark:border-white/10 rounded-2xl p-5 space-y-4">
                        <div className="flex gap-4">
                          <div className="w-10 h-10 rounded-lg bg-black/10 dark:bg-white/10 text-black dark:text-white flex items-center justify-center flex-shrink-0">
                            <Music size={20} />
                          </div>
                          <div className="space-y-1">
                            <h4 className="text-sm font-bold text-gray-900 dark:text-gray-100">TikTok OAuth Login</h4>
                            <p className="text-xs text-[var(--foreground-muted)] leading-relaxed">
                              Sign in with your TikTok account to automatically authorize sidji to manage your direct messages.
                            </p>
                          </div>
                        </div>
                        
                        <button
                          type="button"
                          onClick={() => {
                            // TikTok URL format per official documentation
                            const clientKey = process.env.NEXT_PUBLIC_TIKTOK_CLIENT_KEY || 'YOUR_CLIENT_KEY';
                            const redirectUri = typeof window !== 'undefined' ? `${window.location.protocol}//${window.location.host}/auth/tiktok/callback` : '';
                            const state = Math.random().toString(36).substring(7);
                            const url = `https://www.tiktok.com/v2/auth/authorize/?client_key=${clientKey}&response_type=code&scope=user.info.basic&redirect_uri=${encodeURIComponent(redirectUri)}&state=${state}`;
                            
                            if (clientKey === 'YOUR_CLIENT_KEY') {
                              alert("TikTok Fast Connect requires NEXT_PUBLIC_TIKTOK_CLIENT_KEY in your deployment environment variables.");
                            } else {
                              window.location.href = url;
                            }
                          }}
                          disabled={isSubmitting}
                          className="group w-full py-4 bg-black hover:bg-gray-800 dark:bg-white dark:hover:bg-gray-200 text-white dark:text-black rounded-xl font-bold flex items-center justify-center gap-3 transition-all duration-300 hover:shadow-lg disabled:opacity-50"
                        >
                          <Music size={20} className="fill-current" />
                          <span>Connect with TikTok</span>
                          {!isSubmitting && <ChevronRight size={18} className="group-hover:translate-x-1 transition-transform" />}
                        </button>

                        <div className="flex items-center justify-center gap-4 text-[10px] text-[var(--foreground-muted)] font-medium">
                          <div className="flex items-center gap-1"><Lock size={12} /> Secure Login Flow</div>
                          <div className="flex items-center gap-1"><CheckCircle2 size={12} /> Official API Partner</div>
                        </div>
                      </div>
                    ) : (
                      <div className="bg-gradient-to-br from-blue-500/5 to-purple-500/5 border border-blue-500/10 rounded-2xl p-5 space-y-4">
                        <div className="flex gap-4">
                          <div className="w-10 h-10 rounded-lg bg-blue-500/20 text-blue-500 flex items-center justify-center flex-shrink-0">
                            <Info size={20} />
                          </div>
                          <div className="space-y-1">
                            <h4 className="text-sm font-bold text-blue-400">Meta Embedded Signup</h4>
                            <p className="text-xs text-[var(--foreground-muted)] leading-relaxed">
                              The official way to connect your {type === 'whatsapp' ? 'WhatsApp Business' : type === 'facebook' ? 'Facebook Page' : 'Instagram Professional'} account. Secure, fast, and no technical skills required.
                            </p>
                          </div>
                        </div>
                        
                        <button
                          type="button"
                          onClick={() => handleMetaLogin()}
                          disabled={isSubmitting}
                          className="group w-full py-4 bg-[#1877F2] hover:bg-[#166fe5] text-white rounded-xl font-bold flex items-center justify-center gap-3 transition-all duration-300 hover:shadow-lg hover:shadow-[#1877F2]/20 disabled:opacity-50"
                        >
                          <Facebook size={20} fill="currentColor" />
                          <span>{isSubmitting ? 'Discovering...' : `Connect with Facebook`}</span>
                          {!isSubmitting && <ChevronRight size={18} className="group-hover:translate-x-1 transition-transform" />}
                        </button>

                        <div className="flex items-center justify-center gap-4 text-[10px] text-[var(--foreground-muted)] font-medium">
                          <div className="flex items-center gap-1"><Lock size={12} /> Encrypted</div>
                          <div className="flex items-center gap-1"><CheckCircle2 size={12} /> Official Web API</div>
                        </div>
                      </div>
                    )
                  ) : (
                    <div className="space-y-4">
                      <div className="bg-[var(--background-secondary)] border border-[var(--border)] rounded-2xl p-5 space-y-4">
                        <h4 className="text-sm font-bold flex items-center gap-2">
                          <CheckCircle2 size={16} className="text-green-500" />
                          Authenticated. Select an account:
                        </h4>
                        
                        <div className="space-y-2 max-h-[300px] overflow-y-auto pr-1 custom-scrollbar">
                          {(discoveryResult?.facebook_pages?.length || 0) > 0 && (
                            <div className="space-y-2">
                              <p className="text-[10px] font-bold uppercase tracking-widest text-[var(--foreground-muted)] ml-1">Facebook Pages</p>
                              {discoveryResult?.facebook_pages.map(page => (
                                <button
                                  key={page.id}
                                  onClick={() => handleFacebookConnect(accessToken, page.id)}
                                  className="w-full flex items-center justify-between p-3 rounded-xl border border-[var(--border)] hover:border-[#0866FF] hover:bg-[#0866FF]/5 transition-all group"
                                >
                                  <div className="flex items-center gap-3">
                                    <div className="w-8 h-8 rounded-lg bg-[#0866FF]/10 text-[#0866FF] flex items-center justify-center">
                                      <Facebook size={16} />
                                    </div>
                                    <span className="text-sm font-medium">{page.name}</span>
                                  </div>
                                  <ChevronRight size={16} className="text-[var(--foreground-muted)] group-hover:text-[#0866FF] group-hover:translate-x-1 transition-all" />
                                </button>
                              ))}
                            </div>
                          )}

                          {(discoveryResult?.instagram_users?.length || 0) > 0 && (
                            <div className="space-y-2 mt-4">
                              <p className="text-[10px] font-bold uppercase tracking-widest text-[var(--foreground-muted)] ml-1">Instagram Business</p>
                              {discoveryResult?.instagram_users.map(ig => (
                                <button
                                  key={ig.id}
                                  onClick={() => handleInstagramConnect(accessToken, ig.id)}
                                  className="w-full flex items-center justify-between p-3 rounded-xl border border-[var(--border)] hover:border-[#E1306C] hover:bg-[#E1306C]/5 transition-all group"
                                >
                                  <div className="flex items-center gap-3">
                                    <div className="w-8 h-8 rounded-lg bg-[#E1306C]/10 text-[#E1306C] flex items-center justify-center">
                                      <Instagram size={16} />
                                    </div>
                                    <span className="text-sm font-medium">{ig.name}</span>
                                  </div>
                                  <ChevronRight size={16} className="text-[var(--foreground-muted)] group-hover:text-[#E1306C] group-hover:translate-x-1 transition-all" />
                                </button>
                              ))}
                            </div>
                          )}

                          {(discoveryResult?.whatsapp_numbers?.length || 0) > 0 && (
                            <div className="space-y-2 mt-4">
                              <p className="text-[10px] font-bold uppercase tracking-widest text-[var(--foreground-muted)] ml-1">WhatsApp Business</p>
                              {discoveryResult?.whatsapp_numbers.map(num => (
                                <button
                                  key={num.id}
                                  onClick={() => handleWhatsAppConnect(accessToken, num.id)}
                                  className="w-full flex items-center justify-between p-3 rounded-xl border border-[var(--border)] hover:border-[#25D366] hover:bg-[#25D366]/5 transition-all group"
                                >
                                  <div className="flex items-center gap-3">
                                    <div className="w-8 h-8 rounded-lg bg-[#25D366]/10 text-[#25D366] flex items-center justify-center">
                                      <MessageCircle size={16} />
                                    </div>
                                    <div className="text-left">
                                      <span className="text-sm font-medium block">{num.display_name}</span>
                                      <span className="text-[10px] text-[var(--foreground-muted)]">{num.waba_name}</span>
                                    </div>
                                  </div>
                                  <ChevronRight size={16} className="text-[var(--foreground-muted)] group-hover:text-[#25D366] group-hover:translate-x-1 transition-all" />
                                </button>
                              ))}
                            </div>
                          )}
                        </div>

                        <button 
                          onClick={() => { setStep('init'); setDiscoveryResult(null); }}
                          className="text-xs text-[var(--primary)] hover:underline font-medium w-full text-center mt-2"
                        >
                          Try another account
                        </button>
                      </div>
                    </div>
                  )}
                </div>
              ) : (
                <form onSubmit={handleSubmit} className="animate-fadeIn space-y-4">
                  <div className="space-y-4 bg-[var(--background-secondary)]/20 p-5 rounded-2xl border border-[var(--border)]">
                    <div className="space-y-1.5">
                      <label className="text-[10px] font-bold uppercase tracking-widest text-[var(--foreground-muted)] ml-1">Internal Name</label>
                      <input
                        type="text"
                        required
                        value={name}
                        onChange={(e) => setName(e.target.value)}
                        placeholder={`e.g. My ${type.charAt(0).toUpperCase() + type.slice(1)} Channel`}
                        className="input w-full bg-[var(--background-secondary)]"
                      />
                    </div>

                    <div className="space-y-1.5">
                      <label className="text-[10px] font-bold uppercase tracking-widest text-[var(--foreground-muted)] ml-1">Meta Access Token</label>
                      <input
                        type="password"
                        required
                        value={accessToken}
                        onChange={(e) => setAccessToken(e.target.value)}
                        placeholder="EAAG..."
                        className="input w-full font-mono text-sm bg-[var(--background-secondary)]"
                      />
                    </div>

                    {type === 'whatsapp' && (
                      <div className="grid grid-cols-2 gap-3 pb-2">
                        <div className="space-y-1.5">
                          <label className="text-[10px] font-bold uppercase tracking-widest text-[var(--foreground-muted)] ml-1">Phone ID</label>
                          <input
                            type="text"
                            required
                            value={phoneNumberId}
                            onChange={(e) => setPhoneNumberId(e.target.value)}
                            placeholder="10..."
                            className="input w-full font-mono text-sm bg-[var(--background-secondary)]"
                          />
                        </div>
                        <div className="space-y-1.5">
                          <label className="text-[10px] font-bold uppercase tracking-widest text-[var(--foreground-muted)] ml-1">WABA ID</label>
                          <input
                            type="text"
                            required
                            value={businessAccountId}
                            onChange={(e) => setBusinessAccountId(e.target.value)}
                            placeholder="11..."
                            className="input w-full font-mono text-sm bg-[var(--background-secondary)]"
                          />
                        </div>
                      </div>
                    )}

                    {type === 'instagram' && (
                      <div className="space-y-1.5 pb-2">
                        <label className="text-[10px] font-bold uppercase tracking-widest text-[var(--foreground-muted)] ml-1">Instagram User ID</label>
                        <input
                          type="text"
                          required
                          value={igUserId}
                          onChange={(e) => setIgUserId(e.target.value)}
                          placeholder="From Meta Dashboard"
                          className="input w-full font-mono text-sm bg-[var(--background-secondary)]"
                        />
                      </div>
                    )}

                    {type === 'facebook' && (
                      <div className="space-y-1.5 pb-2">
                        <label className="text-[10px] font-bold uppercase tracking-widest text-[var(--foreground-muted)] ml-1">Facebook Page ID</label>
                        <input
                          type="text"
                          required
                          value={facebookPageId}
                          onChange={(e) => setFacebookPageId(e.target.value)}
                          placeholder="From Meta Dashboard"
                          className="input w-full font-mono text-sm bg-[var(--background-secondary)]"
                        />
                      </div>
                    )}

                    {type === 'tiktok' && (
                      <div className="space-y-1.5 pb-2">
                        <label className="text-[10px] font-bold uppercase tracking-widest text-[var(--foreground-muted)] ml-1">TikTok Account ID</label>
                        <input
                          type="text"
                          required
                          value={tiktokId}
                          onChange={(e) => setTiktokId(e.target.value)}
                          placeholder="From TikTok Developer Dashboard"
                          className="input w-full font-mono text-sm bg-[var(--background-secondary)]"
                        />
                      </div>
                    )}

                    <button
                      type="submit"
                      disabled={isSubmitting || !name || !accessToken}
                      className="btn-primary w-full py-3 h-12 flex items-center justify-center font-bold"
                    >
                      {isSubmitting ? 'Connecting...' : 'Finish Setup'}
                    </button>
                  </div>
                </form>
              )}
            </div>
          </div>
        </div>

        {/* Footer */}
        <div className="p-4 bg-[var(--background-secondary)]/50 border-t border-[var(--border)] text-center">
          <p className="text-[10px] text-[var(--foreground-muted)] bg-[var(--background-tertiary)] inline-block py-1 px-3 rounded-full border border-[var(--border)]">
            Powered by Meta Graph API v24.0
          </p>
        </div>
      </div>
    </div>
  );
}
