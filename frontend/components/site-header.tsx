"use client";

import Link from "next/link";
import Image from "next/image";
import { Github, LogOut, LayoutDashboard } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { ThemeToggle } from "@/components/theme-toggle";
import { useAuth } from "@/context/auth-context";
import { usePathname } from "next/navigation";

export function SiteHeader() {
  const { user, isAuthenticated, isLoading, login, logout } = useAuth();
  const pathname = usePathname();

  // Don't show header on error or callback pages if they have their own UI
  // But usually, it's fine to keep it.
  const isAuthPage = pathname.startsWith("/auth/");

  if (isAuthPage) return null;

  return (
    <header className="sticky top-0 z-50 w-full border-b bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60">
      <div className="container flex h-14 items-center">
        <Link href="/" className="mr-6 flex items-center gap-2">
          <Image
            src="/logo.webp"
            alt="Docker Heatmap Logo"
            width={36}
            height={36}
          />
          <span className="font-bold text-lg hidden sm:inline-block">
            docker-heatmap
          </span>
        </Link>
        <nav className="flex flex-1 items-center justify-end gap-2">
          <Link
            href="https://github.com/sagargujarathi/docker-heatmap"
            target="_blank"
            rel="noopener noreferrer"
            className="text-sm text-muted-foreground hover:text-foreground px-3 py-2 hidden md:block"
          >
            GitHub
          </Link>
          <ThemeToggle />

          {isLoading ? (
            <div className="h-8 w-24 bg-muted animate-pulse rounded-md" />
          ) : isAuthenticated && user ? (
            <div className="flex items-center gap-2 ml-2">
              {pathname !== "/dashboard" && (
                <Link href="/dashboard">
                  <Button
                    variant="ghost"
                    size="sm"
                    className="hidden sm:flex items-center gap-2"
                  >
                    <LayoutDashboard className="h-4 w-4" />
                    Dashboard
                  </Button>
                </Link>
              )}

              <div className="flex items-center gap-3 px-2 border-l ml-1">
                <Avatar className="h-7 w-7">
                  <AvatarImage src={user.avatar_url} />
                  <AvatarFallback className="text-xs">
                    {user.github_username[0].toUpperCase()}
                  </AvatarFallback>
                </Avatar>
                <span className="text-sm font-medium hidden lg:block">
                  {user.github_username}
                </span>
                <Button
                  variant="ghost"
                  size="icon"
                  onClick={logout}
                  className="h-8 w-8 text-muted-foreground hover:text-destructive"
                  title="Logout"
                >
                  <LogOut className="h-4 w-4" />
                </Button>
              </div>
            </div>
          ) : (
            <Button
              size="sm"
              variant="default"
              onClick={login}
              className="ml-2"
            >
              <Github className="mr-2 h-4 w-4" />
              Sign in
            </Button>
          )}
        </nav>
      </div>
    </header>
  );
}
