CREATE DATABASE beverages;
USE beverages;

CREATE TABLE user (
  id INT PRIMARY KEY AUTO_INCREMENT,
  name VARCHAR(100) NOT NULL,
  email VARCHAR(100) UNIQUE NOT NULL,
  password VARCHAR(100) NOT NULL,
  role ENUM('admin', 'user') NOT NULL DEFAULT 'user',
  balance DECIMAL(15,2) NOT NULL DEFAULT 0,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE product (
  id INT PRIMARY KEY AUTO_INCREMENT,
  name VARCHAR(150) NOT NULL,
  price DECIMAL(15,2) NOT NULL,
  stock INT NOT NULL DEFAULT 0,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE `order` (
  id INT PRIMARY KEY AUTO_INCREMENT,
  invoice_number VARCHAR(50) UNIQUE NOT NULL,
  user_id INT NOT NULL,
  total_price DECIMAL(15,2) NOT NULL,
  status ENUM('pending', 'success') NOT NULL DEFAULT 'pending',
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (user_id) REFERENCES user(id)
);

CREATE TABLE order_detail (
  id INT PRIMARY KEY AUTO_INCREMENT,
  order_id INT NOT NULL,
  product_id INT NOT NULL,
  quantity INT NOT NULL,
  price DECIMAL(15,2) NOT NULL,
  FOREIGN KEY (order_id) REFERENCES `order`(id),
  FOREIGN KEY (product_id) REFERENCES product(id)
);


-- 1. Insert Data User
INSERT INTO user (name, email, password, role, balance) VALUES
('Admin Utama', 'admin@beverages.com', 'admin123', 'admin', 0.00),
('Budi Santoso', 'budi.santoso@gmail.com', 'budi123', 'user', 150000.00),
('Siti Aminah', 'siti.aminah@gmail.com', 'siti123', 'user', 75000.00);

-- 2. Insert Data Product (Minuman)
INSERT INTO product (name, price, stock) VALUES
('Kopi Susu Gula Aren', 18000.00, 50),
('Matcha Latte', 24000.00, 30),
('Teh Tarik Cincau', 15000.00, 40),
('Air Mineral Botol', 5000.00, 100);

-- 3. Insert Data Order (Transaksi)
-- Budi (user_id = 2) sukses membeli, Siti (user_id = 3) masih pending
INSERT INTO `order` (invoice_number, user_id, total_price, status) VALUES
('INV-20231025-001', 2, 42000.00, 'success'),
('INV-20231025-002', 3, 15000.00, 'pending');

-- 4. Insert Data Order Detail
-- Detail pesanan Budi (order_id = 1): Beli 1 Kopi Susu dan 1 Matcha Latte
INSERT INTO order_detail (order_id, product_id, quantity, price) VALUES
(1, 1, 1, 18000.00),
(1, 2, 1, 24000.00);

-- Detail pesanan Siti (order_id = 2): Beli 1 Teh Tarik Cincau
INSERT INTO order_detail (order_id, product_id, quantity, price) VALUES
(2, 3, 1, 15000.00);
