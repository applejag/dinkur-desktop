<script>
	import '@skeletonlabs/skeleton/themes/theme-crimson.css';
	import '@skeletonlabs/skeleton/styles/all.css';
	import '../app.postcss';
	import './styles.css';
	import { AppShell, AppBar } from '@skeletonlabs/skeleton';
	import logo from '../assets/images/dinkur-large-512.svg';
	import * as Icon from 'svelte-awesome-icons';
</script>

<AppShell>
	<svelte:fragment slot="header">
		<AppBar>
			<div class="grid grid-cols-[auto_1fr] gap-4">
				<img alt="Dinkur logo" id="logo" src={logo} class="h-6" />
				<span>Dinkur - <em>Time tracking</em></span>
			</div>
		</AppBar>
	</svelte:fragment>

	<slot/>

	<svelte:fragment slot="footer">
		<footer class="p-4">
			<div class="card p-4">
				<form on:submit|preventDefault={void(0)}>
					<label for="entryName" class="label">
						<span>Current entry</span>
						<div class="input-group input-group-divider grid-cols-[1fr_auto]">
							<input id="entryName" type="text" />
							<button type="submit" class="variant-filled-primary"><Icon.PlaySolid/></button>
						</div>
					</label>
				</form>
			</div>
		</footer>
	</svelte:fragment>
</AppShell>