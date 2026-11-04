import { api } from '#lib/server/api.js';

export async function load({ fetch, params }) {
	const product = await api(fetch, `/api/products/${encodeURIComponent(params.slug)}`);
	return { product, info: product };
}
