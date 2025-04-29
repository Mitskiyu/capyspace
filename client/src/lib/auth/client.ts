import { PUBLIC_API_URL } from "$env/static/public";

export async function checkEmail(email: string): Promise<{ exists?: boolean; error?: string }> {
    const url = PUBLIC_API_URL;

    try {
        const res = await fetch(`${url}/api/auth/check-email`, {
            method: "POST",
            body: JSON.stringify({ email }),
            headers: { "Content-Type": "application/json" },
        });

        const data = await res.json();
        if (!res.ok) {
            return { error: data?.error || "Service temporarily unavailable" };
        }

        return { exists: data.data };
    } catch (error) {
        return { error: error instanceof Error ? error.message : "Unknown error" };
    }
}

export async function sendVerification(
    email: string
): Promise<{ success?: boolean; error?: string }> {
    const url = PUBLIC_API_URL;

    try {
        const res = await fetch(`${url}/api/auth/send-verification`, {
            method: "POST",
            body: JSON.stringify({ email }),
            headers: { "Content-Type": "application/json" },
        });

        const data = await res.json();
        if (!res.ok) {
            return { error: data?.error || "Could not send email," };
        }

        return { success: data.data };
    } catch (error) {
        console.error(`Send verification err: ${error}`);
        return { error: "Could not send email, please try again later" };
    }
}
