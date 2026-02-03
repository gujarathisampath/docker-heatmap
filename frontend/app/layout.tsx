import type { Metadata } from "next";
import { Space_Grotesk } from "next/font/google";
import { ThemeProvider } from "next-themes";
import "./globals.css";
import { AuthProvider } from "@/context/auth-context";
import { QueryProvider } from "@/context/query-provider";
import { Toaster } from "@/components/ui/toaster";
import { SiteHeader } from "@/components/site-header";
import { SiteFooter } from "@/components/site-footer";

const spaceGrotesk = Space_Grotesk({
  subsets: ["latin"],
  variable: "--font-sans",
});

export const metadata: Metadata = {
  metadataBase: new URL("https://dockerheatmap.dev"),
  title: {
    default: "Docker Heatmap | Visualize Your Container Activity",
    template: "%s | Docker Heatmap",
  },
  description:
    "Generate stunning, embeddable SVG contribution heatmaps for your Docker Hub activity. Perfect for your GitHub README and developer profiles.",
  keywords: [
    "Docker Hub",
    "Heatmap",
    "Contribution Graph",
    "GitHub README",
    "SVG",
    "Developer Tools",
  ],
  authors: [{ name: "Sagar Gujarathi" }],
  openGraph: {
    type: "website",
    locale: "en_US",
    url: "https://dockerheatmap.dev",
    siteName: "Docker Heatmap",
    title: "Docker Heatmap",
    description: "Visualize your Docker Hub activity like GitHub commits.",
    images: [
      { url: "/logo.webp", width: 800, height: 800, alt: "Docker Heatmap" },
    ],
  },
  twitter: {
    card: "summary",
    title: "Docker Heatmap",
    description: "Visualize your Docker Hub activity like GitHub commits.",
    images: ["/logo.webp"],
  },
  icons: {
    icon: "/logo.webp",
  },
};

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html lang="en" suppressHydrationWarning>
      <body className={`${spaceGrotesk.className} min-h-screen flex flex-col`}>
        <ThemeProvider
          attribute="class"
          defaultTheme="dark"
          enableSystem
          disableTransitionOnChange
        >
          <QueryProvider>
            <AuthProvider>
              <SiteHeader />
              <div className="flex-1">{children}</div>
              <SiteFooter />
              <Toaster />
            </AuthProvider>
          </QueryProvider>
        </ThemeProvider>
      </body>
    </html>
  );
}
