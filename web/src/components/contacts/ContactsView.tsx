'use client';

import { useState, useEffect } from 'react';
import { api } from '@/lib/api';
import { Contact } from '@/lib/types';
import { Avatar } from '../ui/Avatar';

export function ContactsView() {
  const [contacts, setContacts] = useState<Contact[]>([]);
  const [isLoading, setIsLoading] = useState(false);
  const [search, setSearch] = useState('');
  const [page, setPage] = useState(1);
  const [hasMore, setHasMore] = useState(true);
  const [total, setTotal] = useState(0);

  useEffect(() => {
    loadContacts();
  }, [page, search]);

  const loadContacts = async () => {
    try {
      setIsLoading(true);
      const response = await api.getContacts(page, 20, search);
      const data = response.data || [];
      if (page === 1) {
        setContacts(data);
      } else {
        setContacts((prev) => [...prev, ...data]);
      }
      setTotal(response.total);
      setHasMore(contacts.length + data.length < response.total);
    } catch (error) {
      console.error('Failed to load contacts:', error);
    } finally {
      setIsLoading(false);
    }
  };

  const handleSearch = (e: React.ChangeEvent<HTMLInputElement>) => {
    setSearch(e.target.value);
    setPage(1); // Reset to page 1
  };

  return (
    <div className="flex-1 flex flex-col bg-[var(--background)] h-full">
      {/* Header */}
      <div className="p-6 border-b border-[var(--border)] bg-[var(--background-secondary)] flex justify-between items-center">
        <div>
          <h1 className="text-2xl font-bold bg-gradient-to-r from-[var(--foreground)] to-[var(--foreground-secondary)] bg-clip-text text-transparent">
            Contacts
          </h1>
          <p className="text-sm text-[var(--foreground-muted)]">{total} contacts found</p>
        </div>

        {/* Search */}
        <div className="relative w-72">
          <svg
            className="absolute left-3 top-1/2 -translate-y-1/2 w-5 h-5 text-[var(--foreground-muted)]"
            fill="none"
            stroke="currentColor"
            viewBox="0 0 24 24"
          >
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
          </svg>
          <input
            type="text"
            placeholder="Search contacts..."
            value={search}
            onChange={handleSearch}
            className="input w-full pl-10"
          />
        </div>
      </div>

      {/* List */}
      <div className="flex-1 overflow-y-auto p-6">
        {isLoading && page === 1 ? (
          <div className="space-y-4">
            {[1, 2, 3, 4, 5].map((i) => (
              <div key={i} className="h-16 bg-[var(--background-secondary)] rounded-xl animate-pulse" />
            ))}
          </div>
        ) : (
          <div className="grid gap-4">
            {contacts.map((contact) => (
              <div
                key={contact.id}
                className="p-4 bg-[var(--background-secondary)] rounded-xl border border-[var(--border)] hover:border-[var(--primary)] transition-colors flex items-center justify-between group"
              >
                <div className="flex items-center gap-4">
                  <Avatar name={contact.name} src={contact.avatar_url} size="lg" />
                  <div>
                    <h3 className="font-semibold text-[var(--foreground)]">{contact.name}</h3>
                    <div className="flex items-center gap-4 mt-1 text-sm text-[var(--foreground-muted)]">
                      {contact.phone && (
                        <div className="flex items-center gap-1">
                          <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1.5} d="M3 5a2 2 0 012-2h3.28a1 1 0 01.948.684l1.498 4.493a1 1 0 01-.502 1.21l-2.257 1.13a11.042 11.042 0 005.516 5.516l1.13-2.257a1 1 0 011.21-.502l4.493 1.498a1 1 0 01.684.949V19a2 2 0 01-2 2h-1C9.716 21 3 14.284 3 6V5z" />
                          </svg>
                          <span>{contact.phone}</span>
                        </div>
                      )}
                      {contact.email && (
                        <div className="flex items-center gap-1">
                          <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1.5} d="M3 8l7.89 5.26a2 2 0 002.22 0L21 8M5 19h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z" />
                          </svg>
                          <span>{contact.email}</span>
                        </div>
                      )}
                    </div>
                  </div>
                </div>

                <div className="flex gap-2">
                  {contact.whatsapp_id && (
                    <div className="px-2 py-1 bg-green-500/10 text-green-500 rounded text-xs font-mono">
                      WA: {contact.whatsapp_id}
                    </div>
                  )}
                  {contact.instagram_id && (
                    <div className="px-2 py-1 bg-pink-500/10 text-pink-500 rounded text-xs font-mono">
                      IG: {contact.instagram_id}
                    </div>
                  )}
                  {contact.facebook_id && (
                    <div className="px-2 py-1 bg-blue-500/10 text-blue-500 rounded text-xs font-mono">
                      FB: {contact.facebook_id}
                    </div>
                  )}
                  {contact.tiktok_id && (
                    <div className="px-2 py-1 bg-black/10 dark:bg-white/10 text-black dark:text-white rounded text-xs font-mono">
                      TT: {contact.tiktok_id}
                    </div>
                  )}
                </div>
              </div>
            ))}

            {contacts.length === 0 && !isLoading && (
              <div className="text-center py-12">
                <p className="text-[var(--foreground-muted)]">No contacts found.</p>
              </div>
            )}
          </div>
        )}
      </div>
    </div>
  );
}
