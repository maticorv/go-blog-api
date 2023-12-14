# Golang Chi REST API

This is a simple REST API implemented in Golang using the Chi framework. It provides endpoints to perform CRUD operations on posts using the JSONPlaceholder API.

## Endpoints

### GET /v1/posts

Get a list of all posts.

#### Example using curl:

```bash
curl -X GET http://localhost:8080/v1/posts
```

### GET /v1/posts/{id}
Get details of a specific post by ID.

#### Example using curl:

```bash
curl -X GET http://localhost:8080/v1/posts/1
```

### POST /v1/posts
Create a new post.

#### Example using curl:

```bash
curl -X POST -H "Content-Type: application/json" -d '{"title": "New Post", "body": "This is a new post."}' http://localhost:8080/v1/posts
```

### PUT /v1/posts/{id}
Update an existing post by ID.

#### Example using curl:

```bash
curl -X PUT -H "Content-Type: application/json" -d '{"title": "Updated Post", "body": "This post has been updated."}' http://localhost:8080/v1/posts/1
```

### PATCH /v1/posts/{id}
Partially update an existing post by ID.

#### Example using curl:

```bash
curl -X PATCH -H "Content-Type: application/json" -d '{"title": "Updated Title"}' http://localhost:8080/v1/posts/1
```

### DELETE /v1/posts/{id}
Delete an existing post by ID.
```bash
curl -X DELETE http://localhost:8080/v1/posts/1
```



