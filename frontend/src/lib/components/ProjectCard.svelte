<script lang="ts">
	interface Props {
		title: any;
		description: any;
		year: any;
		tags?: any;
		image?: string;
	}

	let {
		title,
		description,
		year,
		tags = [],
		image = '/placeholder.svg?height=300&width=400'
	}: Props = $props();
</script>

<article class="project-card">
	<div class="image-container">
		<img src={image || "/placeholder.svg"} alt={title} />
	</div>
	
	<div class="content">
		<div class="meta">
			<span class="year">{year}</span>
			<div class="tags">
				{#each tags as tag}
					<span class="tag">{tag}</span>
				{/each}
			</div>
		</div>
		
		<h3 class="title">{title}</h3>
		<p class="description">{description}</p>
	</div>
</article>

<style>
	.project-card {
		display: flex;
		gap: 2rem;
		margin-bottom: 4rem;
		padding-bottom: 4rem;
		border-bottom: 1px solid var(--color-border-light);
	}
	
	.project-card:last-child {
		border-bottom: none;
		margin-bottom: 0;
		padding-bottom: 0;
	}
	
	.image-container {
		flex: 0 0 300px;
		height: 200px;
		overflow: hidden;
		border-radius: 2px;
	}
	
	.image-container img {
		width: 100%;
		height: 100%;
		object-fit: cover;
		transition: transform 0.3s ease;
	}
	
	.project-card:hover .image-container img {
		transform: scale(1.02);
	}
	
	.content {
		flex: 1;
		padding-top: 0.5rem;
	}
	
	.meta {
		display: flex;
		justify-content: space-between;
		align-items: center;
		margin-bottom: 1rem;
	}
	
	.year {
		font-size: 0.9rem;
		color: var(--color-text-muted);
		font-weight: 300;
	}
	
	.tags {
		display: flex;
		gap: 0.75rem;
	}
	
	.tag {
		font-size: 0.8rem;
		color: var(--color-text-light);
		font-weight: 400;
		text-transform: uppercase;
		letter-spacing: 0.5px;
	}
	
	.title {
		font-size: 1.5rem;
		margin: 0 0 1rem 0;
		color: var(--color-text);
	}
	
	.description {
		color: var(--color-text-muted);
		font-size: 1rem;
		line-height: 1.7;
		margin: 0;
	}
	
	@media (max-width: 768px) {
		.project-card {
			flex-direction: column;
			gap: 1.5rem;
		}
		
		.image-container {
			flex: none;
			height: 180px;
		}
		
		.meta {
			flex-direction: column;
			align-items: flex-start;
			gap: 0.5rem;
		}
	}
</style>
