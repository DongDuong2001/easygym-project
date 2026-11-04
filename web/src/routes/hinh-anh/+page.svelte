<script>
	import Gallery from '#lib/components/Gallery.svelte';

	let { data } = $props();
</script>

<svelte:head>
	<title>Hình ảnh hoạt động | easyGYM</title>
</svelte:head>

<div class="container page">
	<p class="crumbs"><a href="/">Trang chủ</a> / <span>Hình ảnh hoạt động</span></p>
	<h1 class="heading">Hình ảnh hoạt động</h1>
	<Gallery albums={data.media.albums} />
</div>
