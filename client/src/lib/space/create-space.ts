import { env } from "@/env";

type Result = { ok: true; created: boolean } | { ok: false; error: string };

export async function createSpace(): Promise<Result> {
	const url = env.NEXT_PUBLIC_API_URL;

	try {
		const res = await fetch(`${url}/spaces`, {
			method: "POST",
			credentials: "include",
			headers: { "Content-Type": "application/json" },
		});

		if (res.status === 201) {
			return { ok: true, created: true };
		} else if (res.status === 409) {
			return { ok: false, error: "Space already exists" };
		} else if (res.status === 401) {
			return { ok: false, error: "Unauthorized" };
		} else {
			return { ok: false, error: "Internal server error" };
		}
	} catch (e) {
		return {
			ok: false,
			error: e instanceof Error ? e.message : "Unknown error",
		};
	}
}
