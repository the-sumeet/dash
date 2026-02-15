<script lang="ts">
	import './layout.css';
	import favicon from '$lib/assets/favicon.svg';
	import { onMount } from 'svelte';

	let { children } = $props();

	onMount(() => {
		const mq = window.matchMedia('(prefers-color-scheme: dark)');
		function apply() {
			document.documentElement.classList.toggle('dark', mq.matches);
		}
		apply();
		mq.addEventListener('change', apply);
		const interval = setInterval(apply, 1000);
		return () => {
			mq.removeEventListener('change', apply);
			clearInterval(interval);
		};
	});
</script>

<svelte:head><link rel="icon" href={favicon} /></svelte:head>
{@render children()}
