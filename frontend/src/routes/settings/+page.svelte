<script lang="ts">
	import * as Tabs from '$lib/components/ui/tabs/index.js';
	import { Button } from '$lib/components/ui/button/index.js';
</script>

<div class="h-screen overflow-y-auto p-4">
	<Tabs.Root value="account">
		<Tabs.List>
			<Tabs.Trigger value="apps">Apps</Tabs.Trigger>
			<Tabs.Trigger value="settings">Settings</Tabs.Trigger>
		</Tabs.List>
		<Tabs.Content value="apps" class="flex w-full flex-col items-center">
			<div class="w-full max-w-xl flex-1">
				<Button href="/new" class="w-full" variant="secondary">Add App</Button>
			</div>
		</Tabs.Content>
		<Tabs.Content value="settings">Change your settings here.</Tabs.Content>
	</Tabs.Root>
</div>
