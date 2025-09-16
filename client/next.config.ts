import type { NextConfig } from "next";
import { z } from "zod";

const schema = z.object({
	NEXT_PUBLIC_API_URL: z.url(),
});

const config = schema.safeParse(process.env);

if (!config.success) {
	console.error("Failed to load env:");
	console.error(config.error.issues);
	process.exit(1);
}

const nextConfig: NextConfig = {
	/* config options here */
	experimental: {
		optimizePackageImports: ["@phosphor-icons/react"],
	},
};

export default nextConfig;
