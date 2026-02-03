import { Metadata } from "next";
import { ProfileClient } from "./profile-client";
import { publicApi } from "@/lib/api";

interface PageProps {
  params: Promise<{ username: string }>;
}

export async function generateMetadata({
  params,
}: PageProps): Promise<Metadata> {
  const { username } = await params;

  try {
    const profile = await publicApi.getProfile(username);
    const displayName = profile.user.name || profile.user.github_username;

    return {
      title: `${displayName} (@${profile.user.github_username})`,
      description: `${displayName}'s Docker Hub activity heatmap. Visualize container activity like GitHub commits.`,
      openGraph: {
        title: `${displayName}'s Docker Activity | Docker Heatmap`,
        description: `Visualize ${displayName}'s Docker Hub activity like GitHub commits.`,
        images: [
          {
            url: publicApi.getHeatmapUrl(username, { theme: "github" }),
            width: 800,
            height: 300,
            alt: `${displayName}'s Docker Heatmap`,
          },
        ],
      },
    };
  } catch {
    return {
      title: "Profile Not Found",
      description: "The requested Docker Heatmap profile could not be found.",
    };
  }
}

export default async function ProfilePage({ params }: PageProps) {
  const { username } = await params;
  return <ProfileClient username={username} />;
}
