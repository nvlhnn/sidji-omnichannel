'use client';

import Link from 'next/link';

export default function PrivacyPolicyPage() {
  const lastUpdated = "February 21, 2026";

  return (
    <div className="min-h-screen bg-[var(--background)] text-[var(--foreground)] py-12 px-4 sm:px-6 lg:px-8">
      <div className="max-w-4xl mx-auto">
        {/* Header Section */}
        <div className="text-center mb-16 animate-fadeIn">
          <Link href="/" className="inline-block mb-8">
            <div className="flex items-center gap-3">
              <div className="w-10 h-10 bg-gradient-to-br from-[var(--primary)] to-purple-600 rounded-xl flex items-center justify-center shadow-lg shadow-[var(--primary)]/20">
                <svg className="w-6 h-6 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z" />
                </svg>
              </div>
              <span className="text-2xl font-bold tracking-tighter">Sidji-Omnichannel</span>
            </div>
          </Link>
          <h1 className="text-4xl font-extrabold tracking-tight sm:text-5xl mb-4">
            Privacy Policy
          </h1>
          <p className="text-[var(--foreground-muted)] text-lg">
            Last updated: {lastUpdated}
          </p>
        </div>

        {/* Content Section */}
        <div className="glass rounded-3xl p-8 md:p-12 shadow-2xl space-y-12 animate-slideIn">
          <section>
            <h2 className="text-2xl font-bold mb-4 flex items-center gap-3">
              <span className="w-8 h-8 rounded-lg bg-[var(--primary)]/10 text-[var(--primary)] flex items-center justify-center text-sm">1</span>
              Introduction
            </h2>
            <div className="prose prose-invert max-w-none text-[var(--foreground-secondary)] leading-relaxed space-y-4">
              <p>
                Welcome to Sidji-Omnichannel. We respect your privacy and are committed to protecting your personal data. This privacy policy will inform you as to how we look after your personal data when you visit our website or use our omnichannel messaging platform.
              </p>
              <p>
              Sidji-Omnichannel provides a unified inbox for business communication, integrating services like WhatsApp, Instagram, Facebook, and TikTok. In providing these services, we process data on behalf of our business customers.
              </p>
            </div>
          </section>

          <section>
            <h2 className="text-2xl font-bold mb-4 flex items-center gap-3">
              <span className="w-8 h-8 rounded-lg bg-[var(--primary)]/10 text-[var(--primary)] flex items-center justify-center text-sm">2</span>
              The Data We Collect
            </h2>
            <div className="prose prose-invert max-w-none text-[var(--foreground-secondary)] leading-relaxed space-y-4">
              <p>We may collect, use, store and transfer different kinds of personal data about you which we have grouped together as follows:</p>
              <ul className="list-disc pl-6 space-y-2">
                <li><strong>Identity Data:</strong> Includes first name, last name, username or similar identifier.</li>
                <li><strong>Contact Data:</strong> Includes billing address, email address and telephone numbers.</li>
                <li><strong>Technical Data:</strong> Includes internet protocol (IP) address, login data, browser type and version, time zone setting and location.</li>
                <li><strong>Usage Data:</strong> Includes information about how you use our website, products and services.</li>
                <li><strong>Communication Data:</strong> Includes the content of messages sent through the platform between business agents and their customers via WhatsApp, Instagram, Facebook, and TikTok.</li>
                <li><strong>Third-Party Profile Data:</strong> Includes display names, usernames, and profile pictures from connected social media accounts (Meta, TikTok) used solely for contact identification.</li>
              </ul>
            </div>
          </section>

          <section>
            <h2 className="text-2xl font-bold mb-4 flex items-center gap-3">
              <span className="w-8 h-8 rounded-lg bg-[var(--primary)]/10 text-[var(--primary)] flex items-center justify-center text-sm">3</span>
              Third-Party Integrations
            </h2>
            <div className="prose prose-invert max-w-none text-[var(--foreground-secondary)] leading-relaxed space-y-4">
              <p>
                Our platform integrates with Meta Platforms, Inc. (WhatsApp, Instagram, Facebook) and ByteDance Ltd. (TikTok). By using these channels, you are also subject to their respective privacy policies:
              </p>
              <div className="grid grid-cols-1 md:grid-cols-2 gap-4 mt-6">
                <a href="https://www.whatsapp.com/legal/privacy-policy" target="_blank" className="p-4 rounded-2xl bg-white/5 border border-white/10 hover:bg-white/10 transition-colors flex items-center gap-3">
                  <div className="w-10 h-10 rounded-full bg-[#25D366]/20 flex items-center justify-center">
                    <svg className="w-6 h-6 text-[#25D366]" fill="currentColor" viewBox="0 0 24 24"><path d="M17.472 14.382c-.297-.149-1.758-.867-2.03-.967-.273-.099-.471-.148-.67.15-.197.297-.767.966-.94 1.164-.173.199-.347.223-.644.075-.297-.15-1.255-.463-2.39-1.475-.883-.788-1.48-1.761-1.653-2.059-.173-.297-.018-.458.13-.606.134-.133.298-.347.446-.52.149-.174.198-.298.298-.497.099-.198.05-.371-.025-.52-.075-.149-.669-1.612-.916-2.207-.242-.579-.487-.5-.669-.51-.173-.008-.371-.01-.57-.01-.198 0-.52.074-.792.372-.272.297-1.04 1.016-1.04 2.479 0 1.462 1.067 2.877 1.215 3.076.149.198 2.096 3.2 5.077 4.487.709.306 1.262.489 1.694.625.712.227 1.36.195 1.871.118.571-.085 1.758-.719 2.006-1.413.248-.694.248-1.289.173-1.413-.074-.124-.272-.198-.57-.347m-5.421 7.403h-.004a9.87 9.87 0 01-5.031-1.378l-.361-.214-3.741.982.998-3.648-.235-.374a9.86 9.86 0 01-1.51-5.26c.001-5.45 4.436-9.884 9.888-9.884 2.64 0 5.122 1.03 6.988 2.898a9.825 9.825 0 012.893 6.994c-.003 5.45-4.437 9.884-9.885 9.884m8.413-18.297A11.815 11.815 0 0012.05 0C5.414 0 0 5.414 0 12.05c0 2.123.552 4.197 1.6 6.037L0 24l6.105-1.602a11.834 11.834 0 005.937 1.598h.005c6.637 0 12.05-5.414 12.05-12.05a11.834 11.834 0 00-3.592-8.548z"/></svg>
                  </div>
                  <div>
                    <p className="text-sm font-bold">WhatsApp Privacy</p>
                    <p className="text-xs text-[var(--foreground-muted)]">View Policy</p>
                  </div>
                </a>
                <a href="https://privacycenter.instagram.com/policy" target="_blank" className="p-4 rounded-2xl bg-white/5 border border-white/10 hover:bg-white/10 transition-colors flex items-center gap-3">
                  <div className="w-10 h-10 rounded-full bg-[#E1306C]/20 flex items-center justify-center">
                    <svg className="w-6 h-6 text-[#E1306C]" fill="currentColor" viewBox="0 0 24 24"><path d="M12 0C8.74 0 8.333.015 7.053.072 5.775.132 4.905.333 4.14.63c-.789.306-1.459.717-2.126 1.384S.935 3.35.63 4.14C.333 4.905.131 5.775.072 7.053.012 8.333 0 8.74 0 12s.015 3.667.072 4.947c.06 1.277.261 2.148.558 2.913.306.788.717 1.459 1.384 2.126.667.666 1.336 1.079 2.126 1.384.766.296 1.636.499 2.913.558C8.333 23.988 8.74 24 12 24s3.667-.015 4.947-.072c1.277-.06 2.148-.262 2.913-.558.788-.306 1.459-.718 2.126-1.384.666-.667 1.079-1.335 1.384-2.126.296-.765.499-1.636.558-2.913.06-1.28.072-1.687.072-4.947s-.015-3.667-.072-4.947c-.06-1.277-.262-2.149-.558-2.913-.306-.789-.718-1.459-1.384-2.126C21.066.935 20.397.522 19.608.217c-.765-.297-1.636-.499-2.913-.558C15.667.012 15.26 0 12 0zm0 2.16c3.203 0 3.58.016 4.85.071 1.17.055 1.805.249 2.227.415.562.217.96.477 1.382.896.419.42.679.819.896 1.381.164.422.36 1.057.413 2.227.057 1.27.07 1.646.07 4.85s-.015 3.58-.07 4.85c-.055 1.17-.249 1.805-.413 2.227-.218.562-.477.96-.896 1.382-.419.419-.818.679-1.381.896-.422.164-1.056.36-2.227.413-1.27.057-1.647.07-4.85.07s-3.58-.015-4.85-.07c-1.17-.055-1.805-.249-2.227-.413-.562-.217-.96-.477-1.382-.896-.419-.419-.679-.818-.896-1.381-.164-.422-.36-1.056-.413-2.227-.054-1.27-.07-1.647-.07-4.85s.016-3.58.07-4.85c.055-1.17.249-1.805.413-2.227.217-.562.477-.96.896-1.382.421-.419.82-.679 1.382-.896.422-.164 1.057-.36 2.227-.413 1.27-.054 1.647-.07 4.85-.07zM12 5.837a6.162 6.162 0 100 12.324 6.162 6.162 0 000-12.324zM12 16a4 4 0 110-8 4 4 0 010 8zm6.406-9.145a1.44 1.44 0 11-2.88 0 1.44 1.44 0 012.88 0z"/></svg>
                  </div>
                  <div>
                    <p className="text-sm font-bold">Instagram Privacy</p>
                    <p className="text-xs text-[var(--foreground-muted)]">View Policy</p>
                  </div>
                </a>
                <a href="https://www.tiktok.com/legal/privacy-policy" target="_blank" className="p-4 rounded-2xl bg-white/5 border border-white/10 hover:bg-white/10 transition-colors flex items-center gap-3">
                  <div className="w-10 h-10 rounded-full bg-[#00f2ea]/20 flex items-center justify-center">
                    <svg className="w-6 h-6 text-[#00f2ea]" fill="currentColor" viewBox="0 0 24 24"><path d="M19.59 6.69a4.83 4.83 0 01-3.77-4.25V2h-3.45v13.67a2.89 2.89 0 01-2.88 2.5 2.89 2.89 0 01-2.89-2.89 2.89 2.89 0 012.89-2.89c.28 0 .54.04.79.1v-3.5a6.37 6.37 0 00-.79-.05A6.34 6.34 0 003.15 15.2a6.34 6.34 0 006.34 6.34 6.34 6.34 0 006.34-6.34V9.28a8.28 8.28 0 004.77 1.51V7.34a4.85 4.85 0 01-1.01-.65z"/></svg>
                  </div>
                  <div>
                    <p className="text-sm font-bold">TikTok Privacy</p>
                    <p className="text-xs text-[var(--foreground-muted)]">View Policy</p>
                  </div>
                </a>
              </div>

              <div className="mt-6 p-4 rounded-2xl bg-white/5 border border-white/10">
                <h3 className="font-bold mb-2">TikTok Data Usage</h3>
                <p className="text-sm text-[var(--foreground-muted)] leading-relaxed">When you connect your TikTok account, we access the following data through TikTok&apos;s API:</p>
                <ul className="list-disc pl-6 space-y-1 mt-2 text-sm text-[var(--foreground-muted)]">
                  <li><strong>User Profile:</strong> Display name, username, and avatar URL — used to identify the connected channel.</li>
                  <li><strong>Direct Messages:</strong> Message content received via TikTok DMs — displayed in your unified inbox for customer support purposes only.</li>
                  <li><strong>Open ID:</strong> A unique user identifier — used to route messages to the correct business account.</li>
                </ul>
                <p className="text-sm text-[var(--foreground-muted)] mt-2">We do not store, share, or sell TikTok user data beyond what is necessary for providing customer support functionality. Access tokens are stored securely and are never exposed publicly.</p>
              </div>
            </div>
          </section>

          <section>
            <h2 className="text-2xl font-bold mb-4 flex items-center gap-3">
              <span className="w-8 h-8 rounded-lg bg-[var(--primary)]/10 text-[var(--primary)] flex items-center justify-center text-sm">4</span>
              AI Processing
            </h2>
            <div className="prose prose-invert max-w-none text-[var(--foreground-secondary)] leading-relaxed space-y-4">
              <p>
                Our platform utilizes Artificial Intelligence (AI) to provide features like draft suggestions and automated replies.
              </p>
              <ul className="list-disc pl-6 space-y-2">
                <li>We use Google Gemini and OpenAI to process message content for generative features.</li>
                <li>Your data is used to provide context for the AI but is not used to train global AI models in a way that would expose your private information.</li>
                <li>Vector embeddings are used to search your own knowledge base for relevant context.</li>
              </ul>
            </div>
          </section>

          <section>
            <h2 className="text-2xl font-bold mb-4 flex items-center gap-3">
              <span className="w-8 h-8 rounded-lg bg-[var(--primary)]/10 text-[var(--primary)] flex items-center justify-center text-sm">5</span>
              Data Security
            </h2>
            <div className="prose prose-invert max-w-none text-[var(--foreground-secondary)] leading-relaxed space-y-4">
              <p>
                We have put in place appropriate security measures to prevent your personal data from being accidentally lost, used or accessed in an unauthorized way. We limit access to your personal data to those employees, agents, contractors and other third parties who have a business need to know.
              </p>
            </div>
          </section>

          <section>
            <h2 className="text-2xl font-bold mb-4 flex items-center gap-3">
              <span className="w-8 h-8 rounded-lg bg-[var(--primary)]/10 text-[var(--primary)] flex items-center justify-center text-sm">6</span>
              Contact Us
            </h2>
            <div className="prose prose-invert max-w-none text-[var(--foreground-secondary)] leading-relaxed space-y-4">
              <p>
                If you have any questions about this privacy policy or our privacy practices, please contact us at:
              </p>
              <div className="mt-4 p-6 rounded-2xl bg-[var(--primary)]/5 border border-[var(--primary)]/20 inline-block">
                <p className="font-bold text-[var(--primary)]">Privacy Team</p>
                <p>Email: privacy@Sidji-Omnichannel.com</p>
                <p>Website: www.Sidji-Omnichannel.com</p>
              </div>
            </div>
          </section>
        </div>

        {/* Footer Section */}
        <div className="mt-12 text-center text-[var(--foreground-muted)] text-sm pb-12">
          <p>&copy; {new Date().getFullYear()} Sidji-Omnichannel. All rights reserved.</p>
          <div className="mt-4 flex justify-center gap-6">
            <Link href="/" className="hover:text-[var(--foreground)] transition-colors">Home</Link>
            <Link href="/terms" className="hover:text-[var(--foreground)] transition-colors">Terms of Service</Link>
            <Link href="/login" className="hover:text-[var(--foreground)] transition-colors">Login</Link>
            <Link href="/register" className="hover:text-[var(--foreground)] transition-colors">Get Started</Link>
          </div>
        </div>
      </div>
    </div>
  );
}
