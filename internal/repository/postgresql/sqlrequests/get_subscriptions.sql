SELECT services.name            AS service_name,
       u.uuid                   AS user_uuid,
       subscriptions.sum        AS sum,
       subscriptions.begin_date AS begin_date,
       subscriptions.end_date   AS end_date
FROM effectivemobile.subscriptions as subscriptions
         left join effectivemobile.services as services ON
    subscriptions.service_id = services.id
         left join effectivemobile.users as u ON
    subscriptions.user_id = u.id
order by subscriptions.begin_date desc, service_name, user_uuid
