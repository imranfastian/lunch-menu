-- Database initialization script for Docker
-- This runs automatically when PostgreSQL starts in Docker

-- Create tables (matching GORM structure)
CREATE TABLE IF NOT EXISTS restaurants (
    id SERIAL PRIMARY KEY,
    name VARCHAR(200) NOT NULL,
    description VARCHAR(1000),
    address VARCHAR(500),
    coordinate JSONB,
    homepage VARCHAR(500),
    region VARCHAR(100),
    phone VARCHAR(50),
    email VARCHAR(100),
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS menu_items (
    id SERIAL PRIMARY KEY,
    restaurant_id INTEGER NOT NULL REFERENCES restaurants(id),
    name VARCHAR(200) NOT NULL,
    description VARCHAR(1000),
    price DECIMAL(10,2),
    category VARCHAR(100),
    is_available BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_menu_items_restaurant_id ON menu_items(restaurant_id);
CREATE INDEX IF NOT EXISTS idx_restaurants_region ON restaurants(region);
CREATE INDEX IF NOT EXISTS idx_restaurants_is_active ON restaurants(is_active);
CREATE INDEX IF NOT EXISTS idx_menu_items_is_available ON menu_items(is_available);

-- Insert restaurants from original backend data
INSERT INTO restaurants (name, description, address, coordinate, homepage, region, phone, email, is_active, created_at, updated_at) VALUES
('Restaurang Bikupan', 'Hors restaurant in Uppsala', 'Uppsala', '[59.84220, 17.63403]', 'https://www.hors.se/uppsala/17/3/restaurang-bikupan/', 'Uppsala', '', '', true, NOW(), NOW()),
('Café Biomedicum', 'Hors café at Karolinska Institutet', 'Solna', '[59.34910, 18.02600]', 'https://ki.hors.se/17/5/bio-medicum/', 'Solna', '', '', true, NOW(), NOW()),
('Café Delta', 'Campus café in Solna', 'Solna', '[59.35039, 18.02265]', '', 'Solna', '', '', true, NOW(), NOW()),
('Sven Dufva', 'Traditional Swedish restaurant', 'Uppsala', '[59.84292, 17.64129]', 'http://svendufva.se/', 'Uppsala', '', '', true, NOW(), NOW()),
('Elma', 'Modern restaurant in Uppsala', 'Uppsala', '[59.83928, 17.63800]', 'https://www.elmauppsala.se/', 'Uppsala', '', '', true, NOW(), NOW()),
('Den Glada Restaurangen', 'The Happy Restaurant - traditional Swedish cuisine', 'Solna', '[59.35125, 18.02999]', 'http://www.dengladarestaurangen.se/', 'Solna', '', '', true, NOW(), NOW()),
('Haga gatukök', 'Street kitchen in Solna', 'Solna', '[59.34941, 18.02116]', '', 'Solna', '', '', true, NOW(), NOW()),
('Restaurang Hubben', 'Restaurant Hubben in Uppsala', 'Uppsala', '[59.84346, 17.64162]', 'https://vasakronan.foodbycoor.se/hubben', 'Uppsala', '', '', true, NOW(), NOW()),
('Jöns Jacob', 'Classic restaurant in Solna', 'Solna', '[59.34682, 18.02459]', 'http://gastrogate.com/restaurang/jonsjacob/', 'Solna', '', '', true, NOW(), NOW()),
('Café Erik Jorpes', 'Hors café at Karolinska Institutet', 'Solna', '[59.34849, 18.02719]', 'https://ki.hors.se/17/4/cafe-erik-jorpes/', 'Solna', '', '', true, NOW(), NOW()),
('Hotel von Kraemer', 'Hotel restaurant with lunch service', 'Uppsala', '[59.84827, 17.63534]', 'https://hotelvonkraemer.se/en/restaurant/lunch/', 'Uppsala', '', '', true, NOW(), NOW()),
('Livet Restaurant', 'Modern restaurant in Solna', 'Solna', '[59.34833, 18.03035]', 'https://www.livetbrand.com/har-finns-livet/livet-restaurant-solna/', 'Solna', '', '', true, NOW(), NOW()),
('Mai Thai Express', 'Thai and sushi restaurant', 'Solna', '[59.35093, 18.02430]', 'https://www.maethaiexpress.se/mae-thai-karolinska/', 'Solna', '', '', true, NOW(), NOW()),
('Restaurang Nanna Svartz', 'Hors restaurant at Karolinska Institutet', 'Solna', '[59.34867, 18.0278]', 'https://ki.hors.se/17/3/restaurang-nanna-svartz/', 'Solna', '', '', true, NOW(), NOW()),
('Restaurang Omni', 'Restaurant with daily lunch menu', 'Solna', '[59.34709, 18.03335]', 'https://restaurangomni.se/dagens-lunch/', 'Solna', '', '', true, NOW(), NOW()),
('Bistro Rudbeck', 'Hors bistro in Uppsala', 'Uppsala', '[59.84537, 17.63998]', 'https://www.hors.se/uppsala/17/10/bistro-rudbeck/', 'Uppsala', '', '', true, NOW(), NOW()),
('STHLM Street Lunch', 'Street food in Stockholm style', 'Solna', '[59.34747, 18.02653]', 'http://www.sthlmstreetlunch.com/', 'Solna', '', '', true, NOW(), NOW()),
('Svarta Räfven', 'Hors restaurant at Karolinska Institutet', 'Solna', '[59.34858, 18.0277]', 'https://ki.hors.se/17/7/svarta-rafven/', 'Solna', '', '', true, NOW(), NOW());

-- Insert sample menu items for some restaurants
INSERT INTO menu_items (restaurant_id, name, description, price, category, is_available, created_at, updated_at) VALUES
-- Den Glada Restaurangen (id 6)
(6, 'Dagens kött', 'Daily meat dish with potatoes and vegetables', 125.00, 'Main Course', true, NOW(), NOW()),
(6, 'Dagens fisk', 'Daily fish dish with seasonal sides', 135.00, 'Main Course', true, NOW(), NOW()),
(6, 'Vegetarisk tallrik', 'Vegetarian plate with seasonal vegetables', 115.00, 'Vegetarian', true, NOW(), NOW()),
(6, 'Soppa', 'Daily soup with bread', 95.00, 'Soup', true, NOW(), NOW()),

-- Mai Thai Express (id 13)
(13, 'Pad Thai', 'Traditional Thai stir-fried noodles', 145.00, 'Thai', true, NOW(), NOW()),
(13, 'Green Curry', 'Thai green curry with jasmine rice', 155.00, 'Thai', true, NOW(), NOW()),
(13, 'Sushi Mix', 'Mixed sushi plate with miso soup', 165.00, 'Sushi', true, NOW(), NOW()),
(13, 'Tom Yum Soup', 'Spicy Thai soup with shrimp', 125.00, 'Soup', true, NOW(), NOW()),

-- Restaurang Omni (id 15)
(15, 'Dagens lunch', 'Daily lunch special', 129.00, 'Main Course', true, NOW(), NOW()),
(15, 'Pasta of the day', 'Fresh pasta with daily sauce', 119.00, 'Pasta', true, NOW(), NOW()),
(15, 'Grilled chicken salad', 'Fresh salad with grilled chicken', 139.00, 'Salad', true, NOW(), NOW()),

-- Sven Dufva (id 4)
(4, 'Köttbullar', 'Traditional Swedish meatballs', 145.00, 'Swedish', true, NOW(), NOW()),
(4, 'Pannbiff', 'Pan-fried beef with onions', 155.00, 'Swedish', true, NOW(), NOW()),
(4, 'Fiskpudding', 'Traditional fish pudding', 135.00, 'Swedish', true, NOW(), NOW()),

-- STHLM Street Lunch (id 17)
(17, 'Gourmet Burger', 'Premium burger with fries', 149.00, 'Fast Food', true, NOW(), NOW()),
(17, 'Fish & Chips', 'Beer battered fish with chips', 139.00, 'Fast Food', true, NOW(), NOW()),
(17, 'Pulled Pork Sandwich', 'Slow-cooked pulled pork sandwich', 129.00, 'Fast Food', true, NOW(), NOW());

-- Log completion
\echo 'Database initialized with restaurant data successfully!'