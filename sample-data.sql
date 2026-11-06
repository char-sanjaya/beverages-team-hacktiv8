-- tabel user
INSERT INTO `user` (`id`, `name`, `email`, `password`, `role`, `balance`, `created_at`) VALUES
(1, 'Admin Toko', 'admin@store.com', 'indonesia45', 'admin', 0.00, '2026-01-01 08:00:00'),
(2, 'Budi Santoso', 'budi@gmail.com', 'medankeren123', 'user', 1500000.00, '2026-01-02 10:30:00'),
(3, 'Siti Aminah', 'siti@gmail.com', 'merdeka1945', 'user', 500000.00, '2026-01-03 14:15:00');

-- tabel product
INSERT INTO `product` (`id`, `name`, `price`, `stock`, `created_at`) VALUES
(1, 'Kopi Americano', 45000.00, 10, '2026-01-01 09:00:00'),
(2, 'Cappuccino', 40000.00, 25, '2026-01-01 09:30:00'),
(3, 'Matcha', 30000.00, 15, '2026-01-01 10:00:00');

-- tabel order 
INSERT INTO `order` (`id`, `invoice_number`, `user_id`, `total_price`, `status`, `created_at`) VALUES
(1, 'INV/20260104/001', 2, 130000.00, 'success', '2026-01-04 11:20:00'),
(2, 'INV/20260105/002', 3, 30000.00, 'pending', '2026-01-05 15:45:00');

-- tabel order_detail
INSERT INTO `order_detail` (`id`, `order_id`, `product_id`, `quantity`, `price`) VALUES
-- Detail untuk Order ID 1 (Dibeli oleh Budi Santoso)
(1, 1, 1, 1, 45000.00), 
(2, 1, 2, 1, 40000.00), 

-- Detail untuk Order ID 2 (Dibeli oleh Siti Aminah)
(3, 2, 3, 1, 30000.00); 