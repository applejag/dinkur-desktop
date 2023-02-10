<script lang="ts">
	import { ConnectDinkur, GetActiveEntry } from '$lib/wailsjs/go/main/App';
	import logo from '../assets/images/dinkur-large-512.svg';
	import { Avatar } from '@skeletonlabs/skeleton';

	let resultText: string = 'Please enter your name below 👇';
	let serverAddr: string = 'localhost:59122';

	async function connectToDinkur(): Promise<void> {
		try {
			await ConnectDinkur(serverAddr);
			resultText = `Successfully connected to ${serverAddr}`;
			let entry = await GetActiveEntry();
			if (!entry) {
				resultText = `No active entry.`;
			} else {
				resultText = `Active entry: #${entry.id} '${entry.name}', since: ${entry.start}`;
			}
		} catch (err) {
			resultText = `Error: ${err}`;
		}
	}
</script>

<div class="container mx-auto p-8 space-y-8">
	<h1>Hello Skeleton</h1>
	<p>Lorem ipsum dolor sit amet consectetur adipisicing elit.</p>
	<hr />
	<section class="card p-4">
		<p>Lorem ipsum dolor sit amet consectetur adipisicing elit.</p>
		<Avatar src="https://i.pravatar.cc/" />
	</section>
	<hr />
	<section class="flex space-x-2">
		<a class="btn variant-filled-primary" href="https://kit.svelte.dev/" target="_blank" rel="noreferrer">SvelteKit</a>
		<a class="btn variant-filled-secondary" href="https://tailwindcss.com/" target="_blank" rel="noreferrer">Tailwind</a>
		<a class="btn variant-filled-tertiary" href="https://github.com/" target="_blank" rel="noreferrer">GitHub</a>
	</section>
</div>

<main>
	<img alt="Dinkur logo" id="logo" src={logo} />
	<div class="result" id="result">{resultText}</div>
	<div class="input-box" id="input">
		<form on:submit|preventDefault={connectToDinkur}>
			<label for="serverAddr">Server address:</label>
			<input id="serverAddr" bind:value={serverAddr} class="input" type="text" />
			<button type="submit" class="btn">Connect</button>
		</form>
	</div>
</main>