<script lang="ts">
	import type { dinkur } from './wailsjs/go/models';

	export let entry: dinkur.Entry;
	$: startDate = parseDate(entry.start);
	$: start = timeString(startDate);
	$: endDate = parseDate(entry.end);
	$: end = timeString(endDate);
	$: durationSeconds = calcDurationSeconds(startDate, endDate);
	$: duration = formatDurationSeconds(durationSeconds);
	$: console.log('Start & end:', { start, end });

	function parseDate(dateString: string | null): Date | null {
		if (dateString === null) {
			return null;
		}
		return new Date(dateString);
	}

	function calcDurationSeconds(start: Date | null, end: Date | null): number {
		start = start ?? new Date();
		end = end ?? new Date();
		const startMs = start.getTime();
		const endMs = end.getTime();
		const diffMs = endMs - startMs;
		return diffMs / 1000;
	}

	function formatDurationSeconds(totalSeconds: number): string {
		const seconds = Math.floor(totalSeconds) % 60;
		const totalMinutes = Math.floor(totalSeconds / 60);
		const minutes = totalMinutes % 60;
		const hours = Math.floor(totalMinutes / 60);

		return `${hours.toString().padStart(2, '0')}:${minutes.toString().padStart(2, '0')}:${seconds
			.toString()
			.padStart(2, '0')}`;
	}

	function timeString(date: Date | null): string | null {
		if (date === null) {
			return null;
		}
		const minutes = date.getMinutes().toString();
		const hours = date.getHours().toString();
		return `${hours.padStart(2, '0')}:${minutes.padStart(2, '0')}`;
	}
</script>

<form>
	<div class="input-group input-group-divider grid-cols-[1fr_4em_4em_6em]">
		<input
			type="text"
			bind:value={entry.name}
			name="name"
			placeholder="What to work on..."
			class="input px-2 focus:outline-none"
		/>
		<input
			type="text"
			value={start}
			name="start"
			placeholder="..."
			class="input text-center focus:outline-none"
		/>
		<input
			type="text"
			value={end}
			name="end"
			placeholder="..."
			class="input text-center focus:outline-none"
		/>
		<div class="input text-center focus:outline-none text-secondary-500 align-middle">
			{duration}
		</div>
	</div>
</form>
