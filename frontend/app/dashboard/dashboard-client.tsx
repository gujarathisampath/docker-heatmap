"use client";

import { useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import { useAuth } from "@/context/auth-context";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { useToast } from "@/hooks/use-toast";
import { dockerApi, publicApi } from "@/lib/api";
import { useForm } from "react-hook-form";
import type { ConnectDockerRequest, SVGOptions } from "@/lib/schemas";

import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Loader2, Sparkles } from "lucide-react";
import { DockerConnectionCard } from "@/components/dashboard/docker-connection-card";
import { SVGCustomizationCard } from "@/components/dashboard/svg-customization-card";
import { EmbedCodesCard } from "@/components/dashboard/embed-codes-card";

// Default themes in case API fails
const DEFAULT_THEMES = [
  { id: "github", name: "GitHub Dark" },
  { id: "github-light", name: "GitHub Light" },
  { id: "docker", name: "Docker" },
  { id: "dracula", name: "Dracula" },
  { id: "nord", name: "Nord" },
  { id: "monokai", name: "Monokai" },
  { id: "one-dark", name: "One Dark" },
  { id: "tokyo-night", name: "Tokyo Night" },
  { id: "catppuccin", name: "Catppuccin" },
  { id: "ocean", name: "Ocean" },
  { id: "sunset", name: "Sunset" },
  { id: "forest", name: "Forest" },
  { id: "purple", name: "Purple" },
  { id: "rose", name: "Rose" },
  { id: "minimal", name: "Minimal" },
  { id: "minimal-dark", name: "Minimal Dark" },
];

