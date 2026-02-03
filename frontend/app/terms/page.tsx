import { Metadata } from "next";

export const metadata: Metadata = {
  title: "Terms of Service",
  description:
    "Read the terms and conditions for using the Docker Heatmap service.",
};

export default function TermsPage() {
  return (
    <main className="container max-w-4xl py-12 md:py-24">
      <div className="space-y-6">
        <h1 className="text-4xl font-bold tracking-tight">Terms of Service</h1>
        <p className="text-muted-foreground italic">
          Last update: February 03, 2026
        </p>

        <section className="space-y-4 pt-4">
          <h2 className="text-2xl font-semibold">1. Acceptance of Terms</h2>
          <p>
            By accessing and using Docker Heatmap ("the Service"), you agree to
            be bound by these Terms of Service. If you do not agree to these
            terms, please do not use the Service.
          </p>
        </section>

        <section className="space-y-4">
          <h2 className="text-2xl font-semibold">2. Description of Service</h2>
          <p>
            Docker Heatmap is an open-source tool that allows users to visualize
            their Docker Hub activity through embeddable SVG heatmaps. The
            Service is currently provided for free.
          </p>
        </section>

        <section className="space-y-4">
          <h2 className="text-2xl font-semibold">3. User Responsibilities</h2>
          <p>When using the Service, you agree:</p>
          <ul className="list-disc pl-6 space-y-2">
            <li>To provide accurate information (Docker Hub username).</li>
            <li>To use Read-only Access Tokens whenever possible.</li>
            <li>
              Not to use the Service for any illegal or unauthorized purpose.
            </li>
            <li>
              Not to attempt to disrupt or overwhelm the Service's
              infrastructure (rate limiting is in place).
            </li>
          </ul>
        </section>

        <section className="space-y-4">
          <h2 className="text-2xl font-semibold">4. Intellectual Property</h2>
          <p>
            The Service is open-source. The source code is available on GitHub
            under its respective license. The generated heatmap SVGs are yours
            to use in your profiles and READMEs.
          </p>
        </section>

        <section className="space-y-4">
          <h2 className="text-2xl font-semibold">
            5. Disclaimer of Warranties
          </h2>
          <p>
            The Service is provided "as is" and "as available" without any
            warranties of any kind, either express or implied. We do not
            guarantee that the Service will always be available, accurate, or
            error-free.
          </p>
        </section>

        <section className="space-y-4">
          <h2 className="text-2xl font-semibold">6. Limitation of Liability</h2>
          <p>
            To the maximum extent permitted by law, we shall not be liable for
            any indirect, incidental, special, consequential, or punitive
            damages resulting from your use of or inability to use the Service.
          </p>
        </section>

        <section className="space-y-4">
          <h2 className="text-2xl font-semibold">7. Changes to Terms</h2>
          <p>
            We reserve the right to modify these terms at any time. We will
            notify users of any significant changes by updating the "Last
            update" date at the top of this page.
          </p>
        </section>
      </div>
    </main>
  );
}
