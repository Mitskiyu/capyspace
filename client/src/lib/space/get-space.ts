import { env } from "@/env";
import { spaceSchema, type Space } from "@/lib/types";

type Result =
	| { status: "success"; space: Space }
	| { status: "not_found" }
	| { status: "private" }
	| { status: "error"; message: string };

export async function getSpace(username: string): Promise<Result> {
	const url = env.NEXT_PUBLIC_API_URL;

	try {
		const res = await fetch(`${url}/spaces/${username}`, {
			method: "GET",
			headers: { "Content-Type": "application/json" },
			credentials: "include",
		});

		if (res.status === 404) {
			return { status: "not_found" };
		}

		if (res.status === 403) {
			return { status: "private" };
		}

		if (!res.ok) {
			return {
				status: "error",
				message: "Internal server error",
			};
		}

		const json = await res.json();
		const parsed = spaceSchema.safeParse(json);

		if (!parsed.success) {
			return {
				status: "error",
				message: "Invalid response format",
			};
		}

		return {
			status: "success",
			space: parsed.data,
		};
	} catch (e) {
		return {
			status: "error",
			message: e instanceof Error ? e.message : "Unknown error",
		};
	}
}
