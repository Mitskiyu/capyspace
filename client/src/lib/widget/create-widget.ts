import { env } from "@/env";
import { Widget } from "@/lib/types";

type Result = { ok: true } | { ok: false; error: string };

export async function createWidget(
	widget: Widget,
	spaceId: string,
): Promise<Result> {
	const url = env.NEXT_PUBLIC_API_URL;

	try {
		const res = await fetch(`${url}/spaces/${spaceId}/widgets`, {
			method: "POST",
			credentials: "include",
			headers: { "Content-Type": "application/json" },
			body: JSON.stringify(widget),
		});

		if (res.status === 201) {
			return { ok: true };
		} else if (res.status == 403) {
			return {
				ok: false,
				error: "You do not have permission to edit this space",
			};
		} else if (res.status == 401) {
			return { ok: false, error: "You are not logged in" };
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