export function DashboardClient() {
  const router = useRouter();
  const { user, isLoading: authLoading, isAuthenticated } = useAuth();
  const { toast } = useToast();
  const queryClient = useQueryClient();
  const [copied, setCopied] = useState<string | null>(null);

  // SVG Customization state
  const [svgOptions, setSvgOptions] = useState<SVGOptions>({
    theme: "github",
    days: 365,
    cell_size: 11,
    radius: 2,
    hide_legend: false,
    hide_total: false,
    hide_labels: false,
  });

  // Redirect if not authenticated
  useEffect(() => {
    const token =
      typeof window !== "undefined" ? localStorage.getItem("token") : null;
    if (!authLoading && !isAuthenticated && !token) {
      router.replace("/");
    }
  }, [authLoading, isAuthenticated, router]);

  // Fetch Docker account (404 means no account, not an error)
  const { data: dockerData, isLoading: dockerLoading } = useQuery({
    queryKey: ["docker-account"],
    queryFn: async () => {
      try {
        return await dockerApi.getAccount();
      } catch {
        return { account: null };
      }
    },
    enabled: isAuthenticated,
  });

  // Fetch available themes
  const { data: themesData } = useQuery({
    queryKey: ["themes"],
    queryFn: publicApi.getThemes,
    staleTime: 1000 * 60 * 60, // Cache for 1 hour
  });

  // Connect Docker mutation
  const connectMutation = useMutation({
    mutationFn: dockerApi.connect,
    onSuccess: () => {
      toast({
        title: "Connected!",
        description:
          "Docker Hub account linked successfully. Syncing activity...",
      });
      queryClient.invalidateQueries({ queryKey: ["docker-account"] });
    },
    onError: (error: Error) => {
      toast({
        title: "Connection Failed",
        description: error.message,
        variant: "destructive",
      });
    },
  });

  // Disconnect mutation
  const disconnectMutation = useMutation({
    mutationFn: dockerApi.disconnect,
    onSuccess: () => {
      toast({
        title: "Disconnected",
        description: "Docker Hub account removed.",
      });
      queryClient.invalidateQueries({ queryKey: ["docker-account"] });
    },
  });

  // Sync mutation
  const syncMutation = useMutation({
    mutationFn: dockerApi.sync,
    onSuccess: () => {
      toast({
        title: "Sync Started",
        description: "Your activity is being updated...",
      });
    },
  });

  // Form
  const form = useForm<ConnectDockerRequest>({
    defaultValues: { docker_username: "", access_token: "" },
  });

  const onSubmit = (data: ConnectDockerRequest) => {
    connectMutation.mutate(data);
  };

  const copyToClipboard = async (text: string, type: string) => {
    await navigator.clipboard.writeText(text);
    setCopied(type);
    toast({ title: "Copied!" });
    setTimeout(() => setCopied(null), 2000);
  };

  // Generate custom SVG URL
  const getCustomSvgUrl = () => {
    if (!dockerData?.account?.docker_username) return "";
    return publicApi.getHeatmapUrl(
      dockerData.account.docker_username,
      svgOptions,
    );
  };

  // Loading state
  if (authLoading) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <Loader2 className="h-5 w-5 animate-spin text-muted-foreground" />
      </div>
    );
  }

  // Not authenticated - will redirect
  if (!user) return null;

  const hasDockerAccount = !!dockerData?.account;
  const themes = themesData?.themes || DEFAULT_THEMES;

  return (
    <div className="min-h-screen bg-background">
      <main className="container py-8">
        <div className="mb-8">
          <h1 className="text-2xl font-bold mb-2">Dashboard</h1>
          <p className="text-muted-foreground">
            Connect your Docker Hub account and create embeddable heatmaps for
            your README.
          </p>
        </div>

        {/* Step 1: Connect Docker Hub */}
        <section className="mb-8">
          <div className="flex items-center gap-2 mb-4">
            <div className="flex items-center justify-center w-6 h-6 rounded-full bg-primary text-primary-foreground text-sm font-semibold">
              1
            </div>
            <h2 className="text-lg font-semibold">Connect Docker Hub</h2>
          </div>

          <DockerConnectionCard
            isLoading={dockerLoading}
            hasAccount={hasDockerAccount}
            account={dockerData?.account}
            form={form}
            onSubmit={onSubmit}
            onSync={() => syncMutation.mutate()}
            onDisconnect={() => disconnectMutation.mutate()}
            isConnecting={connectMutation.isPending}
            isSyncing={
              syncMutation.isPending ||
              (dockerData?.account?.sync_in_progress ?? false)
            }
            isDisconnecting={disconnectMutation.isPending}
          />
        </section>

        {hasDockerAccount && (
          <>
            {/* Step 2: Customize Heatmap */}
            <section className="mb-12">
              <div className="flex items-center gap-2 mb-6">
                <div className="flex items-center justify-center w-6 h-6 rounded-full bg-primary text-primary-foreground text-sm font-semibold">
                  2
                </div>
                <h2 className="text-xl font-semibold">Customize & Preview</h2>
              </div>

              <div className="grid gap-8 lg:grid-cols-3">
                <div className="lg:col-span-1">
                  <SVGCustomizationCard
                    options={svgOptions}
                    setOptions={setSvgOptions}
                    themes={themes}
                  />
                </div>

                <Card className="lg:col-span-2">
                  <CardHeader className="pb-4">
                    <CardTitle className="text-base flex items-center gap-2">
                      <Sparkles className="h-4 w-4" />
                      Live Preview
                    </CardTitle>
                    <CardDescription>
                      This is how your heatmap will look in your README
                    </CardDescription>
                  </CardHeader>
                  <CardContent>
                    <div className="rounded-xl border bg-muted/20 p-4 sm:p-8 overflow-x-auto">
                      <div className="min-w-[320px] flex justify-center items-center">
                        {/* Using a regular img tag for the dynamic SVG since Next/Image has issues with dynamic external SVGs and size optimization for small assets */}
                        <img
                          key={JSON.stringify(svgOptions)}
                          src={getCustomSvgUrl()}
                          alt="Your Docker activity heatmap"
                          className="w-full h-auto max-w-full shadow-sm rounded bg-background"
                        />
                      </div>
                    </div>
                  </CardContent>
                </Card>
              </div>
            </section>

            {/* Step 3: Get Embed Codes */}
            <section>
              <div className="flex items-center gap-2 mb-4">
                <div className="flex items-center justify-center w-6 h-6 rounded-full bg-primary text-primary-foreground text-sm font-semibold">
                  3
                </div>
                <h2 className="text-lg font-semibold">Embed in Your README</h2>
              </div>

              <EmbedCodesCard
                customUrl={getCustomSvgUrl()}
                dockerUsername={dockerData.account!.docker_username}
                copied={copied}
                onCopy={copyToClipboard}
              />
            </section>
          </>
        )}
      </main>
    </div>
  );
}
