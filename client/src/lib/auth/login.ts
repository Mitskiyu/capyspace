import { env } from "@/env";
import { z } from "zod";
import { delay } from "../util/delay";

type Result = { ok: true; success: boolean } | { ok: false; error: string };

const schema = z.object({
	message: z.string(),
	user: z.object({
		id: z.uuid(),
		email: z.email(),
	}),
});

export async function login(email: string, password: string): Promise<Result> {
	const url = env.NEXT_PUBLIC_API_URL;

	try {
		const req = fetch(`${url}/login`, {
			method: "POST",
			credentials: "include",
			headers: {
				"Content-Type": "application/json",
			},
			body: JSON.stringify({ email, password }),
		});

		const [res] = await Promise.all([req, delay(250)]);

		if (res.status === 401) {
			return { ok: false, error: "Invalid email or password" };
		}

		if (!res.ok) {
			return { ok: false, error: "Internal server error" };
		}

		const json = await res.json();
		const parsed = schema.safeParse(json);

		if (!parsed.success) {
			return { ok: false, error: `Invalid response format` };
		}

		return { ok: true, success: true };
	} catch (e) {
		return {
			ok: false,
			error: e instanceof Error ? e.message : "Unknown error",
		};
	}
}
