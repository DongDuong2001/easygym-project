<script>
	import { page } from '$app/state';
</script>

<svelte:head>
	<title>Không tìm thấy trang | easyGYM</title>
</svelte:head>

<div class="container empty">
	<h1>{page.status === 404 ? 'Không tìm thấy trang này' : 'Có lỗi xảy ra'}</h1>
	<p>Trang có thể đã đổi địa chỉ. Bạn quay lại cửa hàng để xem tiếp nhé.</p>
	<a class="btn" href="/cua-hang">Về cửa hàng</a>
</div>
