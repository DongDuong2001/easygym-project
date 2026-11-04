import { api } from '#lib/server/api.js';

export async function load({ fetch }) {
	return { media: await api(fetch, '/api/media') };
}
