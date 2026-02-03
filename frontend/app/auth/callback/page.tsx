"use client";

import { useEffect, useRef } from "react";
import { useSearchParams } from "next/navigation";
import Image from "next/image";
import { Loader2 } from "lucide-react";
import { useAuth } from "@/context/auth-context";

function CallbackContent() {
  const searchParams = useSearchParams();
  const { refreshUser } = useAuth();
  const processed = useRef(false);

  useEffect(() => {
    if (processed.current) return;
    processed.current = true;

    const token = searchParams.get("token");

    if (token) {
      localStorage.setItem("token", token);
      refreshUser().then(() => {
        // Use window.location.href for a full page reload to ensure all contexts are updated
        window.location.href = "/dashboard";
      });
    } else {
      window.location.href = "/auth/error?message=no_token";
    }
  }, [searchParams, refreshUser]);

  return (
    <div className="min-h-screen flex items-center justify-center bg-background">
      <div className="text-center">
        <Image
          src="/logo.webp"
          alt="Logo"
          width={80}
          height={80}
          className="mx-auto mb-6"
        />
        <Loader2 className="h-6 w-6 animate-spin mx-auto mb-4" />
        <p className="text-sm text-muted-foreground">Signing you in...</p>
      </div>
    </div>
  );
}

import { Suspense } from "react";

export default function AuthCallbackPage() {
  return (
    <Suspense
      fallback={
        <div className="min-h-screen flex items-center justify-center bg-background">
          <Loader2 className="h-6 w-6 animate-spin" />
        </div>
      }
    >
      <CallbackContent />
    </Suspense>
  );
}
