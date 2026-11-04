<script>
	import CourseCard from '#lib/components/CourseCard.svelte';
	import { catalog } from '#lib/catalog.js';

	const courses = $derived(catalog().courses);
</script>

<svelte:head>
	<title>Khóa học | easyGYM</title>
</svelte:head>

<div class="container page">
	<p class="crumbs"><a href="/">Trang chủ</a> / <span>Khóa học</span></p>
	<h1 class="heading">Khóa học đang mở</h1>
	<p class="lead">Tập trực tiếp, tập nhóm, coach online hoặc tự học kiến thức nền. Chọn cách phù hợp với lịch và chi phí của bạn.</p>
	<h2 class="sr-only">Danh sách khóa học</h2>
	<div class="grid">
		{#each courses as c}<CourseCard course={c} />{/each}
	</div>
</div>
