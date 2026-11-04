// End-to-end walk through the site plus a layout audit: npm run e2e [base-url].
// It places a real order and real sign-ups, so point it at a local or staging database, never production.
import { chromium } from 'playwright';
import assert from 'node:assert/strict';

const base = process.argv[2] ?? 'http://localhost:3000';
const browser = await chromium.launch();
const ctx = await browser.newContext({ viewport: { width: 1440, height: 900 }, locale: 'vi-VN', reducedMotion: 'reduce' });
const page = await ctx.newPage();
const errors = [];
// YouTube embeds log their own noise; only this site's errors count.
const ours = (u) => !/youtube|googlevideo|doubleclick|google\.com|gstatic|ytimg/.test(u);
page.on('console', (m) => m.type() === 'error' && ours(m.location().url) && !/Failed to load resource/.test(m.text()) && errors.push(m.text()));
page.on('response', (r) => r.status() >= 400 && ours(r.url()) && !/khong-co/.test(r.url()) && errors.push(`${r.status()} ${r.url()}`));
page.on('pageerror', (e) => errors.push(String(e)));
const ok = (name) => console.log('✔', name);
const go = (p) => page.goto(base + p, { waitUntil: 'networkidle' });
const cartCount = () => page.locator('header .count').last().innerText();

// Home: offer and a real action on the first screen
await go('/');
assert.match(await page.locator('h1').innerText(), /tập dễ dàng/i);
assert.equal(await page.locator('h1').count(), 1);
const cta = await page.locator('.hero').getByRole('link', { name: 'Chọn gói tập' }).boundingBox();
assert.ok(cta.y + cta.height < 900, 'CTA above the fold');
await page.getByRole('button', { name: 'Bánh Bông Lean', exact: true }).click();
assert.equal(await page.locator('.grid').first().locator('article').count(), 3);
// Benefits ring: ten evenly spaced nodes that stay inside the page, and a selection that follows the pointer
assert.equal(await page.locator('.ring button').count(), 10);
await page.getByRole('button', { name: /Cải thiện tư thế/ }).click();
assert.equal(await page.locator('.ring button.on').count(), 1);
assert.match(await page.locator('.ring button.on').innerText(), /tư thế/i);
const radii = await page.locator('.ring .node').evaluateAll((nodes) => {
	const r = nodes[0].closest('.ring').getBoundingClientRect();
	return nodes.map((n) => {
		const b = n.getBoundingClientRect();
		return Math.round(Math.hypot(b.left + b.width / 2 - (r.left + r.width / 2), b.top + b.height / 2 - (r.top + r.height / 2)));
	});
});
assert.ok(Math.max(...radii) - Math.min(...radii) <= 6, `nodes on one circle: ${radii}`);
ok('home: hero, tabs, benefits ring');

// Hero is the video: it keeps looping, can be paused, and takes most of the width next to the news rail
const live = await browser.newContext({ viewport: { width: 1440, height: 900 }, locale: 'vi-VN' });
const lp = await live.newPage();
await lp.goto(base + '/', { waitUntil: 'domcontentloaded' });
await lp.waitForFunction(() => {
	const v = document.querySelector('.hero video');
	return v && !v.paused && v.currentTime > 0.3;
}, null, { timeout: 15000 });
await lp.waitForTimeout(8000);
const v = await lp.locator('.hero video').evaluate((e) => ({ muted: e.muted, paused: e.paused, loop: e.loop, ended: e.ended, w: e.getBoundingClientRect().width }));
assert.ok(v.muted && !v.paused && v.loop && !v.ended, JSON.stringify(v));
assert.ok(v.w > 1440 * 0.7, `video is the main element: ${v.w}px`);
await lp.getByRole('button', { name: 'Dừng video' }).click();
assert.ok(await lp.locator('.hero video').evaluate((e) => e.paused));
assert.equal(await lp.locator('.hero aside li').count(), 4);
const banner = await lp.locator('.banner img').evaluate((i) => Math.abs(i.getBoundingClientRect().width / i.getBoundingClientRect().height - i.naturalWidth / i.naturalHeight));
assert.ok(banner < 0.05, 'banner is not cropped');
await live.close();
const frozen = await page.locator('.hero video').evaluate((e) => ({ paused: e.paused, at: e.currentTime }));
assert.ok(frozen.paused && frozen.at > 6, `reduced motion holds the logo frame: ${JSON.stringify(frozen)}`);
ok('hero video loops, pauses, sits beside news');

