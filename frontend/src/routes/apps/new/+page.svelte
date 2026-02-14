<script lang="ts">
	import { Button } from '$lib/components/ui/button/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import { Label } from '$lib/components/ui/label/index.js';
	import * as Card from '$lib/components/ui/card/index.js';
	import { Badge } from '$lib/components/ui/badge/index.js';
	import IconPicker from '$lib/components/icon-picker.svelte';

	let appName = $state('');
	let appIcon = $state('');
	let command = $state('');
	let argument = $state('');
	let args = $state<string[]>([]);

	let dragIndex = $state<number | null>(null);
	let dropIndex = $state<number | null>(null);

	function addArgument() {
		const newArg = argument.trim();
		if (newArg === '') return;
		args = [...args, newArg];
		argument = '';
	}

	function removeArgument(index: number) {
		args = args.filter((_, i) => i !== index);
	}

	function handleDragStart(index: number) {
		dragIndex = index;
	}

	function handleDragOver(e: DragEvent, index: number) {
		e.preventDefault();
		dropIndex = index;
	}

	function handleDrop(index: number) {
		if (dragIndex === null || dragIndex === index) return;
		const updated = [...args];
		const [moved] = updated.splice(dragIndex, 1);
		updated.splice(index, 0, moved);
		args = updated;
		dragIndex = null;
		dropIndex = null;
	}

	function handleDragEnd() {
		dragIndex = null;
		dropIndex = null;
	}

	function isValid() {
		return appIcon.trim() !== '' && appName.trim() !== '' && command.trim() !== '';
	}
</script>

<div class="flex h-screen w-screen flex-col items-center overflow-y-auto p-4">
	<div class="flex w-full max-w-2xl flex-1 flex-col gap-3">
		<Card.Root>
			<Card.Header>
				<Card.Title>Command</Card.Title>
			</Card.Header>

			<Card.Content class="flex flex-col gap-4">
				<div class="flex w-full flex-col gap-1.5">
					<Label>Name</Label>
					<Input bind:value={appName} type="text" placeholder="Docker" />
				</div>

				<IconPicker bind:value={appIcon} label="App Icon" />

				<div class="flex w-full flex-col gap-1.5">
					<Label>Command</Label>
					<Input
						class="font-mono"
						bind:value={command}
						type="text"
						placeholder="ls"
						autocapitalize="off"
						autocorrect="off"
						spellcheck="false"
					/>
				</div>

				<div class="flex flex-col gap-1.5">
                    <Label>Arguments</Label>

				<div class="flex w-full items-center gap-2">
					<Input
						class="font-mono"
						bind:value={argument}
						type="text"
						placeholder="-l"
						autocapitalize="off"
						autocorrect="off"
						spellcheck="false"
					/>
					<Button onclick={() => addArgument()} type="submit" variant="outline">Add argument</Button
					>
				</div>
                </div>
				<div class="flex flex-wrap gap-2">
					{#each args as arg, i (i)}
						<button
							type="button"
							draggable="true"
							ondragstart={() => handleDragStart(i)}
							ondragover={(e) => handleDragOver(e, i)}
							ondrop={() => handleDrop(i)}
							ondragend={handleDragEnd}
							class="group cursor-grab active:cursor-grabbing"
						>
							<Badge
								variant="outline"
								class="font-mono transition-opacity select-none {dragIndex === i
									? 'opacity-40'
									: ''} {dropIndex === i && dragIndex !== i ? 'ring-2 ring-primary' : ''}"
							>
								<span class="mr-1">{arg}</span>
								<span
									role="button"
									tabindex="0"
									class="ml-0.5 text-xs text-muted-foreground hover:text-destructive"
									onclick={(e) => {
										e.stopPropagation();
										removeArgument(i);
									}}
									onkeydown={(e) => e.key === 'Enter' && removeArgument(i)}
								>
									✕
								</span>
							</Badge>
						</button>
					{/each}
				</div>
			</Card.Content>
		</Card.Root>

		<!-- Arguments -->
		<Card.Root>
			<Card.Header>
				<Card.Title>Command Arguments</Card.Title>
			</Card.Header>
			<Card.Content class="flex flex-col gap-4"></Card.Content>
		</Card.Root>

		<Card.Root>
			<Card.Header>
				<Card.Title>Final Command</Card.Title>
			</Card.Header>
			<Card.Content>
				<p class="font-mono">
					{command}
					{args.join(' ')}
				</p>
			</Card.Content>
		</Card.Root>

		<Button class="w-full" disabled={!isValid()}>Save</Button>

		<Button class="w-full" variant="secondary" href="/apps">Cancel</Button>
	</div>
</div>
