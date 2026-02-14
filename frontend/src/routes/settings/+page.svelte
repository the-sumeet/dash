<script lang="ts">
	import * as Tabs from '$lib/components/ui/tabs/index.js';
	import * as Card from '$lib/components/ui/card/index.js';
</script>

<div class="h-screen overflow-y-auto p-4">
	<Tabs.Root value="account">
		<Tabs.List>
			<Tabs.Trigger value="apps">Apps</Tabs.Trigger>
			<Tabs.Trigger value="settings">Settings</Tabs.Trigger>
		</Tabs.List>
		<Tabs.Content value="apps">Make changes to your apps here.</Tabs.Content>
		<Tabs.Content value="settings">Change your settings here.</Tabs.Content>
	</Tabs.Root>
</div>
