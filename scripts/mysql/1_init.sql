CREATE DATABASE IF NOT EXISTS binl;
USE binl;

CREATE TABLE IF NOT EXISTS metrics (
  id int not null AUTO_INCREMENT,
  `description` varchar(50) null,
  `value`       DOUBLE not null,
  created_at    timestamp not null default CURRENT_TIMESTAMP,

  PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS customers (
  id int not null AUTO_INCREMENT,
  `name` varchar(100) null,
  created_at    timestamp not null default CURRENT_TIMESTAMP,

  PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS customer_addresses (
  id int not null AUTO_INCREMENT,
  customer_id int not null,
  city varchar(100) not null,

  PRIMARY KEY (id, customer_id),
  INDEX idx_customer_id (customer_id),
    FOREIGN KEY (customer_id)
        REFERENCES customers(id)
        ON DELETE CASCADE
);
