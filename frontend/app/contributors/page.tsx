import { Metadata } from "next";
import Image from "next/image";
import Link from "next/link";
import { Github, ExternalLink } from "lucide-react";
import { Card, CardContent } from "@/components/ui/card";

export const metadata: Metadata = {
  title: "Contributors",
  description: "Meet the amazing developers who helped build Docker Heatmap.",
};

interface Contributor {
  login: string;
  avatar_url: string;
  html_url: string;
  contributions: number;
}

async function getContributors(): Promise<Contributor[]> {
  const res = await fetch(
    "https://api.github.com/repos/sagargujarathi/docker-heatmap/contributors",
    {
      next: { revalidate: 3600 }, // Cache for 1 hour
    },
  );

  if (!res.ok) return [];
  return res.json();
}

export default async function ContributorsPage() {
  const contributors = await getContributors();

  return (
    <main className="container max-w-6xl py-12 md:py-24">
      <div className="flex flex-col items-center text-center space-y-4 mb-16">
        <h1 className="text-4xl font-bold tracking-tight sm:text-5xl">
          Our Contributors
        </h1>
        <p className="text-lg text-muted-foreground max-w-2xl">
          Thank you to everyone who has contributed to making Docker Heatmap
          better. Your efforts help developers visualize their work across the
          community.
        </p>
      </div>

      <div className="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-6">
        {contributors.map((contributor) => (
          <Card
            key={contributor.login}
            className="overflow-hidden group hover:border-primary/50 transition-colors"
          >
            <CardContent className="p-0">
              <Link
                href={contributor.html_url}
                target="_blank"
                rel="noopener noreferrer"
                className="block p-6 text-center space-y-4"
              >
                <div className="relative w-20 h-20 mx-auto">
                  <Image
                    src={contributor.avatar_url}
                    alt={contributor.login}
                    fill
                    className="rounded-full object-cover border-2 border-background shadow-sm"
                  />
                  <div className="absolute -bottom-1 -right-1 bg-primary text-primary-foreground rounded-full p-1 border-2 border-background">
                    <Github className="w-3 h-3" />
                  </div>
                </div>
                <div>
                  <h3 className="font-semibold text-lg group-hover:text-primary transition-colors">
                    @{contributor.login}
                  </h3>
                  <p className="text-xs text-muted-foreground uppercase tracking-wider font-medium">
                    {contributor.contributions}{" "}
                    {contributor.contributions === 1
                      ? "Contribution"
                      : "Contributions"}
                  </p>
                </div>
                <div className="pt-2 flex justify-center text-muted-foreground group-hover:text-primary opacity-0 group-hover:opacity-100 transition-all transform translate-y-1 group-hover:translate-y-0">
                  <ExternalLink className="w-4 h-4" />
                </div>
              </Link>
            </CardContent>
          </Card>
        ))}
      </div>

      {contributors.length === 0 && (
        <div className="text-center py-20 border rounded-xl bg-muted/10">
          <p className="text-muted-foreground">
            No contributors found. Be the first!
          </p>
          <Link
            href="https://github.com/sagargujarathi/docker-heatmap"
            target="_blank"
            rel="noopener noreferrer"
            className="inline-flex items-center gap-2 mt-4 text-primary hover:underline"
          >
            Visit GitHub Repository <ExternalLink className="w-4 h-4" />
          </Link>
        </div>
      )}
    </main>
  );
}
