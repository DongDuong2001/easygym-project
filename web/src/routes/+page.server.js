import { api } from '#lib/server/api.js';

export async function load({ fetch }) {
	const [posts, media] = await Promise.all([api(fetch, '/api/posts'), api(fetch, '/api/media')]);
	return { posts, media };
}
