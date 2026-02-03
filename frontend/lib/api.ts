import type {
  User,
  DockerAccount,
  ConnectDockerRequest,
  UpdateProfileRequest,
  ActivityResponse,
  ProfileData,
  EmbedCodes,
  ThemesResponse,
  SVGOptions,
} from "./schemas";

const API_URL = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080/api";

class ApiError extends Error {
  constructor(
    public status: number,
    message: string,
  ) {
    super(message);
    this.name = "ApiError";
  }
}

async function fetchApi<T>(
  endpoint: string,
  options: RequestInit = {},
): Promise<T> {
  const token =
    typeof window !== "undefined" ? localStorage.getItem("token") : null;

  const headers: Record<string, string> = {
    "Content-Type": "application/json",
  };

  if (token) {
    headers["Authorization"] = `Bearer ${token}`;
  }

  const response = await fetch(`${API_URL}${endpoint}`, {
    ...options,
    headers: {
      ...headers,
      ...(options.headers as Record<string, string>),
    },
  });

  if (!response.ok) {
    let message = "An error occurred";
    try {
      const data = await response.json();
      message = data.error || data.message || message;
    } catch {
      // Ignore JSON parse error
    }
    throw new ApiError(response.status, message);
  }

  return response.json();
}

// Auth API
export const authApi = {
  getAuthUrl: (): Promise<{ auth_url: string }> => {
    return fetchApi("/auth/github");
  },

  getCurrentUser: (): Promise<{ user: User }> => {
    return fetchApi("/user/me");
  },

  logout: (): Promise<{ message: string }> => {
    return fetchApi("/auth/logout", { method: "POST" });
  },
};

// User API
export const userApi = {
  getProfile: (): Promise<{ user: User }> => {
    return fetchApi("/user/me");
  },

  updateProfile: (
    data: UpdateProfileRequest,
  ): Promise<{ user: User; message: string }> => {
    return fetchApi("/user/me", {
      method: "PUT",
      body: JSON.stringify(data),
    });
  },

  getEmbedCodes: (dockerUsername: string): Promise<EmbedCodes> => {
    return fetchApi(`/user/embed?docker_username=${dockerUsername}`);
  },
};

// Docker API
export const dockerApi = {
  connect: (
    data: ConnectDockerRequest,
  ): Promise<{ account: DockerAccount; message: string }> => {
    return fetchApi("/docker/connect", {
      method: "POST",
      body: JSON.stringify(data),
    });
  },

  getAccount: (): Promise<{ account: DockerAccount }> => {
    return fetchApi("/docker/account");
  },

  disconnect: (): Promise<{ message: string }> => {
    return fetchApi("/docker/disconnect", { method: "DELETE" });
  },

  sync: (): Promise<{ message: string }> => {
    return fetchApi("/docker/sync", { method: "POST" });
  },
};

// Helper to build SVG URL with options
function buildSVGUrl(username: string, options?: SVGOptions): string {
  const params = new URLSearchParams();

  if (options) {
    if (options.theme) params.set("theme", options.theme);
    if (options.days) params.set("days", String(options.days));
    if (options.cell_size) params.set("cell_size", String(options.cell_size));
    if (options.radius !== undefined)
      params.set("radius", String(options.radius));
    if (options.hide_legend) params.set("hide_legend", "true");
    if (options.hide_total) params.set("hide_total", "true");
    if (options.hide_labels) params.set("hide_labels", "true");
    if (options.title) params.set("title", options.title);
  }

  const queryString = params.toString();
  return `${API_URL}/heatmap/${username}.svg${queryString ? `?${queryString}` : ""}`;
}

// Public API (no auth required)
export const publicApi = {
  // Get heatmap URL with customization options
  getHeatmapUrl: (username: string, options?: SVGOptions): string => {
    return buildSVGUrl(username, options);
  },

  // Legacy method for backwards compatibility
  getHeatmapUrlSimple: (username: string, days = 365): string => {
    return buildSVGUrl(username, { days });
  },

  getActivityUrl: (username: string, days = 365): string => {
    return `${API_URL}/activity/${username}.json?days=${days}`;
  },

  getActivity: (username: string, days = 365): Promise<ActivityResponse> => {
    return fetchApi(`/activity/${username}?days=${days}`);
  },

  getProfile: (username: string): Promise<ProfileData> => {
    return fetchApi(`/profile/${username}`);
  },

  // Get available themes
  getThemes: (): Promise<ThemesResponse> => {
    return fetchApi("/themes");
  },
};

export { ApiError };
