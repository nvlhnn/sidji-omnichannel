'use client';

import { Suspense, useEffect, useState } from 'react';
import { useRouter, useSearchParams } from 'next/navigation';
import { api } from '@/lib/api';
import { Loader2, AlertCircle, CheckCircle2 } from 'lucide-react';

function TikTokCallbackContent() {
  const router = useRouter();
  const searchParams = useSearchParams();
  const [status, setStatus] = useState<'loading' | 'success' | 'error'>('loading');
  const [errorMessage, setErrorMessage] = useState('');

  useEffect(() => {
    const code = searchParams.get('code');
    const error = searchParams.get('error');
    
    // Check if we came back with an error from TikTok
    if (error) {
       setStatus('error');
       setErrorMessage(`Authentication rejected or failed: ${error}`);
       return;
    }

    if (!code) {
      setStatus('error');
      setErrorMessage('No authorization code provided by TikTok.');
      return;
    }

    // Process the code
    const connectTikTok = async () => {
      try {
        await api.connectTikTok(code);
        setStatus('success');
        // Redirect back to settings after a brief delay
        setTimeout(() => {
          router.push('/settings/channels');
        }, 2000);
      } catch (err: any) {
        console.error('Failed to connect TikTok:', err);
        setStatus('error');
        setErrorMessage(err.response?.data?.error || err.message || 'Failed to exchange authorization code.');
      }
    };

    connectTikTok();
  }, [searchParams, router]);

  return (
    <div className="min-h-screen flex items-center justify-center bg-[#0a0a0f] text-white p-4">
      <div className="max-w-md w-full glass rounded-3xl p-8 text-center space-y-6 animate-fadeIn">
        {status === 'loading' && (
          <div className="flex flex-col items-center gap-4">
            <div className="relative">
              <div className="w-16 h-16 border-4 border-white/10 rounded-full animate-pulse blur-[1px]"></div>
              <Loader2 className="w-16 h-16 text-white absolute top-0 left-0 animate-spin" />
            </div>
            <div className="space-y-2">
              <h2 className="text-xl font-bold tracking-tight">Connecting TikTok...</h2>
              <p className="text-sm text-[var(--foreground-muted)]">Please wait while we finalize your channel securely.</p>
            </div>
          </div>
        )}

        {status === 'success' && (
          <div className="flex flex-col items-center gap-4 animate-scaleIn">
            <div className="w-16 h-16 bg-white text-black rounded-full flex items-center justify-center shadow-[0_0_30px_rgba(255,255,255,0.3)]">
              <CheckCircle2 size={32} />
            </div>
            <div className="space-y-2">
              <h2 className="text-xl font-bold tracking-tight text-white">Successfully Connected!</h2>
              <p className="text-sm text-[var(--foreground-muted)]">Redirecting you back to your Sidji channels...</p>
            </div>
          </div>
        )}

        {status === 'error' && (
          <div className="flex flex-col items-center gap-4 animate-scaleIn">
            <div className="w-16 h-16 bg-red-500/10 text-red-500 rounded-full flex items-center justify-center">
              <AlertCircle size={32} />
            </div>
            <div className="space-y-4 w-full">
              <div>
                <h2 className="text-xl font-bold tracking-tight text-white">Connection Failed</h2>
                <p className="text-sm text-[var(--foreground-muted)] mt-1">{errorMessage}</p>
              </div>
              <button
                onClick={() => router.push('/settings/channels')}
                className="w-full py-3 px-4 bg-white/10 hover:bg-white/20 text-white rounded-xl text-sm font-semibold transition-colors"
              >
                Back to Settings
              </button>
            </div>
          </div>
        )}
      </div>
    </div>
  );
}

export default function TikTokCallbackPage() {
  return (
    <Suspense fallback={
      <div className="min-h-screen flex items-center justify-center bg-[#0a0a0f] text-white p-4">
        <div className="flex flex-col items-center gap-4 animate-pulse">
           <Loader2 className="w-16 h-16 text-white animate-spin" />
           <p className="text-sm font-bold">Loading authorization request...</p>
        </div>
      </div>
    }>
      <TikTokCallbackContent />
    </Suspense>
  );
}
