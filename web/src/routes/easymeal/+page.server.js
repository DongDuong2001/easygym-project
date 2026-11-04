import { api } from '#lib/server/api.js';

export async function load({ fetch }) {
	return { posts: await api(fetch, '/api/posts?cat=easymeal') };
}
