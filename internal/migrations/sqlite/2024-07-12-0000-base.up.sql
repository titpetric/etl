CREATE TABLE `commit` (
    `id` INTEGER PRIMARY KEY AUTOINCREMENT,
    `commit_id` TEXT NOT NULL,
    `repository` TEXT NOT NULL,
    `created_at` TIMESTAMP,
    UNIQUE(`commit_id`, `repository`)
);

CREATE TABLE `commit_output` (
    `id` INTEGER PRIMARY KEY AUTOINCREMENT,
    `commit_id` INTEGER NOT NULL,
    `created_with` TEXT,
    `filename` TEXT,
    `contents` TEXT,
    `created_at` DATETIME,
    UNIQUE(`commit_id`, `filename`)
);

CREATE TABLE `test_func` (
    `id` INTEGER PRIMARY KEY AUTOINCREMENT,
    `commit_id` INTEGER NOT NULL,
    `package_name` TEXT NOT NULL,
    `test_name` TEXT NOT NULL,
    `test_status` TEXT NOT NULL,
    `test_duration` FLOAT NOT NULL,
    `test_coverage` TEXT NOT NULL,
    UNIQUE(`test_name`, `package_name`)
);

CREATE TABLE `test_func_coverage` (
    `id` INTEGER PRIMARY KEY AUTOINCREMENT,
    `test_func_id` INTEGER NOT NULL,
    `symbol_name` TEXT NOT NULL,
    `coverage` INTEGER NOT NULL,
    UNIQUE(`test_func_id`, `symbol_name`)
);
