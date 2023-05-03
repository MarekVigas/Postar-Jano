-- Adminer 4.7.6 PostgreSQL dump

CREATE SEQUENCE owner_id_seq INCREMENT 1 MINVALUE 1 MAXVALUE 2147483647 START 1 CACHE 1;


CREATE TABLE "public"."owners" (
                                   "id" integer DEFAULT nextval('owner_id_seq') NOT NULL,
                                   "name" text NOT NULL,
                                   "surname" text NOT NULL,
                                   "username" text NOT NULL,
                                   "pass" text NOT NULL,
                                   "email" text NOT NULL,
                                   "phone" text NOT NULL,
                                   "photo" text NOT NULL,
                                   "gender" text NOT NULL,
                                   CONSTRAINT "owner_id" PRIMARY KEY ("id")
) WITH (oids = false);


CREATE SEQUENCE events_id_seq INCREMENT 1 MINVALUE 1 MAXVALUE 2147483647 START 1 CACHE 1;

CREATE TABLE "public"."events" (
                                   "id" integer DEFAULT nextval('events_id_seq') NOT NULL,
                                   "title" text NOT NULL,
                                   "owner_id" integer NOT NULL,
                                   "description" text NOT NULL,
                                   "date_from" text NOT NULL,
                                   "location" text NOT NULL,
                                   "min_age" integer NOT NULL,
                                   "max_age" integer NOT NULL,
                                   "info" text,
                                   "photo" text NOT NULL,
                                   "date_to" text NOT NULL,
                                   "time" text,
                                   "mail_info" text,
                                   price text NULL,
                                   "active" boolean DEFAULT false NOT NULL,
                                   "promo_registration" boolean DEFAULT false NOT NULL,
                                   CONSTRAINT "events_id" PRIMARY KEY ("id"),
                                   CONSTRAINT "events_owner_id_fkey" FOREIGN KEY (owner_id) REFERENCES owners(id) ON UPDATE CASCADE ON DELETE RESTRICT NOT DEFERRABLE
) WITH (oids = false);

CREATE SEQUENCE days_id_seq INCREMENT 1 MINVALUE 1 MAXVALUE 2147483647 START 1 CACHE 1;

CREATE TABLE "public"."days" (
                                 "id" integer DEFAULT nextval('days_id_seq') NOT NULL,
                                 "capacity" integer NOT NULL,
                                 "limit_boys" integer,
                                 "limit_girls" integer,
                                 "description" text NOT NULL,
                                 "price" integer NOT NULL,
                                 "event_id" integer NOT NULL,
                                 CONSTRAINT "days_id" PRIMARY KEY ("id"),
                                 CONSTRAINT "days_event_id_fkey" FOREIGN KEY (event_id) REFERENCES events(id) ON UPDATE CASCADE ON DELETE CASCADE NOT DEFERRABLE
) WITH (oids = false);



CREATE SEQUENCE children_id_seq INCREMENT 1 MINVALUE 1 MAXVALUE 2147483647 START 1 CACHE 1;

CREATE TABLE "public"."registrations" (
                                          "id" integer DEFAULT nextval('children_id_seq') NOT NULL,
                                          "name" text NOT NULL,
                                          "surname" text NOT NULL,
                                          "updated_at" timestamp NOT NULL,
                                          "created_at" timestamp NOT NULL,
                                          "deleted_at" timestamp,
                                          "token" text NOT NULL,
                                          "gender" text NOT NULL,
                                          "amount" integer NOT NULL,
                                          "payed" integer,
                                          "date_of_birth" date NOT NULL,
                                          "finished_school" text,
                                          "attended_previous" boolean NOT NULL,
                                          "city" text NOT NULL,
                                          "pills" text,
                                          "notes" text,
                                          "parent_name" text NOT NULL,
                                          "parent_surname" text NOT NULL,
                                          "email" text NOT NULL,
                                          "phone" text NOT NULL,
                                          "attended_activities" text,
                                          "problems" text,
                                          "admin_note" text DEFAULT '' NOT NULL,
                                          "discount" integer,
                                          "promo_code" text,
                                          CONSTRAINT "registrations_id" PRIMARY KEY ("id")
) WITH (oids = false);

CREATE SEQUENCE signups_id_seq INCREMENT 1 MINVALUE 1 MAXVALUE 2147483647 START 1 CACHE 1;

CREATE TABLE "public"."signups" (
                                    "id" integer DEFAULT nextval('signups_id_seq') NOT NULL,
                                    "day_id" integer NOT NULL,
                                    "registration_id" integer NOT NULL,
                                    "state" text NOT NULL,
                                    "updated_at" timestamp NOT NULL,
                                    "created_at" timestamp NOT NULL,
                                    "deleted_at" timestamp,
                                    UNIQUE (day_id, registration_id),
                                    CONSTRAINT "signups_id" PRIMARY KEY ("id"),
                                    CONSTRAINT "signups_day_id_fkey" FOREIGN KEY (day_id) REFERENCES days(id) ON UPDATE CASCADE ON DELETE CASCADE NOT DEFERRABLE,
                                    CONSTRAINT "signups_registration_id_fkey" FOREIGN KEY (registration_id) REFERENCES registrations(id) ON UPDATE CASCADE ON DELETE CASCADE NOT DEFERRABLE
) WITH (oids = false);


CREATE TABLE "public"."promo_codes" (
                                    "id"                     SERIAL PRIMARY KEY,
                                    "created_at"             TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                                    "updated_at"             TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                                    "email"                  TEXT NOT NULL,
                                    "key"                   TEXT NOT NULL UNIQUE,
                                    "available_registrations" integer NOT NULL
);

CREATE INDEX ix_registration_uuid ON "public"."promo_codes"("key");
-- 2020-06-05 22:49:50.223171+00
