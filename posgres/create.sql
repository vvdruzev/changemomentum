
DROP TABLE IF EXISTS Contacts;
create table Contacts
(
    id serial
        constraint contacts_pk
        primary key,
    Firstname varchar(50),
    lastname varchar(50)
);

insert into Contacts (Firstname, lastname) values
    ('ivan','ivanov'),
    ('petr','petrov');




DROP TABLE IF EXISTS Phonenumber;
create table Phonenumber
(
  contact_id  int         not null,
  phonenumber varchar(10) null,
  constraint Phonenumber_Contacts_id_fk
  foreign key (contact_id) references Contacts (id)
    on update cascade
    on delete cascade
);


insert into Phonenumber (contact_id, Phonenumber) VALUES
('1', '952000001'),
('1', '952000002'),
('2', '952000004');

-- SEQUENCE: public."Participants_id_seq"

-- DROP SEQUENCE public."Participants_id_seq";

CREATE SEQUENCE public."Participants_id_seq"
    INCREMENT 1
    START 1
    MINVALUE 1
    MAXVALUE 2147483647
    CACHE 1;


-- Table: public."Participants"

-- DROP TABLE public."Participants";

CREATE TABLE public."Participants"
(
    id integer NOT NULL DEFAULT nextval('"Participants_id_seq"'::regclass),
    "Firstname" character(20) COLLATE pg_catalog."default" NOT NULL,
    "Lastname" character(30) COLLATE pg_catalog."default" NOT NULL,
    command character(20) COLLATE pg_catalog."default",
    data_registration date NOT NULL,
    "UserstokenId" integer NOT NULL,
    CONSTRAINT "Participants_pkey" PRIMARY KEY (id)
);

-- SEQUENCE: public."UsersToken_id_seq"

-- DROP SEQUENCE public."UsersToken_id_seq";

CREATE SEQUENCE public."UsersToken_id_seq"
    INCREMENT 1
    START 1
    MINVALUE 1
    MAXVALUE 2147483647
    CACHE 1;

-- Table: public."UsersToken"

-- DROP TABLE public."UsersToken";

CREATE TABLE public."UsersToken"
(
    id integer NOT NULL DEFAULT nextval('"UsersToken_id_seq"'::regclass),
    login character(20) COLLATE pg_catalog."default" NOT NULL,
    "Firstname" character(20) COLLATE pg_catalog."default" NOT NULL,
    "Lastname" character(20) COLLATE pg_catalog."default" NOT NULL,
    email character(40) COLLATE pg_catalog."default",
    CONSTRAINT "UsersToken_pkey" PRIMARY KEY (id)
)

