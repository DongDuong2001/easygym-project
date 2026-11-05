# Đóng góp

Mọi thay đổi vào `main` đi qua pull request. CI (`.github/workflows/ci.yml`) phải xanh mới được merge.

## Nhánh

Tạo nhánh từ `main`, tên theo dạng `feat/<ten-viec>`: chữ thường, nối bằng dấu gạch ngang.

```bash
git switch main && git pull
git switch -c feat/checkout-momo
```

## Commit

Hai quy tắc, CI kiểm tra cả hai trên mọi commit của PR (bỏ qua commit merge):

1. **Mỗi commit đúng một file.** Sửa ba file thì ba commit. Commit file dùng chung (helper) trước file gọi nó, để mỗi commit vẫn build được.
2. **Conventional Commits:** `<type>(<scope>): <mô tả>`, scope không bắt buộc.

| type | Dùng khi |
|---|---|
| `feat` | thêm tính năng |
| `fix` | sửa lỗi |
| `docs` | tài liệu |
| `style` | định dạng, không đổi logic |
| `refactor` | đổi cấu trúc, không đổi hành vi |
| `perf` | tăng tốc |
| `test` | thêm hoặc sửa test |
| `build` | dependency, công cụ build |
| `ci` | workflow CI |
| `chore` | việc lặt vặt khác |
| `revert` | hoàn tác commit trước |

Thêm `!` sau type khi thay đổi phá vỡ tương thích, ví dụ `feat(api)!: đổi định dạng /api/orders`.

```bash
git add web/src/lib/focus.js
git commit -m "feat(web): add focusError helper" -- web/src/lib/focus.js
git commit -m "fix(checkout): focus first invalid field" -- web/src/routes/thanh-toan/+page.svelte
```

## CI kiểm tra gì

| Job | Bước |
|---|---|
| Branch and commit rules | tên nhánh `feat/...`, commit đúng Conventional Commits, mỗi commit một file |
| Backend (Go) | `gofmt`, `go vet`, `go test -race`, `go build` |
| Frontend (SvelteKit) | `npm ci`, `npm test`, `npm run build` |
| E2E (Playwright) | chạy sau khi Backend và Frontend xanh: migrate + seed vào Postgres 17 tạm, chạy API và site, `npm run e2e` |

Chạy trước trên máy để khỏi chờ CI, các lệnh nằm ở mục "Kiểm tra" trong [README](README.md).

## Bảo vệ nhánh main

Admin repo bật một lần trong GitHub > Settings > Branches > Add rule cho `main`:

- Require a pull request before merging.
- Require status checks to pass: `Branch and commit rules`, `Backend (Go)`, `Frontend (SvelteKit)`, `E2E (Playwright)`.
- Không dùng "Squash merge" (sẽ gộp nhiều file vào một commit). Dùng "Create a merge commit".

## Bản quyền

Mã nguồn thuộc về Ông Đạt, xem [LICENSE](LICENSE). Đóng góp vào repo này đồng nghĩa với việc chuyển quyền sở hữu phần đóng góp cho chủ sở hữu.
