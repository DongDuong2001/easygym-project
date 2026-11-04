// Cart and wishlist live in this browser's localStorage; orders and sign-ups go to the API.
const KEY = 'easygym:v1';
const SAVED = ['cart', 'wishes'];

export const store = $state({
	cart: [],
	wishes: [],
	ready: false
});

/** Call once in the browser after hydration (onMount, not $effect: it reads and writes the store). */
export function loadStore() {
	try {
		Object.assign(store, JSON.parse(localStorage.getItem(KEY) ?? '{}'));
	} catch {
		// corrupt or blocked storage: start empty
	}
	store.ready = true;
}

function save() {
	try {
		localStorage.setItem(KEY, JSON.stringify(Object.fromEntries(SAVED.map((k) => [k, store[k]]))));
	} catch {
		// private mode or the quota is full: keep the in-memory state for this visit
	}
}

/** Course sign-up popup: `opened` counts requests so the layout opens it even for the same course twice. */
export const ui = $state({ opened: 0, course: '' });

export function openRegister(course = '') {
	ui.course = course;
	ui.opened += 1;
}

/** Last message; the layout shows it for a few seconds. */
export const toast = $state({ text: '', at: 0, cart: false });

export function say(text, cart = false) {
	toast.text = text;
	toast.cart = cart;
	toast.at = Date.now();
}

export function addToCart(slug, variant = 0, qty = 1) {
	const line = store.cart.find((l) => l.slug === slug && l.variant === variant);
	if (line) line.qty = Math.min(line.qty + qty, 99);
	else store.cart.push({ slug, variant, qty });
	save();
	say('Đã thêm vào giỏ hàng', true);
}

export function setQty(slug, variant, qty) {
	const same = (l) => l.slug === slug && l.variant === variant;
	store.cart = qty < 1 ? store.cart.filter((l) => !same(l)) : store.cart.map((l) => (same(l) ? { ...l, qty: Math.min(qty, 99) } : l));
	save();
}

export function toggleWish(slug) {
	const on = !store.wishes.includes(slug);
	store.wishes = on ? [...store.wishes, slug] : store.wishes.filter((s) => s !== slug);
	save();
	say(on ? 'Đã thêm vào yêu thích' : 'Đã bỏ khỏi yêu thích');
}

/** POSTs JSON to the API (proxied by hooks.server.js). Throws an Error whose message can be shown to the visitor. */
async function post(url, body) {
	let res;
	try {
		res = await fetch(url, { method: 'POST', headers: { 'content-type': 'application/json' }, body: JSON.stringify(body) });
	} catch {
		throw new Error('Không gửi được, vui lòng kiểm tra kết nối mạng và thử lại');
	}
	const data = await res.json().catch(() => ({}));
	if (!res.ok) throw new Error(data.error ?? 'Có lỗi xảy ra, vui lòng thử lại');
	return data;
}

/** Sends the order; the server re-prices every line. Empties the cart and returns { code, total }. */
export async function placeOrder(order) {
	const res = await post('/api/orders', order);
	store.cart = [];
	save();
	return res;
}

/** Saves a sign-up ({ kind: 'register' | 'newsletter', ... }) and returns its code. */
export async function addLead(lead) {
	return (await post('/api/leads', lead)).code;
}
