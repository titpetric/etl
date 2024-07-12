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
