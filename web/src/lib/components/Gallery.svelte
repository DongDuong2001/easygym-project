<script>
	import Icon from './Icon.svelte';

	let { albums, limit = Infinity } = $props();

	let picked = $state('');
	let box = $state();
	let index = $state(0);

	const tab = $derived(picked || albums[0]?.id);
	const album = $derived(albums.find((a) => a.id === tab) ?? { images: [] });
	const shown = $derived(album.images.slice(0, limit));

	function open(i) {
		index = i;
		box.showModal();
	}

	const step = (by) => (index = (index + by + album.images.length) % album.images.length);

	function keys(e) {
		if (e.key === 'ArrowRight') step(1);
		if (e.key === 'ArrowLeft') step(-1);
	}
</script>

<div class="tabs" role="group" aria-label="Chọn album ảnh">
	{#each albums as a}
		<button onclick={() => (picked = a.id)} aria-pressed={tab === a.id}>{a.title}</button>
	{/each}
</div>

<div class="photos">
	{#each shown as src, i (src)}
		<button onclick={() => open(i)} aria-label={`Xem ảnh ${i + 1} của album ${album.title}`}>
			<img {src} alt="" loading="lazy" />
		</button>
	{/each}
</div>

<dialog bind:this={box} onkeydown={keys} aria-label={`Album ${album.title}`}>
	<img src={album.images[index]} alt={`Ảnh ${index + 1} trong album ${album.title}`} />
	<p>{index + 1} / {album.images.length}</p>
	<button class="close" onclick={() => box.close()} aria-label="Đóng"><Icon name="x" size={24} /></button>
	<button class="prev" onclick={() => step(-1)} aria-label="Ảnh trước"><Icon name="left" size={24} /></button>
	<button class="next" onclick={() => step(1)} aria-label="Ảnh sau"><Icon name="right" size={24} /></button>
</dialog>

<style>
	.photos {
		display: grid;
		grid-template-columns: repeat(4, minmax(0, 1fr));
		gap: var(--gap);
	}

	.photos button {
		display: block;
		aspect-ratio: 4 / 5;
		overflow: hidden;
		background: var(--navy);
	}

	.photos img {
		width: 100%;
		height: 100%;
		object-fit: cover;
		transition: transform 0.35s;
	}

	.photos button:hover img {
		transform: scale(1.03);
	}

	dialog {
		width: 100vw;
		max-width: none;
		height: 100dvh;
		max-height: none;
		padding: 0;
		background: rgb(13 13 63 / 0.96);
		color: var(--on-dark);
		border: 0;
	}

	dialog::backdrop {
		background: rgb(13 13 63 / 0.96);
	}

	dialog img {
		width: 100%;
		height: 100%;
		padding: 56px 16px;
		object-fit: contain;
	}

	dialog p {
		position: absolute;
		left: 50%;
		bottom: 16px;
		font-size: var(--t-sm);
		transform: translateX(-50%);
	}

	dialog button {
		position: absolute;
		display: grid;
		place-items: center;
		width: 48px;
		height: 48px;
		background: var(--surface);
		color: var(--ink);
	}

	.close {
		top: 8px;
		right: 8px;
	}

	.prev,
	.next {
		top: 50%;
		transform: translateY(-50%);
	}

	.prev {
		left: 8px;
	}

	.next {
		right: 8px;
	}

	@media (max-width: 720px) {
		.photos {
			grid-template-columns: repeat(2, minmax(0, 1fr));
		}
	}
</style>
