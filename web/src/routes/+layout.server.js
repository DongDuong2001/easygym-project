import { api } from '#lib/server/api.js';

// Product and course cards are needed on almost every page (header search, footer, cart), so load them once here.
export async function load({ fetch }) {
	return { catalog: await api(fetch, '/api/catalog') };
}
