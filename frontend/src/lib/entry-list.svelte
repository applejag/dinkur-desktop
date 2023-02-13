<script lang="ts">
	import Entry from './entry.svelte';

	import type { dinkur } from '$lib/wailsjs/go/models';

	export let date = new Date();
	export let entries: dinkur.Entry[] = [];

	function dayString(day: number): string {
		switch (day) {
			case 1:
				return 'Monday';
			case 2:
				return 'Tuesday';
			case 3:
				return 'Wednesday';
			case 4:
				return 'Thursday';
			case 5:
				return 'Friday';
			case 6:
				return 'Saturday';
			case 7:
				return 'Sunday';
			default:
				return 'Day ' + day;
		}
	}
</script>

<section>
	<header class="pb-2 text-center">
		<h2 class="text-primary-500">{dayString(date.getDay())}</h2>
		<h6 class="text-tertiary-500">{date.toLocaleDateString()}</h6>
	</header>

	<div>
		<ol class="list">
			{#each entries as entry}
				<li>
					<span class="flex-auto">
						<Entry {entry} />
					</span>
				</li>
			{/each}
		</ol>
	</div>
	{#if entries.length == 0}
		<div class="card variant-soft-tertiary-200 p-4">
			<p>There are no entries on this day.</p>
		</div>
	{/if}
</section>

<style>
	li {
		padding-left: 0;
		padding-right: 0;
		padding-bottom: 0;
	}
</style>
