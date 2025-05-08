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

    if (code.length !== 8) {
        return false;
    }

    if (!/^\d{8}$/.test(code)) {
        return false;
    }

    return true;
}

export function validatePassword(password: string): boolean {
    if (!password) {
        return false;
    }

    if (password.length < 8 || password.length > 255) {
        return false;
    }

    if (password.trim() === "") {
        return false;
    }

    return true;
}
