// Pure helpers for the catalogue and cart. No DOM, so they run in Node tests.

export const formatVND = (n) => n.toLocaleString('vi-VN') + 'đ';

/** Lowercase and strip Vietnamese accents so "banh bong" matches "Bánh Bông". */
export const fold = (s) =>
	s.normalize('NFD').replace(/[̀-ͯ]/g, '').replace(/đ/gi, 'd').toLowerCase().trim();

/** Cheapest option. Works for product variants and course plans. */
export const fromPrice = (list) => Math.min(...list.map((x) => x.price));

export const SORTS = [
	{ value: '', label: 'Phổ biến' },
	{ value: 'gia-tang', label: 'Giá thấp đến cao' },
	{ value: 'gia-giam', label: 'Giá cao đến thấp' }
];

/** The filter keys live in the URL so a search can be shared or bookmarked. */
export const FILTER_KEYS = ['q', 'dm', 'sap-xep'];

export function readFilters(params) {
	return Object.fromEntries(FILTER_KEYS.map((k) => [k, params.get(k) ?? '']));
}

export function filterProducts(list, f) {
	const words = fold(f.q ?? '').split(/\s+/).filter(Boolean);
	const out = list.filter((p) => words.every((w) => fold(`${p.name} ${p.tag}`).includes(w)) && (!f.dm || p.category === f.dm));
	const by = {
		'gia-tang': (a, b) => fromPrice(a.variants) - fromPrice(b.variants),
		'gia-giam': (a, b) => fromPrice(b.variants) - fromPrice(a.variants)
	}[f['sap-xep']] ?? ((a, b) => b.reviews - a.reviews);
	return out.sort(by);
}

/** Header search box: names that start with the first word come first. */
export function suggest(list, q, limit = 5) {
	const first = fold(q).split(/\s+/)[0];
	if (!first) return [];
	const rank = (p) => (fold(p.name).startsWith(first) ? 0 : 1);
	return filterProducts(list, { q })
		.sort((a, b) => rank(a) - rank(b))
		.slice(0, limit);
}

/** Lines are {slug, variant, qty}. Lines whose product or variant no longer exists are dropped. */
export function priceCart(lines, catalog) {
	const items = lines
		.map((l) => {
			const p = catalog.find((x) => x.slug === l.slug);
			const v = p?.variants[l.variant];
			return v && { ...l, product: p, label: v.label, price: v.price, total: v.price * l.qty };
		})
		.filter(Boolean);
	return {
		items,
		count: items.reduce((s, i) => s + i.qty, 0),
		total: items.reduce((s, i) => s + i.total, 0)
	};
}

/** Vietnamese mobile or landline: 10 digits starting with 0, spaces and dots allowed. */
export const isPhone = (s) => /^0\d{9}$/.test(s.replace(/[\s.-]/g, ''));
