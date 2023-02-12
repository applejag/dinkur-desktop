<script lang="ts">
	import { ConnectDinkur, GetActiveEntry } from '$lib/wailsjs/go/app/App';
	import { Avatar } from '@skeletonlabs/skeleton';

	let resultVariant: string = 'variant-ghost-surface'
	let resultText: string = 'Please connect to Dinkur daemon 👇';
	let serverAddr: string = 'localhost:59122';

	async function connectToDinkur(): Promise<void> {
		try {
			await ConnectDinkur(serverAddr);
			resultText = `Successfully connected to ${serverAddr}`;
			let entry = await GetActiveEntry();
			if (!entry) {
				resultText = `No active entry.`;
				resultVariant = 'variant-ghost-tertiary';
			} else {
				resultText = `Active entry: #${entry.id} '${entry.name}', since: ${entry.start}`;
				resultVariant = 'variant-ghost-success';
			}
		} catch (err) {
			resultText = `Error: ${err}`;
			resultVariant = 'variant-ghost-error';
		}
	}
</script>

<main class="p-4">
	<div class="card p-4 mb-4 {resultVariant}" id="result">{resultText}</div>
	<div class="card p-4">
		<form on:submit|preventDefault={connectToDinkur}>
			<label for="serverAddr" class="label">
				<span>Server address</span>
				<div class="input-group input-group-divider grid-cols-[1fr_auto]">
					<input id="serverAddr" bind:value={serverAddr} type="text" />
					<button type="submit" class="variant-filled-primary">Connect</button>
				</div>
			</label>
		</form>
	</div>
</main>
