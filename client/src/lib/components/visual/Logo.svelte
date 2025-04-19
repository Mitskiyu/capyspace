<script lang="ts">
    import Icon from "./icons/Icon.svelte";

    let { title, fontSize, iconSize } = $props();
</script>

<div class="flex flex-row gap-x-3">
    <Icon {iconSize} />
    <h1 class={`font-head text-primary ${fontSize}`}>{title}</h1>
</div>
