<script>
	import { onMount } from 'svelte';

	let mounted = $state(false);

	onMount(() => {
		mounted = true;
	});
</script>

<svelte:head>
	<title>Seth Steen-Fuentes</title>
	<meta name="description" content="Personal website of Seth Steen-Fuentes" />
</svelte:head>

<section class="intro">
	<h1>Hi, I'm Seth.</h1>

	<p class="bio">
		I'm a software engineer working in the biopharma industry.
		<br>
		<br>
		This site is currently a work in progress, but will soon be the home of a blog
		and portfolio of my work.
	</p>

	<div class="links">
		<a href="/about">About me</a>
		<span class="separator">|</span>
		<a href="mailto:seth@steenfuentes.dev">Email</a>
		<span class="separator">|</span>
		<a href="https://github.com/steenfuentes" target="_blank" rel="noopener">GitHub</a>
	</div>
</section>

<section class="recent-posts">
	<h2>Recent Writing</h2>
	<p class="coming-soon">Coming soon...</p>
</section>

<style>
	.intro {
		margin: 4rem 0 4rem;
		padding-bottom: 2rem;
		border-bottom: 1px solid var(--color-border-light);
	}

	h1 {
		font-size: 2rem;
		margin-bottom: 1.5rem;
		font-weight: 600;
		color: var(--color-text);
	}

	.bio {
		font-size: 1.125rem;
		line-height: 1.7;
		color: var(--color-text-muted);
		margin-bottom: 2rem;
		max-width: 600px;
	}

	.links {
		display: flex;
		align-items: center;
		gap: 0.75rem;
		font-size: 1rem;
	}

	.links a {
		text-decoration: underline;
	}

	.separator {
		color: var(--color-border);
		font-size: 0.875rem;
	}

	.recent-posts {
		margin: 4rem 0;
	}

	.recent-posts h2 {
		font-size: 1.5rem;
		margin-bottom: 1rem;
		font-weight: 600;
		color: var(--color-text);
	}

	.coming-soon {
		color: var(--color-text-light);
		font-style: italic;
	}

	@media (max-width: 768px) {
		.intro {
			margin: 2rem 0 4rem;
		}

		h1 {
			font-size: 1.75rem;
		}

		.bio {
			font-size: 1rem;
		}

		.links {
			font-size: 0.95rem;
		}
	}
</style>
