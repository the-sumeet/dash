<script lang="ts">
	import { GetApps, RunApp, RefreshApp } from '$bindings/changeme/appservice.js';
	import type { App } from '$bindings/changeme/models.js';
	import Icon from '@iconify/svelte';
	import { Button } from '$lib/components/ui/button/index.js';
	import RefreshCw from '@lucide/svelte/icons/refresh-cw';
	import * as Card from '$lib/components/ui/card/index.js';
	import { marked } from 'marked';
	import { Browser } from '@wailsio/runtime';
	import { onMount } from 'svelte';

	let apps = $state<App[]>([]);
	let activeIndex = $state(0);
	let items = $state<{ title?: string; subtitle?: string; markdown?: string; href?: string }[]>([]);
	let loading = $state(false);
	let error = $state('');

	async function loadApps() {
		const updated = await GetApps();
		const changed = updated.length !== apps.length || updated.some((a, i) => a.name !== apps[i]?.name);
		apps = updated;
		if (changed && apps.length > 0) {
			activeIndex = Math.min(activeIndex, apps.length - 1);
			runCommand(activeIndex);
		}
	}

	onMount(() => {
		const handler = () => {
			if (!document.hidden) loadApps();
		};
		document.addEventListener('visibilitychange', handler);
		return () => document.removeEventListener('visibilitychange', handler);
	});

	async function runCommand(index: number) {
		loading = true;
		items = [];
		error = '';
		try {
			const output = await RunApp(index);
			items = JSON.parse(output);
		} catch (e: any) {
			error = e?.message ?? 'Command failed';
		} finally {
			loading = false;
		}
	}

	async function refreshCommand(index: number) {
		loading = true;
		items = [];
		error = '';
		try {
			const output = await RefreshApp(index);
			items = JSON.parse(output);
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
					onclick={() => refreshCommand(activeIndex)}
				>
					<RefreshCw class="h-3.5 w-3.5 {loading ? 'animate-spin' : ''}" />
				</Button>
			</div>
		{/if}
		<div class="flex flex-1 flex-col gap-2 overflow-auto p-2">
			{#if apps.length === 0}
				<p class="p-3 text-sm text-muted-foreground">No apps configured.</p>
			{:else if loading}
				<p class="p-3 text-sm text-muted-foreground">Running...</p>
			{:else if error}
				<pre class="p-3 text-sm whitespace-pre-wrap text-destructive">{error}</pre>
			{:else}
				{#each items as item}
					<Card.Root
						class="py-3 {item.href ? 'cursor-pointer transition-colors hover:bg-muted/50' : ''}"
						onclick={item.href ? () => Browser.OpenURL(item.href) : undefined}
					>
						{#if item.title?.trim() || item.subtitle?.trim()}
							<Card.Header class="px-3">
								{#if item.title?.trim()}
									<Card.Title class="text-sm font-bold">{item.title}</Card.Title>
								{/if}
								{#if item.subtitle?.trim()}
									<Card.Description>{item.subtitle}</Card.Description>
								{/if}
							</Card.Header>
						{/if}
						{#if item.markdown?.trim()}
							<Card.Content class="px-3 ">
								<div
									class="border rounded-xl p-2 max-h-48 overflow-y-auto prose prose-sm max-w-none dark:prose-invert"
									onclick={(e) => {
										e.stopPropagation();
										const a = (e.target as HTMLElement).closest('a');
										if (a?.href) {
											e.preventDefault();
											Browser.OpenURL(a.href);
										}
									}}
								>
									{@html marked(item.markdown)}
								</div>
							</Card.Content>
						{/if}
					</Card.Root>
				{/each}
			{/if}
		</div>
	</div>
</div>
