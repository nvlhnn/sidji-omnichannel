import type { Metadata } from "next";
import { Inter } from "next/font/google";
import "./globals.css";

const inter = Inter({ subsets: ["latin"] });

export const metadata: Metadata = {
  title: "Sidji - Unified Inbox",
  description: "Manage WhatsApp and Instagram conversations in one place",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en">
      <body className={inter.className}>
        {children}
        <script
          dangerouslySetInnerHTML={{
             __html: `
              window.fbAsyncInit = function() {
                FB.init({
                  appId            : '1283159606968560',
                  autoLogAppEvents : true,
                  xfbml            : true,
                  version          : 'v21.0'
                });
              };
             `
          }}
        />
        <script async defer crossOrigin="anonymous" src="https://connect.facebook.net/en_US/sdk.js"></script>
      </body>
    </html>
  );
}
