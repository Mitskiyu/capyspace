<script lang="ts">
    import { Eye, EyeOff } from "@lucide/svelte";

    let password = $state<string>("");
    let confirmPassword = $state<string>("");
    let showPassword = $state<boolean>(false);
    let showConfirmPassword = $state<boolean>(false);
</script>

<div class="relative flex w-full items-center justify-center">
    <input
        type={showPassword ? "text" : "password"}
        placeholder="Enter your password"
        bind:value={password}
        class="focus:outline-overlay2 outline-overlay3/40 bg-background3/40 relative h-9 w-11/12 rounded-lg px-2 py-1 outline-1"
    />
    <button
        type="button"
        onclick={() => (showPassword = !showPassword)}
        class="absolute top-1/2 right-5 -translate-y-1/2 transform hover:cursor-pointer focus:outline-none"
    >
        {#if showPassword}
            <EyeOff size="20" />
        {:else}
            <Eye size="20" />
        {/if}
    </button>
</div>

<div class="relative flex w-full items-center justify-center">
    <input
        type="password"
        placeholder="Confirm your password"
        bind:value={confirmPassword}
        class="focus:outline-overlay2 outline-overlay3/40 bg-background3/40 h-9 w-11/12 rounded-lg px-2 py-1 outline-1"
    />
    <button
        type="button"
        onclick={() => (showConfirmPassword = !showConfirmPassword)}
        class="absolute top-1/2 right-5 -translate-y-1/2 transform hover:cursor-pointer focus:outline-none"
    >
        {#if showConfirmPassword}
            <EyeOff size="20" />
        {:else}
            <Eye size="20" />
        {/if}
    </button>
</div>
