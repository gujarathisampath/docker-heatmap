"use client";

import { useQuery } from "@tanstack/react-query";
import { publicApi } from "@/lib/api";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { ArrowLeft, Loader2, Github, Share2, Copy, Check } from "lucide-react";
import Link from "next/link";
import { useState } from "react";
import { ProfileData } from "@/lib/schemas";
import { useToast } from "@/hooks/use-toast";

interface ProfileClientProps {
  username: string;
}

export function ProfileClient({ username }: ProfileClientProps) {
  const { toast } = useToast();
  const [copied, setCopied] = useState(false);
  const {
    data: profile,
    isLoading,
    error,
  } = useQuery<ProfileData>({
    queryKey: ["profile", username],
    queryFn: () => publicApi.getProfile(username),
  });

  const handleCopyLink = () => {
    navigator.clipboard.writeText(window.location.href);
    setCopied(true);
    toast({
      title: "Link copied!",
      description: "Profile link has been copied to your clipboard.",
    });
    setTimeout(() => setCopied(false), 2000);
  };

  if (isLoading) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <Loader2 className="h-5 w-5 animate-spin text-muted-foreground" />
      </div>
    );
  }

  if (error || !profile) {
    return <NotFoundState />;
  }

  return (
    <div className="min-h-screen bg-background">
      <main className="container py-10">
        <div className="flex flex-col md:flex-row md:items-end justify-between gap-6 mb-12">
          {/* Profile Header */}
          <div className="flex items-start gap-5">
            <Avatar className="h-20 w-20 border-2 border-muted shadow-sm">
              <AvatarImage src={profile.user.avatar_url} />
              <AvatarFallback className="text-2xl font-bold bg-primary/10">
                {profile.user.github_username[0].toUpperCase()}
              </AvatarFallback>
            </Avatar>
            <div className="space-y-1">
              <h1 className="text-3xl font-bold tracking-tight">
                {profile.user.name || profile.user.github_username}
              </h1>
              <div className="flex flex-wrap items-center gap-3 text-muted-foreground">
                <Link
                  href={`https://github.com/${profile.user.github_username}`}
                  target="_blank"
                  className="flex items-center gap-1.5 hover:text-primary transition-colors bg-muted/50 px-2 py-0.5 rounded-md text-sm border"
                >
                  <Github className="h-3.5 w-3.5" />
                  {profile.user.github_username}
                </Link>
                <span className="text-muted-foreground/30 hidden sm:inline">
                  |
                </span>
                <span className="flex items-center gap-1.5 bg-primary/5 text-primary px-2 py-0.5 rounded-md text-sm border border-primary/10 font-medium">
                  {profile.docker.username} (Docker Hub)
                </span>
              </div>
              {profile.user.bio && (
                <p className="text-base text-muted-foreground mt-3 max-w-xl leading-relaxed">
                  {profile.user.bio}
                </p>
              )}
            </div>
          </div>

          <div className="flex items-center gap-3">
            <Button
              variant="outline"
              size="sm"
              className="gap-2 shrink-0 shadow-sm"
              onClick={handleCopyLink}
            >
              {copied ? (
                <Check className="h-4 w-4 text-green-500" />
              ) : (
                <Copy className="h-4 w-4" />
              )}
              {copied ? "Copied" : "Copy Profile Link"}
            </Button>
            <Button
              variant="default"
              size="sm"
              className="gap-2 shrink-0 shadow-sm"
              asChild
            >
              <Link
                href={`https://hub.docker.com/u/${profile.docker.username}`}
                target="_blank"
              >
                <Share2 className="h-4 w-4" />
                View Docker Hub
              </Link>
            </Button>
          </div>
        </div>

        {/* Heatmap */}
        <Card className="overflow-hidden border-muted/50 shadow-sm">
          <CardHeader className="bg-muted/30 border-b py-4">
            <CardTitle className="text-sm font-semibold flex items-center gap-2">
              <div className="w-2 h-2 rounded-full bg-primary" />
              Docker Hub Activity (Last Year)
            </CardTitle>
          </CardHeader>
          <CardContent className="p-0">
            <div className="p-6 md:p-10 overflow-x-auto flex flex-col items-center">
              <img
                src={publicApi.getHeatmapUrl(username, { theme: "github" })}
                alt={`${username}'s Docker activity`}
                className="max-w-none rounded shadow-sm bg-background border p-4"
              />
              <p className="text-xs text-muted-foreground mt-4 italic">
                Activity data is fetched directly from Docker Hub API.
              </p>
            </div>
          </CardContent>
        </Card>

        {/* CTA */}
        <div className="mt-16 text-center border-t pt-16 pb-12">
          <h2 className="text-2xl font-bold mb-4">
            Want your own Docker activity heatmap?
          </h2>
          <p className="text-muted-foreground mb-8 max-w-md mx-auto">
            Join hundreds of developers who showcase their container activity on
            their GitHub READMEs.
          </p>
          <Link href="/">
            <Button size="lg" className="rounded-full px-8 shadow-lg">
              Get started for free
            </Button>
          </Link>
        </div>
      </main>
    </div>
  );
}

function NotFoundState() {
  return (
    <div className="min-h-screen flex flex-col items-center justify-center gap-6 p-4 text-center">
      <div className="w-20 h-20 bg-muted rounded-full flex items-center justify-center mb-2">
        <ArrowLeft className="h-10 w-10 text-muted-foreground/50" />
      </div>
      <div className="space-y-2">
        <h1 className="text-3xl font-bold tracking-tight">Profile not found</h1>
        <p className="text-muted-foreground max-w-xs mx-auto text-lg leading-relaxed">
          The user you&apos;re looking for doesn&apos;t exist, or they
          haven&apos;t connected their Docker Hub account yet.
        </p>
      </div>
      <Link href="/">
        <Button variant="outline" className="rounded-full px-6">
          <ArrowLeft className="mr-2 h-4 w-4" />
          Back to homepage
        </Button>
      </Link>
    </div>
  );
}
