<script>
	import { img, stats } from '#lib/site.js';

	let { data } = $props();
</script>

<svelte:head>
	<title>Giới thiệu | easyGYM</title>
</svelte:head>

<div class="container page">
	<p class="crumbs"><a href="/">Trang chủ</a> / <span>Giới thiệu</span></p>

	<div class="layout">
		<aside>
			<img src={img('hinh-anh-thuc-te/511269006_3716249228507661_3973768484470522313_n.webp')} alt="Đạt easyGYM trong buổi tập" />
			<ul class="stats">
				{#each stats as s}
					<li><b>{s.value}</b>{s.label}</li>
				{/each}
			</ul>
			<a class="btn block" href="/khoa-tap">Xem khóa học</a>
		</aside>
		<div>
			<h1>Giới thiệu</h1>
			<div class="prose">{@html data.about}</div>
		</div>
	</div>
</div>

<style>
	.layout {
		display: grid;
		grid-template-columns: minmax(0, 0.8fr) minmax(0, 1fr);
		gap: clamp(24px, 4vw, 64px);
		align-items: start;
	}

	aside {
		position: sticky;
		top: calc(var(--bar) + 16px);
		display: grid;
		gap: 16px;
	}

	aside img {
		width: 100%;
		aspect-ratio: 4 / 4.4;
		object-fit: cover;
		background: var(--soft);
	}

	h1 {
		margin-bottom: 16px;
		font-size: clamp(36px, 5vw, 64px);
		text-transform: uppercase;
	}

	.stats {
		display: grid;
		grid-template-columns: repeat(4, minmax(0, 1fr));
		gap: 12px;
		color: var(--muted);
		font-size: var(--t-xs);
	}

	.stats b {
		display: block;
		color: var(--ink);
		font-size: var(--t-xl);
		font-weight: 700;
		line-height: 1.1;
	}

	@media (max-width: 860px) {
		.layout {
			grid-template-columns: minmax(0, 1fr);
		}

		aside {
			position: static;
		}

		.stats {
			grid-template-columns: repeat(2, minmax(0, 1fr));
		}
	}
</style>
