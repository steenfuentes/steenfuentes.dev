<script>
	export let data;
	$: post = data.post;
</script>

<svelte:head>
	<title>{post.title} - Seth Steen-Fuentes</title>
</svelte:head>

<article class="blog-post">
	<header class="post-header">
		<time class="date">{new Date(post.created_at).toLocaleDateString('en-US', { year: 'numeric', month: 'long', day: 'numeric' })}</time>
		<h1 class="post-title">{post.title}</h1>
	</header>
	
	<div class="post-content">
		{@html post.content}
	</div>
	
	<footer class="post-footer">
		<a href="/blog" class="back-link">← Back to all posts</a>
	</footer>
</article>

<style>
	.blog-post {
		max-width: 650px;
		margin: 0 auto;
	}
	
	.post-header {
		margin-bottom: 3rem;
		padding-bottom: 2rem;
		border-bottom: 1px solid var(--color-border-light);
	}
	
	.date {
		font-size: 0.875rem;
		color: var(--color-text-light);
		display: block;
		margin-bottom: 1rem;
		font-family: var(--font-mono);
	}
	
	.post-title {
		font-size: 2.25rem;
		margin: 0;
		line-height: 1.2;
		font-weight: 600;
		color: var(--color-text);
	}
	
	.post-content {
		font-size: 1.125rem;
		line-height: 1.8;
		color: var(--color-text);
	}
	
	.post-content :global(h2) {
		font-size: 1.5rem;
		margin: 3rem 0 1rem 0;
		font-weight: 600;
	}
	
	.post-content :global(h3) {
		font-size: 1.25rem;
		margin: 2.5rem 0 1rem 0;
		font-weight: 600;
	}
	
	.post-content :global(p) {
		margin: 0 0 1.5rem 0;
	}
	
	.post-content :global(blockquote) {
		margin: 2rem 0;
		padding-left: 1.5rem;
		border-left: 3px solid var(--color-border);
		font-style: italic;
		color: var(--color-text-muted);
	}
	
	.post-content :global(blockquote p) {
		margin: 0;
	}
	
	.post-content :global(ul),
	.post-content :global(ol) {
		margin: 1.5rem 0;
		padding-left: 1.5rem;
	}
	
	.post-content :global(li) {
		margin-bottom: 0.5rem;
	}
	
	.post-content :global(code) {
		font-family: var(--font-mono);
		font-size: 0.9em;
		background: var(--color-bg-subtle);
		padding: 0.2em 0.4em;
		border-radius: 3px;
	}
	
	.post-content :global(pre) {
		background: var(--color-bg-subtle);
		padding: 1rem;
		border-radius: 4px;
		overflow-x: auto;
		margin: 1.5rem 0;
	}
	
	.post-content :global(pre code) {
		background: none;
		padding: 0;
		font-size: 0.875rem;
	}
	
	.post-content :global(img) {
		max-width: 100%;
		height: auto;
		margin: 2rem 0;
	}
	
	.post-content :global(a) {
		text-decoration: underline;
	}
	
	.post-footer {
		margin-top: 4rem;
		padding-top: 2rem;
		border-top: 1px solid var(--color-border-light);
	}
	
	.back-link {
		font-size: 0.95rem;
	}
	
	@media (max-width: 768px) {
		.post-title {
			font-size: 1.75rem;
		}
		
		.post-content {
			font-size: 1rem;
		}
		
		.post-content :global(h2) {
			font-size: 1.375rem;
		}
		
		.post-content :global(h3) {
			font-size: 1.125rem;
		}
	}
</style>