// Coach achievements, real logos
await go('/');
assert.equal(await page.locator('.awards .titles li').count(), 11);
assert.equal(await page.locator('.awards .medals li').count(), 3);
assert.equal(await page.locator('footer .socials a svg').count(), 5);
ok('coach achievements + social logos');

// Everything brought over from easygym.vn
await go('/tin-tuc');
assert.equal(await page.locator('main article').count(), 22);
await page.getByRole('button', { name: 'easyMEAL', exact: true }).click();
assert.equal(await page.locator('main article').count(), 4);
await page.locator('main article h3 a').first().click();
await page.waitForURL('**/tin-tuc/*');
assert.ok((await page.locator('.prose').innerText()).length > 400, 'article has its body');
await go('/easymeal');
assert.equal(await page.locator('main .posts article').count(), 4);
await go('/faq');
assert.equal(await page.locator('details').count(), 55);
await page.locator('details summary').first().click();
assert.ok(await page.locator('details[open] .prose').isVisible());
await go('/hinh-anh');
assert.equal(await page.locator('.photos button').count(), 8);
await page.getByRole('button', { name: 'Feedback khách hàng' }).click();
assert.equal(await page.locator('.photos button').count(), 23);
await page.locator('.photos button').first().click();
assert.match(await page.locator('dialog p').innerText(), /1 \/ 23/);
await page.keyboard.press('ArrowRight');
assert.match(await page.locator('dialog p').innerText(), /2 \/ 23/);
await page.keyboard.press('Escape');
await go('/video');
assert.equal(await page.locator('.videos iframe').count(), 0, 'no player before a click');
assert.equal(await page.locator('.videos button[aria-label^="Phát video"]').count(), 5);
await page.locator('.videos button').first().click();
assert.equal(await page.locator('.videos iframe').getAttribute('src'), 'https://www.youtube.com/embed/0ylyvZuqgW8?autoplay=1&rel=0');
await page.locator('.videos button').first().click();
assert.equal(await page.locator('.videos iframe').count(), 1, 'one player at a time');
assert.match(await page.locator('.videos iframe').getAttribute('src'), /cN0ODQBfz1Q/);
await go('/san-pham/chocolean');
assert.ok((await page.locator('.below .prose').first().innerText()).length > 2000, 'full product copy');
assert.equal(await page.locator('.thumbs button').count(), 15);
await go('/san-pham/kim-chi-organic');
assert.doesNotMatch(await page.locator('main').innerText(), /ung thư/i);
await go('/san-pham/goat-kefir-a');
const asked = await page.locator('details summary').allInnerTexts();
assert.equal(asked.length, 5);
assert.ok(asked.every((q) => /goat kefir a/i.test(q)), `questions belong to this product: ${asked}`);
await go('/khoa-tap/coach-online');
assert.match(await page.locator('.body').innerText(), /Quy trình/);
await page.locator('header nav').getByRole('link', { name: 'Sản phẩm', exact: true }).hover();
assert.ok(await page.getByRole('link', { name: 'Bánh Bông Lean', exact: true }).first().isVisible(), 'dropdown opens');
ok('news, easyMEAL, FAQ, gallery, video, product and course copy, menu');

// Header search ignores accents and goes to the product
await page.getByPlaceholder('Bạn đang tìm gì?').fill('chocolean');
await page.locator('.hits a').first().click();
await page.waitForURL('**/san-pham/chocolean');
ok('header search suggestion');

// Search by Enter lands on a filtered shop, and that URL holds up on reload
await page.getByPlaceholder('Bạn đang tìm gì?').fill('banh matcha');
await page.keyboard.press('Enter');
await page.waitForURL('**/cua-hang?q=*');
await page.waitForSelector('main article');
assert.equal(await page.locator('main article').count(), 1);
await page.reload({ waitUntil: 'networkidle' });
assert.equal(await page.locator('main article').count(), 1, 'query survives reload');
ok('search results and reload');

// Shop filters and sort
await go('/cua-hang?dm=kefir');
assert.equal(await page.locator('main article').count(), 4);
await page.getByRole('button', { name: 'Tất cả', exact: true }).click();
assert.equal(await page.locator('main article').count(), 11);
await page.locator('.sort select').selectOption('gia-giam');
await page.waitForURL('**sap-xep=gia-giam');
assert.match(await page.locator('main article h3').first().innerText(), /áo easygym/i);
await page.getByPlaceholder('Tìm theo tên sản phẩm').fill('khong co gi');
await page.getByRole('button', { name: 'Xóa bộ lọc' }).click();
assert.equal(await page.locator('main article').count(), 11);
ok('shop filter, sort, empty state');

