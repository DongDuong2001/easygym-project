<script>
	import Icon from './Icon.svelte';
	import { catalog } from '#lib/catalog.js';
	import { isPhone } from '#lib/shop.js';
	import { addLead, focusError } from '#lib/store.svelte.js';

	const courses = $derived(catalog().courses);

	// `course` fixes the course (course pages); `preset` only pre-selects it in the dropdown (popup).
	let { course = '', plan = '', preset = '', id = 'dk' } = $props();

	let form = $state({ name: '', phone: '', course: preset, note: '' });
	let tried = $state(false);
	let code = $state('');
	let sending = $state(false);
	let failed = $state('');

	const nameErr = $derived(tried && !form.name.trim() ? 'Vui lòng nhập họ tên' : '');
	const phoneErr = $derived(tried && !isPhone(form.phone) ? 'Số điện thoại gồm 10 số, bắt đầu bằng 0' : '');

	async function submit(e) {
		e.preventDefault();
		tried = true;
		if (!form.name.trim() || !isPhone(form.phone)) return focusError(e.currentTarget);
		if (sending) return;
		sending = true;
		failed = '';
		try {
			code = await addLead({ kind: 'register', name: form.name, phone: form.phone, course: course || form.course, plan, note: form.note, source: location.pathname });
		} catch (err) {
			failed = err.message;
		} finally {
			sending = false;
		}
	}
</script>

{#if code}
	<div class="done" role="status">
		<Icon name="check" size={28} />
		<p><b>Đã nhận đăng ký {code}</b></p>
		<p>easyGYM sẽ gọi lại cho bạn trong giờ làm việc.</p>
	</div>
{:else}
	<form onsubmit={submit} novalidate>
		<label class="field">
			Họ và tên
			<input bind:value={form.name} autocomplete="name" aria-invalid={Boolean(nameErr)} aria-describedby={`${id}-ten`} />
			<span class="err" id={`${id}-ten`}>{nameErr}</span>
		</label>
		<label class="field">
			Số điện thoại
			<input type="tel" bind:value={form.phone} autocomplete="tel" inputmode="tel" aria-invalid={Boolean(phoneErr)} aria-describedby={`${id}-sdt`} />
			<span class="err" id={`${id}-sdt`}>{phoneErr}</span>
		</label>
		{#if !course}
			<label class="field">
				Chọn khóa học
				<select bind:value={form.course}>
					<option value="">Chưa chọn, cần tư vấn</option>
					{#each courses as c}<option value={c.slug}>{c.name}</option>{/each}
				</select>
			</label>
		{/if}
		<label class="field">
			Ghi chú
			<textarea rows="2" bind:value={form.note}></textarea>
		</label>
		{#if failed}<p class="err" role="alert">{failed}</p>{/if}
		<button class="btn block submit" disabled={sending}>{sending ? 'Đang gửi...' : 'Đăng ký ngay'}</button>
	</form>
{/if}

<style>
	form {
		display: grid;
		gap: 10px;
	}

	.submit {
		margin-top: 6px;
	}

	.err:empty {
		display: none;
	}

	:global(.dark) .submit {
		background: var(--surface);
		border-color: var(--surface);
		color: var(--ink);
	}

	:global(.dark) .submit:hover {
		background: var(--soft);
		border-color: var(--soft);
	}

	.done {
		display: grid;
		gap: 6px;
		padding: 24px;
		border: 1px solid currentColor;
	}
</style>
