import { z } from "zod";

// User schema
export const userSchema = z.object({
  id: z.number(),
  github_id: z.number(),
  github_username: z.string(),
  email: z.string().nullable(),
  avatar_url: z.string(),
  name: z.string().nullable(),
  bio: z.string().nullable(),
  public_profile: z.boolean(),
  created_at: z.string(),
  updated_at: z.string(),
});

export type User = z.infer<typeof userSchema>;

// Docker Account schema
export const dockerAccountSchema = z.object({
  id: z.number(),
  docker_username: z.string(),
  is_active: z.boolean().optional(),
  auto_refresh: z.boolean().optional(),
  last_sync_at: z.string().nullable(),
  last_sync_error: z.string().nullable().optional(),
  sync_in_progress: z.boolean().optional(),
});

export type DockerAccount = z.infer<typeof dockerAccountSchema>;

// Activity schema
export const activityEventSchema = z.object({
  date: z.string(),
  count: z.number(),
  level: z.number(),
  pushes: z.number().optional(),
  pulls: z.number().optional(),
  builds: z.number().optional(),
});

export type ActivityEvent = z.infer<typeof activityEventSchema>;

// SVG Theme schema
export const themeSchema = z.object({
  id: z.string(),
  name: z.string(),
  bg_color: z.string(),
  text_color: z.string(),
  colors: z.array(z.string()),
});

export type Theme = z.infer<typeof themeSchema>;

// SVG Options for customization
export interface SVGOptions {
  theme?: string;
  days?: number;
  cell_size?: number;
  radius?: number;
  hide_legend?: boolean;
  hide_total?: boolean;
  hide_labels?: boolean;
  title?: string;
}

// API Request schemas
export const connectDockerSchema = z.object({
  docker_username: z.string().min(1, "Username is required"),
  access_token: z.string().min(1, "Access token is required"),
});

export type ConnectDockerRequest = z.infer<typeof connectDockerSchema>;

export const updateProfileSchema = z.object({
  name: z.string().optional(),
  bio: z.string().optional(),
  public_profile: z.boolean().optional(),
});

export type UpdateProfileRequest = z.infer<typeof updateProfileSchema>;

// API Response types
export interface ActivityResponse {
  username: string;
  days: number;
  totals: {
    activities: number;
    pushes: number;
    pulls: number;
    builds: number;
  };
  activity: ActivityEvent[];
}

export interface ProfileData {
  user: {
    github_username: string;
    name: string | null;
    bio: string | null;
    avatar_url: string;
  };
  docker: {
    username: string;
    last_sync_at: string | null;
  };
  stats: {
    total_activities: number;
  };
  available_themes?: string[];
}

export interface EmbedCodes {
  svg_url: string;
  json_url: string;
  markdown: string;
  html: string;
  html_link: string;
}

export interface ThemesResponse {
  themes: Theme[];
}