// Product: variant changes the price, cart stores that variant
await go('/san-pham/goat-kefir-s');
assert.match(await page.locator('.price').first().innerText(), /75\.000đ/);
await page.getByLabel('500ml').check({ force: true });
assert.match(await page.locator('.price').first().innerText(), /130\.000đ/);
await page.getByRole('button', { name: 'Tăng số lượng' }).click();
await page.getByRole('button', { name: 'Thêm vào giỏ' }).click();
assert.equal(await cartCount(), '2');
assert.match(await page.locator('.toast').innerText(), /giỏ hàng/i);
ok('product variant + add to cart');

// Moving to a related product resets the picked variant
await page.locator('.more article h3 a').first().click();
await page.waitForURL(/san-pham\/(?!goat-kefir-s)/);
assert.equal(await page.locator('.opts input:checked').inputValue(), '0');
ok('variant resets on navigation');

// Quick add from a card, wishlist
await go('/cua-hang?dm=khoe');
await page.getByRole('button', { name: /Thêm Thạch Jellean/ }).click();
await page.getByRole('button', { name: 'Yêu thích Chocolean' }).click();
await go('/yeu-thich');
assert.equal(await page.locator('main article').count(), 1);
ok('quick add + wishlist');

// Cart maths and persistence: 2 x 130k + 1 x 89k = 349k
await go('/gio-hang');
assert.match(await page.locator('aside').innerText(), /349\.000đ/);
await page.getByRole('button', { name: 'Tăng số lượng' }).first().click();
assert.match(await page.locator('aside').innerText(), /479\.000đ/);
await page.reload({ waitUntil: 'networkidle' });
assert.match(await page.locator('aside').innerText(), /479\.000đ/, 'cart persists');
await page.getByRole('button', { name: 'Xóa Thạch Jellean' }).click();
assert.match(await page.locator('aside').innerText(), /390\.000đ/);
ok('cart maths, persistence, remove');

// Checkout validates, then places the order and empties the cart
await page.getByRole('link', { name: 'Thanh toán' }).click();
await page.waitForURL('**/thanh-toan');
await page.getByRole('button', { name: 'Đặt hàng' }).click();
assert.match(await page.locator('main').innerText(), /Vui lòng nhập họ tên/);
assert.match(await page.locator('main').innerText(), /Số điện thoại gồm 10 số/);
await page.locator('main').getByLabel('Họ và tên').fill('Nguyễn Văn Demo');
await page.locator('main').getByLabel('Số điện thoại').fill('123');
await page.getByLabel('Địa chỉ giao hàng').fill('1 Đường Demo');
await page.getByRole('button', { name: 'Đặt hàng' }).click();
assert.match(await page.locator('main').innerText(), /Số điện thoại gồm 10 số/);
await page.locator('main').getByLabel('Số điện thoại').fill('0900000000');
const [orderRes] = await Promise.all([page.waitForResponse('**/api/orders'), page.getByRole('button', { name: 'Đặt hàng' }).click()]);
assert.equal(orderRes.status(), 201);
const order = await orderRes.json();
assert.ok(order.total > 0, 'server priced the order');
assert.match(await page.locator('main').innerText(), new RegExp(`Đã nhận đơn ${order.code}`));
assert.equal(await page.locator('header .count').count(), 1, 'only the wishlist badge is left');
ok('checkout validation + order');

// Course cards match easygym.vn: no price, a detail link and a sign-up popup with the course pre-selected
await go('/khoa-tap');
assert.doesNotMatch(await page.locator('main').innerText(), /\d\.\d{3}\.\d{3}đ|\d{3}\.000đ/, 'no price on cards');
assert.equal(await page.getByRole('link', { name: 'Xem chi tiết' }).count(), 4);
await page.getByRole('button', { name: 'Đăng ký ngay' }).nth(2).click();
const popup = page.locator('dialog[open]');
await popup.waitFor();
assert.equal(await popup.locator('select').inputValue(), 'coach-online');
await popup.getByRole('button', { name: 'Đăng ký ngay' }).click();
assert.match(await popup.innerText(), /Vui lòng nhập họ tên/);
await popup.getByLabel('Họ và tên').fill('Nguyễn Văn Demo');
await popup.getByLabel('Số điện thoại').fill('0900 000 000');
await popup.getByLabel('Ghi chú').fill('Tập buổi tối');
const [popupRes] = await Promise.all([page.waitForResponse('**/api/leads'), popup.getByRole('button', { name: 'Đăng ký ngay' }).click()]);
assert.match(await popup.innerText(), /Đã nhận đăng ký DK\d+/);
const viaPopup = popupRes.request().postDataJSON();
assert.deepEqual([viaPopup.kind, viaPopup.course, viaPopup.note], ['register', 'coach-online', 'Tập buổi tối']);
await popup.getByRole('button', { name: 'Đóng' }).click();

