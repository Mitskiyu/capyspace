import { PUBLIC_API_URL } from "$env/static/public";

export async function checkEmail(email: string): Promise<{ exists?: boolean; error?: string }> {
    const url = PUBLIC_API_URL;

    try {
        const res = await fetch(`${url}/auth/check-email`, {
            method: "POST",
            body: JSON.stringify({ email }),
            headers: { "Content-Type": "application/json" },
        });

        const data = await res.json();
        if (!res.ok) {
            return { error: data?.error || "Could not check email, try again later" };
        }

        return { exists: data.data };
    } catch (error) {
        console.error(`Check email err: ${error}`);
        return { error: "Could not check email, try again later" };
    }
}

export async function sendVerificationCode(
    email: string
): Promise<{ success?: boolean; error?: string }> {
    const url = PUBLIC_API_URL;

    try {
        const res = await fetch(`${url}/auth/send-verification`, {
            method: "POST",
            body: JSON.stringify({ email }),
            headers: { "Content-Type": "application/json" },
        });

        const data = await res.json();
        if (!res.ok) {
            return { error: data?.error || "Could not send email, try again later" };
        }

        return { success: data.data };
    } catch (error) {
        console.error(`Send verification err: ${error}`);
        return { error: "Could not send email, try again later" };
    }
}

export async function checkVerificationCode(
    email: string,
    code: string
): Promise<{ verified?: boolean; error?: string }> {
    const url = PUBLIC_API_URL;

    try {
        const res = await fetch(`${url}/auth/check-verification`, {
            method: "POST",
            body: JSON.stringify({ email, code }),
            headers: { "Content-Type": "application/json" },
        });

        const data = await res.json();
        if (!res.ok) {
            return { error: data?.error || "Could not verify, try again later" };
        }

        return { verified: data.data };
    } catch (error) {
        console.error(error);
        return { error: "Could not verify, try again later" };
    }
}

export async function createUser(
    email: string,
    password: string
): Promise<{ success?: boolean; error?: string }> {
    const url = PUBLIC_API_URL;

    try {
        const res = await fetch(`${url}/auth/create-user`, {
            method: "POST",
            body: JSON.stringify({ email, password }),
            headers: { "Content-Type": "application/json" },
        });

        const data = await res.json();
        if (!res.ok) {
            return { error: data?.error || "Could not sign up, try again later" };
        }

        return { success: data.data };
    } catch (error) {
        console.error(error);
        return { error: "Could not sign up, try again later" };
    }
}
