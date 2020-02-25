
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