<script lang="ts">
	import '../app.css';
	import Header from '$lib/components/Header.svelte';
	interface Props {
		children?: import('svelte').Snippet;
	}

	let { children }: Props = $props();
</script>

<Header />

<main class="container">
	{@render children?.()}
</main>

<footer class="footer">
	<div class="container">
		<div class="footer-content">
			<p>&copy; {new Date().getUTCFullYear()} Seth Steen-Fuentes.</p>
			<div class="social-links">
				<a href="mailto:hello@steenfuentes.dev">Email</a>
				<a href="https://github.com/steenfuentes">GitHub</a>
				<a href="https://twitter.com/steenfuentes">Twitter</a>
			</div>
		</div>
	</div>
</footer>

<style>
	main {
		width: 100%;
		flex: 1;
	}

	.footer {
		margin-top: 6rem;
		padding: 3rem 0;
		border-top: 1px solid var(--color-border-light);
		background-color: var(--color-surface);
	}
	
	.footer-content {
		display: flex;
		justify-content: space-between;
		align-items: center;
	}
	
	.footer-content p {
		margin: 0;
		color: var(--color-text-muted);
		font-size: 0.9rem;
	}
	
	.social-links {
		display: flex;
		gap: 2rem;
	}
	
	.social-links a {
		color: var(--color-accent);
		font-size: 0.9rem;
		transition: color 0.2s ease;
	}
	
	.social-links a:hover {
		color: var(--color-accent-hover);
	}
	
	@media (max-width: 768px) {
		.footer-content {
			flex-direction: column;
			gap: 1.5rem;
			text-align: center;
		}
	}
</style>
