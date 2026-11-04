import { error } from '@sveltejs/kit';
import { env } from '$env/dynamic/private';

export const apiBase = () => env.API_URL || 'http://localhost:8080';

/** GET from the Go API while rendering on the server. A 404 there becomes this site's 404 page. */
export async function api(fetch, path) {
	let res;
	try {
		res = await fetch(apiBase() + path);
	} catch {
		error(503, 'Hệ thống đang bảo trì, vui lòng thử lại sau ít phút');
	}
	if (res.status === 404) error(404, 'Không tìm thấy trang');
	if (!res.ok) error(502, 'Hệ thống đang bận, vui lòng thử lại sau');
	return res.json();
}
