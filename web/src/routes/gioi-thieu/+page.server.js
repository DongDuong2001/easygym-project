import { api } from '#lib/server/api.js';

export async function load({ fetch }) {
	return { about: (await api(fetch, '/api/pages/gioi-thieu')).body };
}
