CREATE TABLE IF NOT EXISTS expenses (
  id SERIAL PRIMARY KEY,
  title TEXT,
  amount FLOAT,
  note TEXT,
  tags TEXT[]
);

INSERT INTO expenses(id, title, amount, note, tags)
  VALUES
    (121,'watermelon',54,'Big C promotion discount 9 bath','{"food","beverage"}'),
    (123, 'iPhone 14 Pro Max 1TB', 66900, 'birthday gift from my love', '{"gadget"}')