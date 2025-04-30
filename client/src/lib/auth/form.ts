export function validateEmail(email: string): boolean {
    if (!email) {
        return false;
    }

    // standard html type email
    const regex =
        /^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$/;
    if (!regex.test(email)) {
        return false;
    }

    return true;
}

export function validateVerificationCode(code: string): boolean {
    if (!code) {
        return false;
    }

    if (code.length !== 6) {
        return false;
    }

    if (!/^\d{6}$/.test(code)) {
        return false;
    }

    return true;
}
