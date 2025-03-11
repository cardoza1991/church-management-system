'use client';

import ContactForm from '@/components/contacts/ContactForm';

export default function NewContactPage() {
  return (
    <div>
      <h1 className="text-3xl font-bold mb-6">Add New Contact</h1>
      <ContactForm />
    </div>
  );
}