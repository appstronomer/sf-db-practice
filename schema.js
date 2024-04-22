db.posts.drop()

db.posts.insertOne({
    ID: 0,
    Title: 'Статья',
    Content: 'Содержание статьи',
    CreatedAt: 0,
    PublishedAt: 0,
    AuthorID: 0,
    AuthorName: 'Дмитрий'
})
