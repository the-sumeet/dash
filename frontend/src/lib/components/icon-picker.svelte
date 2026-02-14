<script lang="ts">
	import Icon from '@iconify/svelte';
	import { Input } from '$lib/components/ui/input/index.js';
	import { Label } from '$lib/components/ui/label/index.js';
	import { Button } from '$lib/components/ui/button/index.js';
	import X from '@lucide/svelte/icons/x';
	let {
		value = $bindable(''),
		label = 'Icon'
	}: {
		value: string;
		label?: string;
	} = $props();

	let query = $state('');
	let results = $state<string[]>([]);
	let loading = $state(false);
	let open = $state(false);
	let debounceTimer: ReturnType<typeof setTimeout>;

	function handleInput(q: string) {
		query = q;
		clearTimeout(debounceTimer);
		if (q.trim().length < 2) {
			results = [];
			return;
		}
		loading = true;
		debounceTimer = setTimeout(() => searchIcons(q.trim()), 300);
	}

	async function searchIcons(q: string) {
		try {
			const res = await fetch(
				`https://api.iconify.design/search?query=${encodeURIComponent(q)}&limit=48`
			);
			const data = await res.json();
			results = data.icons ?? [];
		} catch {
			results = [];
		} finally {
			loading = false;
		}
	}

	function selectIcon(icon: string) {
		value = icon;
		open = false;
		query = '';
		results = [];
	}

	function clearIcon() {
		value = '';
	}
</script>

<div class="flex flex-col gap-1.5">
	<Label>{label}</Label>

	{#if value}
		<div class="flex items-center gap-2">
			<div class="flex items-center gap-2 rounded-md border border-input px-3 py-2 text-sm">
				<Icon icon={value} width={20} height={20} />
				<span class="font-mono text-xs text-muted-foreground">{value}</span>
			</div>

			<Button onclick={clearIcon} variant="destructive" size="icon" aria-label="Submit">
				<X />
			</Button>
		</div>
	{/if}

	{#if !value || open}
		<Input
			type="text"
			placeholder="Search icons... (e.g. docker, terminal, github)"
			value={query}
			oninput={(e) => handleInput(e.currentTarget.value)}
			onfocus={() => (open = true)}
			autocapitalize="off"
			autocorrect="off"
			spellcheck="false"
		/>

		{#if loading}
			<p class="text-xs text-muted-foreground">Searching...</p>
		{/if}

		{#if results.length > 0}
			<div
				class="grid max-h-48 grid-cols-6 gap-1 overflow-y-auto rounded-md border border-input p-2"
			>
				{#each results as icon (icon)}
					<Button
						onclick={() => selectIcon(icon)}
						variant="outline"
						size="icon"
						aria-label="Submit"
					>
						<Icon {icon} width={24} height={24} />
					</Button>

					<!-- <button
						type="button"
						class="hover:bg-accent flex flex-col items-center gap-1 rounded-md p-2 transition-colors {value ===
						icon
							? 'bg-accent ring-primary ring-2'
							: ''}"
						onclick={() => selectIcon(icon)}
						title={icon}
					>
						<Icon {icon} width={24} height={24} />
					</button> -->
				{/each}
			</div>
		{:else if query.length >= 2 && !loading}
			<p class="text-xs text-muted-foreground">No icons found</p>
		{/if}
	{/if}
</div>
