"use client";

import { Button } from "@/components/ui/button";
import { useAuth } from "@/context/auth-context";
import { Github, ArrowRight, Star } from "lucide-react";
import Link from "next/link";

interface HeroCTAProps {
  stars: number | null;
}

export function HeroCTA({ stars }: HeroCTAProps) {
  const { isAuthenticated, login, isLoading } = useAuth();

  if (isLoading) {
    return (
      <div className="flex flex-col sm:flex-row gap-3 mt-2 h-[44px]">
        <div className="w-[200px] h-full bg-muted animate-pulse rounded-md" />
        <div className="w-[140px] h-full bg-muted animate-pulse rounded-md" />
      </div>
    );
  }

  return (
    <div className="flex flex-col sm:flex-row gap-3 mt-2">
      {isAuthenticated ? (
        <Link href="/dashboard" className="w-full sm:w-auto">
          <Button size="lg" className="w-full">
            Go to Dashboard
            <ArrowRight className="ml-2 h-4 w-4" />
          </Button>
        </Link>
      ) : (
        <Button size="lg" onClick={login} className="w-full sm:w-auto">
          <Github className="mr-2 h-4 w-4" />
          Get started with GitHub
        </Button>
      )}
      <Link
        href="https://github.com/sagargujarathi/docker-heatmap"
        target="_blank"
        className="w-full sm:w-auto"
      >
        <Button size="lg" variant="outline" className="w-full gap-2">
          <Star className="h-4 w-4 fill-current" />
          Star on GitHub
          {stars !== null && (
            <span className="ml-1 px-1.5 py-0.5 rounded-full bg-muted text-xs font-semibold tabular-nums">
              {stars.toLocaleString()}
            </span>
          )}
        </Button>
      </Link>
    </div>
  );
}
