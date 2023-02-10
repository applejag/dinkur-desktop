<script>
	import '@skeletonlabs/skeleton/themes/theme-gold-nouveau.css';
	import '@skeletonlabs/skeleton/styles/all.css';
	import '../app.postcss';
	import './styles.css';
	import { AppShell, AppBar } from '@skeletonlabs/skeleton';
</script>

<AppShell>
	<svelte:fragment slot="header">
		<AppBar>Dinkur</AppBar>
	</svelte:fragment>

	<slot/>
</AppShell>