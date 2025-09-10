import { env } from "@/env";
import { z } from "zod";
import { delay } from "../util/delay";

type Result = { ok: true; exists: boolean } | { ok: false; error: string };

const schema = z.object({
	exists: z.boolean(),
});

export async function checkEmail(email: string): Promise<Result> {
	const url = env.NEXT_PUBLIC_API_URL;

	try {
		const req = fetch(`${url}/users/check/email`, {
			method: "POST",
			headers: {
				"Content-Type": "application/json",
			},
			body: JSON.stringify({ email }),
		});

		const [res] = await Promise.all([req, delay(400)]);

		if (!res.ok) {
			return { ok: false, error: `Internal server error` };
		}

		const json = await res.json();
		const parsed = schema.safeParse(json);

		if (!parsed.success) {
			return { ok: false, error: "Invalid response format" };
		}

		return { ok: true, exists: parsed.data.exists };
	} catch (e) {
		return {
			ok: false,
			error: e instanceof Error ? e.message : "Unknown error",
		};
	}
}
