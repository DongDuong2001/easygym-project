<script>
	import { onMount } from 'svelte';
	import Benefits from '#lib/components/Benefits.svelte';
	import CourseCard from '#lib/components/CourseCard.svelte';
	import Gallery from '#lib/components/Gallery.svelte';
	import Icon from '#lib/components/Icon.svelte';
	import PostCard from '#lib/components/PostCard.svelte';
	import ProductCard from '#lib/components/ProductCard.svelte';
	import SignupForm from '#lib/components/SignupForm.svelte';
	import VideoGrid from '#lib/components/VideoGrid.svelte';
	import MapEmbed from '#lib/components/MapEmbed.svelte';
	import { catalog } from '#lib/catalog.js';
	import { filterProducts } from '#lib/shop.js';
	import { brand, founder, heroVideo, img, stats } from '#lib/site.js';

	let { data } = $props();

	const CATEGORIES = $derived(catalog().categories);
	const courses = $derived(catalog().courses);
	const products = $derived(catalog().products);

	const ticker = ['Hơn 10 năm đồng hành và phát triển', '2.900+ học viên đã đồng hành', 'easyGYM: tập luyện mỗi ngày', 'easyMEAL: ăn sạch, sống dễ', 'Mở cửa 8:00 - 20:00 mỗi ngày'];

	let tab = $state('');
	let video = $state();
	let paused = $state(true);

	const trending = $derived(filterProducts(products, { dm: tab }).slice(0, 8));

	// Reduced-motion visitors get the finished logo frame instead of a moving hero.
	onMount(() => {
		if (!matchMedia('(prefers-reduced-motion: reduce)').matches) {
			video.play().catch(() => {});
			return;
		}
		const toEnd = () => (video.currentTime = video.duration - 0.1);
		if (video.readyState) toEnd();
		else video.addEventListener('loadedmetadata', toEnd, { once: true });
	});
</script>

