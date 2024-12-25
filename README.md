### 1. Install Docker

If you haven't installed Docker yet, install it by running:

```bash
curl -sSL https://get.docker.com | sh
sudo usermod -aG docker $(whoami)
exit
```


---

### Cách hoạt động
1. API nhận nhiều file từ form-data (key là `files`).
2. Tối đa 12 file có thể được tải lên.
3. Mỗi file được kiểm tra:
   - PDF: Lưu trực tiếp.
   - Hình ảnh: Resize nếu cần.
   - Các định dạng không hỗ trợ sẽ bị từ chối.
4. Kết quả trả về dạng JSON:
   - **Thành công**: Trả về danh sách các đường dẫn file.
   - **Thất bại**: Trả về thông báo lỗi chi tiết.

---

### Request mẫu
Sử dụng cURL để kiểm tra:

```bash
curl -X POST -F "files=@image1.jpg" -F "files=@image2.pdf" http://localhost:8080/upload
```

---

### Phản hồi JSON
**Thành công:**
```json
{
    "data": [
        "uploads/0440bfbc/ca3e/448a/a393/7309b787a88e.jpg",
        "uploads/1234abcd/ef56/7890/gh12/ijklmnopqrst.pdf"
    ],
    "status": "success",
    "message": "ok"
}
```

**Lỗi:**
```json
{
    "status": "error",
    "message": "Error processing upload"
}
```