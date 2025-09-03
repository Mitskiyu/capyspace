export async function checkEmail(email: string): Promise<boolean> {
	const url = process.env.NEXT_PUBLIC_API_URL;
	if (!url) {
		throw new Error("NEXT_PUBLIC_API_URL not set");
	}

	try {
		const res = await fetch(`${url}/check-email`, {
			method: "POST",
			headers: {
				"Content-Type": "application/json",
			},
			body: JSON.stringify({ email }),
		});

		if (!res.ok) {
			throw new Error(`Request failed with status ${res.status}`);
		}

		const data: { exists: boolean } = await res.json();
		return data.exists;
	} catch (e) {
		console.error("Error checking email: ", e);
		throw e;
	}
}
