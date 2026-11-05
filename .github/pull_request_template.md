## Mô tả

Thay đổi gì và vì sao. Đóng issue nào (nếu có): Fixes #

## Loại thay đổi

- [ ] `fix`: sửa lỗi
- [ ] `feat`: tính năng mới
- [ ] Breaking change (`!`): làm hỏng hành vi cũ
- [ ] `refactor`: đổi cấu trúc, không đổi hành vi
- [ ] `docs`: tài liệu

## Checklist

- [ ] Nhánh tên `feat/<ten-viec>`, mỗi commit một file, message theo Conventional Commits (xem [CONTRIBUTING.md](../CONTRIBUTING.md))
- [ ] Đã tự review diff
- [ ] Đã chạy trên máy: `cd backend && go vet ./... && go test ./...` và `cd web && npm test && npm run build`
- [ ] Đổi giao diện: đã chạy `npm run e2e` với DB local
- [ ] Đổi DB: có migration mới trong `backend/db/migrations/` với phần `-- +goose Down` hoàn tác được, không sửa migration cũ
- [ ] Đã cập nhật README hoặc CONTRIBUTING nếu cần

## Ảnh chụp (nếu đổi giao diện)

Điện thoại (375px) và máy tính (1440px).
