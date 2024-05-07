CREATE OR REPLACE VIEW public.registrations_with_event
 AS
 SELECT registrations.id,
    registrations.name,
    registrations.surname,
    registrations.updated_at,
    registrations.created_at,
    registrations.deleted_at,
    registrations.token,
    registrations.gender,
    registrations.amount,
    registrations.payed,
    registrations.date_of_birth,
    registrations.finished_school,
    registrations.attended_previous,
    registrations.city,
    registrations.pills,
    registrations.notes,
    registrations.parent_name,
    registrations.parent_surname,
    registrations.email,
    registrations.phone,
    registrations.attended_activities,
    registrations.problems,
    registrations.admin_note,
    registrations.discount,
    registrations.promo_code,
    registrations.specific_symbol,
    registrations.notification_sent_at,
    events.title AS event_name,
    events.payment_reference
   FROM registrations
     LEFT JOIN signups ON registrations.id = signups.registration_id
     LEFT JOIN days ON signups.day_id = days.id
     LEFT JOIN events ON days.event_id = events.id;