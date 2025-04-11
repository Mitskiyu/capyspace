<script lang="ts">
    import "../app.css";
    import * as Layout from "$lib/components/layout/";

    let { children } = $props();
</script>

<div class="mx-auto mt-8 min-h-screen max-w-7xl px-4">
    <Layout.Header />
    <main>
        {@render children()}
    </main>
</div>
