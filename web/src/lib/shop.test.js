import { test } from 'node:test';
import assert from 'node:assert/strict';
import { filterProducts, fold, formatVND, fromPrice, isPhone, priceCart, suggest } from './shop.js';
// The same data the backend seeds into the database.
import seed from '../../../backend/db/seed/seed.json' with { type: 'json' };

const CATEGORIES = seed.categories.map((c) => ({ value: c.slug, label: c.label }));
const { courses, products } = seed;

const cat = [
	{ slug: 'a', name: 'Bánh A', tag: '', category: 'x', reviews: 1, variants: [{ label: 'x', price: 100000 }, { label: 'y', price: 400000 }] },
	{ slug: 'b', name: 'Đồ B', tag: '', category: 'y', reviews: 5, variants: [{ label: 'z', price: 50000 }] }
];

test('fold strips Vietnamese accents', () => {
	assert.equal(fold('Bánh Bông Đường'), 'banh bong duong');
});

test('formatVND uses dots and the đ sign', () => {
	assert.equal(formatVND(8700000), '8.700.000đ');
});

test('search ignores accents and word order', () => {
	assert.deepEqual(filterProducts(products, { q: 'matcha banh' }).map((p) => p.slug), ['bong-lean-matcha']);
	assert.equal(filterProducts(products, { q: 'khong co mon nay' }).length, 0);
});

test('category filter and sorting', () => {
	assert.ok(filterProducts(products, { dm: 'banh' }).every((p) => p.category === 'banh'));
	assert.deepEqual(filterProducts(cat, {}).map((p) => p.slug), ['b', 'a']);
	assert.deepEqual(filterProducts(cat, { 'sap-xep': 'gia-giam' }).map((p) => p.slug), ['a', 'b']);
	assert.deepEqual(filterProducts(cat, { 'sap-xep': 'gia-tang' }).map((p) => p.slug), ['b', 'a']);
});

test('filtering never reorders the catalogue itself', () => {
	const before = products.map((p) => p.slug);
	filterProducts(products, { 'sap-xep': 'gia-giam' });
	assert.deepEqual(products.map((p) => p.slug), before);
});

test('suggest puts names starting with the word first', () => {
	assert.equal(suggest(products, 'kim')[0].slug, 'kim-chi-organic');
	assert.deepEqual(suggest(products, ''), []);
});

test('priceCart uses the chosen variant and drops stale lines', () => {
	const c = priceCart(
		[
			{ slug: 'a', variant: 1, qty: 2 },
			{ slug: 'b', variant: 0, qty: 1 },
			{ slug: 'gone', variant: 0, qty: 9 },
			{ slug: 'a', variant: 7, qty: 1 }
		],
		cat
	);
	assert.equal(c.items.length, 2);
	assert.equal(c.count, 3);
	assert.equal(c.total, 850000);
});

test('isPhone accepts Vietnamese numbers only', () => {
	assert.ok(isPhone('0905152949'));
	assert.ok(isPhone('0905 152 949'));
	assert.ok(!isPhone('905152949'));
	assert.ok(!isPhone('09051529499'));
	assert.ok(!isPhone('abc'));
});

test('catalogue is consistent', () => {
	assert.equal(new Set(products.map((p) => p.slug)).size, products.length);
	assert.equal(new Set(courses.map((c) => c.slug)).size, courses.length);
	for (const p of products) {
		assert.ok(CATEGORIES.some((c) => c.value === p.category), p.slug);
		assert.equal(p.variants[0].price, fromPrice(p.variants), `${p.slug} lists its cheapest variant first`);
	}
	for (const c of courses) assert.ok(c.plans.length > 0, c.slug);
});
