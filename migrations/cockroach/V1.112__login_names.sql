WITH doa AS (
    SELECT l1.user_id, array_agg(l1.login_name)::STRING as login_names, l2.login_name as preferred_login_name
    FROM zitadel.projections.login_names l1
    JOIN (SELECT user_id, login_name FROM zitadel.projections.login_names WHERE is_primary) l2
        ON l1.user_id = l2.user_id
    WHERE l1.user_id in (
        SELECT id FROM notification.notify_users u
           LEFT JOIN zitadel.projections.login_names n
                ON n.user_id = u.id
        WHERE u.preferred_login_name <> n.login_name AND n.is_primary
    )
    GROUP BY l1.user_id, l2.login_name
)
UPDATE notification.notify_users SET preferred_login_name = doa.preferred_login_name, login_names = doa.login_names FROM doa WHERE doa.user_id = notification.notify_users.id;

WITH doa AS (
    SELECT l1.user_id, array_agg(l1.login_name) as login_names, l2.login_name as preferred_login_name
    FROM zitadel.projections.login_names l1
    JOIN (SELECT user_id, login_name FROM zitadel.projections.login_names WHERE is_primary) l2
      ON l1.user_id = l2.user_id
    WHERE l1.user_id in (
        SELECT id FROM auth.users u
            LEFT JOIN zitadel.projections.login_names n
                ON n.user_id = u.id
        WHERE u.preferred_login_name <> n.login_name AND n.is_primary
    )
    GROUP BY l1.user_id, l2.login_name
)
UPDATE auth.users SET preferred_login_name = doa.preferred_login_name, login_names = doa.login_names FROM doa WHERE doa.user_id = auth.users.id;

