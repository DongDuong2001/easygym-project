<script>
	import VideoGrid from '#lib/components/VideoGrid.svelte';

	let { data } = $props();
</script>

<svelte:head>
	<title>Video | easyGYM</title>
</svelte:head>

<div class="container page">
	<p class="crumbs"><a href="/">Trang chủ</a> / <span>Video</span></p>
	<h1 class="heading">Video</h1>
	<VideoGrid videos={data.media.videos} />
</div>
