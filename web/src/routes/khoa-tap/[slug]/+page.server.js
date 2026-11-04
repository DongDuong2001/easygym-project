import { api } from '#lib/server/api.js';

export async function load({ fetch, params }) {
	const course = await api(fetch, `/api/courses/${encodeURIComponent(params.slug)}`);
	return { course, body: course.body };
}
