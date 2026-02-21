'use client';

import Link from 'next/link';

export default function TermsOfServicePage() {
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
            Terms of Service
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
              Acceptance of Terms
            </h2>
            <div className="prose prose-invert max-w-none text-[var(--foreground-secondary)] leading-relaxed space-y-4">
              <p>
                By accessing or using Sidji-Omnichannel (&quot;the Service&quot;), you agree to be bound by these Terms of Service. If you do not agree to these terms, please do not use the Service.
              </p>
              <p>
                Sidji-Omnichannel is a unified customer communication platform that integrates messaging channels including WhatsApp, Instagram, Facebook Messenger, and TikTok into a single dashboard for business use.
              </p>
            </div>
          </section>

          <section>
            <h2 className="text-2xl font-bold mb-4 flex items-center gap-3">
              <span className="w-8 h-8 rounded-lg bg-[var(--primary)]/10 text-[var(--primary)] flex items-center justify-center text-sm">2</span>
              Description of Service
            </h2>
            <div className="prose prose-invert max-w-none text-[var(--foreground-secondary)] leading-relaxed space-y-4">
              <p>Sidji-Omnichannel provides the following services:</p>
              <ul className="list-disc pl-6 space-y-2">
                <li><strong>Unified Inbox:</strong> Centralized management of customer messages from WhatsApp, Instagram, Facebook, and TikTok.</li>
                <li><strong>Channel Integration:</strong> OAuth-based connection to third-party messaging platforms.</li>
                <li><strong>AI-Powered Features:</strong> Automated reply suggestions and smart responses powered by AI.</li>
                <li><strong>Team Collaboration:</strong> Multi-agent support with conversation assignment and internal notes.</li>
                <li><strong>Contact Management:</strong> Unified customer profiles across all connected channels.</li>
              </ul>
            </div>
          </section>

          <section>
            <h2 className="text-2xl font-bold mb-4 flex items-center gap-3">
              <span className="w-8 h-8 rounded-lg bg-[var(--primary)]/10 text-[var(--primary)] flex items-center justify-center text-sm">3</span>
              User Responsibilities
            </h2>
            <div className="prose prose-invert max-w-none text-[var(--foreground-secondary)] leading-relaxed space-y-4">
              <p>By using the Service, you agree to:</p>
              <ul className="list-disc pl-6 space-y-2">
                <li>Provide accurate and complete registration information.</li>
                <li>Maintain the security of your account credentials.</li>
                <li>Comply with all applicable laws and regulations, including data protection laws.</li>
                <li>Not use the Service for spam, harassment, or any unlawful purpose.</li>
                <li>Respect the terms of service of connected third-party platforms (Meta, TikTok).</li>
                <li>Obtain necessary consent from your customers before processing their data through the platform.</li>
              </ul>
            </div>
          </section>

          <section>
            <h2 className="text-2xl font-bold mb-4 flex items-center gap-3">
              <span className="w-8 h-8 rounded-lg bg-[var(--primary)]/10 text-[var(--primary)] flex items-center justify-center text-sm">4</span>
              Third-Party Integrations
            </h2>
            <div className="prose prose-invert max-w-none text-[var(--foreground-secondary)] leading-relaxed space-y-4">
              <p>
                The Service integrates with third-party platforms including Meta Platforms, Inc. (WhatsApp, Instagram, Facebook) and ByteDance Ltd. (TikTok). Your use of these integrations is also subject to the respective terms and policies of each platform.
              </p>
              <ul className="list-disc pl-6 space-y-2">
                <li>We access third-party APIs on your behalf using OAuth authorization.</li>
                <li>Access tokens are stored securely and used solely to provide the Service.</li>
                <li>We do not control the availability or functionality of third-party platforms.</li>
                <li>Changes to third-party APIs may affect the availability of certain features.</li>
              </ul>
            </div>
          </section>

          <section>
            <h2 className="text-2xl font-bold mb-4 flex items-center gap-3">
              <span className="w-8 h-8 rounded-lg bg-[var(--primary)]/10 text-[var(--primary)] flex items-center justify-center text-sm">5</span>
              Data Handling
            </h2>
            <div className="prose prose-invert max-w-none text-[var(--foreground-secondary)] leading-relaxed space-y-4">
              <p>
                We take data security seriously. Please refer to our <Link href="/privacy" className="text-[var(--primary)] hover:underline">Privacy Policy</Link> for detailed information on how we collect, use, and protect your data.
              </p>
              <ul className="list-disc pl-6 space-y-2">
                <li>Message data is stored securely and only accessible to the organization that owns the connected channel.</li>
                <li>We do not sell, share, or use customer message data for advertising or unrelated purposes.</li>
                <li>AI features process message content for generating responses but do not use your data to train global models.</li>
                <li>You may request deletion of your data at any time by contacting our support team.</li>
              </ul>
            </div>
          </section>

          <section>
            <h2 className="text-2xl font-bold mb-4 flex items-center gap-3">
              <span className="w-8 h-8 rounded-lg bg-[var(--primary)]/10 text-[var(--primary)] flex items-center justify-center text-sm">6</span>
              Limitation of Liability
            </h2>
            <div className="prose prose-invert max-w-none text-[var(--foreground-secondary)] leading-relaxed space-y-4">
              <p>
                The Service is provided &quot;as is&quot; without warranty of any kind. We shall not be liable for any indirect, incidental, special, consequential, or punitive damages resulting from your use of the Service.
              </p>
              <p>
                We do not guarantee uninterrupted availability of the Service or third-party integrations. Service may be temporarily unavailable due to maintenance, updates, or factors beyond our control.
              </p>
            </div>
          </section>

          <section>
            <h2 className="text-2xl font-bold mb-4 flex items-center gap-3">
              <span className="w-8 h-8 rounded-lg bg-[var(--primary)]/10 text-[var(--primary)] flex items-center justify-center text-sm">7</span>
              Termination
            </h2>
            <div className="prose prose-invert max-w-none text-[var(--foreground-secondary)] leading-relaxed space-y-4">
              <p>
                We reserve the right to suspend or terminate your access to the Service at any time if you violate these Terms of Service. You may also terminate your account at any time by contacting our support team.
              </p>
              <p>
                Upon termination, your data will be retained for a reasonable period to comply with legal obligations, after which it will be securely deleted.
              </p>
            </div>
          </section>

          <section>
            <h2 className="text-2xl font-bold mb-4 flex items-center gap-3">
              <span className="w-8 h-8 rounded-lg bg-[var(--primary)]/10 text-[var(--primary)] flex items-center justify-center text-sm">8</span>
              Changes to Terms
            </h2>
            <div className="prose prose-invert max-w-none text-[var(--foreground-secondary)] leading-relaxed space-y-4">
              <p>
                We may update these Terms of Service from time to time. We will notify you of any changes by updating the &quot;Last updated&quot; date at the top of this page. Continued use of the Service after any changes constitutes acceptance of the new terms.
              </p>
            </div>
          </section>

          <section>
            <h2 className="text-2xl font-bold mb-4 flex items-center gap-3">
              <span className="w-8 h-8 rounded-lg bg-[var(--primary)]/10 text-[var(--primary)] flex items-center justify-center text-sm">9</span>
              Contact Us
            </h2>
            <div className="prose prose-invert max-w-none text-[var(--foreground-secondary)] leading-relaxed space-y-4">
              <p>
                If you have any questions about these Terms of Service, please contact us at:
              </p>
              <div className="mt-4 p-6 rounded-2xl bg-[var(--primary)]/5 border border-[var(--primary)]/20 inline-block">
                <p className="font-bold text-[var(--primary)]">Sidji-Omnichannel Support</p>
                <p>Email: support@sidji-omnichannel.com</p>
                <p>Website: www.sidji-omnichannel.com</p>
              </div>
            </div>
          </section>
        </div>

        {/* Footer Section */}
        <div className="mt-12 text-center text-[var(--foreground-muted)] text-sm pb-12">
          <p>&copy; {new Date().getFullYear()} Sidji-Omnichannel. All rights reserved.</p>
          <div className="mt-4 flex justify-center gap-6">
            <Link href="/" className="hover:text-[var(--foreground)] transition-colors">Home</Link>
            <Link href="/privacy" className="hover:text-[var(--foreground)] transition-colors">Privacy Policy</Link>
            <Link href="/login" className="hover:text-[var(--foreground)] transition-colors">Login</Link>
            <Link href="/register" className="hover:text-[var(--foreground)] transition-colors">Get Started</Link>
          </div>
        </div>
      </div>
    </div>
  );
}
