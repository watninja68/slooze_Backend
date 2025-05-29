-- +goose Up
-- Seed Restaurants
INSERT INTO restaurants (name, address, country_id)
VALUES
  ('Bombay Bites', '123 Marine Drive, Mumbai', 0),
  ('Spice Villa',  '45 MG Road, Bangalore',   0),
  ('Liberty Diner','200 Broadway, New York',   1),
  ('Route 66 Grill','555 Sunset Blvd, LA',     1);

-- Seed Menu Items
INSERT INTO menu_items (restaurant_id, name, description, price)
VALUES
  ((SELECT id FROM restaurants WHERE name = 'Bombay Bites'), 'Paneer Tikka',    'Grilled cottage cheese cubes', 250.00),
  ((SELECT id FROM restaurants WHERE name = 'Bombay Bites'), 'Veg Biryani',      'Aromatic rice with veggies',   180.00),
  ((SELECT id FROM restaurants WHERE name = 'Spice Villa'),  'Masala Dosa',      'Rice crepe with potato masala', 120.00),
  ((SELECT id FROM restaurants WHERE name = 'Spice Villa'),  'Sambar Vada',      'Lentil donuts in sambar',       100.00),
  ((SELECT id FROM restaurants WHERE name = 'Liberty Diner'),'Cheeseburger',     'Beef patty with cheese',        8.50),
  ((SELECT id FROM restaurants WHERE name = 'Liberty Diner'),'Fries',            'Crispy potato fries',           3.00),
  ((SELECT id FROM restaurants WHERE name = 'Route 66 Grill'),'BBQ Ribs',         'Slow-cooked pork ribs',         15.00),
  ((SELECT id FROM restaurants WHERE name = 'Route 66 Grill'),'Coleslaw',         'Creamy cabbage salad',          4.00);

-- Seed Payment Methods
INSERT INTO payment_methods (user_id, method_type, details, is_default)
VALUES
  ((SELECT id FROM users WHERE email = 'nick@shield.com'),    'credit_card', '{"card":"**** **** **** 1111"}', TRUE),
  ((SELECT id FROM users WHERE email = 'carol@shield.com'),   'upi',         '{"vpa":"carol@upi"}',             TRUE),
  ((SELECT id FROM users WHERE email = 'steve@shield.com'),   'credit_card', '{"card":"**** **** **** 2222"}', TRUE),
  ((SELECT id FROM users WHERE email = 'thanos@shield.com'),  'upi',         '{"vpa":"thanos@upi"}',            TRUE),
  ((SELECT id FROM users WHERE email = 'thor@shield.com'),    'credit_card', '{"card":"**** **** **** 3333"}', TRUE),
  ((SELECT id FROM users WHERE email = 'travis@shield.com'),  'upi',         '{"vpa":"travis@upi"}',            TRUE);

-- Seed Orders
-- Captain Marvel places an order
INSERT INTO orders (user_id, restaurant_id, status, total_price)
VALUES
  (
    (SELECT id FROM users WHERE email = 'carol@shield.com'),
    (SELECT id FROM restaurants WHERE name = 'Bombay Bites'),
    'placed',
    430.00
  );

-- Thanos creates but then cancels an order
INSERT INTO orders (user_id, restaurant_id, status, total_price)
VALUES
  (
    (SELECT id FROM users WHERE email = 'thanos@shield.com'),
    (SELECT id FROM restaurants WHERE name = 'Spice Villa'),
    'cancelled',
    220.00
  );

-- Seed Order Items
-- Items for Carol’s order
INSERT INTO order_items (order_id, menu_item_id, quantity, price_at_order)
VALUES
  (
    (SELECT id FROM orders WHERE user_id = (SELECT id FROM users WHERE email = 'carol@shield.com') AND status='placed'),
    (SELECT id FROM menu_items WHERE name = 'Paneer Tikka'),
    1,
    250.00
  ),
  (
    (SELECT id FROM orders WHERE user_id = (SELECT id FROM users WHERE email = 'carol@shield.com') AND status='placed'),
    (SELECT id FROM menu_items WHERE name = 'Veg Biryani'),
    1,
    180.00
  );

-- Items for Thanos’s cancelled order
INSERT INTO order_items (order_id, menu_item_id, quantity, price_at_order)
VALUES
  (
    (SELECT id FROM orders WHERE user_id = (SELECT id FROM users WHERE email = 'thanos@shield.com') AND status='cancelled'),
    (SELECT id FROM menu_items WHERE name = 'Masala Dosa'),
    1,
    120.00
  ),
  (
    (SELECT id FROM orders WHERE user_id = (SELECT id FROM users WHERE email = 'thanos@shield.com') AND status='cancelled'),
    (SELECT id FROM menu_items WHERE name = 'Sambar Vada'),
    1,
    100.00
  );

-- Seed Payments
INSERT INTO payments (order_id, payment_method_id, amount, status, transaction_id)
VALUES
  (
    (SELECT id FROM orders WHERE user_id = (SELECT id FROM users WHERE email = 'carol@shield.com') AND status='placed'),
    (SELECT id FROM payment_methods WHERE user_id = (SELECT id FROM users WHERE email = 'carol@shield.com') AND is_default),
    430.00,
    'successful',
    'TXN12345CAROL'
  );

-- +goose Down
-- Remove payments
DELETE FROM payments WHERE transaction_id = 'TXN12345CAROL';

-- Remove order items
DELETE FROM order_items
 WHERE order_id IN (
    SELECT id FROM orders
    WHERE user_id = (SELECT id FROM users WHERE email = 'carol@shield.com') AND status='placed'
 )
 OR order_id IN (
    SELECT id FROM orders
    WHERE user_id = (SELECT id FROM users WHERE email = 'thanos@shield.com') AND status='cancelled'
 );

-- Remove orders
DELETE FROM orders
 WHERE user_id IN (
    SELECT id FROM users
    WHERE email IN ('carol@shield.com','thanos@shield.com')
 );

-- Remove payment methods
DELETE FROM payment_methods
 WHERE user_id IN (
    SELECT id FROM users
    WHERE email IN ('nick@shield.com','carol@shield.com','steve@shield.com','thanos@shield.com','thor@shield.com','travis@shield.com')
 );

-- Remove menu items
DELETE FROM menu_items
 WHERE restaurant_id IN (
    SELECT id FROM restaurants
    WHERE name IN ('Bombay Bites','Spice Villa','Liberty Diner','Route 66 Grill')
 );

-- Remove restaurants
DELETE FROM restaurants
 WHERE name IN ('Bombay Bites','Spice Villa','Liberty Diner','Route 66 Grill');
