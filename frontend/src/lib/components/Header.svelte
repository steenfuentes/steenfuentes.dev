<script>
	import { page } from '$app/stores';

	let currentPath = $derived($page.url.pathname);
</script>

<header class="header">
	<div class="container">
		<nav class="nav">
			<a href="/" class="logo">
				<span class="name">Seth</span>
				<span class="domain">steenfuentes.dev</span>
			</a>

			<div class="nav-links">
				<a href="/" class:active={currentPath === '/'}>Home</a>
				<a href="/blog" class:active={currentPath.startsWith('/blog')}>Writing</a>
				<a href="/about" class:active={currentPath === '/about'}>About</a>
			</div>
		</nav>
	</div>
</header>

<style>
	.header {
		padding: 2rem 0;
		margin-bottom: 3rem;
		border-bottom: 1px solid var(--color-border-light);
		background: var(--color-surface);
	}

	.nav {
		display: flex;
		flex-wrap: wrap;
		gap: 1rem;
		justify-content: space-between;
		align-items: center;
	}

	.logo {
		display: flex;
		align-items: baseline;
		gap: 0.5rem;
		text-decoration: none;
	}

	.logo:hover {
		text-decoration: none;
	}

	.logo:hover .name {
		color: var(--color-accent-hover);
	}

	.name {
		font-size: 1.25rem;
		font-weight: 600;
		color: var(--color-accent);
		transition: color 0.2s ease;
	}

	.domain {
		font-size: 1rem;
		color: var(--color-text-light);
		font-weight: 400;
	}

	.nav-links {
		display: flex;
		gap: 0.5rem;
		align-items: center;
	}

	.nav-links a {
		padding: 0.25rem 0.75rem;
		border-radius: 2rem;
		font-size: 0.95rem;
		color: var(--color-text-muted);
		transition: color 0.2s ease, background-color 0.2s ease;
		text-decoration: none;
	}

	.nav-links a:hover {
		color: var(--color-accent-hover);
		background: var(--color-bg);
	}

	.nav-links a.active {
		color: var(--color-accent);
		background: var(--color-bg-subtle);
		font-weight: 600;
	}

	@media (max-width: 768px) {
		.nav-links {
			gap: 0.25rem;
		}

		.nav-links a {
			font-size: 0.9rem;
		}
	}
</style>
