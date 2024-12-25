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


# Cách sử dụng

### 1. **Node.js (sử dụng thư viện `axios`)**

Để thực hiện upload file trong Node.js, bạn có thể sử dụng thư viện `axios` và `form-data`.

Cài đặt `axios` và `form-data`:

```bash
npm install axios form-data
```

Ví dụ mã nguồn:

```javascript
const axios = require('axios');
const FormData = require('form-data');
const fs = require('fs');

// Tạo form-data
const form = new FormData();
form.append('files', fs.createReadStream('image1.jpg'));
form.append('files', fs.createReadStream('image2.pdf'));

// Gửi POST request
axios.post('http://localhost:8080/upload', form, {
    headers: {
        ...form.getHeaders()
    }
})
.then(response => {
    console.log('Upload successful:', response.data);
})
.catch(error => {
    console.error('Error during upload:', error);
});
```

### 2. **PHP (sử dụng `cURL`)**

Trong PHP, bạn có thể sử dụng `cURL` để upload file như sau:

Ví dụ mã nguồn:

```php
<?php
$curl = curl_init();

// Tạo dữ liệu POST với file
$data = array(
    'files' => new CURLFile('image1.jpg'),
    'files' => new CURLFile('image2.pdf')
);

// Cấu hình cURL
curl_setopt_array($curl, array(
    CURLOPT_URL => "http://localhost:8080/upload",
    CURLOPT_RETURNTRANSFER => true,
    CURLOPT_POST => true,
    CURLOPT_POSTFIELDS => $data,
));

$response = curl_exec($curl);

if(curl_errno($curl)) {
    echo 'Error:' . curl_error($curl);
} else {
    echo 'Upload successful: ' . $response;
}

curl_close($curl);
?>
```

### 3. **Python (sử dụng `requests`)**

Trong Python, bạn có thể sử dụng thư viện `requests` để thực hiện upload file.

Cài đặt `requests`:

```bash
pip install requests
```

Ví dụ mã nguồn:

```python
import requests

# Tạo các file cần upload
files = {
    'files': open('image1.jpg', 'rb'),
    'files': open('image2.pdf', 'rb')
}

# Gửi POST request
response = requests.post('http://localhost:8080/upload', files=files)

# In kết quả
if response.status_code == 200:
    print('Upload successful:', response.text)
else:
    print('Error during upload:', response.status_code)

# Đóng file
files['files'].close()
```
