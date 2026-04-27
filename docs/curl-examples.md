# Sample cURL Requests

## Login Admin
curl -X POST http://localhost:8080/admin/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin123"}'

## Public Posts
curl "http://localhost:8080/posts?page=1&limit=10&search=welcome"

## Public Post by Slug
curl "http://localhost:8080/posts/welcome-to-cms"

## Create Post (Admin)
curl -X POST http://localhost:8080/admin/posts \
  -H "Authorization: Bearer <TOKEN>" \
  -H "Content-Type: application/json" \
  -d '{
    "title":"My First Post",
    "content":"<h1>Hello</h1>",
    "excerpt":"short text",
    "cover_image":"https://cdn.example.com/cover.jpg",
    "status":"draft",
    "meta_title":"My First Post",
    "meta_description":"meta"
  }'

## Reorder Navigation (drag-drop)
curl -X POST http://localhost:8080/admin/navigation/reorder \
  -H "Authorization: Bearer <TOKEN>" \
  -H "Content-Type: application/json" \
  -d '{"location":"header","ids":["<NAV_ID_1>","<NAV_ID_2>"]}'

## Upload File
curl -X POST http://localhost:8080/admin/upload \
  -H "Authorization: Bearer <TOKEN>" \
  -F "file=@./image.jpg"
