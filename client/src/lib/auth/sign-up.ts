import { env } from "@/env";
import { delay } from "../util/delay";

type Result = { ok: true; created: boolean } | { ok: false; error: string };

export async function signUp(
	email: string,
	password: string,
	username: string,
): Promise<Result> {
	const url = env.NEXT_PUBLIC_API_URL;

	try {
		const req = fetch(`${url}/register`, {
			method: "POST",
			headers: { "Content-Type": "application/json" },
			body: JSON.stringify({ email, password, username }),
		});

		const [res] = await Promise.all([req, delay(300)]);

		if (res.status === 201) {
			return { ok: true, created: true };
		} else if (res.status === 400) {
			return { ok: false, error: "Invalid Request" };
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
