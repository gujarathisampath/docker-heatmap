import { Metadata } from "next";
import { Github } from "lucide-react";
import Link from "next/link";
import Image from "next/image";
import { HeatmapPreview } from "@/components/landing/heatmap-preview";
import { Feature, Step } from "@/components/landing/sections";
import { HeroCTA } from "@/components/landing/hero-cta";

export const metadata: Metadata = {
  title: "Docker Heatmap | Visualize Your Container Activity",
  description:
    "Generate stunning, embeddable SVG contribution heatmaps for your Docker Hub activity. Perfect for your GitHub README and developer profiles.",
};

async function getGitHubStars() {
  try {
    const res = await fetch(
      "https://api.github.com/repos/sagargujarathi/docker-heatmap",
      {
        next: { revalidate: 3600 }, // Cache stars for 1 hour
      },
    );
    if (!res.ok) return null;
    const data = await res.json();
    return data.stargazers_count as number;
  } catch {
    return null;
  }
}

export default async function HomePage() {
  const stars = await getGitHubStars();

  return (
    <div className="min-h-screen bg-background flex flex-col">
      <main className="flex-1">
        {/* Hero */}
        <section className="container flex flex-col items-center justify-center gap-6 py-24 md:py-32 text-center">
          <Image
            src="/logo.webp"
            alt="Docker Heatmap Large Logo"
            width={180}
            height={180}
            className="mb-6 animate-in fade-in zoom-in duration-1000"
            priority
          />
          <div className="flex flex-col items-center gap-3">
            <div className="inline-flex items-center rounded-full border border-primary/20 bg-primary/5 px-3 py-1 text-xs text-primary font-medium">
              Open source · Free forever
            </div>
          </div>

          <h1 className="text-4xl sm:text-5xl md:text-6xl font-bold text-balance max-w-3xl">
            Docker Hub activity heatmaps for your README
          </h1>

          <p className="text-lg text-muted-foreground max-w-xl text-balance">
            Generate embeddable contribution graphs that automatically update.
            Show your container activity like GitHub shows commits.
          </p>

          <HeroCTA stars={stars} />
        </section>

        {/* Preview */}
        <section className="container pb-16">
          <div className="rounded-lg border bg-card p-6 md:p-8">
            <div className="text-xs text-muted-foreground mb-4 flex items-center justify-between">
              <span>@dockeruser</span>
              <span>1,247 contributions</span>
            </div>
            <HeatmapPreview />
          </div>
        </section>

        {/* Features */}
        <section className="container py-16 border-t">
          <div className="grid gap-8 md:grid-cols-3">
            <Feature title="Embeddable SVG">
              One URL. Works in GitHub README, websites, or anywhere that
              renders images.
            </Feature>
            <Feature title="Encrypted tokens">
              Your Docker Hub access token is encrypted with AES-256. Never
              stored in plaintext.
            </Feature>
            <Feature title="Manual Sync">
              Keep your activity updated with a single click from your
              dashboard.
            </Feature>
          </div>
        </section>

        {/* How it works */}
        <section className="container py-16 border-t">
          <h2 className="text-2xl font-bold mb-8">How it works</h2>
          <div className="grid gap-6 md:grid-cols-3">
            <Step num={1} title="Sign in with GitHub">
              Quick OAuth login. No password to remember.
            </Step>
            <Step num={2} title="Connect Docker Hub">
              Enter your username and access token.
            </Step>
            <Step num={3} title="Copy embed URL">
              Paste the SVG URL into your README. Done.
            </Step>
          </div>
        </section>
      </main>
    </div>
  );
}
