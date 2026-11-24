<svelte:head>
	<title>About - Seth Steen-Fuentes</title>
	<meta name="description" content="About Seth Steen-Fuentes" />
</svelte:head>

<section class="about">
	<h1>About</h1>

	<div class="bio">
		<p>
			I'm a software engineer based in southern California. I enjoy working on challenging problems
			and building tools that make people's lives easier.


		</p>

		<p>
			I work in the biopharma industry and am part of a small team in a fast-paced
			org, where forays into a multitude of domains including data engineering, web development, devops, AI, etc., is the norm.
		</p>

		<p>
			Hailing from a background in the physical and life sciences, I have a special interest in building systems that propagate productivity and knowledge.
		</p>

		<p>
			Outside of programming, I enjoy gaming, reading, and playing music.
		</p>

		<h2>Contact</h2>
		<p>
			You can reach me at <a href="mailto:seth@steenfuentes.dev">seth@steenfuentes.dev</a> or
			find me on <a href="https://github.com/steenfuentes" target="_blank" rel="noopener">GitHub</a>.
		</p>
	</div>
</section>

<style>
	.about {
		max-width: 650px;
		margin: 0 auto;
	}

	h1 {
		font-size: 2rem;
		margin: 0 0 2rem 0;
		font-weight: 600;
		color: var(--color-text);
		padding-bottom: 2rem;
		border-bottom: 1px solid var(--color-border-light);
	}

	h2 {
		font-size: 1.5rem;
		margin: 3rem 0 1rem 0;
		font-weight: 600;
		color: var(--color-text);
	}

	.bio p {
		font-size: 1.125rem;
		line-height: 1.7;
		margin: 0 0 1.5rem 0;
		color: var(--color-text-muted);
	}

	.bio a {
		text-decoration: underline;
	}

	@media (max-width: 768px) {
		h1 {
			font-size: 1.75rem;
		}

		.bio p {
			font-size: 1rem;
		}
	}
</style>
