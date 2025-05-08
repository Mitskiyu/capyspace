<script lang="ts">
    import { X } from "@lucide/svelte";
    import {
        validateEmail,
        validateVerificationCode,
        validatePassword,
        checkEmail,
        sendVerificationCode,
        checkVerificationCode,
        createUser,
        signIn,
    } from "$lib/auth";
    import SignIn from "../functional/SignIn.svelte";
    import SignUp from "../functional/SignUp.svelte";
    import Google from "../visual/icons/Google.svelte";
    import Icon from "../visual/icons/Icon.svelte";

    type AuthState = "initial" | "verify" | "signup" | "signin";

    let { isModal } = $props();
    let authState = $state<AuthState>("initial");

    // Verification code message
    let message = $state<string>("Code will be sent to your inbox");

    // Form states
    let email = $state<string>("");
    let password = $state<string>("");
    let confirmPassword = $state<string>("");
    let verificationCode = $state<string>("");

    let validEmail = $derived<boolean>(validateEmail(email));
    let validPassword = $derived<boolean>(validatePassword(password));
    let validVerificationCode = $derived<boolean>(validateVerificationCode(verificationCode));

    // User-facing error
    let err = $state<string>("");

    const handleSubmit = async (e: SubmitEvent): Promise<void> => {
        e.preventDefault();

        if (authState === "initial") {
            // Check if a user with this email exists
            const { exists, error } = await checkEmail(email);
            if (error) {
                err = error;
                return;
            }

            if (!exists) {
                // User does not exist, go to sign up
                authState = "signup";
                return;
            } else {
                // User exists, go to sign in
                authState = "signin";
                return;
            }
        }

        if (authState === "signup") {
            // Send the email and continue to verification
            sendVerificationCode(email);
            authState = "verify";
            return;
        }

        if (authState === "verify") {
            // Check if the code is correct
            const { verified, error } = await checkVerificationCode(email, verificationCode);
            if (error) {
                err = error;
                return;
            }

            if (!verified) {
                err = "Code is invalid or expired";
            } else {
                // Code is valid, create the user
                const { success, error } = await createUser(email, password, verificationCode);
                if (error) {
                    err = error;
                    return;
                }

                if (!success) {
                    err = "Could not sign up, try again later";
                    return;
                } else {
                    console.log("yay");
                    // TODO:
                    // Issue session & go to onboarding
                }
            }

            return;
        }

        if (authState === "signin") {
            const { success, error } = await signIn(email, password);
            if (error) {
                err = error;
                return;
            }

            if (!success) {
                err = "Could not sign in, try again later";
                return;
            } else {
                // Go to the space
            }
        }
    };
</script>

<div
    class={[
        "bg-background2 text-text font-body outline-overlay2/40 h-auto w-80 rounded-xl px-4 py-3 outline-1",
        { "min-h-[20em]": authState === "initial" },
        { "min-h-[30em]": authState === "verify" },
        { "min-h-[30em]": authState === "signup" },
        { "min-h-[28em]": authState === "signin" },
    ]}
>
    <!-- Show close button if modal -->
    {#if isModal}
        <button class="flex w-full justify-end hover:cursor-pointer focus:outline-none">
            <X size="20" />
        </button>
    {/if}

    <div class="flex w-full flex-col items-center justify-between">
        <!-- Top section -->
        <div class="focus:outline-none">
            <Icon iconSize="32" />
        </div>
        <h3 class="text-text mt-4 text-lg">Welcome to Capyspace</h3>
        <h4 class="text-subtext text-base">
            {#if authState === "initial"}Sign in or sign up{/if}
            {#if authState === "verify"}Sign up{/if}
            {#if authState === "signin"}Sign in{/if}
            {#if authState === "signup"}Sign up{/if}
        </h4>

        <div class="mt-6 flex h-full w-full flex-col items-center justify-between gap-y-4">
            <!-- Google sign in -->
            <button
                class="bg-background0 hover:bg-background0/80 flex h-9 w-11/12 items-center justify-center gap-x-2 rounded-lg p-2 hover:cursor-pointer focus:outline-1 focus:outline-none"
            >
                <div class="flex items-center justify-center">
                    <Google iconSize="20" />
                </div>
                <span class="font-medium">Continue with Google</span>
            </button>

            <!-- Line -->
            <div class="my-1 flex w-11/12 items-center gap-3">
                <div class="bg-background0/80 h-px flex-grow"></div>
                <span class="text-subtext text-sm font-medium">OR</span>
                <div class="bg-background0/80 h-px flex-grow"></div>
            </div>

            <!-- Form -->
            <form onsubmit={handleSubmit} class="flex w-full flex-col items-center gap-y-2.5">
                <!-- Error -->
                <div class="-mt-3 flex w-11/12 items-center justify-center text-center">
                    {#if err}
                        <span class="text-error text-sm">{err}</span>
                    {/if}
                </div>

                <input
                    type="text"
                    placeholder="Enter your email"
                    bind:value={email}
                    oninput={() => {
                        err = "";
                        if (authState !== "initial") {
                            authState = "initial";
                        }
                    }}
                    class="focus:outline-overlay2 outline-overlay3/40 bg-background3/40 h-9 w-11/12 rounded-lg px-2 py-1 outline-1"
                />
                {#if authState === "verify"}
                    <div class="flex w-full flex-col items-center justify-center">
                        <input
                            type="text"
                            inputmode="numeric"
                            placeholder="Enter verification code"
                            bind:value={verificationCode}
                            oninput={() => (err = "")}
                            class="focus:outline-overlay2 outline-overlay3/40 bg-background3/40 h-9 w-11/12 rounded-lg px-2 py-1 outline-1"
                        />
                        <span class="text-subtext mt-1 text-sm">{message}</span>
                    </div>
                {:else if authState === "signup"}
                    <SignUp bind:password bind:confirmPassword bind:err />
                {:else if authState === "signin"}
                    <SignIn bind:password />
                {/if}

                <!-- Submit -->
                <button
                    type="submit"
                    disabled={(authState === "initial" && !validEmail) ||
                        (authState === "signup" &&
                            (!validPassword || password !== confirmPassword)) ||
                        (authState === "verify" && !validVerificationCode) ||
                        (authState === "signin" && !password)}
                    class={[
                        "bg-background3 hover:bg-overlay1 focus:outline-overlay1 mt-2 h-9 w-11/12 rounded-lg focus:outline-1",
                        (authState === "initial" && !validEmail) ||
                        (authState === "signup" &&
                            (!validPassword || password !== confirmPassword)) ||
                        (authState === "verify" && !validVerificationCode) ||
                        (authState === "signin" && !password)
                            ? "cursor-not-allowed opacity-60"
                            : "hover:cursor-pointer",
                    ]}
                >
                    <span class="p-1 text-base">
                        {#if authState === "initial"}
                            Continue with email
                        {:else if authState === "verify"}
                            Create account
                        {:else if authState === "signup"}
                            Continue
                        {:else}
                            Sign in
                        {/if}
                    </span>
                </button>

                <!-- Optional -->
                {#if authState === "verify"}
                    <button
                        class="text-info/80 hover:text-info text-sm hover:cursor-pointer focus:outline-none"
                    >
                        Resend in 60s
                    </button>
                {:else if authState === "signin"}
                    <button
                        class="text-info/80 hover:text-info text-sm hover:cursor-pointer focus:outline-none"
                    >
                        Forgot your password?
                    </button>
                {/if}
            </form>
            <div class="h-2"></div>
        </div>
    </div>
</div>
