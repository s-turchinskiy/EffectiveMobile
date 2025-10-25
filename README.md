Про тесты ничего не написано в ТЗ, так что их не добавлял

http://localhost:8080/api/subscription/create
Body:
{
"service_name" : "Yandex Plus",
"price": 400,
"user_id": "60601fee-2bf1-4721-ae6f-7636e79a0cbc",
"start_date": "03-2025",
"end_date": "04-2025"
}

http://localhost:8080/api/subscription/update
Body:
{
"id" : "1",
"service_name" : "Yandex Plus",
"price": 500,
"user_id": "60601fee-2bf1-4721-ae6f-7636e79a0cbc",
"start_date": "04-2025",
"end_date": "05-2025"
}

http://localhost:8080/api/subscription/read

http://localhost:8080/api/subscription/delete?id=18

http://localhost:8080/api/subscriptions/sum
Body:
{
"period": "03-2025",
"service_name" : "Yandex Plus2",
"user_id": "60601fee-2bf1-4721-ae6f-7636e79a0cbc"
}

Документация: http://localhost:8080/swagger/index.html