-- Adminer 4.7.6 PostgreSQL dump

CREATE SEQUENCE owner_id_seq INCREMENT 1 MINVALUE 1 MAXVALUE 2147483647 START 1 CACHE 1;

CREATE TABLE "public"."owners" (
                                   "id" integer DEFAULT nextval('owner_id_seq') NOT NULL,
                                   "name" text NOT NULL,
                                   "surname" text NOT NULL,
                                   "username" text NOT NULL,
                                   "pass" text NOT NULL,
                                   CONSTRAINT "owner_id" PRIMARY KEY ("id")
) WITH (oids = false);


CREATE SEQUENCE events_id_seq INCREMENT 1 MINVALUE 1 MAXVALUE 2147483647 START 1 CACHE 1;

CREATE TABLE "public"."events" (
                                   "id" integer DEFAULT nextval('events_id_seq') NOT NULL,
                                   "name" text NOT NULL,
                                   "owner_id" integer NOT NULL,
                                   CONSTRAINT "events_id" PRIMARY KEY ("id"),
                                   CONSTRAINT "events_owner_id_fkey" FOREIGN KEY (owner_id) REFERENCES owners(id) ON UPDATE CASCADE ON DELETE RESTRICT NOT DEFERRABLE
) WITH (oids = false);

CREATE TABLE "public"."days" (
                                 "id" integer NOT NULL,
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
                                          "updated_at" timestamptz NOT NULL,
                                          "created_at" timestamptz NOT NULL,
                                          CONSTRAINT "registrations_id" PRIMARY KEY ("id")
) WITH (oids = false);


CREATE TABLE "public"."signups" (
                                    "id" integer NOT NULL,
                                    "day_id" integer NOT NULL,
                                    "registration_id" integer NOT NULL,
                                    "state" text NOT NULL,
                                    "updated_at" timestamp NOT NULL,
                                    "created_at" timestamp NOT NULL,
                                    CONSTRAINT "signups_id" PRIMARY KEY ("id"),
                                    CONSTRAINT "signups_day_id_fkey" FOREIGN KEY (day_id) REFERENCES days(id) ON UPDATE CASCADE ON DELETE CASCADE NOT DEFERRABLE,
                                    CONSTRAINT "signups_registration_id_fkey" FOREIGN KEY (registration_id) REFERENCES registrations(id) ON UPDATE CASCADE ON DELETE CASCADE NOT DEFERRABLE
) WITH (oids = false);


-- 2020-06-05 22:49:50.223171+00