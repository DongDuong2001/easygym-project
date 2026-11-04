import { api } from '#lib/server/api.js';

export async function load({ fetch }) {
	return { groups: await api(fetch, '/api/faq') };
}
