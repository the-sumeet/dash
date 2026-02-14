<script lang="ts">
	import { Button } from '$lib/components/ui/button/index.js';
	import * as Card from '$lib/components/ui/card/index.js';
	import { GetApps, DeleteApp } from '$bindings/changeme/appservice.js';
	import type { App } from '$bindings/changeme/models.js';
	import Icon from '@iconify/svelte';
    import Trash from '@lucide/svelte/icons/trash';

	let apps = $state<App[]>([]);
	let loading = $state(true);

	async function loadApps() {
		loading = true;
		apps = await GetApps();
		loading = false;
	}

	async function handleDelete(index: number) {
		await DeleteApp(index);
		await loadApps();
	}

	loadApps();
</script>

<div class="flex h-screen w-screen flex-col items-center overflow-y-auto p-4">
	<div class="flex w-full max-w-2xl flex-1 flex-col gap-3">
		{#if loading}
			<p class="text-sm text-muted-foreground">Loading...</p>
		{:else if apps.length === 0}
        			<Button href="/apps/new" class="w-full" variant="secondary">Add App</Button>

		{:else}
			<Button href="/apps/new" class="w-full" variant="secondary">Add App</Button>

			{#each apps as app, i}
				<Card.Root class="group">
					<Card.Content class="flex items-center gap-3 p-4">
						<Icon icon={app.icon} class="h-8 w-8 shrink-0" />
						<div class="min-w-0 flex-1">
							<p class="font-medium">{app.name}</p>
							<p class="font-mono text-sm text-muted-foreground">
								{app.command}
								{app.args.join(' ')}
							</p>
						</div>
					</Card.Content>
                    <Card.Footer class="hidden group-hover:flex">
                        <Button variant="destructive" size="icon" aria-label="Delete" onclick={() => handleDelete(i)}>
							<Trash />
						</Button>
                    </Card.Footer>
				</Card.Root>
			{/each}
		{/if}
	</div>
</div>
