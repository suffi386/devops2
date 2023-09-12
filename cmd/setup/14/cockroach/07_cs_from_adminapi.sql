INSERT INTO projections.current_states (
    projection_name
    , instance_id
    , event_date
    , "position"
    , last_updated
) SELECT 
    cs.view_name
    , cs.instance_id
    , e.creation_date
    , e.position
    , cs.last_successful_spooler_run
FROM 
    adminapi.current_sequences cs
JOIN eventstore.events e ON
    e.instance_id = cs.instance_id 
    AND e.event_sequence = cs.current_sequence 
    AND cs.current_sequence = (
        SELECT 
            MAX(cs2.current_sequence)
        FROM
            adminapi.current_sequences cs2
        WHERE
            cs.view_name = cs2.view_name
            AND cs.instance_id = cs2.instance_id
    )
;