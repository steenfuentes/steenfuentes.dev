<script lang="ts">
	interface Props {
		title: any;
		excerpt: any;
		date: any;
		slug: any;
		readTime?: string;
	}

	let {
		title,
		excerpt,
		date,
		slug,
		readTime = '5 min read'
	}: Props = $props();
</script>

<article class="blog-card">
	<a href="/blog/{slug}" class="card-link">
		<div class="meta">
			<time class="date">{date}</time>
			<span class="read-time">{readTime}</span>
		</div>
		
		<h2 class="title">{title}</h2>
		<p class="excerpt">{excerpt}</p>
		
		<span class="read-more">Read more →</span>
	</a>
</article>

<style>
	.blog-card {
		margin-bottom: 3rem;
		padding-bottom: 3rem;
		border-bottom: 1px solid var(--color-border-light);
	}
	
	.blog-card:last-child {
		border-bottom: none;
		margin-bottom: 0;
		padding-bottom: 0;
	}
	
	.card-link {
		display: block;
		border: none;
		transition: transform 0.2s ease;
	}
	
	.card-link:hover {
		transform: translateX(0.5rem);
	}
	
	.meta {
		display: flex;
		justify-content: space-between;
		align-items: center;
		margin-bottom: 1rem;
	}
	
	.date {
		font-size: 0.9rem;
		color: var(--color-text-muted);
		font-weight: 300;
	}
	
	.read-time {
		font-size: 0.8rem;
		color: var(--color-text-light);
		text-transform: uppercase;
		letter-spacing: 0.5px;
	}
	
	.title {
		font-size: 1.75rem;
		margin: 0 0 1rem 0;
		color: var(--color-text);
		line-height: 1.3;
	}
	
	.excerpt {
		color: var(--color-text-muted);
		font-size: 1rem;
		line-height: 1.7;
		margin: 0 0 1.5rem 0;
	}
	
	.read-more {
		color: var(--color-accent);
		font-size: 0.9rem;
		font-weight: 500;
	}
	
	@media (max-width: 768px) {
		.meta {
			flex-direction: column;
			align-items: flex-start;
			gap: 0.25rem;
		}
		
		.title {
			font-size: 1.5rem;
		}
	}
</style>
