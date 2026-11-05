# Bảo mật

## Báo lỗ hổng

Không mở issue công khai. Gửi email cho người bảo trì repo, kèm:

- mô tả lỗ hổng và mức ảnh hưởng,
- các bước tái hiện,
- endpoint hoặc trang bị ảnh hưởng.

Người bảo trì trả lời trong 48 giờ và báo tiến độ đến khi sửa xong.

## Dữ liệu khách hàng

- Bảng `leads` và `orders` chứa tên, số điện thoại, địa chỉ khách. Chỉ xem qua Supabase Studio, không xuất ra file chưa mã hóa.
- Các schema của app (`catalog`, `content`, `leads`, `orders`, `platform`) không được thêm vào "Exposed schemas" của Supabase, để anon key không đọc được.
- Backup hằng ngày (`.github/workflows/db-backup.yml`) được mã hóa bằng GPG trước khi rời GitHub. Khóa riêng chỉ chủ sở hữu giữ.

## API

- POST chỉ nhận `Content-Type: application/json`, body tối đa 64 KB, mọi trường được validate phía Go.
- Giá đơn hàng luôn tính lại từ `catalog`, không tin giá trình duyệt gửi lên.
- Chưa có: giới hạn số lần gửi form theo IP, trang admin và đăng nhập (xem "Chưa làm" trong [README](README.md)).

## Secret

- Không commit `.env`. Mẫu ở `.env.example`.
- `DATABASE_URL` production phải có `sslmode=require`.
- Secret của CI nằm trong GitHub > Settings > Secrets and variables > Actions, danh sách ở [CONTRIBUTING.md](CONTRIBUTING.md).
