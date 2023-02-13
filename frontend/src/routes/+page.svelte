<script lang="ts">
	import { building } from '$app/environment';
	import EntryList from '$lib/entry-list.svelte';
	import { GetEntriesForDay } from '$lib/wailsjs/go/app/App';
	import type { dinkur } from '$lib/wailsjs/go/models';

	let date: Date = new Date();
	let entriesPromise = getEntries();

	async function getEntries(): Promise<dinkur.Entry[]> {
		if (building) {
			return [];
		}
		return GetEntriesForDay(date);
	}
</script>

<main class="p-4">
	{#if building}
		<div class="placeholder animate-pulse" />
	{/if}
	{#await entriesPromise}
		<div>
			<span>Loading...</span>
		</div>
	{:then entries}
		<EntryList {date} {entries} />
	{/await}
</main>
