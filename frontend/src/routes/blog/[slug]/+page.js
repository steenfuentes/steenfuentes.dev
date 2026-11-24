import { error } from '@sveltejs/kit';

export async function load({ params }) {
	try {
		// Static sample posts data
		const posts = {
			'welcome-to-my-blog': {
				title: 'Welcome to My Blog',
				created_at: '2024-01-01T00:00:00Z',
				reading_time: 5,
				content: '# Welcome to My Blog\n\nThis is a sample blog post to demonstrate the blog functionality.\n\n## About This Blog\n\nThis blog will cover various topics related to web development, technology, and more.\n\n### What to Expect\n\n- Technical tutorials\n- Industry insights\n- Personal projects\n- And more!\n\nStay tuned for more content!',
				slug: 'welcome-to-my-blog'
			},
			'my-second-post': {
				title: 'My Second Post',
				created_at: '2024-01-02T00:00:00Z',
				reading_time: 3,
				content: '# My Second Post\n\nThis is another sample post to show multiple posts in the blog.\n\n## Content Goes Here\n\nLorem ipsum dolor sit amet, consectetur adipiscing elit.',
				slug: 'my-second-post'
			}
		};
		
		const post = posts[params.slug];
		
		if (!post) {
			throw error(404, 'Post not found');
		}

		return {
			post: {
				title: post.title,
				date: new Date(post.created_at).toLocaleDateString('en-US', { 
					year: 'numeric', 
					month: 'long', 
					day: 'numeric' 
				}),
				readTime: `${post.reading_time} min read`,
				content: post.content,
				slug: post.slug
			}
		};
	} catch (err) {
		if (err.status === 404) {
			throw err;
		}
		console.error('Failed to load post:', err);
		throw error(500, 'Failed to load post');
	}
}