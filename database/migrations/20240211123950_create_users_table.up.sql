CREATE TABLE users (
  id_user VARCHAR(255) PRIMARY KEY,
  name VARCHAR(255) NOT NULL,
  email VARCHAR(255) UNIQUE NOT NULL,
  email_verified_at TIMESTAMP,
  password VARCHAR(255) NOT NULL,
  password_string VARCHAR(255) NOT NULL,
  level ENUM('client', 'admin', 'owner') NOT NULL,
  remember_token VARCHAR(100),
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4;