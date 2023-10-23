DROP TABLE IF EXISTS products;
CREATE TABLE products (
  product_id INT AUTO_INCREMENT NOT NULL,
  user_id INT,
  product_name      VARCHAR(256) NOT NULL,
  product_description     TEXT(4096) NOT NULL,
  product_images JSON,
  product_price   DECIMAL(5,2) NOT NULL,
  compressed_product_images JSON,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`product_id`),
  FOREIGN KEY (user_id) REFERENCES users(user_id)
);


CREATE TABLE users (
  user_id INT AUTO_INCREMENT NOT NULL,
  user_name      VARCHAR(64) NOT NULL,
  mobile  VARCHAR(15),
  latitude FLOAT(24) DEFAULT 0.0,
  longitude FLOAT(24) DEFAULT 0.0,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`user_id`)
);

