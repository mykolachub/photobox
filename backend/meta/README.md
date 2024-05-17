# Metadata Service

Routes:

- `POST /api/meta/` upload image
- `GET /api/meta/:file_id` get image
- `GET /api/meta/` search (by name, keywords)
- `DELETE /api/meta/:file_id` delete image

Middleware Protection:

- `middle.Protect` check user existance for any action
- `middle.StorageAvailability` check storage availability for uploads
