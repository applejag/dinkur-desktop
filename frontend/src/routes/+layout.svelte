<script>
	import 'carbon-components-svelte/css/g90.css';
	import './styles.css';
	import {
		Header,
		HeaderNav,
		HeaderNavItem,
		HeaderNavMenu,
		SideNav,
		SideNavItems,
		SideNavMenu,
		SideNavMenuItem,
		SideNavLink,
		SideNavDivider,
		SkipToContent,
		Content,
		Grid,
		Row,
		Column,
	} from "carbon-components-svelte";
	import Fade from "carbon-icons-svelte/lib/Fade.svelte";
  
	let isSideNavOpen = false;
</script>

<Header company="Dinkur" platformName="Time tracking" bind:isSideNavOpen>
	<svelte:fragment slot="skip-to-content">
		<SkipToContent />
	</svelte:fragment>
	<HeaderNav>
		<HeaderNavItem href="/" text="Link 1" />
		<HeaderNavItem href="/" text="Link 2" />
		<HeaderNavItem href="/" text="Link 3" />
		<HeaderNavMenu text="Menu">
			<HeaderNavItem href="/" text="Link 1" />
			<HeaderNavItem href="/" text="Link 2" />
			<HeaderNavItem href="/" text="Link 3" />
		</HeaderNavMenu>
		<HeaderNavItem href="/" text="Link 4" />
	</HeaderNav>
</Header>
  
<SideNav bind:isOpen={isSideNavOpen} fixed>
	<SideNavItems>
		<SideNavLink icon={Fade} text="Link 1" href="/" isSelected />
		<SideNavLink icon={Fade} text="Link 2" href="/" />
		<SideNavLink icon={Fade} text="Link 3" href="/" />
		<SideNavMenu icon={Fade} text="Menu">
			<SideNavMenuItem href="/" text="Link 1" />
			<SideNavMenuItem href="/" text="Link 2" />
			<SideNavMenuItem href="/" text="Link 3" />
		</SideNavMenu>
		<SideNavDivider />
		<SideNavLink icon={Fade} text="Link 4" href="/" />
	</SideNavItems>
</SideNav>
  
<Content>
	<Grid>
		<Row>
			<Column>
				<slot/>
			</Column>
		</Row>
	</Grid>
</Content>