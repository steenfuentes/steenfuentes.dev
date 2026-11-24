<script>
	export let data;
	$: posts = data.posts || [];
</script>

<svelte:head>
	<title>Writing - Seth Steen-Fuentes</title>
	<meta name="description" content="Writing about software, technology, and other interests." />
</svelte:head>

<section class="blog-header">
	<h1>Writing</h1>
</section>

<section class="blog-posts">
	{#if posts.length === 0}
		<p class="no-posts">No posts yet. Check back soon.</p>
	{:else}
		{#each posts as post}
			<article class="post-item">
				<a href="/blog/{post.slug}" class="post-link">
					<time class="post-date">{new Date(post.created_at).toLocaleDateString('en-US', { year: 'numeric', month: 'short', day: 'numeric' })}</time>
					<h2 class="post-title">{post.title}</h2>
					{#if post.description}
						<p class="post-description">{post.description}</p>
					{/if}
				</a>
			</article>
		{/each}
	{/if}
</section>

<style>
	.blog-header {
		margin-bottom: 3rem;
		padding-bottom: 2rem;
		border-bottom: 1px solid var(--color-border-light);
	}
	
	.blog-header h1 {
		font-size: 2rem;
		margin: 0;
		font-weight: 600;
		color: var(--color-text);
	}
	
	.blog-posts {
		margin-bottom: 4rem;
	}
	
	.no-posts {
		color: var(--color-text-light);
		font-style: italic;
	}
	
	.post-item {
		margin-bottom: 2rem;
		padding-bottom: 2rem;
		border-bottom: 1px solid var(--color-border-light);
	}
	
	.post-item:last-child {
		border-bottom: none;
		padding-bottom: 0;
	}
	
	.post-link {
		display: block;
		text-decoration: none;
		color: inherit;
	}
	
	.post-link:hover {
		text-decoration: none;
	}
	
	.post-link:hover .post-title {
		color: var(--color-accent);
	}
	
	.post-date {
		font-size: 0.875rem;
		color: var(--color-text-muted);
		display: block;
		margin-bottom: 0.5rem;
		font-family: var(--font-mono);
	}
	
	.post-title {
		font-size: 1.5rem;
		margin: 0 0 0.5rem 0;
		font-weight: 600;
		line-height: 1.3;
		transition: color 0.2s ease;
	}
	
	.post-description {
		color: var(--color-text-muted);
		line-height: 1.6;
		margin: 0;
	}
	
	@media (max-width: 768px) {
		.blog-header h1 {
			font-size: 1.75rem;
		}
		
		.post-title {
			font-size: 1.25rem;
		}
	}
</style>