<section class="hero">
	<div class="stage dark">
		<video bind:this={video} bind:paused src={heroVideo} muted loop playsinline preload="auto" aria-hidden="true"></video>
		<div class="caption">
			<div>
				<h1>Tập dễ dàng, bền vững</h1>
				<p>Hơn 10 năm đồng hành và phát triển</p>
			</div>
			<a class="btn light" href="/khoa-tap">Chọn gói tập</a>
		</div>
		<button class="toggle" onclick={() => (paused = !paused)} aria-label={paused ? 'Phát video' : 'Dừng video'}>
			<Icon name={paused ? 'play' : 'pause'} size={18} />
		</button>
	</div>

	<aside aria-labelledby="su-kien">
		<h2 id="su-kien">Sự kiện</h2>
		<ul>
			{#each data.posts.slice(0, 4) as p}
				<li>
					<a href={`/tin-tuc/${p.slug}`}>
						<img src={p.image} alt="" width="84" height="84" />
						<span>{p.title}<small>{p.date}</small></span>
					</a>
				</li>
			{/each}
		</ul>
		<a class="banner" href="/khoa-tap" aria-label="Sức khỏe hôm nay, tương lai ngày mai. Xem các khóa học easyGYM">
			<img src={img('slide/688c0f9b-f04d-4e27-888b-f1f87e60e56d.webp')} alt="" />
		</a>
	</aside>
</section>

<div class="ticker">
	<div class="run">
		{#each [0, 1] as n}
			<ul aria-hidden={n === 1 ? 'true' : undefined}>
				{#each ticker as t}<li>{t}</li>{/each}
				{#each ticker as t}<li>{t}</li>{/each}
			</ul>
		{/each}
	</div>
</div>

<section class="section container">
	<h2 class="heading">Sản phẩm</h2>
	<div class="tabs" role="group" aria-label="Lọc theo nhóm sản phẩm">
		<button onclick={() => (tab = '')} aria-pressed={tab === ''}>Tất cả</button>
		{#each CATEGORIES as c}
			<button onclick={() => (tab = c.value)} aria-pressed={tab === c.value}>{c.label}</button>
		{/each}
	</div>
	<div class="grid">
		{#each trending as p (p.slug)}<ProductCard product={p} />{/each}
	</div>
	<p class="more"><a class="btn line" href={tab ? `/cua-hang?dm=${tab}` : '/cua-hang'}>Xem tất cả sản phẩm</a></p>
</section>

<section class="container promos" aria-label="Nhóm sản phẩm">
	<article>
		<div>
			<p class="kicker">easyMEAL</p>
			<h2>Ăn sạch, sống dễ</h2>
			<p>Bánh Bông Lean từ bột nguyên cám, mật ong rừng và trứng gà ta.</p>
			<a class="btn" href="/easymeal">Khám phá easyMEAL</a>
		</div>
		<img src={img('food/screenshot-2026-05-28-233247.webp')} alt="Bánh Bông Lean vị socola" loading="lazy" />
	</article>
	<article>
		<div>
			<p class="kicker">Lên men tự nhiên</p>
			<h2>Goat Kefir</h2>
			<p>Lên men từ sữa dê thuần khiết, đóng chai từ 300ml đến 1000ml.</p>
			<a class="btn" href="/cua-hang?dm=kefir">Mua ngay</a>
		</div>
		<img src={img('food/screenshot-2026-05-28-225109.webp')} alt="Các chai Goat Kefir xếp cạnh nhau" loading="lazy" />
	</article>
</section>

<section class="section container">
	<h2 class="heading">Khóa học đang mở</h2>
	<div class="grid">
		{#each courses as c}<CourseCard course={c} />{/each}
	</div>
</section>

<section class="founder container">
	<figure>
		<img src={img('hinhanh/jimjsafkasmc.webp')} alt="Đạt easyGYM, poster Shape Your Body" loading="lazy" />
		<figcaption><b>{founder.name}</b>{founder.role}</figcaption>
	</figure>
	<div class="story">
		<h2>{founder.title}</h2>
		<ul class="intro-list">
			{#each founder.intro as t}<li>{t}</li>{/each}
		</ul>
		<ul class="path">
			{#each founder.path as t}<li>{t}</li>{/each}
		</ul>
		<p>{founder.now}</p>
		<div class="lists">
			<h3>Kinh nghiệm</h3>
			<ul>
				{#each founder.work as t}<li><Icon name="check" size={16} /> {t}</li>{/each}
			</ul>
		</div>
		<a class="btn" href="/gioi-thieu">Đôi lời giới thiệu</a>
	</div>
</section>

<section class="awards dark">
	<div class="container">
		<h2 class="heading">Thành tích huấn luyện viên</h2>
		<ul class="medals">
			{#each founder.medals as m}
				<li><b>{m.value}</b>{m.label}</li>
			{/each}
		</ul>
		<ul class="titles">
			{#each founder.titles as t}
				<li><Icon name="trophy" size={20} /> {t}</li>
			{/each}
		</ul>
	</div>
</section>

<Benefits />

<section class="media dark">
	<div class="container">
		<p class="eyebrow">Video thực tế</p>
		<h2 class="heading">Tại easyGYM</h2>
		<VideoGrid videos={data.media.videos} />

		<h2 class="heading second">Hình ảnh hoạt động</h2>
		<Gallery albums={data.media.albums} limit={8} />
		<p class="more"><a class="btn light" href="/hinh-anh">Xem tất cả hình ảnh</a></p>
	</div>
</section>

<section class="section container">
	<h2 class="heading">Tin tức</h2>
	<div class="news">
		{#each data.posts.slice(0, 3) as p}<PostCard post={p} />{/each}
	</div>
	<p class="more"><a class="btn line" href="/tin-tuc">Xem tất cả tin tức</a></p>
</section>

<section class="join dark" id="dang-ky">
	<div class="container join-in">
		<div>
			<h2>Đăng ký trải nghiệm ngay hôm nay</h2>
			<ul class="stats">
				{#each stats as s}
					<li><b>{s.value}</b>{s.label}</li>
				{/each}
			</ul>
			<ul class="info">
				<li><Icon name="clock" /> {brand.hours}</li>
				<li><Icon name="pin" /> <a href={brand.map} target="_blank" rel="noreferrer">{brand.address}</a></li>
			</ul>
		</div>
		<SignupForm id="dk-home" />
	</div>
	<div class="container where">
		<MapEmbed height={380} />
	</div>
</section>

<style>
	.hero {
		display: grid;
		grid-template-columns: minmax(0, 1fr) clamp(300px, 24vw, 400px);
		min-height: clamp(480px, 72vh, 780px);
		background: var(--ink);
	}

	.stage {
		position: relative;
		display: flex;
		align-items: flex-end;
		overflow: hidden;
		color: var(--on-dark);
		isolation: isolate;
	}

	video {
		position: absolute;
		inset: 0;
		z-index: -2;
		width: 100%;
		height: 100%;
		object-fit: cover;
	}

	/* Only the bottom edge is shaded, so the video stays the main thing on screen. */
	.stage::after {
		content: '';
		position: absolute;
		inset: 50% 0 0;
		z-index: -1;
		background: linear-gradient(180deg, rgb(13 13 63 / 0) 0%, rgb(13 13 63 / 0.92) 100%);
	}

	.caption {
		display: flex;
		align-items: flex-end;
		justify-content: space-between;
		gap: 16px 24px;
		width: 100%;
		padding: 24px var(--gutter);
	}

	h1 {
		font-size: clamp(36px, 5vw, 80px);
		font-weight: 700;
		line-height: 1.05;
		text-transform: uppercase;
	}

	.caption p {
		margin-top: 4px;
		font-size: clamp(16px, 1.5vw, 20px);
		font-weight: 500;
	}

	.caption .btn {
		flex: none;
	}

	.toggle {
		position: absolute;
		top: 12px;
		right: 12px;
		display: grid;
		place-items: center;
		width: 44px;
		height: 44px;
		background: rgb(13 13 63 / 0.7);
		color: var(--on-dark);
		border-radius: 50%;
	}

	.toggle:hover {
		background: var(--ink);
	}

	aside {
		display: flex;
		flex-direction: column;
		min-width: 0;
		background: var(--surface);
		border-left: 1px solid var(--line);
	}

	aside h2 {
		padding: 12px 16px;
		background: var(--brand);
		color: var(--on-dark);
		font-size: var(--t-md);
		font-weight: 600;
		text-transform: uppercase;
	}

	aside ul {
		display: grid;
		flex: 1;
		grid-auto-rows: 1fr;
	}

	aside li a {
		display: flex;
		align-items: center;
		gap: 12px;
		height: 100%;
		padding: 10px 16px;
		border-bottom: 1px solid var(--line);
		font-size: var(--t-sm);
		font-weight: 600;
		line-height: 1.3;
	}

	aside li a:hover {
		background: var(--soft);
	}

	aside li img {
		flex: none;
		object-fit: cover;
		background: var(--soft);
	}

	aside li span {
		display: -webkit-box;
		overflow: hidden;
		-webkit-box-orient: vertical;
		-webkit-line-clamp: 4;
		line-clamp: 4;
	}

	aside small {
		display: block;
		margin-top: 4px;
		color: var(--muted);
		font-size: var(--t-xs);
		font-weight: 500;
	}

	/* Native ratio: the banner carries text, so it must never crop. */
	.banner img {
		width: 100%;
		aspect-ratio: 1916 / 821;
	}

	.ticker {
		overflow: hidden;
		border-bottom: 1px solid var(--line);
		font-size: var(--t-xs);
		font-weight: 600;
		text-transform: uppercase;
	}

	.run {
		display: flex;
		width: max-content;
		animation: run 40s linear infinite;
	}

	.run ul {
		display: flex;
	}

	.run li {
		display: flex;
		align-items: center;
		gap: 16px;
		padding: 7px 8px;
		white-space: nowrap;
	}

	.run li::after {
		content: '';
		width: 7px;
		height: 7px;
		margin-left: 8px;
		border: 1px solid currentColor;
		border-radius: 50%;
	}

	@keyframes run {
		to {
			transform: translateX(-50%);
		}
	}

	.more {
		margin-top: 32px;
		text-align: center;
	}

	.promos {
		display: grid;
		grid-template-columns: repeat(2, minmax(0, 1fr));
		gap: var(--gap);
	}

	.promos article {
		display: grid;
		grid-template-columns: 1fr minmax(0, 0.8fr);
		align-items: center;
		gap: 16px;
		min-height: 280px;
		padding-left: clamp(20px, 3vw, 40px);
		background: var(--soft);
	}

	.promos img {
		width: 100%;
		height: 100%;
		object-fit: cover;
	}

	.kicker,
	.eyebrow {
		font-size: var(--t-xs);
		font-weight: 600;
		text-transform: uppercase;
	}

	.promos h2 {
		margin: 6px 0 8px;
		font-size: clamp(24px, 2.6vw, 32px);
		font-weight: 600;
	}

	.promos article p:not(.kicker) {
		margin-bottom: 16px;
		color: var(--muted);
		font-size: var(--t-sm);
	}

	.founder {
		display: grid;
		grid-template-columns: minmax(0, 0.8fr) minmax(0, 1fr);
		gap: clamp(24px, 4vw, 64px);
		align-items: start;
	}

	.founder figure {
		position: sticky;
		top: calc(var(--bar) + 16px);
		background: var(--ink);
	}

	.founder figure img {
		width: 100%;
	}

	figcaption {
		padding: 12px 16px;
		color: var(--on-dark-2);
		font-size: var(--t-sm);
	}

	figcaption b {
		display: block;
		color: var(--on-dark);
		font-size: var(--t-lg);
		font-weight: 600;
		text-transform: uppercase;
	}

	.story {
		display: grid;
		gap: 14px;
		justify-items: start;
	}

	.story h2 {
		font-size: clamp(30px, 3.6vw, 48px);
		font-weight: 600;
		text-transform: uppercase;
	}

	.intro-list {
		font-weight: 600;
	}

	.path {
		padding-left: 18px;
		list-style: disc;
		color: var(--muted);
	}

	.story > p {
		color: var(--muted);
	}

	.lists {
		width: 100%;
		padding-block: 14px;
		border-block: 1px solid var(--line);
	}

	.lists h3 {
		margin-bottom: 8px;
		font-size: var(--t-sm);
		font-weight: 600;
		text-transform: uppercase;
	}

	.lists li {
		display: flex;
		gap: 8px;
		padding-block: 3px;
	}

	.lists :global(svg) {
		flex: none;
		margin-top: 5px;
		color: var(--brand);
	}

	.awards {
		margin-top: var(--section);
		padding-block: var(--section);
		background: var(--navy);
		color: var(--on-dark);
	}

	.medals {
		display: grid;
		grid-template-columns: repeat(3, minmax(0, 1fr));
		gap: var(--gap);
		margin-bottom: var(--gap);
	}

	.medals li {
		padding: clamp(16px, 2.4vw, 28px);
		background: var(--ink);
		color: var(--on-dark-2);
		font-size: var(--t-sm);
		text-align: center;
	}

	.medals b {
		display: block;
		color: var(--on-dark);
		font-size: clamp(44px, 6vw, 80px);
		font-weight: 700;
		line-height: 1;
	}

	.titles {
		display: grid;
		grid-template-columns: repeat(3, minmax(0, 1fr));
		column-gap: 32px;
	}

	.titles li {
		display: flex;
		align-items: center;
		gap: 12px;
		min-height: 56px;
		padding-block: 8px;
		border-bottom: 1px solid var(--brand);
		font-weight: 500;
	}

	.titles :global(svg) {
		flex: none;
		color: var(--on-dark-2);
	}

	.media {
		padding-block: var(--section);
		background: var(--ink);
		color: var(--on-dark);
	}

	.eyebrow {
		margin-bottom: 4px;
		text-align: center;
		color: var(--on-dark-2);
	}

	.second {
		margin-top: var(--section);
	}

	.news {
		display: grid;
		grid-template-columns: repeat(3, minmax(0, 1fr));
		gap: 24px var(--gap);
	}

	.join {
		padding-block: var(--section);
		background: var(--navy);
		color: var(--on-dark);
	}

	.where {
		margin-top: 48px;
	}

	.join-in {
		display: grid;
		grid-template-columns: 1.2fr minmax(0, 1fr);
		gap: 40px clamp(32px, 6vw, 96px);
		align-items: start;
	}

	.join h2 {
		font-size: var(--t-xl);
		font-weight: 600;
		text-transform: uppercase;
	}

	.stats {
		display: grid;
		grid-template-columns: repeat(4, minmax(0, 1fr));
		gap: 16px;
		margin-block: 24px;
		color: var(--on-dark-2);
		font-size: var(--t-sm);
	}

	.stats b {
		display: block;
		color: var(--on-dark);
		font-size: clamp(26px, 3vw, 40px);
		font-weight: 700;
		line-height: 1.1;
	}

	.info {
		display: grid;
		gap: 10px;
		color: var(--on-dark-2);
	}

	.info li {
		display: flex;
		align-items: center;
		gap: 10px;
	}

	.info a:hover {
		text-decoration: underline;
	}

	.info :global(svg) {
		flex: none;
	}

	@media (max-width: 1024px) {
		.promos,
		.join-in {
			grid-template-columns: 1fr;
		}

		/* The whole video stays in frame on smaller screens; its logo would crop in a tall box. */
		.hero {
			grid-template-columns: minmax(0, 1fr);
			min-height: 0;
		}

		.stage {
			flex-direction: column;
			align-items: stretch;
			background: var(--ink);
		}

		video {
			position: static;
			z-index: auto;
			height: auto;
			aspect-ratio: 16 / 9;
		}

		.stage::after {
			display: none;
		}

		.caption {
			padding-block: 16px;
		}

		aside {
			border-left: 0;
		}

		.banner {
			display: none;
		}

		.titles {
			grid-template-columns: repeat(2, minmax(0, 1fr));
		}
	}

	@media (max-width: 860px) {
		.founder {
			grid-template-columns: minmax(0, 1fr);
		}

		.founder figure {
			position: static;
			max-width: 420px;
		}

		.news {
			grid-template-columns: minmax(0, 1fr);
		}
	}

	@media (max-width: 720px) {
		.caption {
			flex-direction: column;
			align-items: stretch;
		}

		.promos article {
			grid-template-columns: 1fr 0.7fr;
			min-height: 220px;
			padding-left: 16px;
		}

		.stats {
			grid-template-columns: repeat(2, minmax(0, 1fr));
		}

		.titles {
			grid-template-columns: minmax(0, 1fr);
		}
	}
</style>
