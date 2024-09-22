CREATE TABLE IF NOT EXISTS guest (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    surname VARCHAR(255) NOT NULL,
    phone VARCHAR(255) NOT NULL,
    email VARCHAR(255),
    country VARCHAR(255)
);

CREATE UNIQUE INDEX IF NOT EXISTS i_guest_phone_unique ON guest (phone) WHERE length(phone) > 0;
CREATE UNIQUE INDEX IF NOT EXISTS i_guest_email_unique ON guest (email) WHERE length(email) > 0;