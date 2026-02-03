"use client";

import { Suspense } from "react";
import { useSearchParams } from "next/navigation";
import { Button } from "@/components/ui/button";
import { ArrowLeft } from "lucide-react";
import Link from "next/link";

const errorMessages: Record<string, string> = {
  missing_params: "Missing required parameters. Please try signing in again.",
  invalid_state: "Invalid OAuth state. The link may have expired.",
  auth_failed: "Authentication with GitHub failed. Please try again.",
  token_failed: "Failed to generate authentication token.",
  no_token: "No authentication token received.",
  default: "An unexpected error occurred during authentication.",
};

function ErrorContent() {
  const searchParams = useSearchParams();
  const errorCode = searchParams.get("message") || "default";
  const errorMessage = errorMessages[errorCode] || errorMessages.default;

  return (
    <div className="min-h-screen flex items-center justify-center bg-background p-4">
      <div className="max-w-md w-full text-center">
        <h1 className="text-xl font-bold mb-2">Authentication Error</h1>
        <p className="text-muted-foreground mb-6">{errorMessage}</p>
        <Link href="/">
          <Button>
            <ArrowLeft className="mr-2 h-4 w-4" />
            Back to Home
          </Button>
        </Link>
        <p className="text-xs text-muted-foreground mt-4">Error: {errorCode}</p>
      </div>
    </div>
  );
}

export default function AuthErrorPage() {
  return (
    <Suspense
      fallback={
        <div className="min-h-screen flex items-center justify-center bg-background">
          <p className="text-muted-foreground">Loading...</p>
        </div>
      }
    >
      <ErrorContent />
    </Suspense>
  );
}
