<script lang="ts">
	import { GetApps } from '$bindings/changeme/appservice.js';
	import type { App } from '$bindings/changeme/models.js';
	import Icon from '@iconify/svelte';
	import { Button } from '$lib/components/ui/button/index.js';
	let apps = $state<App[]>([]);
	let activeIndex = $state(0);

	async function loadApps() {
		apps = await GetApps();
	}

	loadApps();
</script>

<div class="flex h-screen w-screen">
	<!-- Apps -->
	<div class="flex flex-col gap-1 overflow-auto border-r p-2">
		{#each apps as app, i}
			<Button class="{activeIndex == i ? "bg-primary/30" : ""}" variant="secondary" size="icon" aria-label="Submit" onclick={() => (activeIndex = i)}>
				<Icon icon={app.icon} class="h-4 w-4" />
			</Button>

			<!-- <button
				class="flex items-center justify-center rounded-lg p-2 transition-colors {activeIndex === i
					? 'bg-primary text-primary-foreground'
					: 'hover:bg-muted'}"
				title={app.name}
				onclick={() => (activeIndex = i)}
			>
			</button> -->
		{/each}
	</div>

	<!-- Content -->
	<div class="flex-1 overflow-auto p-4">
		{#if apps.length === 0}
			<p class="text-sm text-muted-foreground">No apps configured.</p>
		{:else if apps[activeIndex]}
			<p class="text-sm text-muted-foreground">{apps[activeIndex].name}</p>
		{/if}
	</div>
</div>
