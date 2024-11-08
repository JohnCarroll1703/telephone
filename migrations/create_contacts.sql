CREATE TABLE contacts (
    contact_id INTEGER PRIMARY KEY ,
    phone_number VARCHAR(11),
);


CREATE TABLE users (
    id INTEGER PRIMARY KEY,
    name VARCHAR(100),
    email VARCHAR(100),
);

CREATE TABLE user_contacts(
    user_contacts_id INTEGER PRIMARY KEY,
    contact_id INTEGER FOREIGN KEY REFERENCES contacts(contact_id),
    is_fav BOOLEAN NOT NULL,
    id INTEGER FOREIGN KEY REFERENCES users(id)
);

SELECT * FROM user_contact_relation;
