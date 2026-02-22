'use client';

import { Suspense, useEffect, useState } from 'react';
import { useRouter, useSearchParams } from 'next/navigation';
import { api } from '@/lib/api';
import { Loader2, AlertCircle, CheckCircle2 } from 'lucide-react';

function GoogleCallbackContent() {
  const router = useRouter();
  const searchParams = useSearchParams();
  const [status, setStatus] = useState<'loading' | 'success' | 'error'>('loading');
  const [errorMessage, setErrorMessage] = useState('');

  useEffect(() => {
    const code = searchParams.get('code');
    const error = searchParams.get('error');

    if (error) {
       setStatus('error');
       setErrorMessage(`Authentication rejected or failed: ${error}`);
       return;
    }

    if (!code) {
      setStatus('error');
      setErrorMessage('No authorization code provided by Google.');
      return;
    }

    const completeLogin = async () => {
      try {
        await api.loginWithGoogle(code);
        setStatus('success');
        // Redirect back to inbox securely
        setTimeout(() => {
          router.push('/inbox');
        }, 1500);
      } catch (err: any) {
        console.error('Failed to login via Google:', err);
        setStatus('error');
        setErrorMessage(err.response?.data?.error || err.message || 'Failed to authenticate.');
      }
    };

    completeLogin();
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
              <h2 className="text-xl font-bold tracking-tight">Authenticating with Google...</h2>
              <p className="text-sm text-[var(--foreground-muted)]">Please wait while we log you in securely.</p>
            </div>
          </div>
        )}

        {status === 'success' && (
          <div className="flex flex-col items-center gap-4 animate-scaleIn">
            <div className="w-16 h-16 bg-white text-black rounded-full flex items-center justify-center shadow-[0_0_30px_rgba(255,255,255,0.3)]">
              <CheckCircle2 size={32} />
            </div>
            <div className="space-y-2">
              <h2 className="text-xl font-bold tracking-tight text-white">Login Successful!</h2>
              <p className="text-sm text-[var(--foreground-muted)]">Redirecting you to your inbox...</p>
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
                <h2 className="text-xl font-bold tracking-tight text-white">Authentication Failed</h2>
                <p className="text-sm text-[var(--foreground-muted)] mt-1">{errorMessage}</p>
              </div>
              <button
                onClick={() => router.push('/login')}
                className="w-full py-3 px-4 bg-white/10 hover:bg-white/20 text-white rounded-xl text-sm font-semibold transition-colors"
              >
                Back to Login
              </button>
            </div>
          </div>
        )}
      </div>
    </div>
  );
}

export default function GoogleCallbackPage() {
  return (
    <Suspense fallback={
      <div className="min-h-screen flex items-center justify-center bg-[#0a0a0f] text-white p-4">
        <div className="flex flex-col items-center gap-4 animate-pulse">
           <Loader2 className="w-16 h-16 text-white animate-spin" />
           <p className="text-sm font-bold">Verifying Google login...</p>
        </div>
      </div>
    }>
      <GoogleCallbackContent />
    </Suspense>
  );
}
