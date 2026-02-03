"use client";

import { Check, Loader2, RefreshCw, Clock } from "lucide-react";
import { useForm } from "react-hook-form";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Card, CardContent } from "@/components/ui/card";
import { ConnectDockerRequest } from "@/lib/schemas";

interface DockerConnectionCardProps {
  isLoading: boolean;
  hasAccount: boolean;
  account?: {
    docker_username: string;
    last_sync_at: string | null;
  } | null;
  form: ReturnType<typeof useForm<ConnectDockerRequest>>;
  onSubmit: (data: ConnectDockerRequest) => void;
  onSync: () => void;
  onDisconnect: () => void;
  isConnecting: boolean;
  isSyncing: boolean;
  isDisconnecting: boolean;
}

export function DockerConnectionCard({
  isLoading,
  hasAccount,
  account,
  form,
  onSubmit,
  onSync,
  onDisconnect,
  isConnecting,
  isSyncing,
  isDisconnecting,
}: DockerConnectionCardProps) {
  const formatSyncTime = (dateStr: string) => {
    const date = new Date(dateStr);
    return date.toLocaleString(undefined, {
      month: "short",
      day: "numeric",
      hour: "2-digit",
      minute: "2-digit",
    });
  };

  return (
    <Card>
      <CardContent className="pt-6">
        {isLoading ? (
          <div className="py-8 flex justify-center">
            <Loader2 className="h-5 w-5 animate-spin text-muted-foreground" />
          </div>
        ) : hasAccount && account ? (
          <div className="flex items-center justify-between">
            <div className="flex items-center gap-4">
              <div className="flex items-center justify-center w-10 h-10 rounded-full bg-green-500/10">
                <Check className="h-5 w-5 text-green-500" />
              </div>
              <div>
                <p className="font-medium">{account.docker_username}</p>
                <div className="flex items-center gap-1.5 text-sm text-muted-foreground mt-0.5">
                  <Clock className="h-3.5 w-3.5" />
                  <span>
                    {account.last_sync_at
                      ? `Last synced ${formatSyncTime(account.last_sync_at)}`
                      : "Syncing..."}
                  </span>
                </div>
              </div>
            </div>
            <div className="flex flex-col items-end gap-1">
              <div className="flex items-center gap-2">
                <Button
                  variant="outline"
                  size="sm"
                  onClick={onSync}
                  disabled={isSyncing}
                >
                  <RefreshCw
                    className={`h-4 w-4 mr-2 ${isSyncing ? "animate-spin" : ""}`}
                  />
                  Sync Now
                </Button>
                <Button
                  variant="ghost"
                  size="sm"
                  onClick={onDisconnect}
                  disabled={isDisconnecting}
                  className="text-destructive hover:text-destructive hover:bg-destructive/10"
                >
                  Disconnect & Delete Data
                </Button>
              </div>
              <p className="text-[10px] text-muted-foreground mr-2">
                Permanently removes all your activity data from our servers.
              </p>
            </div>
          </div>
        ) : (
          <div>
            <div className="mb-6 rounded-lg bg-primary/5 border border-primary/20 p-4">
              <p className="text-sm text-primary font-medium mb-1">
                Security Note
              </p>
              <p className="text-sm text-muted-foreground">
                We only need to read your activity. Please create a{" "}
                <strong>Read-only</strong> Personal Access Token (PAT) for the
                best security.
              </p>
            </div>

            <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-4">
              <div className="grid gap-4 sm:grid-cols-2">
                <div className="space-y-2">
                  <Label htmlFor="docker_username" className="text-sm">
                    Docker Hub Username
                  </Label>
                  <Input
                    id="docker_username"
                    placeholder="username"
                    className="h-9"
                    {...form.register("docker_username")}
                  />
                </div>
                <div className="space-y-2">
                  <Label htmlFor="access_token" className="text-sm">
                    Read-only Access Token (PAT)
                  </Label>
                  <Input
                    id="access_token"
                    type="password"
                    placeholder="dckr_pat_..."
                    className="h-9"
                    {...form.register("access_token")}
                  />
                </div>
              </div>
              <div className="flex items-center justify-between">
                <p className="text-xs text-muted-foreground">
                  Create a token at{" "}
                  <a
                    href="https://hub.docker.com/settings/security"
                    target="_blank"
                    rel="noopener noreferrer"
                    className="underline underline-offset-2 hover:text-foreground"
                  >
                    hub.docker.com/settings/security
                  </a>
                </p>
                <Button type="submit" size="sm" disabled={isConnecting}>
                  {isConnecting && (
                    <Loader2 className="h-4 w-4 mr-2 animate-spin" />
                  )}
                  Connect
                </Button>
              </div>
            </form>
          </div>
        )}
      </CardContent>
    </Card>
  );
}
