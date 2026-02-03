"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";

export function SiteFooter() {
  const pathname = usePathname();
  const isAuthPage = pathname.startsWith("/auth/");

  if (isAuthPage) return null;

  return (
    <footer className="border-t py-6 bg-background">
      <div className="container flex flex-col sm:flex-row items-center justify-between gap-4 text-sm text-muted-foreground">
        <div className="flex items-center gap-2">
          <span className="font-semibold text-foreground">docker-heatmap</span>
          <span>© {new Date().getFullYear()}</span>
        </div>
        <nav className="flex items-center gap-6">
          <Link
            href="https://github.com/sagargujarathi/docker-heatmap"
            target="_blank"
            rel="noopener noreferrer"
            className="hover:text-foreground transition-colors"
          >
            GitHub
          </Link>
          <Link
            href="/contributors"
            className="hover:text-foreground transition-colors"
          >
            Contributors
          </Link>
          <Link
            href="/privacy"
            className="hover:text-foreground transition-colors"
          >
            Privacy
          </Link>
          <Link
            href="/terms"
            className="hover:text-foreground transition-colors"
          >
            Terms
          </Link>
        </nav>
      </div>
    </footer>
  );
}
