<script>
	import { createEventDispatcher } from 'svelte';
	
	export let post = {
		title: '',
		slug: '',
		content: '',
		excerpt: '',
		tags: [],
		published: false
	};
	
	const dispatch = createEventDispatcher();
	
	let formData = { ...post };
	let tagInput = Array.isArray(post.tags) ? post.tags.join(', ') : '';
	
	function generateSlug() {
		if (!formData.slug && formData.title) {
			formData.slug = formData.title
				.toLowerCase()
				.replace(/[^a-z0-9]+/g, '-')
				.replace(/^-+|-+$/g, '');
		}
	}
	
	function handleSubmit() {
		const tags = tagInput
			.split(',')
			.map(tag => tag.trim())
			.filter(tag => tag.length > 0);
		
		dispatch('save', {
			...formData,
			tags
		});
	}
	
	function insertMarkdown(before, after = '') {
		const textarea = document.getElementById('content');
		const start = textarea.selectionStart;
		const end = textarea.selectionEnd;
		const selectedText = formData.content.substring(start, end);
		
		const newText = before + selectedText + after;
		formData.content = 
			formData.content.substring(0, start) + 
			newText + 
			formData.content.substring(end);
		
		setTimeout(() => {
			textarea.focus();
			textarea.setSelectionRange(
				start + before.length,
				start + before.length + selectedText.length
			);
		}, 0);
	}
</script>

<form on:submit|preventDefault={handleSubmit} class="post-form">
	<div class="form-group">
		<label for="title">Title</label>
		<input
			id="title"
			type="text"
			bind:value={formData.title}
			on:input={generateSlug}
			required
		/>
	</div>
	
	<div class="form-group">
		<label for="slug">Slug</label>
		<input
			id="slug"
			type="text"
			bind:value={formData.slug}
			required
		/>
	</div>
	
	<div class="form-group">
		<label for="excerpt">Excerpt</label>
		<textarea
			id="excerpt"
			bind:value={formData.excerpt}
			rows="3"
			placeholder="Brief description of the post..."
		></textarea>
	</div>
	
	<div class="form-group">
		<label for="content">Content</label>
		<div class="editor-toolbar">
			<button type="button" on:click={() => insertMarkdown('**', '**')}>Bold</button>
			<button type="button" on:click={() => insertMarkdown('*', '*')}>Italic</button>
			<button type="button" on:click={() => insertMarkdown('## ')}>Heading</button>
			<button type="button" on:click={() => insertMarkdown('[', '](url)')}>Link</button>
			<button type="button" on:click={() => insertMarkdown('![alt](', ')')}>Image</button>
			<button type="button" on:click={() => insertMarkdown('```\n', '\n```')}>Code</button>
			<button type="button" on:click={() => insertMarkdown('> ')}>Quote</button>
			<button type="button" on:click={() => insertMarkdown('- ')}>List</button>
		</div>
		<textarea
			id="content"
			bind:value={formData.content}
			rows="20"
			placeholder="Write your post in Markdown..."
			required
		></textarea>
	</div>
	
	<div class="form-group">
		<label for="tags">Tags (comma-separated)</label>
		<input
			id="tags"
			type="text"
			bind:value={tagInput}
			placeholder="technology, design, tutorial"
		/>
	</div>
	
	<div class="form-group checkbox">
		<label>
			<input
				type="checkbox"
				bind:checked={formData.published}
			/>
			Publish immediately
		</label>
	</div>
	
	<div class="form-actions">
		<button type="submit" class="btn-primary">Save Post</button>
		<a href="/admin" class="btn-secondary">Cancel</a>
	</div>
</form>

<style>
	.post-form {
		display: flex;
		flex-direction: column;
		gap: 1.5rem;
	}
	
	.form-group {
		display: flex;
		flex-direction: column;
		gap: 0.5rem;
	}
	
	label {
		font-weight: 600;
		color: var(--color-text);
	}
	
	input[type="text"],
	textarea {
		padding: 0.75rem;
		border: 1px solid var(--color-border);
		border-radius: 4px;
		font-family: inherit;
		font-size: 1rem;
		background-color: var(--color-surface);
		color: var(--color-text);
		transition: border-color 0.2s ease;
	}
	
	input[type="text"]:focus,
	textarea:focus {
		border-color: var(--color-accent);
	}
	
	.editor-toolbar {
		display: flex;
		gap: 0.5rem;
		padding: 0.5rem;
		background-color: var(--color-bg-subtle);
		border: 1px solid var(--color-border);
		border-bottom: none;
		border-radius: 4px 4px 0 0;
	}
	
	.editor-toolbar button {
		padding: 0.5rem 0.75rem;
		background-color: var(--color-surface);
		color: var(--color-text);
		border: 1px solid var(--color-border);
		border-radius: 4px;
		cursor: pointer;
		font-size: 0.875rem;
		transition: background-color 0.2s ease;
	}
	
	.editor-toolbar button:hover {
		background-color: var(--color-bg);
	}
	
	#content {
		border-radius: 0 0 4px 4px;
		font-family: 'Courier New', monospace;
	}
	
	.checkbox label {
		display: flex;
		align-items: center;
		gap: 0.5rem;
	}
	
	input[type="checkbox"] {
		accent-color: var(--color-accent);
		width: 1.25rem;
		height: 1.25rem;
	}
	
	.form-actions {
		display: flex;
		gap: 1rem;
		margin-top: 1rem;
	}
	
	.btn-primary,
	.btn-secondary {
		padding: 0.75rem 1.5rem;
		border-radius: 4px;
		font-weight: 500;
		transition: background-color 0.2s ease;
		border: none;
		cursor: pointer;
		font-size: 1rem;
	}
	
	.btn-primary {
		background-color: var(--color-accent);
		color: var(--color-surface);
	}
	
	.btn-secondary {
		background-color: var(--color-bg-subtle);
		color: var(--color-text);
	}
	
	.btn-primary:hover {
		background-color: var(--color-accent-hover);
	}

	.btn-secondary:hover {
		background-color: var(--color-surface);
		color: var(--color-accent-hover);
	}
</style>
