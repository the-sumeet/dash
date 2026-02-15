<script lang="ts">
	import { GetApps, RunApp } from '$bindings/changeme/appservice.js';
	import type { App } from '$bindings/changeme/models.js';
	import Icon from '@iconify/svelte';
	import { Button } from '$lib/components/ui/button/index.js';
	import RefreshCw from '@lucide/svelte/icons/refresh-cw';

	let apps = $state<App[]>([]);
	let activeIndex = $state(0);
	let output = $state('');
	let loading = $state(false);
	let error = $state('');

	async function loadApps() {
		apps = await GetApps();
		if (apps.length > 0) {
			runCommand(0);
		}
	}

	async function runCommand(index: number) {
		loading = true;
		output = '';
		error = '';
		try {
			output = await RunApp(index);
			console.log('Command output:', output);
		} catch (e: any) {
			error = e?.message ?? 'Command failed';
		} finally {
			loading = false;
		}
	}

	function selectApp(index: number) {
		activeIndex = index;
		runCommand(index);
	}

	loadApps();
</script>

<div class="flex h-screen w-screen">
	<!-- Apps -->
	<div class="flex flex-col gap-1 overflow-auto border-r p-2">
		{#each apps as app, i}
			<Button
				class={activeIndex == i ? 'bg-primary/30' : ''}
				variant="secondary"
				size="icon"
				title={app.name}
				onclick={() => selectApp(i)}
			>
				<Icon icon={app.icon} class="h-4 w-4" />
			</Button>
		{/each}
	</div>

	<!-- Content -->
	<div class="flex flex-1 flex-col overflow-hidden">
		{#if apps.length > 0}
			<div class="flex items-center justify-between border-b px-3 py-2">
				<span class="text-sm font-medium">{apps[activeIndex]?.name}</span>
				<Button
					variant="ghost"
					size="icon"
					class="h-7 w-7"
					disabled={loading}
					onclick={() => runCommand(activeIndex)}
				>
					<RefreshCw class="h-3.5 w-3.5 {loading ? 'animate-spin' : ''}" />
				</Button>
			</div>
		{/if}
		<div class="flex-1 overflow-auto p-3">
			{#if apps.length === 0}
				<p class="text-sm text-muted-foreground">No apps configured.</p>
			{:else if loading}
				<p class="text-sm text-muted-foreground">Running...</p>
			{:else if error}
				<pre class="whitespace-pre-wrap text-sm text-destructive">{error}</pre>
			{:else}
				<pre class="whitespace-pre-wrap font-mono text-sm">{output}</pre>
			{/if}
		</div>
	</div>
</div>
