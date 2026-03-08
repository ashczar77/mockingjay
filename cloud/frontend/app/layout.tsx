import type { Metadata } from "next";
import "./globals.css";

export const metadata: Metadata = {
  title: "MockingJay - Voice AI Testing",
  description: "Monitor your voice AI test performance",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en">
      <body>{children}</body>
    </html>
  );
}
