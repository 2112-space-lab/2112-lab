"use client"

import React from 'react';
import { ClerkProvider } from '@clerk/clerk-react'; // Correct import for clerk-react
import AppWrappers from './AppWrappers';
import { useRouter } from 'next/navigation';

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  const router = useRouter();

  return (
    <ClerkProvider
      publishableKey={process.env.NEXT_PUBLIC_CLERK_PUBLISHABLE_KEY} // Use Clerk's publishable key from .env
      routerPush={(to) => router.push(to)} // Provide routerPush implementation
      routerReplace={(to) => router.replace(to)} // Provide routerReplace implementation
    >
      <html lang="en">
        <body className="dark" id="root">
          <AppWrappers>{children}</AppWrappers>
        </body>
      </html>
    </ClerkProvider>
  );
}
