'use client';

import React from 'react';
import DashboardLayout from '@/components/layout/DashboardLayout';
import { ContactsView } from '@/components/contacts/ContactsView';

export default function ContactsPage() {
  return (
    <DashboardLayout>
      <ContactsView />
    </DashboardLayout>
  );
}
