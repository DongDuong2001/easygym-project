import { json } from '@sveltejs/kit';
import { apiBase } from '#lib/server/api.js';

/**
 * The browser posts sign-ups and orders to /api/* on this origin; forward them to the Go API,
 * so the API stays private and there is no CORS to configure.
 */
export async function handle({ event, resolve }) {
	if (!event.url.pathname.startsWith('/api/')) return resolve(event);

	const { request } = event;
	const hasBody = !['GET', 'HEAD'].includes(request.method);
	try {
		const res = await fetch(apiBase() + event.url.pathname + event.url.search, {
			method: request.method,
			headers: {
				'content-type': request.headers.get('content-type') ?? '',
				authorization: request.headers.get('authorization') ?? '',
				'x-forwarded-for': event.getClientAddress()
			},
			body: hasBody ? await request.arrayBuffer() : undefined
		});
		return new Response(res.body, { status: res.status, headers: { 'content-type': res.headers.get('content-type') ?? 'application/json' } });
	} catch {
		return json({ error: 'Hệ thống đang bảo trì, vui lòng thử lại sau ít phút' }, { status: 503 });
	}
}
