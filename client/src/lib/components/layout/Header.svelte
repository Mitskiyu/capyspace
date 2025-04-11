<script lang="ts">
    import { page } from "$app/state";
    import Logo from "$lib/components/visual/Logo.svelte";

    const nav = [
        { href: "/", label: "Home" },
        { href: "/product", label: "Product" },
        { href: "/about", label: "About" },
        { href: "/pricing", label: "Pricing" },
    ];
</script>

<header
    class="bg-background1 relative flex h-16 w-full items-center justify-between rounded-xl px-4 py-3"
>
    <Logo title="CAPYSPACE" fontSize="text-4xl" iconSize="40" />
    <nav
        class="font-body absolute top-1/2 left-1/2 flex -translate-x-1/2 -translate-y-1/2 gap-x-6 text-lg"
    >
        {#each nav as link (link.href)}
            {@const isActive = page.url.pathname === link.href}
            <a
                href={link.href}
                class="hover:text-primary transition-colors duration-300
                {isActive ? 'text-primary' : 'text-secondary'}">{link.label}</a
            >
        {/each}
    </nav>
    <button
        class="bg-primary font-body text-background0 h-full w-32 rounded-lg text-base hover:cursor-pointer"
    >
        Get Started
    </button>
</header>
