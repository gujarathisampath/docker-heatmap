"use client";

import Link from "next/link";
import { Check, Copy, ExternalLink, Link as LinkIcon } from "lucide-react";
import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Label } from "@/components/ui/label";

interface EmbedCodesCardProps {
  customUrl: string;
  dockerUsername: string;
  copied: string | null;
  onCopy: (text: string, type: string) => void;
}

export function EmbedCodesCard({
  customUrl,
  dockerUsername,
  copied,
  onCopy,
}: EmbedCodesCardProps) {
  const markdownCode = `![Docker Activity](${customUrl})`;
  const htmlCode = `<img src="${customUrl}" alt="Docker Activity Heatmap" />`;

  const codes = [
    {
      label: "Markdown (for GitHub README)",
      value: markdownCode,
      type: "md",
      icon: "📝",
    },
    {
      label: "Image URL",
      value: customUrl,
      type: "url",
      icon: "🔗",
    },
    {
      label: "HTML",
      value: htmlCode,
      type: "html",
      icon: "🌐",
    },
  ];

  return (
    <Card>
      <CardHeader className="pb-4">
        <CardTitle className="text-base flex items-center gap-2">
          <LinkIcon className="h-4 w-4" />
          Embed Codes
        </CardTitle>
        <CardDescription>
          Copy and paste into your GitHub README or website
        </CardDescription>
      </CardHeader>
      <CardContent className="space-y-4">
        {codes.map(({ label, value, type, icon }) => (
          <div key={type} className="space-y-2">
            <Label className="text-sm text-muted-foreground flex items-center gap-2">
              <span>{icon}</span> {label}
            </Label>
            <div className="flex gap-2">
              <code className="flex-1 text-xs bg-muted px-3 py-2.5 rounded-md overflow-x-auto whitespace-nowrap font-mono">
                {value}
              </code>
              <Button
                variant="outline"
                size="icon"
                className="h-9 w-9 shrink-0"
                onClick={() => onCopy(value, type)}
              >
                {copied === type ? (
                  <Check className="h-4 w-4 text-green-500" />
                ) : (
                  <Copy className="h-4 w-4" />
                )}
              </Button>
            </div>
          </div>
        ))}

        <div className="pt-6 border-t space-y-4">
          <div className="flex flex-col sm:flex-row items-center gap-3">
            <Button
              variant="secondary"
              className="w-full sm:w-auto gap-2"
              asChild
            >
              <Link href={`/profile/${dockerUsername}`}>
                <ExternalLink className="h-4 w-4" />
                View Public Profile
              </Link>
            </Button>
            <Button
              variant="outline"
              className="w-full sm:w-auto gap-2"
              onClick={() =>
                onCopy(
                  `${window.location.origin}/profile/${dockerUsername}`,
                  "profile",
                )
              }
            >
              {copied === "profile" ? (
                <Check className="h-4 w-4 text-green-500" />
              ) : (
                <Copy className="h-4 w-4" />
              )}
              {copied === "profile" ? "Copied Link" : "Copy Profile Link"}
            </Button>
          </div>
          <p className="text-xs text-muted-foreground italic">
            Your public profile page showcases your heatmap and total
            contributions to everyone.
          </p>
        </div>
      </CardContent>
    </Card>
  );
}