// Reopening starts a clean form
await page.getByRole('button', { name: 'Đăng ký ngay' }).first().click();
assert.equal(await page.locator('dialog[open] select').inputValue(), 'tap-luyen-1-1');
assert.equal(await page.locator('dialog[open] input').first().inputValue(), '');
await page.keyboard.press('Escape');

// The detail page is where prices live
await page.getByRole('link', { name: 'Xem chi tiết' }).first().click();
await page.waitForURL('**/khoa-tap/tap-luyen-1-1');
await page.getByLabel('6 tháng').check({ force: true });
assert.match(await page.locator('.price').innerText(), /45\.700\.000đ/);
await page.locator('#dang-ky').getByRole('button', { name: 'Đăng ký ngay' }).click();
assert.match(await page.locator('main').innerText(), /Vui lòng nhập họ tên/);
await page.locator('#dang-ky').getByLabel('Họ và tên').fill('Nguyễn Văn Demo');
await page.locator('#dang-ky').getByLabel('Số điện thoại').fill('0900 000 000');
const [leadRes] = await Promise.all([page.waitForResponse('**/api/leads'), page.locator('#dang-ky').getByRole('button', { name: 'Đăng ký ngay' }).click()]);
assert.match(await page.locator('main').innerText(), /Đã nhận đăng ký DK\d+/);
assert.equal(leadRes.request().postDataJSON().plan, '6 tháng');
ok('course cards, popup, detail page');

await go('/lien-he#dang-ky');
assert.ok(await page.locator('#dang-ky select').isVisible());
assert.match(await page.locator('.map iframe').getAttribute('src'), /q=10\.8475319,106\.6541837&.*output=embed/);
assert.match(await page.locator('.map a', { hasText: 'Chỉ đường' }).getAttribute('href'), /destination=10\.8475319,106\.6541837$/);
ok('contact form + map');

const missing = await page.goto(base + '/san-pham/khong-co');
assert.equal(missing.status(), 404);
assert.match(await page.locator('main').innerText(), /Không tìm thấy/);
ok('404 page');

// ---- Audit ----
const PAGES = ['/', '/cua-hang', '/san-pham/goat-kefir-s', '/san-pham/chocolean', '/san-pham/kim-chi-organic', '/khoa-tap', '/khoa-tap/coach-online', '/gioi-thieu', '/lien-he', '/gio-hang', '/yeu-thich', '/tin-tuc', '/tin-tuc/review-goat-kefir-sua-len-men-tu-nhien-tu-sua-de-a2', '/easymeal', '/faq', '/hinh-anh', '/video'];

