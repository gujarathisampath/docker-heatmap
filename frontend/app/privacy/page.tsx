import { Metadata } from "next";

export const metadata: Metadata = {
  title: "Privacy Policy",
  description:
    "Learn about how Docker Heatmap collects, uses, and protects your data.",
};

export default function PrivacyPage() {
  return (
    <main className="container max-w-4xl py-12 md:py-24">
      <div className="space-y-6">
        <h1 className="text-4xl font-bold tracking-tight">Privacy Policy</h1>
        <p className="text-muted-foreground italic">
          Last update: February 03, 2026
        </p>

        <section className="space-y-4 pt-4">
          <h2 className="text-2xl font-semibold">1. Introduction</h2>
          <p>
            Docker Heatmap ("we", "us", or "our") respects your privacy and is
            committed to protecting your personal data. This privacy policy will
            inform you as to how we look after your personal data when you visit
            our website and tell you about your privacy rights and how the law
            protects you.
          </p>
        </section>

        <section className="space-y-4">
          <h2 className="text-2xl font-semibold">2. The Data We Collect</h2>
          <p>We collect and process the following data:</p>
          <ul className="list-disc pl-6 space-y-2">
            <li>
              <strong>GitHub Account Information:</strong> To manage your
              account, we collect your GitHub username, avatar URL, and internal
              GitHub ID when you sign in.
            </li>
            <li>
              <strong>Docker Hub Information:</strong> To generate your activity
              heatmap, we store your Docker Hub username and the Read-only
              access token you provide.
            </li>
            <li>
              <strong>Activity Data:</strong> We fetch and store metadata about
              your public Docker Hub activity (pushes, tags, etc.) to generate
              the visualizations.
            </li>
          </ul>
        </section>

        <section className="space-y-4">
          <h2 className="text-2xl font-semibold">3. How We Use Your Data</h2>
          <p>We use your data only for the following purposes:</p>
          <ul className="list-disc pl-6 space-y-2">
            <li>To provide and maintain our Service.</li>
            <li>To generate and serve your custom Docker activity heatmaps.</li>
            <li>To manage your account and authentication via GitHub.</li>
          </ul>
        </section>

        <section className="space-y-4">
          <h2 className="text-2xl font-semibold">4. Data Security</h2>
          <p>
            We prioritize the security of your data. Your Docker Hub Access
            Tokens are **encrypted at rest** using industry-standard AES-256
            encryption. We specifically recommend providing **Read-only** tokens
            to minimize security risks.
          </p>
        </section>

        <section className="space-y-4">
          <h2 className="text-2xl font-semibold">5. Third-Party Services</h2>
          <p>We interact with the following third-party services:</p>
          <ul className="list-disc pl-6 space-y-2">
            <li>
              <strong>GitHub:</strong> For authentication and user profiles.
            </li>
            <li>
              <strong>Docker Hub:</strong> To fetch your repository activity
              metadata.
            </li>
          </ul>
        </section>

        <section className="space-y-4">
          <h2 className="text-2xl font-semibold">6. Contact Us</h2>
          <p>
            If you have any questions about this Privacy Policy, you can contact
            us via our GitHub repository.
          </p>
        </section>
      </div>
    </main>
  );
}
