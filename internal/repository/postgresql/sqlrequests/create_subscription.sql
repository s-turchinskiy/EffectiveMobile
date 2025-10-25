INSERT INTO effectivemobile.subscriptions (user_id, service_id, sum, begin_date, end_date)
VALUES ((select id from effectivemobile.users where uuid = $1),
        (select id from effectivemobile.services where name = $2),
        $3,
        $4,
        $5)