const audit = () => {
		const lum = (c) => {
			const a = c.map((v) => ((v /= 255) <= 0.03928 ? v / 12.92 : ((v + 0.055) / 1.055) ** 2.4));
			return 0.2126 * a[0] + 0.7152 * a[1] + 0.0722 * a[2];
		};
		const nums = (s) => s.match(/[\d.]+/g).map(Number);
		const bgOf = (e) => {
			for (let n = e; n && n !== document.documentElement; n = n.parentElement) {
				const s = getComputedStyle(n);
				if (s.backgroundImage !== 'none') return null;
				const b = nums(s.backgroundColor);
				if (b.length < 4 || b[3] > 0.9) return b.slice(0, 3);
			}
			return [255, 255, 255];
		};
		const shown = (e) => {
			const r = e.getBoundingClientRect();
			const s = getComputedStyle(e);
			return r.width > 0 && r.height > 0 && s.visibility !== 'hidden' && s.display !== 'none' && s.opacity !== '0';
		};
		const contrast = [];
		for (const e of document.querySelectorAll('body *')) {
			if (e.closest('.ghost, .hero') || !shown(e) || ![...e.childNodes].some((n) => n.nodeType === 3 && n.textContent.trim())) continue;
			const b = bgOf(e);
			if (!b) continue;
			const c = nums(getComputedStyle(e).color).slice(0, 3);
			const [x, y] = [lum(c), lum(b)];
			const q = (Math.max(x, y) + 0.05) / (Math.min(x, y) + 0.05);
			if (q < 4.5) contrast.push(`${q.toFixed(2)} "${e.textContent.trim().slice(0, 30)}"`);
		}
		const targets = [...document.querySelectorAll('a,button,input,select,label:has(input[type=radio])')].filter((e) => shown(e) && !e.matches('.sr-only') && !e.closest('.hits, .skip, .crumbs, footer, p') && e.getAttribute('aria-hidden') !== 'true' && e.tabIndex !== -1);
		const small = targets.filter((e) => {
			const r = e.getBoundingClientRect();
			return r.height < 44 || (r.width < 44 && !e.closest('nav'));
		});
		const unnamed = [...document.querySelectorAll('a,button')].filter((e) => shown(e) && !e.closest('[aria-hidden=true]') && e.getAttribute('aria-hidden') !== 'true' && !(e.innerText || '').trim() && !e.getAttribute('aria-label'));
		// Anything floating over titles or prices at mid-screen, the way the old contact dock did.
		const dock = [...document.querySelectorAll('h1, .price')].filter(shown).filter((e) => {
			e.scrollIntoView({ block: 'center' });
			const r = e.getBoundingClientRect();
			const top = document.elementFromPoint(r.left + 2, r.top + r.height / 2);
			return top && !e.contains(top) && !top.contains(e);
		});
		scrollTo(0, 0);
		return {
			overflow: document.documentElement.scrollWidth - innerWidth,
			contrast: [...new Set(contrast)],
			small: small.map((e) => `${(e.innerText || e.getAttribute('aria-label') || e.className).toString().trim().slice(0, 26)} ${Math.round(e.getBoundingClientRect().width)}x${Math.round(e.getBoundingClientRect().height)}`),
			unnamed: unnamed.length,
			noAlt: [...document.images].filter((i) => i.getAttribute('alt') === null).length,
			h1: document.querySelectorAll('h1').length,
			covered: dock.map((e) => e.textContent.trim().slice(0, 30))
		};
};

let bad = 0;
const widths = [
	[375, 812],
	[768, 1024],
	[1440, 900],
	[1920, 1080]
];
for (const [w, h] of widths) {
	const c = await browser.newContext({ viewport: { width: w, height: h }, locale: 'vi-VN', isMobile: w < 700, hasTouch: w < 700, reducedMotion: 'reduce' });
	const pg = await c.newPage();
	pg.on('pageerror', (e) => errors.push(String(e)));
	for (const p of PAGES) {
		await pg.goto(base + p, { waitUntil: 'networkidle' });
		const a = await pg.evaluate(audit);
		const issues = [];
		if (a.overflow > 0) issues.push(`overflow ${a.overflow}px`);
		if (a.contrast.length) issues.push(`contrast ${a.contrast.join(', ')}`);
		if (w < 700 && a.small.length) issues.push(`tap ${a.small.join(', ')}`);
		if (a.unnamed) issues.push(`${a.unnamed} unnamed controls`);
		if (a.noAlt) issues.push(`${a.noAlt} images without alt`);
		if (a.h1 !== 1) issues.push(`${a.h1} h1`);
		if (a.covered.length) issues.push(`covered ${a.covered.join(', ')}`);
		if (p === '/' && w >= 1440) {
			const lines = await pg.locator('h1').evaluate((e) => Math.round(e.getBoundingClientRect().height / parseFloat(getComputedStyle(e).lineHeight)));
			if (lines !== 1) issues.push(`hero title on ${lines} lines`);
		}
		if (issues.length) {
			bad += issues.length;
			console.log(`✘ ${w}px ${p}: ${issues.join(' | ')}`);
		}
	}
	await c.close();
}
assert.match(await page.locator('meta[name=viewport]').getAttribute('content'), /^width=device-width, initial-scale=1$/, 'zoom is allowed');
console.log(bad ? `audit: ${bad} issue(s)` : '✔ audit clean at 375, 768, 1440, 1920');

if (errors.length) console.log('errors:\n' + [...new Set(errors)].join('\n'));
await browser.close();
process.exit(bad || errors.length ? 1 : 0);
