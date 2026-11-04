import { api } from '#lib/server/api.js';

export async function load({ fetch, params }) {
	const [post, posts] = await Promise.all([api(fetch, `/api/posts/${encodeURIComponent(params.slug)}`), api(fetch, '/api/posts')]);
	return { post, body: post.body, posts };
}
