<script lang="ts">
    import { X } from "@lucide/svelte";
    import { validateEmail, checkEmail } from "$lib/auth";
    import SignIn from "../functional/SignIn.svelte";
    import SignUp from "../functional/SignUp.svelte";
    import Google from "../visual/icons/Google.svelte";
    import Icon from "../visual/icons/Icon.svelte";
    import { sendVerification } from "$lib/auth/client";

    type AuthState = "initial" | "verify" | "signup" | "signin";

    let { isModal } = $props();
    let authState = $state<AuthState>("initial");

    // user facing error
    let err = $state<string>("");

    let message = $state<string>("");

    // states for submit button
    let email = $state<string>("");
    let validEmail = $derived<boolean>(validateEmail(email));

    const handleSubmit = async (e: SubmitEvent): Promise<void> => {
        e.preventDefault();

        if (authState === "initial") {
            // check if email exists
            const { exists, error } = await checkEmail(email);

            if (error) {
                err = error;
                return;
            }

            if (!exists) {
                // try sending code
                const { success, error } = await sendVerification(email);

                // if error, stay on initial state
                if (error) {
                    err = error;
                    return;
                }

                // go to verify
                if (success) {
                    authState = "verify";
                    message = "We sent a code to your inbox";
                }
            } else {
                // go to signin
                authState = "signin";
            }
        }

        if (authState === "verify") {
            // TODO:
            // check if code correct
            // success -> go to signup
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
    <!-- close button if modal -->
    {#if isModal}
        <button class="flex w-full justify-end hover:cursor-pointer focus:outline-none">
            <X size="20" />
        </button>
    {/if}

    <div class="flex w-full flex-col items-center justify-between">
        <!-- top section -->
        <div class="focus:outline-none">
            <Icon iconSize="32" />
        </div>
        <h3 class="text-text mt-4 text-lg">Welcome to Capyspace</h3>
        <h4 class="text-subtext text-base">
            {#if authState === "initial"}Sign in or sign up{/if}
            {#if authState === "verify"}Verify your email{/if}
            {#if authState === "signin"}Sign in{/if}
            {#if authState === "signup"}Sign up{/if}
        </h4>

        <div class="mt-6 flex h-full w-full flex-col items-center justify-between gap-y-4">
            <!-- oauth -->
            <button
                class="bg-background0 hover:bg-background0/80 flex h-9 w-11/12 items-center justify-center gap-x-2 rounded-lg p-2 hover:cursor-pointer focus:outline-1 focus:outline-none"
            >
                <div class="flex items-center justify-center">
                    <Google iconSize="20" />
                </div>
                <span class="font-medium">Continue with Google</span>
            </button>

            <!-- line -->
            <div class="my-1 flex w-11/12 items-center gap-3">
                <div class="bg-background0/80 h-px flex-grow"></div>
                <span class="text-subtext text-sm font-medium">OR</span>
                <div class="bg-background0/80 h-px flex-grow"></div>
            </div>

            <!-- form -->
            <form onsubmit={handleSubmit} class="flex w-full flex-col items-center gap-y-2.5">
                <!-- error message -->
                {#if err}
                    <span class="text-error -mt-4 text-sm">{err}</span>
                {/if}

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
                            pattern="[0-9]*"
                            maxlength="6"
                            placeholder="Enter verification code"
                            oninput={() => (err = "")}
                            class="focus:outline-overlay2 outline-overlay3/40 bg-background3/40 h-9 w-11/12 rounded-lg px-2 py-1 outline-1"
                        />
                        <span class="text-subtext mt-1 text-sm">{message}</span>
                    </div>
                {:else if authState === "signup"}
                    <SignUp />
                {:else if authState === "signin"}
                    <SignIn />
                {/if}

                <!-- submit button -->
                <button
                    type="submit"
                    disabled={!validEmail}
                    class={[
                        "bg-background3 hover:bg-overlay1 focus:outline-overlay1 mt-2 h-9 w-11/12 rounded-lg focus:outline-1",
                        !validEmail ? "cursor-not-allowed opacity-60" : "hover:cursor-pointer",
                    ]}
                >
                    <span class="p-1 text-base">
                        {#if authState === "initial"}
                            Continue with email
                        {:else if authState === "verify"}
                            Verify code
                        {:else if authState === "signup"}
                            Sign up
                        {:else}
                            Sign in
                        {/if}
                    </span>
                </button>

                <!-- resend code -->
                {#if authState === "verify"}
                    <button
                        class="text-info/80 hover:text-info text-sm hover:cursor-pointer focus:outline-none"
                    >
                        Resend verification code
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
