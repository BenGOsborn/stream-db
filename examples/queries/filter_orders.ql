SELECT id, customer_id, status
FROM orders, customers
WHERE status = 'pending' AND total > 100 OR item_count < 5;
