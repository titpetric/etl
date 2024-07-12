CREATE TABLE "commit" (
    "id" INTEGER PRIMARY KEY AUTOINCREMENT,
    "commit_id" TEXT NOT NULL,
    "repository" TEXT NOT NULL,
    "created_at" TIMESTAMP,
    UNIQUE("commit_id", "repository")
);
