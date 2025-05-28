-- +goose Up

-- 00001: Create Countries Table
CREATE TABLE countries (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE
);

-- 00002: Create Roles Table
CREATE TABLE roles (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL UNIQUE
);

-- 00003: Create Users Table
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    role_id INT NOT NULL,
    country_id INT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (role_id) REFERENCES roles(id),
    FOREIGN KEY (country_id) REFERENCES countries(id)
);

-- 00004: Create Restaurants Table
CREATE TABLE restaurants (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    address TEXT,
    country_id INT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (country_id) REFERENCES countries(id)
);

-- 00005: Create Menu Items Table
CREATE TABLE menu_items (
    id SERIAL PRIMARY KEY,
    restaurant_id INT NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    price DECIMAL(10, 2) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (restaurant_id) REFERENCES restaurants(id) ON DELETE CASCADE
);

-- 00006: Create Orders Table
CREATE TABLE orders (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    restaurant_id INT NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'created', -- (created, placed, cancelled)
    total_price DECIMAL(10, 2) NOT NULL DEFAULT 0.00,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (restaurant_id) REFERENCES restaurants(id)
);

-- 00007: Create Order Items Table
CREATE TABLE order_items (
    id SERIAL PRIMARY KEY,
    order_id INT NOT NULL,
    menu_item_id INT NOT NULL,
    quantity INT NOT NULL,
    price_at_order DECIMAL(10, 2) NOT NULL,
    FOREIGN KEY (order_id) REFERENCES orders(id) ON DELETE CASCADE,
    FOREIGN KEY (menu_item_id) REFERENCES menu_items(id)
);

-- 00008: Create Payment Methods Table
CREATE TABLE payment_methods (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    method_type VARCHAR(100) NOT NULL, -- (e.g., 'credit_card', 'upi')
    details TEXT NOT NULL, -- Could be encrypted or store a token
    is_default BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- 00009: Create Payments Table
CREATE TABLE payments (
    id SERIAL PRIMARY KEY,
    order_id INT NOT NULL,
    payment_method_id INT NOT NULL,
    amount DECIMAL(10, 2) NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'pending', -- (pending, successful, failed)
    transaction_id VARCHAR(255),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (order_id) REFERENCES orders(id),
    FOREIGN KEY (payment_method_id) REFERENCES payment_methods(id)
);

-- 00010: Seed Initial Data
-- Insert Countries
INSERT INTO countries (name) VALUES ('India'), ('America');

-- Insert Roles
INSERT INTO roles (name) VALUES ('ADMIN'), ('MANAGER'), ('MEMBER');

-- Insert Users (Replace 'hashed_pw' with actual hashed passwords)
INSERT INTO users (name, email, password_hash, role_id, country_id) VALUES
('Nick Fury', 'nick@shield.com', 'hashed_pw', (SELECT id FROM roles WHERE name = 'ADMIN'), (SELECT id FROM countries WHERE name = 'India')),
('Captain Marvel', 'carol@shield.com', 'hashed_pw', (SELECT id FROM roles WHERE name = 'MANAGER'), (SELECT id FROM countries WHERE name = 'India')),
('Captain America', 'steve@shield.com', 'hashed_pw', (SELECT id FROM roles WHERE name = 'MANAGER'), (SELECT id FROM countries WHERE name = 'America')),
('Thanos', 'thanos@shield.com', 'hashed_pw', (SELECT id FROM roles WHERE name = 'MEMBER'), (SELECT id FROM countries WHERE name = 'India')),
('Thor', 'thor@shield.com', 'hashed_pw', (SELECT id FROM roles WHERE name = 'MEMBER'), (SELECT id FROM countries WHERE name = 'India')),
('Travis', 'travis@shield.com', 'hashed_pw', (SELECT id FROM roles WHERE name = 'MEMBER'), (SELECT id FROM countries WHERE name = 'America'));

-- +goose Down
DROP TABLE IF EXISTS payments;
DROP TABLE IF EXISTS payment_methods;
DROP TABLE IF EXISTS order_items;
DROP TABLE IF EXISTS orders;
DROP TABLE IF EXISTS menu_items;
DROP TABLE IF EXISTS restaurants;
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS roles;
DROP TABLE IF EXISTS countries;
