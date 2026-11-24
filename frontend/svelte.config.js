import adapter from '@sveltejs/adapter-node';
import { vitePreprocess } from '@sveltejs/vite-plugin-svelte';

/** @type {import('@sveltejs/kit').Config} */
const config = {
	// Consult https://svelte.dev/docs/kit/integrations
	// for more information about preprocessors
	preprocess: vitePreprocess(),

	kit: {
		// adapter-node for deployment to Cloud Run and other Node.js environments
		adapter: adapter({
			// default build output directory
			out: 'build',
			// precompress assets for better performance
			precompress: true,
			// environment variables prefix
			envPrefix: 'VITE_'
		})
	}
};

export default config;
