import { Metadata } from "next";
import { DashboardClient } from "./dashboard-client";

export const metadata: Metadata = {
  title: "Dashboard",
  description:
    "Manage your Docker Hub connection and customize your activity heatmap.",
};

export default function DashboardPage() {
  return <DashboardClient />;
}
