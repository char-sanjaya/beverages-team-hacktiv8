create database beverages
use beverages

CREATE TABLE user (
  id int PRIMARY KEY AUTO_INCREMENT,
  name varchar(100) NOT NULL,
  email varchar(100) UNIQUE NOT NULL,
  password varchar(100) NOT NULL,
  role ENUM('admin', 'user') NOT NULL DEFAULT 'user',
  balance decimal(15,2) NOT NULL DEFAULT 0,
  created_at timestamp DEFAULT (now())
);

CREATE TABLE products (
  id int PRIMARY KEY AUTO_INCREMENT,
  name varchar(150) NOT NULL,
  price decimal(15,2) NOT NULL,
  stock int NOT NULL DEFAULT 0,
  created_at timestamp DEFAULT (now())
);

CREATE TABLE orders (
  id int PRIMARY KEY AUTO_INCREMENT,
  invoice_number varchar(50) UNIQUE NOT NULL,
  user_id int NOT NULL,
  total_price decimal(15,2) NOT NULL,
  status ENUM('pending', 'succes') NOT NULL DEFAULT 'pending',
  created_at timestamp DEFAULT (now()),

  foreign key (user_id) references user(id)
);

CREATE TABLE order_details (
  id int PRIMARY KEY AUTO_INCREMENT,
  order_id int NOT NULL,
  product_id int NOT NULL,
  quantity int NOT NULL,
  price decimal(15,2) NOT NULL,

  foreign key (order_id) references orders(id),
  foreign key (product_id) references product(id)
);
